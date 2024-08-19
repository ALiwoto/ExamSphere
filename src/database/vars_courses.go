package database

import (
	"time"

	"github.com/ALiwoto/ssg/ssg"
)

var (
	coursesInfoMap = func() *ssg.SafeEMap[int, CourseInfo] {
		m := ssg.NewSafeEMap[int, CourseInfo]()
		m.SetExpiration(time.Hour * 3)
		m.SetInterval(time.Hour * 12)
		m.EnableChecking()

		return m
	}()
)

var (
	valueCourseNotFound = &CourseInfo{}
)
