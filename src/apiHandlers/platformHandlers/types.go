package platformHandlers

import "ExamSphere/src/core/utils/logging"

type GetPlatformLogsResult struct {
	Logs []*logging.LogEntry `json:"logs"`
} // @name GetPlatformLogsResult
