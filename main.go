package main

import (
	"OnlineExams/src/core/appConfig"
	"OnlineExams/src/core/utils/logging"
	"OnlineExams/src/masterServer"
)

// @title           ExamSphere API
// @version         1.0
// @description     This is the API for the ExamSphere system
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

	err = masterServer.RunServer()
	if err != nil {
		logging.Fatal("Error in running server:", err)
	}
}
