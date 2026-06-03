package service

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"

	"training_eval_system/config"
	"training_eval_system/internal/dto/request"
	"training_eval_system/internal/model"
	"training_eval_system/internal/repository"
)

type ReportService struct {
	reportRepo *repository.ReportRepo
	subRepo    *repository.SubmissionRepo
	evalRepo   *repository.EvalRepo
	courseRepo *repository.CourseRepo
	taskRepo   *repository.TaskRepo
	userRepo   *repository.UserRepo
}

func NewReportService(
	reportRepo *repository.ReportRepo,
	subRepo *repository.SubmissionRepo,
	evalRepo *repository.EvalRepo,
	courseRepo *repository.CourseRepo,
	taskRepo *repository.TaskRepo,
	userRepo *repository.UserRepo,
) *ReportService {
	return &ReportService{
		reportRepo: reportRepo,
		subRepo:    subRepo,
		evalRepo:   evalRepo,
		courseRepo: courseRepo,
		taskRepo:   taskRepo,
		userRepo:   userRepo,
	}
}

func (s *ReportService) List(req *request.ReportQuery) ([]model.Report, error) {
	return s.reportRepo.List(req.CourseID, req.TaskID)
}

func (s *ReportService) Generate(req *request.GenerateReportReq, userID uint) (*model.Report, error) {
	reportDir := filepath.Join(config.AppConfig.Upload.Dir, "reports")
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return nil, fmt.Errorf("create report dir: %w", err)
	}

	report := &model.Report{
		CourseID:    &req.CourseID,
		TaskID:      &req.TaskID,
		Type:        req.Type,
		Format:      req.Format,
		GeneratedBy: userID,
	}

	var err error
	switch {
	case req.Type == "individual" && req.Format == "pdf":
		err = s.generateIndividualPDF(report, reportDir)
	case req.Type == "class" && req.Format == "excel":
		err = s.generateClassExcel(report, reportDir)
	case req.Type == "course" && req.Format == "excel":
		err = s.generateCourseExcel(report, reportDir)
	default:
		return nil, fmt.Errorf("不支持的报告类型: %s/%s", req.Type, req.Format)
	}

	if err != nil {
		return nil, err
	}

	if err := s.reportRepo.Create(report); err != nil {
		return nil, err
	}

	return report, nil
}

func (s *ReportService) GetByID(id uint) (*model.Report, error) {
	return s.reportRepo.FindByID(id)
}

func (s *ReportService) GetExportPath(id uint) (string, string, error) {
	report, err := s.reportRepo.FindByID(id)
	if err != nil {
		return "", "", fmt.Errorf("报表不存在")
	}
	if report.FileURL == "" {
		return "", "", fmt.Errorf("报表文件未生成")
	}

	absPath := filepath.Join(config.AppConfig.Upload.Dir, filepath.Base(report.FileURL))
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("报表文件已丢失")
	}

	var contentType string
	switch report.Format {
	case "pdf":
		contentType = "application/pdf"
	case "excel":
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return absPath, contentType, nil
}

func (s *ReportService) generateIndividualPDF(report *model.Report, dir string) error {
	submissions, _ := s.subRepo.List(*report.TaskID, "")
	if len(submissions) == 0 {
		return fmt.Errorf("该任务下无提交记录")
	}

	task, err := s.taskRepo.FindByID(*report.TaskID)
	if err != nil {
		return fmt.Errorf("任务不存在")
	}
	course, err := s.courseRepo.FindByID(*report.CourseID)
	if err != nil {
		return fmt.Errorf("课程不存在")
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 20)

	for _, sub := range submissions {
		scores, _ := s.evalRepo.GetScores(sub.ID)
		vr, _ := s.evalRepo.GetVerificationBySubmission(sub.ID)

		pdf.AddPage()
		s.addPDFHeader(pdf, course.Name, task.Title)
		s.addPDFStudentInfo(pdf, sub, vr)
		s.addPDFScoresTable(pdf, scores)
		s.addPDFFooter(pdf)
	}

	fileName := fmt.Sprintf("individual_%d_%d.pdf", *report.TaskID, time.Now().UnixNano())
	filePath := filepath.Join(dir, fileName)
	if err := pdf.OutputFileAndClose(filePath); err != nil {
		return fmt.Errorf("生成PDF失败: %w", err)
	}

	report.FileURL = fmt.Sprintf("/uploads/reports/%s", fileName)
	report.Title = fmt.Sprintf("%s - %s 个人评价报告", course.Name, task.Title)
	return nil
}

func (s *ReportService) addPDFHeader(pdf *gofpdf.Fpdf, courseName, taskName string) {
	pdf.SetFont("Helvetica", "B", 16)
	pdf.CellFormat(190, 10, "软件实训评价报告", "", 1, "C", false, 0, "")
	pdf.SetFont("Helvetica", "", 11)
	pdf.CellFormat(190, 7, fmt.Sprintf("课程: %s", courseName), "", 1, "C", false, 0, "")
	pdf.CellFormat(190, 7, fmt.Sprintf("任务: %s", taskName), "", 1, "C", false, 0, "")
	pdf.Ln(4)
}

func (s *ReportService) addPDFStudentInfo(pdf *gofpdf.Fpdf, sub model.Submission, vr *model.VerificationResult) {
	passText := "未核查"
	if vr != nil && vr.OverallPass != nil {
		if *vr.OverallPass == 1 {
			passText = "通过"
		} else {
			passText = "不通过"
		}
	}

	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(190, 8, fmt.Sprintf("学生: %s  |  提交时间: %s  |  核查: %s",
		sub.StudentName, sub.SubmitTime.Format("2006-01-02 15:04"), passText), "", 1, "L", false, 0, "")
	pdf.Ln(4)
}

func (s *ReportService) addPDFScoresTable(pdf *gofpdf.Fpdf, scores []model.EvalScore) {
	headers := []string{"评价指标", "LLM评分", "教师评分", "最终得分", "评语"}
	widths := []float64{45, 25, 25, 25, 70}

	pdf.SetFont("Helvetica", "B", 10)
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Helvetica", "", 10)
	totalFinal := 0.0
	for _, sc := range scores {
		row := []string{
			sc.IndicatorName,
			formatScore(sc.LLMScore),
			formatScore(sc.TeacherScore),
			formatScore(sc.FinalScore),
			formatComment(sc.LLMComment, sc.TeacherComment),
		}
		for i, v := range row {
			pdf.CellFormat(widths[i], 7, truncateText(v, 40), "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
		if sc.FinalScore != nil {
			totalFinal += *sc.FinalScore
		}
	}

	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(120, 8, "加权总分", "1", 0, "R", false, 0, "")
	pdf.CellFormat(70, 8, fmt.Sprintf("%.1f", totalFinal), "1", 1, "C", false, 0, "")
}

func (s *ReportService) addPDFFooter(pdf *gofpdf.Fpdf) {
	pdf.Ln(5)
	pdf.SetFont("Helvetica", "I", 9)
	pdf.SetTextColor(128, 128, 128)
	pdf.CellFormat(190, 6, fmt.Sprintf("报告生成时间: %s  |  系统自动生成", time.Now().Format("2006-01-02 15:04:05")),
		"", 1, "C", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
}

func (s *ReportService) generateClassExcel(report *model.Report, dir string) error {
	return s.generateExcel(report, dir, *report.TaskID)
}

func (s *ReportService) generateCourseExcel(report *model.Report, dir string) error {
	return s.generateExcel(report, dir, 0)
}

func (s *ReportService) generateExcel(report *model.Report, dir string, taskID uint) error {
	var submissions []model.Submission
	var err error
	if taskID > 0 {
		submissions, err = s.subRepo.List(taskID, "")
	} else {
		submissions, err = s.subRepo.List(0, "")
	}
	if err != nil || len(submissions) == 0 {
		return fmt.Errorf("无提交记录")
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "实训成绩统计"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"序号", "学生姓名", "任务名称", "提交时间", "核查结果", "加权总分"}
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	dataStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	passStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Color: "008000"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	failStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Color: "FF0000"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	f.SetRowStyle(sheet, 1, 1, headerStyle)

	totalScore := 0.0
	count := 0

	for idx, sub := range submissions {
		row := idx + 2
		scores, _ := s.evalRepo.GetScores(sub.ID)
		vr, _ := s.evalRepo.GetVerificationBySubmission(sub.ID)

		finalTotal := 0.0
		for _, sc := range scores {
			if sc.FinalScore != nil {
				finalTotal += *sc.FinalScore
			}
		}

		passText := "未核查"
		passStylePtr := dataStyle
		if vr != nil && vr.OverallPass != nil {
			if *vr.OverallPass == 1 {
				passText = "通过"
				passStylePtr = passStyle
			} else {
				passText = "不通过"
				passStylePtr = failStyle
			}
		}

		values := []interface{}{
			idx + 1, sub.StudentName, sub.TaskTitle,
			sub.SubmitTime.Format("2006-01-02 15:04"),
			passText, math.Round(finalTotal*10) / 10,
		}

		for i, v := range values {
			cell, _ := excelize.CoordinatesToCellName(i+1, row)
			f.SetCellValue(sheet, cell, v)
			if i == 4 {
				f.SetCellStyle(sheet, cell, cell, passStylePtr)
			} else {
				f.SetCellStyle(sheet, cell, cell, dataStyle)
			}
		}

		totalScore += finalTotal
		count++
	}

	if count > 0 {
		summaryRow := count + 3
		avgScore := math.Round(totalScore/float64(count)*10) / 10
		f.SetCellValue(sheet, fmt.Sprintf("E%d", summaryRow), fmt.Sprintf("平均分: %.1f", avgScore))
	}

	f.SetColWidth(sheet, "A", "F", 18)

	prefix := "course"
	if taskID > 0 {
		prefix = fmt.Sprintf("task_%d", taskID)
	}
	fileName := fmt.Sprintf("%s_%d.xlsx", prefix, time.Now().UnixNano())
	filePath := filepath.Join(dir, fileName)
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("生成Excel失败: %w", err)
	}

	report.FileURL = fmt.Sprintf("/uploads/reports/%s", fileName)
	report.Title = "实训成绩统计报表"
	return nil
}

func formatScore(score *float64) string {
	if score == nil {
		return "-"
	}
	return fmt.Sprintf("%.1f", *score)
}

func formatComment(llmComment, teacherComment string) string {
	if teacherComment != "" {
		return teacherComment
	}
	if llmComment != "" {
		return llmComment
	}
	return "-"
}

func truncateText(s string, max int) string {
	runes := []rune(s)
	if len(runes) > max {
		return string(runes[:max]) + "..."
	}
	return s
}
