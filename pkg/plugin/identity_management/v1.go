package identity_management

import (
	"context"

	"github.com/openkcm/plugin-sdk/api"
	"github.com/openkcm/plugin-sdk/api/service/identitymanagement"
	"github.com/openkcm/plugin-sdk/pkg/plugin"
	grpcidentitymanagementv1 "github.com/openkcm/plugin-sdk/proto/plugin/identity_management/v1"
)

type V1 struct {
	plugin.Facade
	grpcidentitymanagementv1.IdentityManagementPluginClient
}

func (v1 *V1) Version() uint {
	return 1
}

func (v1 *V1) ServiceInfo() api.Info {
	return v1.Info
}

func (v1 *V1) GetGroup(ctx context.Context, req *identitymanagement.GetGroupRequest) (*identitymanagement.GetGroupResponse, error) {
	in := &grpcidentitymanagementv1.GetGroupRequest{
		GroupName:   req.GroupName,
		AuthContext: AuthContextToGRPC(&req.AuthContext),
	}
	grpcResp, err := v1.IdentityManagementPluginClient.GetGroup(ctx, in)
	if err != nil {
		return nil, err
	}
	return &identitymanagement.GetGroupResponse{
		Group: FromGRPCGroup(grpcResp.GetGroup()),
	}, nil
}

func (v1 *V1) ListGroups(ctx context.Context, req *identitymanagement.ListGroupsRequest) (*identitymanagement.ListGroupsResponse, error) {
	in := &grpcidentitymanagementv1.ListGroupsRequest{
		AuthContext: AuthContextToGRPC(&req.AuthContext),
	}
	grpcResp, err := v1.IdentityManagementPluginClient.ListGroups(ctx, in)
	if err != nil {
		return nil, err
	}
	return &identitymanagement.ListGroupsResponse{
		Groups: FromGRPCGroups(grpcResp.GetGroups()),
	}, nil
}

func (v1 *V1) ListGroupUsers(ctx context.Context, req *identitymanagement.ListGroupUsersRequest) (*identitymanagement.ListGroupUsersResponse, error) {
	in := &grpcidentitymanagementv1.ListGroupUsersRequest{
		GroupId:     req.GroupID,
		AuthContext: AuthContextToGRPC(&req.AuthContext),
	}
	grpcResp, err := v1.IdentityManagementPluginClient.ListGroupUsers(ctx, in)
	if err != nil {
		return nil, err
	}
	return &identitymanagement.ListGroupUsersResponse{
		Users: FromGRPCUsers(grpcResp.GetUsers()),
	}, nil
}

func (v1 *V1) ListUserGroups(ctx context.Context, req *identitymanagement.ListUserGroupsRequest) (*identitymanagement.ListUserGroupsResponse, error) {
	in := &grpcidentitymanagementv1.ListUserGroupsRequest{
		UserId:      req.UserID,
		AuthContext: AuthContextToGRPC(&req.AuthContext),
	}
	grpcResp, err := v1.IdentityManagementPluginClient.ListUserGroups(ctx, in)
	if err != nil {
		return nil, err
	}
	return &identitymanagement.ListUserGroupsResponse{
		Groups: FromGRPCGroups(grpcResp.GetGroups()),
	}, nil
}

func AuthContextToGRPC(v *identitymanagement.AuthContext) *grpcidentitymanagementv1.AuthContext {
	if v == nil {
		return nil
	}
	return &grpcidentitymanagementv1.AuthContext{
		Data: v.Data,
	}
}

func FromGRPCGroup(v *grpcidentitymanagementv1.Group) identitymanagement.Group {
	if v == nil {
		return identitymanagement.Group{}
	}
	return identitymanagement.Group{
		ID:   v.Id,
		Name: v.Name,
	}
}

func FromGRPCGroups(groups []*grpcidentitymanagementv1.Group) []identitymanagement.Group {
	var wrapperGroups []identitymanagement.Group
	for _, group := range groups {
		wrapperGroups = append(wrapperGroups, FromGRPCGroup(group))
	}
	return wrapperGroups
}

func FromGRPCUser(v *grpcidentitymanagementv1.User) identitymanagement.User {
	if v == nil {
		return identitymanagement.User{}
	}
	return identitymanagement.User{
		ID:    v.Id,
		Name:  v.Name,
		Email: v.Email,
	}
}

func FromGRPCUsers(users []*grpcidentitymanagementv1.User) []identitymanagement.User {
	var wrapperUsers []identitymanagement.User
	for _, user := range users {
		wrapperUsers = append(wrapperUsers, FromGRPCUser(user))
	}
	return wrapperUsers
}
