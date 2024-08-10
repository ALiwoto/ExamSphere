package appValues

func (r UserRole) ToString() string {
	return string(r)
}

func (r UserRole) IsInvalid() bool {
	switch r {
	case UserRoleOwner, UserRoleAdmin, UserRoleStudent, UserRoleTeacher:
		return false
	default:
		return true
	}
}
