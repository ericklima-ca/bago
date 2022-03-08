package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type loginPayload interface{}

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty"`
	Name     string             `json:"name,omitempty"`
	Lastname string             `json:"lastname,omitempty"`
	Email    string             `json:"email,omitempty"`
	Role     string             `json:"role,omitempty"`
	// password string
}

func TryAuthenticate(_ loginPayload) (User, bool) {
	return User{
		ID:       primitive.ObjectID{1},
		Name:     "Erick",
		Lastname: "Lima",
		Email:    "ericklima@email.com",
		Role:     "Log",
	}, true
}
