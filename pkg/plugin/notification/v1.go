package notification

import (
	"context"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/api/service/notification"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	grpcnotification1 "github.com/openkcm/plugin-sdk/proto/plugin/notification/v1"
)

type V1 struct {
	plugin.Facade
	grpcnotification1.NotificationServicePluginClient
}

func (v1 *V1) Version() uint {
	return 1
}

func (v1 *V1) ServiceInfo() api.Info {
	return v1.Info
}

func (v1 *V1) Send(ctx context.Context, req *notification.SendNotificationRequest) (*notification.SendNotificationResponse, error) {
	in := &grpcnotification1.SendNotificationRequest{
		NotificationType: grpcnotification1.NotificationType(req.Type),
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
