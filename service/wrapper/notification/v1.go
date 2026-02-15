package notification

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	grpcnotification1 "github.com/openkcm/plugin-sdk/proto/plugin/notification/v1"
	"github.com/openkcm/plugin-sdk/service/api/notification"
)

type V1 struct {
	plugin.Facade
	grpcnotification1.NotificationPluginClient
}

func (v1 *V1) Version() uint {
	return 1
}

func (v1 *V1) ServiceInfo() api.Info {
	return v1.Info
}

func (v1 *V1) Send(ctx context.Context, req *notification.SendNotificationRequest) (*notification.SendNotificationResponse, error) {
	pbReq := mapNotificationRequestToProto(req)

	if err := protovalidate.Validate(pbReq); err != nil {
		return nil, fmt.Errorf("failed validation: %v", err)
	}

	pbResp, err := v1.NotificationPluginClient.Send(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	return mapNotificationResponseToDomain(pbResp), nil
}

// Mapping Functions

func mapNotificationRequestToProto(req *notification.SendNotificationRequest) *grpcnotification1.SendNotificationRequest {
	pbReq := &grpcnotification1.SendNotificationRequest{
		PreferredChannel: grpcnotification1.DeliveryChannel(req.PreferredChannel),
		Recipients:       make([]*grpcnotification1.Recipient, 0, len(req.Recipients)),
	}

	// Map Recipients array using a clean switch for the oneof
	for _, r := range req.Recipients {
		pbRecipient := &grpcnotification1.Recipient{}
		switch {
		case r.EmailAddress != nil:
			pbRecipient.Target = &grpcnotification1.Recipient_EmailAddress{EmailAddress: *r.EmailAddress}
		case r.PhoneNumber != nil:
			pbRecipient.Target = &grpcnotification1.Recipient_PhoneNumber{PhoneNumber: *r.PhoneNumber}
		case r.DeviceToken != nil:
			pbRecipient.Target = &grpcnotification1.Recipient_DeviceToken{DeviceToken: *r.DeviceToken}
		case r.UserID != nil:
			pbRecipient.Target = &grpcnotification1.Recipient_UserId{UserId: *r.UserID}
		}
		pbReq.Recipients = append(pbReq.Recipients, pbRecipient)
	}

	// Map Content oneof using a clean switch
	pbReq.Content = &grpcnotification1.Content{}
	switch {
	case req.Content.Raw != nil:
		pbReq.Content.Payload = &grpcnotification1.Content_Raw{
			Raw: &grpcnotification1.RawMessage{
				Subject:  req.Content.Raw.Subject,
				Body:     req.Content.Raw.Body,
				Metadata: req.Content.Raw.Metadata,
			},
		}
	case req.Content.Template != nil:
		pbReq.Content.Payload = &grpcnotification1.Content_Template{
			Template: &grpcnotification1.TemplateMessage{
				TemplateId: req.Content.Template.TemplateID,
				Parameters: req.Content.Template.Parameters,
			},
		}
	}

	return pbReq
}

func mapNotificationResponseToDomain(pbResp *grpcnotification1.SendNotificationResponse) *notification.SendNotificationResponse {
	domainResp := &notification.SendNotificationResponse{
		TrackingID:      pbResp.TrackingId,
		PartialFailures: make([]notification.DeliveryFailure, 0, len(pbResp.PartialFailures)),
	}

	// Map partial failures back to domain objects
	for _, pf := range pbResp.PartialFailures {
		domainFailure := notification.DeliveryFailure{
			ErrorReason: pf.ErrorReason,
			Recipient:   notification.Recipient{},
		}

		if pf.Recipient != nil {
			// Type switch to unpack the gRPC oneof into domain pointers
			switch target := pf.Recipient.Target.(type) {
			case *grpcnotification1.Recipient_EmailAddress:
				domainFailure.Recipient.EmailAddress = &target.EmailAddress
			case *grpcnotification1.Recipient_PhoneNumber:
				domainFailure.Recipient.PhoneNumber = &target.PhoneNumber
			case *grpcnotification1.Recipient_DeviceToken:
				domainFailure.Recipient.DeviceToken = &target.DeviceToken
			case *grpcnotification1.Recipient_UserId:
				domainFailure.Recipient.UserID = &target.UserId
			}
		}

		domainResp.PartialFailures = append(domainResp.PartialFailures, domainFailure)
	}

	return domainResp
}
