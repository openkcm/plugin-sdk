package certificate_issuer

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/api/service/certificateissuer"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	grpccertificateissuerv1 "github.com/openkcm/plugin-sdk/proto/plugin/certificate_issuer/v1"
)

type V1 struct {
	plugin.Facade
	grpccertificateissuerv1.CertificateIssuerPluginClient
}

func (v1 *V1) Version() uint {
	return 1
}

func (v1 *V1) ServiceInfo() api.Info {
	return v1.Info
}

func (v1 *V1) IssueCertificate(ctx context.Context, req *certificateissuer.IssueCertificateRequest) (*certificateissuer.IssueCertificateResponse, error) {
	pbReq := mapRequestToProto(req)

	if err := protovalidate.Validate(pbReq); err != nil {
		return nil, fmt.Errorf("failed validation: %v", err)
	}

	pbResp, err := v1.CertificateIssuerPluginClient.IssueCertificate(ctx, pbReq)
	if err != nil {
		return nil, handleGRPCError(err)
	}

	return mapResponseToDomain(pbResp), nil
}

func mapRequestToProto(req *certificateissuer.IssueCertificateRequest) *grpccertificateissuerv1.IssueCertificateRequest {
	pbReq := &grpccertificateissuerv1.IssueCertificateRequest{
		PreferredFormat: grpccertificateissuerv1.IssueCertificateRequest_CertificateFormat(req.PreferredFormat),
		Subject: &grpccertificateissuerv1.Subject{
			CommonName: req.Subject.CommonName,

			//Country:            req.Subject.Country,
			//Organization:       req.Subject.Organization,
			//OrganizationalUnit: req.Subject.OrganizationalUnit,
			//Locality:           req.Subject.Locality,
			//Province:           req.Subject.Province,
		},
	}

	//if req.Subject.AdvancedAttributes != nil {
	//	pbReq.Subject.AdvancedAttributes = &grpccertificateissuerv1.AdvancedSubjectAttributes{
	//		SerialNumber:  req.Subject.AdvancedAttributes.SerialNumber,
	//		StreetAddress: req.Subject.AdvancedAttributes.StreetAddress,
	//		PostalCode:    req.Subject.AdvancedAttributes.PostalCode,
	//	}
	//}

	if req.PrivateKey != nil {
		pbReq.PrivateKey = &grpccertificateissuerv1.PrivateKey{
			Data:   req.PrivateKey.Data,
			Format: grpccertificateissuerv1.PrivateKey_KeyFormat(req.PrivateKey.Format),
		}
	}

	switch {
	case req.Lifetime.Duration != nil:
		pbReq.Lifetime = &grpccertificateissuerv1.CertificateLifetime{
			Lifetime: &grpccertificateissuerv1.CertificateLifetime_Duration{
				Duration: durationpb.New(*req.Lifetime.Duration),
			},
		}
	case req.Lifetime.NotAfter != nil:
		pbReq.Lifetime = &grpccertificateissuerv1.CertificateLifetime{
			Lifetime: &grpccertificateissuerv1.CertificateLifetime_NotAfter{
				NotAfter: timestamppb.New(*req.Lifetime.NotAfter),
			},
		}
	case req.Lifetime.Relative != nil:
		pbReq.Lifetime = &grpccertificateissuerv1.CertificateLifetime{
			Lifetime: &grpccertificateissuerv1.CertificateLifetime_Relative{
				Relative: &grpccertificateissuerv1.RelativeValidity{
					Value: req.Lifetime.Relative.Value,
					Unit:  grpccertificateissuerv1.RelativeValidity_ValidityUnit(req.Lifetime.Relative.Unit),
				},
			},
		}
	}

	return pbReq
}

func mapResponseToDomain(pbResp *grpccertificateissuerv1.IssueCertificateResponse) *certificateissuer.IssueCertificateResponse {
	return &certificateissuer.IssueCertificateResponse{
		CertificateData: pbResp.CertificateData,
		Format:          certificateissuer.CertificateFormat(pbResp.Format),
		CAChain:         pbResp.CaChain,
	}
}

func handleGRPCError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	// Iterate over the "Any" details packed into the gRPC status
	for _, detail := range st.Details() {
		switch t := detail.(type) {
		case *grpccertificateissuerv1.SupportedKeyFormatsError:
			domainErr := &certificateissuer.SupportedKeyFormatsError{
				RejectedFormat:   t.RejectedFormat,
				Reason:           t.Reason,
				SupportedFormats: make([]certificateissuer.KeyFormat, len(t.SupportedFormats)),
			}
			for i, format := range t.SupportedFormats {
				domainErr.SupportedFormats[i] = certificateissuer.KeyFormat(format)
			}
			return domainErr
		}
	}

	return err
}
