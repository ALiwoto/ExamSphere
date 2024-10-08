package database

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// CreateNewCourse creates a new course in the database using the
// plpgsql function create_course_info.
func CreateNewCourse(data *NewCourseData) (*CourseInfo, error) {
	info := &CourseInfo{
		CourseName:        strings.TrimSpace(data.CourseName),
		TopicId:           data.TopicId,
		CourseDescription: strings.TrimSpace(data.CourseDescription),
		AddedBy:           data.AddedBy,
		CreatedAt:         time.Now(),
	}

	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT create_course_info(
			p_course_name := $1,
			p_topic_id := $2,
			p_course_description := $3,
			p_added_by := $4
		)`,
		info.CourseName,
		info.TopicId,
		info.CourseDescription,
		info.AddedBy,
	).Scan(&info.CourseId)
	if err != nil {
		return nil, err
	}

	coursesInfoMap.Add(info.CourseId, info)

	return info, nil
}

// EditCourseInfo edits a course in the database.
func EditCourseInfo(data *EditCourseInfoData) (*CourseInfo, error) {
	_, err := DefaultContainer.db.Exec(context.Background(),
		`
		UPDATE course_info
			SET course_name = $1,
				course_description = $2,
				topic_id = $3
			WHERE course_id = $4;
		`,
		data.CourseName,
		data.CourseDescription,
		data.TopicId,
		data.CourseId,
	)
	if err != nil {
		return nil, err
	}

	info := coursesInfoMap.Get(data.CourseId)
	if info != nil {
		info.CourseName = data.CourseName
		info.CourseDescription = data.CourseDescription
		info.TopicId = data.TopicId
		return info, nil
	}

	return GetCourseInfo(data.CourseId)
}

// GetCourseInfo gets a course from the database.
func GetCourseInfo(courseId int) (*CourseInfo, error) {
	info := coursesInfoMap.Get(courseId)
	if info != nil && info != valueCourseNotFound && info.CourseId == courseId {
		return info, nil
	}

	info = &CourseInfo{}
	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT course_id, 
			course_name,
			course_description,
			topic_id,
			created_at,
			added_by
		FROM course_info WHERE course_id = $1`,
		courseId,
	).Scan(
		&info.CourseId,
		&info.CourseName,
		&info.CourseDescription,
		&info.TopicId,
		&info.CreatedAt,
		&info.AddedBy,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCourseNotFound
		}

		return nil, err
	}

	coursesInfoMap.Add(info.CourseId, info)
	coursesInfoByNameMap.Add(strings.ToLower(info.CourseName), info)
	return info, nil
}

// GetCourseByName gets a course from the database by its name.
func GetCourseByName(courseName string) (*CourseInfo, error) {
	courseName = strings.TrimSpace(strings.ToLower(courseName))
	info := coursesInfoByNameMap.Get(courseName)
	if info != nil && info != valueCourseNotFound &&
		strings.EqualFold(info.CourseName, courseName) {
		return info, nil
	}

	rows, err := DefaultContainer.db.Query(context.Background(),
		`SELECT course_id, 
			course_name,
			course_description,
			topic_id,
			created_at,
			added_by
		FROM course_info WHERE LOWER(course_name) = $1`,
		courseName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, ErrCourseNotFound
	}

	info = &CourseInfo{}
	err = rows.Scan(
		&info.CourseId,
		&info.CourseName,
		&info.CourseDescription,
		&info.TopicId,
		&info.CreatedAt,
		&info.AddedBy,
	)
	if err != nil {
		return nil, err
	}

	coursesInfoMap.Add(info.CourseId, info)
	coursesInfoByNameMap.Add(courseName, info)
	return info, nil
}

// SearchCourseByName searches for courses in the database.
func SearchCourseByName(data *SearchCourseByNameData) ([]*CourseInfo, error) {
	if data.Limit == 0 {
		data.Limit = DefaultPaginationLimit
	}

	rows, err := DefaultContainer.db.Query(context.Background(),
		`SELECT course_id,
			course_name,
			course_description,
			topic_id,
			created_at,
			added_by
		FROM course_info WHERE course_name ILIKE '%' || $1 || '%'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;`,
		data.CourseName,
		data.Limit,
		data.Offset,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*CourseInfo
	for rows.Next() {
		info := &CourseInfo{}
		err = rows.Scan(
			&info.CourseId,
			&info.CourseName,
			&info.CourseDescription,
			&info.TopicId,
			&info.CreatedAt,
			&info.AddedBy,
		)
		if err != nil {
			return nil, err
		}

		courses = append(courses, info)
	}

	return courses, nil
}

// GetCreatedCoursesByUser gets all courses created by a user.
func GetCreatedCoursesByUser(userId string) ([]*CourseInfo, error) {
	rows, err := DefaultContainer.db.Query(context.Background(),
		`SELECT course_id,
			course_name,
			course_description,
			topic_id,
			created_at,
			added_by
		FROM course_info WHERE added_by = $1
		ORDER BY created_at DESC;`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*CourseInfo
	for rows.Next() {
		info := &CourseInfo{}
		err = rows.Scan(
			&info.CourseId,
			&info.CourseName,
			&info.CourseDescription,
			&info.TopicId,
			&info.CreatedAt,
			&info.AddedBy,
		)
		if err != nil {
			return nil, err
		}

		courses = append(courses, info)
	}

	return courses, nil
}

// GetAllUserCourses calls user_courses view to get all courses
// a user has ever participated in.
func GetAllUserCourses(userId string) ([]*UserParticipatedCourse, error) {
	rows, err := DefaultContainer.db.Query(context.Background(),
		`SELECT course_id, course_name
		FROM user_courses 
		WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*UserParticipatedCourse
	for rows.Next() {
		currentCourse := &UserParticipatedCourse{}
		err = rows.Scan(
			&currentCourse.CourseId,
			&currentCourse.CourseName,
		)
		if err != nil {
			return nil, err
		}

		courses = append(courses, currentCourse)
	}

	return courses, nil
}

// GetAllParticipantsOfCourse gets all participants of a course.
func GetAllParticipantsOfCourse(courseId int) ([]*CourseParticipantInfo, error) {
	rows, err := DefaultContainer.db.Query(context.Background(),
		`SELECT user_id, full_name
		FROM course_participants
		WHERE course_id = $1`,
		courseId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []*CourseParticipantInfo
	for rows.Next() {
		participant := &CourseParticipantInfo{}
		err = rows.Scan(
			&participant.UserId,
			&participant.FullName,
		)
		if err != nil {
			return nil, err
		}

		participants = append(participants, participant)
	}

	return participants, nil
}
