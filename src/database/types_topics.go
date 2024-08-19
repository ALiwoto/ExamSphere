package database

import "time"

// TopicInfo is a struct that holds the information about a topic.
type TopicInfo struct {
	TopicId   int    `json:"topic_id"`
	TopicName string `json:"topic_name"`
}

// UserTopicStat is a struct that holds the user's statistics
// about a specific topic.
type UserTopicStat struct {
	UserId       string    `json:"user_id"`
	TopicId      int       `json:"topic_id"`
	CurrentExp   int       `json:"current_exp"`
	TotalExp     int       `json:"total_exp"`
	CurrentLevel int       `json:"current_level"`
	LastVisited  time.Time `json:"last_visited"`
}
