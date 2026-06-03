package router

import (
	"github.com/gin-gonic/gin"

	"training_eval_system/internal/handler"
	"training_eval_system/internal/middleware"
)

func Setup(
	auth *handler.AuthHandler,
	user *handler.UserHandler,
	course *handler.CourseHandler,
	task *handler.TaskHandler,
	submission *handler.SubmissionHandler,
	evaluation *handler.EvaluationHandler,
	upload *handler.UploadHandler,
	report *handler.ReportHandler,
	config *handler.ConfigHandler,
	verify *handler.VerifyHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	api := r.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", auth.Login)
			authGroup.POST("/register", auth.Register)
			authGroup.POST("/logout", auth.Logout)
		}

		api.POST("/upload", middleware.AuthRequired(), upload.Upload)

		users := api.Group("/users")
		users.Use(middleware.AuthRequired())
		{
			users.GET("/profile", user.GetProfile)
			users.GET("/list", middleware.RoleRequired("admin"), user.List)
			users.POST("", middleware.RoleRequired("admin"), user.Create)
			users.PUT("/:id", middleware.RoleRequired("admin"), user.Update)
			users.DELETE("/:id", middleware.RoleRequired("admin"), user.Delete)
		}

		courses := api.Group("/courses")
		courses.Use(middleware.AuthRequired())
		{
			courses.GET("", course.List)
			courses.GET("/my", middleware.RoleRequired("student"), course.MyCourses)
			courses.POST("/join", middleware.RoleRequired("student"), course.JoinByCode)
			courses.DELETE("/:id/leave", middleware.RoleRequired("student"), course.LeaveCourse)
			courses.GET("/:id", course.GetByID)
			courses.POST("", middleware.RoleRequired("admin", "teacher"), course.Create)
			courses.PUT("/:id", middleware.RoleRequired("admin", "teacher"), course.Update)
			courses.DELETE("/:id", middleware.RoleRequired("admin"), course.Delete)
		}

		tasks := api.Group("/tasks")
		tasks.Use(middleware.AuthRequired())
		{
			tasks.GET("", task.List)
			tasks.GET("/:id", task.GetByID)
			tasks.POST("", middleware.RoleRequired("admin", "teacher"), task.Create)
			tasks.PUT("/:id", middleware.RoleRequired("admin", "teacher"), task.Update)
			tasks.DELETE("/:id", middleware.RoleRequired("admin"), task.Delete)
		}

		submissions := api.Group("/submissions")
		submissions.Use(middleware.AuthRequired())
		{
			submissions.GET("", submission.List)
			submissions.GET("/:id", submission.GetByID)
			submissions.GET("/:id/download/:idx", submission.DownloadFile)
			submissions.POST("", middleware.RoleRequired("student"), submission.Create)
			submissions.POST("/resubmit", middleware.RoleRequired("student"), submission.Resubmit)
		}

		verifyGroup := api.Group("/verify")
		verifyGroup.Use(middleware.AuthRequired())
		{
			verifyGroup.GET("/:submissionId", verify.GetVerification)
			verifyGroup.POST("/:submissionId", middleware.RoleRequired("admin", "teacher"), verify.Verify)
		}

		eval := api.Group("/eval")
		eval.Use(middleware.AuthRequired())
		{
			eval.GET("/indicators", evaluation.ListIndicators)
			eval.POST("/indicators", middleware.RoleRequired("admin"), evaluation.CreateIndicator)
			eval.PUT("/indicators/:id", middleware.RoleRequired("admin"), evaluation.UpdateIndicator)
			eval.DELETE("/indicators/:id", middleware.RoleRequired("admin"), evaluation.DeleteIndicator)

			eval.GET("/task-indicators/:taskId", evaluation.GetTaskIndicators)
			eval.PUT("/task-indicators/:taskId", middleware.RoleRequired("admin", "teacher"), evaluation.SaveTaskIndicators)

			eval.GET("/course-indicators/:courseId", evaluation.GetCourseIndicators)
			eval.PUT("/course-indicators/:courseId", middleware.RoleRequired("admin", "teacher"), evaluation.SaveCourseIndicators)

			eval.GET("/score/:submissionId", evaluation.GetScores)
			eval.POST("/score/teacher", middleware.RoleRequired("admin", "teacher"), evaluation.SubmitTeacherScore)
			eval.POST("/score/llm/:submissionId", middleware.RoleRequired("admin", "teacher"), evaluation.ScoreByLLM)
		}

		reports := api.Group("/reports")
		reports.Use(middleware.AuthRequired())
		{
			reports.GET("", report.List)
			reports.POST("/generate", report.Generate)
			reports.GET("/export/:id", report.Export)
		}

		configApi := api.Group("/config")
		configApi.Use(middleware.AuthRequired(), middleware.RoleRequired("admin"))
		{
			configApi.GET("/llm", config.GetLlmConfig)
			configApi.PUT("/llm", config.UpdateLlmConfig)
		}
	}

	return r
}
