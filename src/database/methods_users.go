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

// CanSearchUser returns true if and only if the current user has
// the permission to search for a user.
func (i *UserInfo) CanSearchUser() bool {
	if i == nil || i.Role == appValues.UserRoleUnknown {
		// looks like an uninitialized user to me, just in case
		return false
	}

	return i.Role == appValues.UserRoleOwner || i.Role == appValues.UserRoleAdmin
}
