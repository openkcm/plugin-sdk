package extauthzplugin

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"google.golang.org/grpc"

	authzpluginv1 "github.com/openkcm/plugin-sdk/proto/kms/plugin/extauthz/v1"
)

type clientMock struct {
	error bool
	allow bool
}

func (m *clientMock) Check(ctx context.Context, in *authzpluginv1.CheckRequest, opts ...grpc.CallOption) (*authzpluginv1.CheckResponse, error) {
	if m.error {
		return nil, fmt.Errorf("error")
	}
	if m.allow {
		return &authzpluginv1.CheckResponse{Allowed: true}, nil
	} else {
		return &authzpluginv1.CheckResponse{Allowed: false}, nil
	}
}

func TestClientCheck(t *testing.T) {
	// create test cases
	tests := []struct {
		name      string
		request   CheckRequest
		mockErr   bool
		mockAllow bool
		error     bool
		want      CheckResponse
	}{
		{
			name:  "zero values",
			error: false,
		}, {
			name:    "mock error",
			mockErr: true,
			error:   true,
		}, {
			name:      "mock reject",
			mockAllow: false,
			error:     false,
			want:      CheckResponse{Allowed: false},
		}, {
			name:      "mock allow",
			mockAllow: true,
			error:     false,
			want:      CheckResponse{Allowed: true},
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			clnt := GRPCClient{client: &clientMock{
				error: tc.mockErr,
				allow: tc.mockAllow,
			}}

			// Act
			got, err := clnt.Check(tc.request)

			// Assert
			if tc.error && err != nil { // expected error and got it
				return
			} else if tc.error && err == nil { // expected error but did not get it
				t.Errorf("expected error value: %v, got: %s", tc.error, err)
			} else if !tc.error && err != nil { // got unexpected error
				t.Errorf("expected error value: %v, got: %s", tc.error, err)
			} else if !tc.error && err == nil {
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("expected: %+v, got: %+v", tc.want, got)
				}
			}
		})
	}
}

type authzMock struct {
	error bool
	allow bool
}

func (m *authzMock) Check(CheckRequest) (CheckResponse, error) {
	if m.error {
		return CheckResponse{}, fmt.Errorf("error")
	}
	if m.allow {
		return CheckResponse{Allowed: true}, nil
	} else {
		return CheckResponse{Allowed: false}, nil
	}
}

func TestServerCheck(t *testing.T) {
	req := &authzpluginv1.CheckRequest{
		Subject: "subject",
		Object:  "object",
		Action:  "action",
	}

	// create test cases
	tests := []struct {
		name      string
		request   *authzpluginv1.CheckRequest
		mockErr   bool
		mockAllow bool
		error     bool
		want      *authzpluginv1.CheckResponse
	}{
		{
			name:    "mock error",
			request: req,
			mockErr: true,
			error:   true,
		}, {
			name:    "mock reject",
			request: req,
			error:   false,
			want:    &authzpluginv1.CheckResponse{Allowed: false},
		}, {
			name:      "mock allow",
			request:   req,
			mockAllow: true,
			error:     false,
			want:      &authzpluginv1.CheckResponse{Allowed: true},
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			srv := GRPCServer{Impl: &authzMock{
				error: tc.mockErr,
				allow: tc.mockAllow,
			}}

			// Act
			got, err := srv.Check(context.Background(), tc.request)

			// Assert
			if tc.error && err != nil { // expected error and got it
				return
			} else if tc.error && err == nil { // expected error but did not get it
				t.Errorf("expected error value: %v, got: %s", tc.error, err)
			} else if !tc.error && err != nil { // got unexpected error
				t.Errorf("expected error value: %v, got: %s", tc.error, err)
			} else if !tc.error && err == nil {
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("expected: %+v, got: %+v", tc.want, got)
				}
			}
		})
	}
}
