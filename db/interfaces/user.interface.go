package interfaces

import "github.com/google/uuid"

type UserInterface interface {
	GetMail() string
	GetName() string
	GetId() uuid.UUID
}
