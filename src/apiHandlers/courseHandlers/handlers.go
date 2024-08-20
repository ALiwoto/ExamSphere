package courseHandlers

import (
	"ExamSphere/src/apiHandlers"
	"ExamSphere/src/database"

	"github.com/gofiber/fiber/v2"
)

// CreateCourseV1 godoc
// @Summary Create a new course
// @Description Allows a user to create a new course.
// @Tags Course
// @Accept json
// @Produce json
// @Param data body CreateCourseData true "Data needed to create a new course"
// @Param Authorization header string true "Authorization token"
// @Success 200 {object} apiHandlers.EndpointResponse{result=CreateCourseResult}
// @Router /api/v1/course/create [post]
func CreateCourseV1(c *fiber.Ctx) error {
	claimInfo := apiHandlers.GetJWTClaimsInfo(c)
	if claimInfo == nil {
		return apiHandlers.SendErrInvalidJWT(c)
	}

	userInfo := database.GetUserInfoByAuthHash(
		claimInfo.UserId, claimInfo.AuthHash,
	)
	if userInfo == nil {
		return apiHandlers.SendErrInvalidAuth(c)
	} else if !userInfo.CanCreateNewCourse() {
		return apiHandlers.SendErrPermissionDenied(c)
	}

	data := &CreateCourseData{}
	if err := c.BodyParser(data); err != nil {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	courseInfo, err := database.GetCourseByName(data.CourseName)
	if err == nil && courseInfo != nil {
		return apiHandlers.SendErrCourseAlreadyExists(c)
	}

	courseInfo, err = database.CreateNewCourse(&database.NewCourseData{
		CourseName:        data.CourseName,
		CourseDescription: data.CourseDescription,
		AddedBy:           userInfo.UserId,
	})
	if err != nil {
		return apiHandlers.SendErrInternalServerError(c)
	}

	return apiHandlers.SendResult(c, &CreateCourseResult{
		CourseId:          courseInfo.CourseId,
		CourseName:        courseInfo.CourseName,
		CourseDescription: courseInfo.CourseDescription,
		AddedBy:           courseInfo.AddedBy,
	})
}

// GetCourseInfoV1 godoc
// @Summary Get course information
// @Description Allows a user to get information about a course by its id.
// @Tags Course
// @Accept json
// @Produce json
// @Param id query int true "Course ID"
// @Param Authorization header string true "Authorization token"
// @Success 200 {object} apiHandlers.EndpointResponse{result=GetCourseInfoResult}
// @Router /api/v1/course/info [get]
func GetCourseInfoV1(c *fiber.Ctx) error {
	claimInfo := apiHandlers.GetJWTClaimsInfo(c)
	if claimInfo == nil {
		return apiHandlers.SendErrInvalidJWT(c)
	}

	userInfo := database.GetUserInfoByAuthHash(
		claimInfo.UserId, claimInfo.AuthHash,
	)
	if userInfo == nil {
		return apiHandlers.SendErrInvalidAuth(c)
	}

	courseId := c.QueryInt("id")

	courseInfo, err := database.GetCourseInfo(courseId)
	if err != nil {
		return apiHandlers.SendErrCourseNotFound(c)
	}

	return apiHandlers.SendResult(c, &GetCourseInfoResult{
		CourseId:          courseInfo.CourseId,
		CourseName:        courseInfo.CourseName,
		CourseDescription: courseInfo.CourseDescription,
		CreatedAt:         courseInfo.CreatedAt,
		AddedBy:           courseInfo.AddedBy,
	})
}

// SearchCourseV1 godoc
// @Summary Search for courses
// @Description Allows a user to search for courses by their name.
// @Tags Course
// @Accept json
// @Produce json
// @Param data body SearchCourseData true "Data needed to search for courses"
// @Param Authorization header string true "Authorization token"
// @Success 200 {object} apiHandlers.EndpointResponse{result=SearchCourseResult}
// @Router /api/v1/course/search [post]
func SearchCourseV1(c *fiber.Ctx) error {
	claimInfo := apiHandlers.GetJWTClaimsInfo(c)
	if claimInfo == nil {
		return apiHandlers.SendErrInvalidJWT(c)
	}

	userInfo := database.GetUserInfoByAuthHash(
		claimInfo.UserId, claimInfo.AuthHash,
	)
	if userInfo == nil {
		return apiHandlers.SendErrInvalidAuth(c)
	}

	data := &SearchCourseData{}
	if err := c.BodyParser(data); err != nil {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	courses, err := database.SearchCourseByName(data.CourseName)
	if err != nil {
		return apiHandlers.SendErrInternalServerError(c)
	}

	var coursesInfo []*SearchedCourseInfo
	for _, course := range courses {
		coursesInfo = append(coursesInfo, &SearchedCourseInfo{
			CourseId:          course.CourseId,
			CourseName:        course.CourseName,
			CourseDescription: course.CourseDescription,
			CreatedAt:         course.CreatedAt,
			AddedBy:           course.AddedBy,
		})
	}

	return apiHandlers.SendResult(c, &SearchCourseResult{
		Courses: coursesInfo,
	})
}
