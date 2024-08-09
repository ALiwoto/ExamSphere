package main

import (
	"OnlineExams/src/core/appConfig"
	"OnlineExams/src/core/utils/logging"
)

func main() {
	f := logging.LoadLogger(true)
	if f != nil {
		defer f()
	}

	err := appConfig.LoadConfig()
	if err != nil {
		logging.Fatal("Error in loading config:", err)
		return
	}
}
