package database

import "time"

type CourseInfo struct {
	CourseId          int       `json:"course_id"`
	CourseName        string    `json:"course_name"`
	CourseDescription string    `json:"course_description"`
	CreatedAt         time.Time `json:"created_at"`
	AddedBy           string    `json:"added_by"`
}

type UserParticipatedCourse struct {
	CourseId   int    `json:"course_id"`
	CourseName string `json:"course_name"`
}

// CourseParticipantInfo is a minimal information about someone
// who has participated in a course.
type CourseParticipantInfo struct {
	UserId   string `json:"user_id"`
	FullName string `json:"full_name"`
}

type NewCourseData struct {
	CourseName        string `json:"course_name"`
	CourseDescription string `json:"course_description"`
	AddedBy           string `json:"added_by"`
}

type EditCourseInfoData struct {
	CourseId          int    `json:"course_id"`
	CourseName        string `json:"course_name"`
	CourseDescription string `json:"course_description"`
}
