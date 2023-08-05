package shared

import (
	"time"
)

type UserSignupData struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type ErrorResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Errors     interface{} `json:"errors"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type AuthData struct {
}

type Message struct {
	MsgType     string            `json:"msgType"`
	ClientID    string            `json:"clientID"`
	ClientCount int               `json:"clientCount"`
	ClientsList map[string]string `json:"clientsList"`
	Action      string            `json:"action"`
	Message     string            `json:"message"`
	Sender      string            `json:"sender"`
	Private     bool              `json:"private"`
	To          []string          `json:"to"`
	Mentioned   bool              `json:"mentioned"`
	Date        time.Time         `json:"date"`
}
