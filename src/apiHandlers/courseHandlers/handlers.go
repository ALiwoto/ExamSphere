package courseHandlers

import (
	"ExamSphere/src/apiHandlers"
	"ExamSphere/src/database"

	"github.com/gofiber/fiber/v2"
)

// CreateCourseV1 godoc
// @Summary Create a new course
// @Description Allows a user to create a new course.
// @ID createCourseV1
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
// @ID getCourseInfoV1
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
// @ID searchCourseV1
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

// GetCreatedCoursesV1 godoc
// @Summary Get created courses
// @Description Allows a user to get all courses created by a user.
// @ID getCreatedCoursesV1
// @Tags Course
// @Accept json
// @Produce json
// @Param data body GetCreatedCoursesData true "Data needed to get created courses"
// @Param Authorization header string true "Authorization token"
// @Success 200 {object} apiHandlers.EndpointResponse{result=GetCreatedCoursesResult}
// @Router /api/v1/course/createdCourses [post]
func GetCreatedCoursesV1(c *fiber.Ctx) error {
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

	data := &GetCreatedCoursesData{}
	if err := c.BodyParser(data); err != nil {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	if data.UserId == "" {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	courses, err := database.GetCreatedCoursesByUser(data.UserId)
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

	return apiHandlers.SendResult(c, &GetCreatedCoursesResult{
		Courses: coursesInfo,
	})
}

// GetUserCoursesV1 godoc
// @Summary Get user courses
// @Description Allows a user to get all courses participated by a user.
// @ID getUserCoursesV1
// @Tags Course
// @Accept json
// @Produce json
// @Param data body GetUserCoursesData true "Data needed to get user courses"
// @Param Authorization header string true "Authorization token"
// @Success 200 {object} apiHandlers.EndpointResponse{result=GetUserCoursesResult}
// @Router /api/v1/course/userCourses [post]
func GetUserCoursesV1(c *fiber.Ctx) error {
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

	data := &GetUserCoursesData{}
	if err := c.BodyParser(data); err != nil {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	if data.UserId == "" {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	courses, err := database.GetAllUserCourses(data.UserId)
	if err != nil {
		return apiHandlers.SendErrInternalServerError(c)
	}

	var coursesInfo []*UserParticipatedCourseInfo
	for _, course := range courses {
		coursesInfo = append(coursesInfo, &UserParticipatedCourseInfo{
			CourseId:   course.CourseId,
			CourseName: course.CourseName,
		})
	}

	return apiHandlers.SendResult(c, &GetUserCoursesResult{
		Courses: coursesInfo,
	})
}

// GetCourseParticipantsV1 godoc
// @Summary Get course participants
// @Description Allows a user to get all participants of a course.
// @ID getCourseParticipantsV1
// @Tags Course
// @Accept json
// @Produce json
// @Param data body GetCourseParticipantsData true "Data needed to get course participants"
// @Param Authorization header string true "Authorization token"
// @Success 200 {object} apiHandlers.EndpointResponse{result=GetCourseParticipantsResult}
// @Router /api/v1/course/courseParticipants [post]
func GetCourseParticipantsV1(c *fiber.Ctx) error {
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

	data := &GetCourseParticipantsData{}
	if err := c.BodyParser(data); err != nil {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	if data.CourseId == 0 {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	participants, err := database.GetAllParticipantsOfCourse(data.CourseId)
	if err != nil {
		return apiHandlers.SendErrInternalServerError(c)
	}

	var participantsInfo []*CourseParticipantInfo
	for _, participant := range participants {
		participantsInfo = append(participantsInfo, &CourseParticipantInfo{
			UserId:   participant.UserId,
			FullName: participant.FullName,
		})
	}

	return apiHandlers.SendResult(c, &GetCourseParticipantsResult{
		Participants: participantsInfo,
	})
}