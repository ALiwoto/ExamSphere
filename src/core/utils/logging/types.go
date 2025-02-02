package logging

type LogType string

type LogEntry struct {
	LogId   string  `json:"log_id"`
	Type    LogType `json:"log_type"`
	Date    string  `json:"date"`
	Message string  `json:"message"`
	Details string  `json:"details"`
}

type LogStorage interface {
	StoreLog(log *LogEntry) error
	GetAllLogEntries() ([]*LogEntry, error)
	GetLogsByType(logType LogType) ([]*LogEntry, error)
	GetLogById(logId string) (*LogEntry, error)
	DeleteLogById(logId string) error
	ClearLogs() error
	GetStorageName() string
}

type memoryLogStorage struct {
	logs []*LogEntry
}
