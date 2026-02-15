package notification

import (
	"context"

	"github.com/openkcm/plugin-sdk/api"
)

type Notification interface {
	ServiceInfo() api.Info

	Send(ctx context.Context, req *SendNotificationRequest) (*SendNotificationResponse, error)
}

// Enums translated to pure Go types

type DeliveryChannel int

const (
	DeliveryChannelUnspecified DeliveryChannel = iota
	DeliveryChannelEmail
	DeliveryChannelSMS
	DeliveryChannelPush
	DeliveryChannelInApp
)

func (d DeliveryChannel) String() string {
	switch d {
	case DeliveryChannelEmail:
		return "EMAIL"
	case DeliveryChannelSMS:
		return "SMS"
	case DeliveryChannelPush:
		return "PUSH"
	case DeliveryChannelInApp:
		return "IN_APP"
	default:
		return "UNSPECIFIED"
	}
}

// Domain Models

// Recipient uses pointers to represent the oneof field.
// Only one of these should be non-nil.
type Recipient struct {
	EmailAddress *string
	PhoneNumber  *string
	DeviceToken  *string
	UserID       *string
}

type RawMessage struct {
	Subject  string
	Body     string
	Metadata map[string]string
}

type TemplateMessage struct {
	TemplateID string
	Parameters map[string]string
}

// Content uses pointers for the oneof field.
type Content struct {
	Raw      *RawMessage
	Template *TemplateMessage
}

type SendNotificationRequest struct {
	Recipients       []Recipient
	Content          Content
	PreferredChannel DeliveryChannel
}

type DeliveryFailure struct {
	Recipient   Recipient
	ErrorReason string
}

type SendNotificationResponse struct {
	TrackingID      string
	PartialFailures []DeliveryFailure
}
