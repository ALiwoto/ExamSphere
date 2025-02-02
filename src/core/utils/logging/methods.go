package logging

import "github.com/ALiwoto/ssg/ssg"

func (m *memoryLogStorage) StoreLog(entry *LogEntry) error {
	if entry == nil {
		return ErrInvalidLogEntry
	}

	if entry.LogId == "" {
		entry.LogId = ssg.ToBase10(logIdGenerator.Next())
	}

	m.logs = append(m.logs, entry)
	return nil
}

func (m *memoryLogStorage) GetAllLogEntries() ([]*LogEntry, error) {
	return m.logs, nil
}

func (m *memoryLogStorage) GetLogsByType(logType LogType) ([]*LogEntry, error) {
	var logs []*LogEntry
	for _, log := range m.logs {
		if log.Type == logType {
			logs = append(logs, log)
		}
	}

	return logs, nil
}

func (m *memoryLogStorage) GetLogById(logId string) (*LogEntry, error) {
	for _, log := range m.logs {
		if log.LogId == logId {
			return log, nil
		}
	}

	return nil, ErrLogNotFound
}

func (m *memoryLogStorage) DeleteLogById(logId string) error {
	for i, log := range m.logs {
		if log.LogId == logId {
			m.logs = append(m.logs[:i], m.logs[i+1:]...)
			return nil
		}
	}

	return ErrLogNotFound
}

func (m *memoryLogStorage) GetStorageName() string {
	return "MemoryLogStorage"
}

func (m *memoryLogStorage) ClearLogs() error {
	m.logs = make([]*LogEntry, 0)
	return nil
}
