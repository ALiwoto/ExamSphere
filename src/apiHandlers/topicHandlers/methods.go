package topicHandlers

func (d *CreateNewTopicData) IsValid() bool {
	if len(d.TopicName) < MinTopicNameLength ||
		len(d.TopicName) > MaxTopicNameLength {
		return false
	}

	return true
}
