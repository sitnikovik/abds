package user

import "time"

// User represents a user account in the EnergoSbyt system.
type User struct {
	// createdAt stores the timestamp when the user was created.
	createdAt time.Time

	// updatedAt stores the timestamp when the user was last updated.
	updatedAt time.Time

	// fio stores the full name of the user (family, first name, patronymic).
	fio string

	// phone stores the current phone number of the user.
	phone string

	// phonePrev stores the previous phone number of the user.
	phonePrev string

	// id stores the unique identifier of the user.
	id uint64
}

// NewUser creates a new User instance with the given parameters.
func NewUser(
	id uint64,
	fio string,
	phone string,
	phonePrev string,
	createdAt time.Time,
	updatedAt time.Time,
) User {
	return User{
		createdAt: createdAt,
		updatedAt: updatedAt,
		fio:       fio,
		phone:     phone,
		phonePrev: phonePrev,
		id:        id,
	}
}

// ID returns the unique identifier of the user.
func (u User) ID() uint64 {
	return u.id
}
