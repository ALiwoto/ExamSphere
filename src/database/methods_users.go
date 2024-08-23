package database

import "ExamSphere/src/core/appValues"

// CanCreateRole returns true if and only if the current user has
// the permission to create a new user with the specified role.
func (i *UserInfo) CanCreateRole(targetRole appValues.UserRole) bool {
	if i == nil || i.Role == appValues.UserRoleUnknown {
		// looks like an uninitialized user to me, just in case
		return false
	} else if targetRole == appValues.UserRoleOwner ||
		targetRole == appValues.UserRoleUnknown {
		// there can be only 1 owner role; and it has
		// to be set in config file
		return false
	}

	switch i.Role {
	case appValues.UserRoleOwner:
		// owner can create any role
		return true
	case appValues.UserRoleAdmin:
		// admins can create any role below themselves
		return targetRole != appValues.UserRoleOwner && targetRole != appValues.UserRoleAdmin
	}

	return false
}

// CanEditUser returns true if and only if the current user has
// the permission to edit the specified user's information.
func (i *UserInfo) CanEditUser(targetRole appValues.UserRole) bool {
	if i == nil || i.Role == appValues.UserRoleUnknown {
		// looks like an uninitialized user to me, just in case
		return false
	} else if targetRole == appValues.UserRoleOwner {
		// owner's info can never ben edited. If we want to edit it
		// it should be done inside of the config file.
		return false
	}

	switch i.Role {
	case appValues.UserRoleOwner:
		// owner can edit any role
		return true
	case appValues.UserRoleAdmin:
		// admins can edit any role below themselves
		return targetRole != appValues.UserRoleOwner
	}

	return false
}

// CanSearchUser returns true if and only if the current user has
// the permission to search for a user.
func (i *UserInfo) CanSearchUser() bool {
	if i == nil || i.Role == appValues.UserRoleUnknown {
		// looks like an uninitialized user to me, just in case
		return false
	}

	return i.Role == appValues.UserRoleOwner || i.Role == appValues.UserRoleAdmin
}

// CanGetUserInfo returns true if and only if the current user has
// the permission to get information about the specified user.
func (i *UserInfo) CanGetUserInfo(targetUser *UserInfo) bool {
	if i == nil || i.Role == appValues.UserRoleUnknown {
		// looks like an uninitialized user to me, just in case
		return false
	} else if targetUser == nil || targetUser.Role == appValues.UserRoleUnknown {
		// looks like an uninitialized user to me, just in case
		return false
	}

	if i.Role == appValues.UserRoleOwner {
		// owner can get info of anyone
		return true
	} else if i.Role == appValues.UserRoleAdmin {
		// admins can get info of anyone below themselves
		return targetUser.Role != appValues.UserRoleOwner
	}

	return false
}

// CanBanUser returns true if and only if the current user has
// the permission to ban the specified user.
func (i *UserInfo) CanBanUser(targetUser *UserInfo) bool {
	if i == nil || i.Role == appValues.UserRoleUnknown {
		// looks like an uninitialized user to me, just in case
		return false
	} else if targetUser == nil || targetUser.Role == appValues.UserRoleUnknown {
		// looks like an uninitialized user to me, just in case
		return false
	}

	if i.Role == appValues.UserRoleOwner {
		// owner can ban anyone
		return true
	} else if i.Role == appValues.UserRoleAdmin {
		// admins can ban anyone below themselves
		return targetUser.Role != appValues.UserRoleOwner &&
			targetUser.Role != appValues.UserRoleAdmin
	}

	return false
}

// CanChangePassword returns true if and only if the current user has
// the permission to change the password of the specified user.
func (i *UserInfo) CanChangePassword(targetUser *UserInfo) bool {
	if i == nil || i.Role == appValues.UserRoleUnknown {
		// looks like an uninitialized user to me, just in case
		return false
	} else if targetUser == nil ||
		targetUser.Role == appValues.UserRoleUnknown ||
		targetUser.Role == appValues.UserRoleOwner {
		// looks like an uninitialized user to me, just in case
		// The 'owner' user can never have its password changed from
		// the application. It should be done inside of the config file.
		return false
	}

	if i.Role == appValues.UserRoleOwner {
		// owner can change anyone's password
		return true
	} else if i.Role == appValues.UserRoleAdmin {
		// admins can change anyone's password below themselves
		return targetUser.Role != appValues.UserRoleAdmin
	}

	return false
}

// CanCreateNewTopic returns true if and only if the current user has
// the permission to create a new topic.
func (i *UserInfo) CanCreateNewTopic() bool {
	if i == nil || i.Role == appValues.UserRoleUnknown {
		// looks like an uninitialized user to me, just in case
		return false
	}

	return i.Role == appValues.UserRoleOwner ||
		i.Role == appValues.UserRoleAdmin
}

// CanSearchTopic returns true if and only if the current user has
// the permission to search for a topic.
func (i *UserInfo) CanSearchTopic() bool {
	return i != nil && i.Role != appValues.UserRoleUnknown
}

// CanGetTopicInfo returns true if and only if the current user has
// the permission to get information about a topic.
// All users can get information about a topic, as long as they are
// authenticated (i.e. not unknown, guest user, etc...).
func (i *UserInfo) CanGetTopicInfo() bool {
	return i != nil && i.Role != appValues.UserRoleUnknown
}

// CanCreateNewCourse returns true if and only if the current user has
// the permission to create a new course.
func (i *UserInfo) CanCreateNewCourse() bool {
	return i != nil && i.Role == appValues.UserRoleAdmin ||
		i.Role == appValues.UserRoleOwner
}

// CanCreateNewExam returns true if and only if the current user has
// the permission to create a new exam.
// Owners, admins, and teachers can create new exams.
func (i *UserInfo) CanCreateNewExam() bool {
	return i != nil && i.Role == appValues.UserRoleAdmin ||
		i.Role == appValues.UserRoleOwner ||
		i.Role == appValues.UserRoleTeacher
}

// CanGetExamInfo returns true if and only if the current user has
// the permission to get information about an exam.
// All users can get information about an exam, as long as they are
// authenticated (i.e. not unknown, guest user, etc...).
func (i *UserInfo) CanGetExamInfo() bool {
	return i != nil && i.Role != appValues.UserRoleUnknown
}

// CanGetExamQuestions returns true if and only if the current user has
// the permission to get questions of an exam.
func (i *UserInfo) CanGetExamQuestions() bool {
	return i != nil && i.Role != appValues.UserRoleUnknown
}

// CanPeekExamQuestions returns true if and only if the current user has
// the permission to peek at the questions of an exam.
// *peeking* here means getting exam questions without participating in
// the itself exam.
func (i *UserInfo) CanPeekExamQuestions(createdBy string) bool {
	if i == nil || i.Role == appValues.UserRoleUnknown {
		// looks like an uninitialized user to me, just in case
		return false
	}

	if i.UserId == createdBy {
		return true
	}

	return i.Role == appValues.UserRoleOwner ||
		i.Role == appValues.UserRoleAdmin
}

// CanTryToScoreExam returns true if and only if the current user has
// the permission to *try* to score an exam.
func (i *UserInfo) CanTryToScoreExam() bool {
	if i == nil || i.Role == appValues.UserRoleUnknown {
		// looks like an uninitialized user to me, just in case
		return false
	}

	return i.Role == appValues.UserRoleOwner ||
		i.Role == appValues.UserRoleAdmin ||
		i.Role == appValues.UserRoleTeacher
}

// CanTryToScoreExam returns true if and only if the current user has
// the permission to forcefully insert score an exam.
func (i *UserInfo) CanForceScoreExam() bool {
	if i == nil || i.Role == appValues.UserRoleUnknown {
		// looks like an uninitialized user to me, just in case
		return false
	}

	return i.Role == appValues.UserRoleOwner ||
		i.Role == appValues.UserRoleAdmin ||
		i.Role == appValues.UserRoleTeacher
}

//---------------------------------------------------------

func (d *UpdateUserData) IsEmpty() bool {
	return d == nil ||
		(d.FullName == "" &&
			d.Email == "")
}
