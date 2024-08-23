package courseHandlers

import "time"

type CreateCourseData struct {
	CourseName        string `json:"course_name"`
	TopicId           int    `json:"topic_id"`
	CourseDescription string `json:"course_description"`
} // @name CreateCourseData

type CreateCourseResult struct {
	CourseId          int       `json:"course_id"`
	TopicId           int       `json:"topic_id"`
	CourseName        string    `json:"course_name"`
	CourseDescription string    `json:"course_description"`
	AddedBy           string    `json:"added_by"`
	CreatedAt         time.Time `json:"created_at"`
} // @name CreateCourseResult

type EditCourseData struct {
	CourseId          int    `json:"course_id"`
	TopicId           int    `json:"topic_id"`
	CourseName        string `json:"course_name"`
	CourseDescription string `json:"course_description"`
} // @name EditCourseData

type EditCourseResult struct {
	CourseId          int       `json:"course_id"`
	CourseName        string    `json:"course_name"`
	CourseDescription string    `json:"course_description"`
	AddedBy           string    `json:"added_by"`
	CreatedAt         time.Time `json:"created_at"`
} // @name EditCourseResult

type GetCourseInfoResult struct {
	CourseId          int       `json:"course_id"`
	CourseName        string    `json:"course_name"`
	CourseDescription string    `json:"course_description"`
	TopicId           int       `json:"topic_id"`
	AddedBy           string    `json:"added_by"`
	CreatedAt         time.Time `json:"created_at"`
} // @name GetCourseInfoResult

type SearchCourseData struct {
	CourseName string `json:"course_name"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
} // @name SearchCourseData

type SearchCourseResult struct {
	Courses []*SearchedCourseInfo `json:"courses"`
} // @name SearchCourseResult

type SearchedCourseInfo struct {
	CourseId          int       `json:"course_id"`
	CourseName        string    `json:"course_name"`
	CourseDescription string    `json:"course_description"`
	TopicId           int       `json:"topic_id"`
	CreatedAt         time.Time `json:"created_at"`
	AddedBy           string    `json:"added_by"`
} // @name SearchedCourseInfo

type GetCreatedCoursesData struct {
	UserId string `json:"user_id"`
} // @name GetCreatedCoursesData

type GetCreatedCoursesResult struct {
	Courses []*SearchedCourseInfo `json:"courses"`
} // @name GetCreatedCoursesResult

type GetUserCoursesData struct {
	UserId string `json:"user_id"`
} // @name GetUserCoursesData

type GetUserCoursesResult struct {
	Courses []*UserParticipatedCourseInfo `json:"courses"`
} // @name GetUserCoursesResult

type UserParticipatedCourseInfo struct {
	CourseId   int    `json:"course_id"`
	CourseName string `json:"course_name"`
} // @name UserParticipatedCourse

type GetCourseParticipantsData struct {
	CourseId int `json:"course_id"`
} // @name GetCourseParticipantsData

type GetCourseParticipantsResult struct {
	Participants []*CourseParticipantInfo `json:"participants"`
} // @name GetCourseParticipantsResult

type CourseParticipantInfo struct {
	UserId   string `json:"user_id"`
	FullName string `json:"full_name"`
} // @name CourseParticipantInfo
