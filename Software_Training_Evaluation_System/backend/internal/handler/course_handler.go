package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"training_eval_system/internal/dto/request"
	"training_eval_system/internal/middleware"
	"training_eval_system/internal/model"
	"training_eval_system/internal/service"
	"training_eval_system/pkg/response"
)

type CourseHandler struct {
	courseService *service.CourseService
}

func NewCourseHandler(courseService *service.CourseService) *CourseHandler {
	return &CourseHandler{courseService: courseService}
}

func (h *CourseHandler) List(c *gin.Context) {
	var req request.CourseQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		req = request.CourseQuery{}
	}
	var courses []model.Course
	var err error
	if middleware.IsTeacher(c) {
		courses, err = h.courseService.ListByTeacher(middleware.GetUserID(c), &req)
	} else {
		courses, err = h.courseService.List(&req)
	}
	if err != nil {
		response.InternalError(c, "获取课程列表失败")
		return
	}
	response.Success(c, courses)
}

func (h *CourseHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的课程ID")
		return
	}
	course, err := h.courseService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "课程不存在")
		return
	}
	response.Success(c, course)
}

func (h *CourseHandler) Create(c *gin.Context) {
	var req request.CreateCourseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请检查输入信息")
		return
	}
	teacherID := middleware.GetUserID(c)
	if err := h.courseService.Create(&req, teacherID); err != nil {
		response.InternalError(c, "创建课程失败")
		return
	}
	response.Created(c, nil)
}

func (h *CourseHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的课程ID")
		return
	}
	var req request.UpdateCourseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请检查输入信息")
		return
	}
	if err := h.courseService.Update(uint(id), &req); err != nil {
		response.InternalError(c, "更新失败")
		return
	}
	response.Success(c, nil)
}

func (h *CourseHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的课程ID")
		return
	}
	if err := h.courseService.Delete(uint(id)); err != nil {
		response.InternalError(c, "删除失败")
		return
	}
	response.Success(c, nil)
}

func (h *CourseHandler) JoinByCode(c *gin.Context) {
	var req request.JoinCourseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请输入课程代码")
		return
	}
	studentID := middleware.GetUserID(c)
	if err := h.courseService.JoinByCode(req.Code, studentID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "加入课程成功")
}

func (h *CourseHandler) MyCourses(c *gin.Context) {
	studentID := middleware.GetUserID(c)
	courses, err := h.courseService.GetMyCourses(studentID)
	if err != nil {
		response.InternalError(c, "获取课程列表失败")
		return
	}
	response.Success(c, courses)
}

func (h *CourseHandler) LeaveCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的课程ID")
		return
	}
	studentID := middleware.GetUserID(c)
	if err := h.courseService.LeaveCourse(uint(id), studentID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "已退出课程")
}
