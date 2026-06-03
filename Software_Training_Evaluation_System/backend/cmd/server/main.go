package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"training_eval_system/config"
	"training_eval_system/internal/handler"
	"training_eval_system/internal/model"
	"training_eval_system/internal/repository"
	"training_eval_system/internal/router"
	"training_eval_system/internal/service"
)

func main() {
	if err := config.InitConfig("config/config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	gin.SetMode(config.AppConfig.Server.Mode)

	dsn := config.AppConfig.Database.DSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Course{},
		&model.Task{},
		&model.Submission{},
		&model.EvalIndicator{},
		&model.TaskIndicator{},
		&model.CourseIndicator{},
		&model.EvalScore{},
		&model.VerificationResult{},
		&model.Report{},
		&model.SystemConfig{},
		&model.CourseEnrollment{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	userRepo := repository.NewUserRepo(db)
	courseRepo := repository.NewCourseRepo(db)
	taskRepo := repository.NewTaskRepo(db)
	subRepo := repository.NewSubmissionRepo(db)
	evalRepo := repository.NewEvalRepo(db)
	reportRepo := repository.NewReportRepo(db)
	configRepo := repository.NewConfigRepo(db)

	if dbCfg, err := configRepo.GetByKey("llm_provider"); err == nil {
		var llmCfg config.LLMConfig
		if json.Unmarshal([]byte(dbCfg.ConfigValue), &llmCfg) == nil {
			config.SetLLMConfig(llmCfg)
			log.Println("LLM config loaded from database")
		}
	}

	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	courseService := service.NewCourseService(courseRepo)
	taskService := service.NewTaskService(taskRepo)
	submissionService := service.NewSubmissionService(subRepo, evalRepo, taskRepo)
	fileService := service.NewFileService()
	llmService := service.NewLLMService()
	evalService := service.NewEvaluationService(evalRepo, subRepo, taskRepo, llmService, fileService)
	reportService := service.NewReportService(reportRepo, subRepo, evalRepo, courseRepo, taskRepo, userRepo)
	verifyService := service.NewVerifyService(subRepo, taskRepo, evalRepo, llmService, fileService)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	courseHandler := handler.NewCourseHandler(courseService)
	taskHandler := handler.NewTaskHandler(taskService)
	submissionHandler := handler.NewSubmissionHandler(submissionService, fileService)
	evalHandler := handler.NewEvaluationHandler(evalService)
	uploadHandler := handler.NewUploadHandler(fileService)
	reportHandler := handler.NewReportHandler(reportService)
	configHandler := handler.NewConfigHandler(configRepo)
	verifyHandler := handler.NewVerifyHandler(verifyService)

	r := router.Setup(
		authHandler,
		userHandler,
		courseHandler,
		taskHandler,
		submissionHandler,
		evalHandler,
		uploadHandler,
		reportHandler,
		configHandler,
		verifyHandler,
	)

	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
