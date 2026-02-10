package notification

import (
	"context"

	"github.com/openkcm/plugin-sdk/api/service/notification"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	notification1 "github.com/openkcm/plugin-sdk/proto/plugin/notification/v1"
)

type V1 struct {
	plugin.Facade
	notification1.NotificationServicePluginClient
}

func (v1 *V1) Send(ctx context.Context, req *notification.SendNotificationRequest) (*notification.SendNotificationResponse, error) {
	in := &notification1.SendNotificationRequest{
		NotificationType: notification1.NotificationType(req.Type),
		Recipients:       req.Recipients,
		Subject:          req.Subject,
		Body:             req.Body,
	}
	grpcResp, err := v1.SendNotification(ctx, in)
	if err != nil {
		return nil, err
	}
	return &notification.SendNotificationResponse{
		Success: grpcResp.GetSuccess(),
		Message: grpcResp.GetMessage(),
	}, nil
}
