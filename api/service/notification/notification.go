package notification

import "context"

type Notification interface {
	Send(ctx context.Context, req *SendNotificationRequest) (*SendNotificationResponse, error)
}

type Type int32

const (
	Unspecified Type = iota
	Email
	Text
	Web
)

type SendNotificationRequest struct {
	// V1 Fields
	Type       Type
	Recipients []string
	Subject    string
	Body       string
}

type SendNotificationResponse struct {
	// V1 Fields
	Success bool
	Message string
}
