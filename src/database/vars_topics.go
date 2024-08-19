package database

import (
	"time"

	"github.com/ALiwoto/ssg/ssg"
)

var (
	topicsInfoMap = func() *ssg.SafeEMap[int, TopicInfo] {
		m := ssg.NewSafeEMap[int, TopicInfo]()
		m.SetExpiration(time.Hour * 3)
		m.SetInterval(time.Hour * 12)
		m.EnableChecking()

		return m
	}()
)

var (
	valueTopicNotFound = &TopicInfo{}
)
