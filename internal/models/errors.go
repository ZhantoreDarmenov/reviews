package models

import "errors"

var (
	ErrInvalidVerificationCode = errors.New("invalid verification code")
	ErrWorkoutProgramNotFound  = errors.New("workout program not found")
	ErrExerciseNotFound        = errors.New("exercise not found")
	ErrFoodNotFound            = errors.New("food not found")

	ErrDayNotFound    = errors.New("day not found")
	ErrInviteNotFound = errors.New("invite not found")
)
