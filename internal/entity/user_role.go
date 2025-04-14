package entity

import (
	"errors"
)

type UserRole string

const (
	UserRoleEmployee  UserRole = "employee"
	UserRoleModerator UserRole = "moderator"
)

var validRolesMap = map[UserRole]struct{}{
	UserRoleEmployee:  {},
	UserRoleModerator: {},
}

func (r UserRole) IsValidRole() bool {
	_, exists := validRolesMap[r]
	return exists
}

func (r UserRole) ValidateRole() error {
	if !r.IsValidRole() {
		return errors.New("invalid user role")
	}
	return nil
}
