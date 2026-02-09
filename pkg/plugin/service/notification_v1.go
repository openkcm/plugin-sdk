package service

import (
	"context"

	"github.com/openkcm/plugin-sdk/api/service/notification"
	"github.com/openkcm/plugin-sdk/pkg/catalog"
	notificationv1 "github.com/openkcm/plugin-sdk/proto/plugin/notification/v1"
)

var _ notification.Notification = (*hashicorpNotificationV1Plugin)(nil)

type hashicorpNotificationV1Plugin struct {
	plugin     catalog.Plugin
	grpcClient notificationv1.NotificationServiceClient
}

func NewNotificationV1Plugin(plugin catalog.Plugin) notification.Notification {
	return &hashicorpNotificationV1Plugin{
		plugin:     plugin,
		grpcClient: notificationv1.NewNotificationServiceClient(plugin.ClientConnection()),
	}
}

func (h *hashicorpNotificationV1Plugin) Send(ctx context.Context, req *notification.SendNotificationRequest) (*notification.SendNotificationResponse, error) {
	in := &notificationv1.SendNotificationRequest{
		NotificationType: notificationv1.NotificationType(req.Type),
		Recipients:       req.Recipients,
		Subject:          req.Subject,
		Body:             req.Body,
	}
	grpcResp, err := h.grpcClient.SendNotification(ctx, in)
	if err != nil {
		return nil, err
	}
	return &notification.SendNotificationResponse{
		Success: grpcResp.GetSuccess(),
		Message: grpcResp.GetMessage(),
	}, nil
}
