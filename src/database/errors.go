package database

import "errors"

var (
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrInternalDatabaseError = errors.New("internal database error")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrOperationNotAllowed   = errors.New("operation not allowed")
	ErrCourseNotFound        = errors.New("course not found")
	ErrTopicNotFound         = errors.New("topic not found")
	ErrUserTopicStatNotFound = errors.New("user topic stat not found")
	ErrExamNotFound          = errors.New("exam not found")
	ErrExamQuestionNotFound  = errors.New("exam question not found")
	ErrGivenExamNotFound     = errors.New("given exam not found")
)
