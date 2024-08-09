package database

import (
	"time"

	"github.com/ALiwoto/ssg/ssg"
)

var (
	usersInfoMap = func() *ssg.SafeEMap[string, UserInfo] {
		m := ssg.NewSafeEMap[string, UserInfo]()
		m.SetExpiration(time.Hour * 3)
		m.SetInterval(time.Hour * 12)
		m.EnableChecking()

		return m
	}()
)
