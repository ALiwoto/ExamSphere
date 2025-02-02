package database

import "errors"

var (
	ErrUserAlreadyExists                         = errors.New("user already exists")
	ErrInternalDatabaseError                     = errors.New("internal database error")
	ErrUserNotFound                              = errors.New("user not found")
	ErrInvalidPassword                           = errors.New("invalid password")
	ErrOperationNotAllowed                       = errors.New("operation not allowed")
	ErrCourseNotFound                            = errors.New("course not found")
	ErrTopicNotFound                             = errors.New("topic not found")
	ErrUserTopicStatNotFound                     = errors.New("user topic stat not found")
	ErrExamNotFound                              = errors.New("exam not found")
	ErrExamQuestionTimeExceeded                  = errors.New("exam question time exceeded")
	ErrInvalidParticipantUserId                  = errors.New("invalid participant user id")
	ErrExamQuestionNotFound                      = errors.New("exam question not found")
	ErrGivenExamNotFound                         = errors.New("given exam not found")
	ErrGivenAnswerNotFound                       = errors.New("given answer not found")
	ErrInvalidAnswer                             = errors.New("invalid answer")
	ErrSampleExamCannotHavePointer               = errors.New("sample exam cannot have pointer")
	ErrInvalidPointerToExamId                    = errors.New("invalid pointer to exam id")
	ErrPointedExamMustBeSample                   = errors.New("pointed exam must be a sample exam")
	ErrNonPointerQuestionCannotHavePointerCount  = errors.New("non-pointer question cannot have pointer count")
	ErrNonPointerQuestionCannotHavePointedToExam = errors.New("non-pointer question cannot have pointed to exam")
	ErrInvalidPointerCount                       = errors.New("invalid pointer count")
)
