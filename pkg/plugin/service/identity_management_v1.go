package service

import (
	"context"

	"github.com/openkcm/plugin-sdk/api/service/identitymanagement"
	"github.com/openkcm/plugin-sdk/pkg/catalog"
	identity_managementv1 "github.com/openkcm/plugin-sdk/proto/plugin/identity_management/v1"
)

var _ identitymanagement.IdentityManagement = (*hashicorpIdentityManagementV1Plugin)(nil)

type hashicorpIdentityManagementV1Plugin struct {
	plugin     catalog.Plugin
	grpcClient identity_managementv1.IdentityManagementServiceClient
}

func NewIdentityManagementV1Plugin(plugin catalog.Plugin) identitymanagement.IdentityManagement {
	return &hashicorpIdentityManagementV1Plugin{
		plugin:     plugin,
		grpcClient: identity_managementv1.NewIdentityManagementServiceClient(plugin.ClientConnection()),
	}
}

func (h *hashicorpIdentityManagementV1Plugin) GetGroup(ctx context.Context, req *identitymanagement.GetGroupRequest) (*identitymanagement.GetGroupResponse, error) {
	in := &identity_managementv1.GetGroupRequest{
		GroupName:   req.GroupName,
		AuthContext: AuthContextToGRPC(&req.AuthContext),
	}
	grpcResp, err := h.grpcClient.GetGroup(ctx, in)
	if err != nil {
		return nil, err
	}
	return &identitymanagement.GetGroupResponse{
		Group: FromGRPCGroup(grpcResp.GetGroup()),
	}, nil
}

func (h *hashicorpIdentityManagementV1Plugin) GetAllGroups(ctx context.Context, req *identitymanagement.GetAllGroupsRequest) (*identitymanagement.GetAllGroupsResponse, error) {
	in := &identity_managementv1.GetAllGroupsRequest{
		AuthContext: AuthContextToGRPC(&req.AuthContext),
	}
	grpcResp, err := h.grpcClient.GetAllGroups(ctx, in)
	if err != nil {
		return nil, err
	}
	return &identitymanagement.GetAllGroupsResponse{
		Groups: FromGRPCGroups(grpcResp.GetGroups()),
	}, nil
}

func (h *hashicorpIdentityManagementV1Plugin) GetUsersForGroup(ctx context.Context, req *identitymanagement.GetUsersForGroupRequest) (*identitymanagement.GetUsersForGroupResponse, error) {
	in := &identity_managementv1.GetUsersForGroupRequest{
		GroupId:     req.GroupID,
		AuthContext: AuthContextToGRPC(&req.AuthContext),
	}
	grpcResp, err := h.grpcClient.GetUsersForGroup(ctx, in)
	if err != nil {
		return nil, err
	}
	return &identitymanagement.GetUsersForGroupResponse{
		Users: FromGRPCUsers(grpcResp.GetUsers()),
	}, nil
}

func (h *hashicorpIdentityManagementV1Plugin) GetGroupsForUser(ctx context.Context, req *identitymanagement.GetGroupsForUserRequest) (*identitymanagement.GetGroupsForUserResponse, error) {
	in := &identity_managementv1.GetGroupsForUserRequest{
		UserId:      req.UserID,
		AuthContext: AuthContextToGRPC(&req.AuthContext),
	}
	grpcResp, err := h.grpcClient.GetGroupsForUser(ctx, in)
	if err != nil {
		return nil, err
	}
	return &identitymanagement.GetGroupsForUserResponse{
		Groups: FromGRPCGroups(grpcResp.GetGroups()),
	}, nil
}

func AuthContextToGRPC(v *identitymanagement.AuthContext) *identity_managementv1.AuthContext {
	if v == nil {
		return nil
	}
	return &identity_managementv1.AuthContext{
		Data: v.Data,
	}
}

func FromGRPCGroup(v *identity_managementv1.Group) identitymanagement.Group {
	if v == nil {
		return identitymanagement.Group{}
	}
	return identitymanagement.Group{
		ID:   v.Id,
		Name: v.Name,
	}
}

func FromGRPCGroups(groups []*identity_managementv1.Group) []identitymanagement.Group {
	var wrapperGroups []identitymanagement.Group
	for _, group := range groups {
		wrapperGroups = append(wrapperGroups, FromGRPCGroup(group))
	}
	return wrapperGroups
}

func FromGRPCUser(v *identity_managementv1.User) identitymanagement.User {
	if v == nil {
		return identitymanagement.User{}
	}
	return identitymanagement.User{
		ID:    v.Id,
		Name:  v.Name,
		Email: v.Email,
	}
}

func FromGRPCUsers(users []*identity_managementv1.User) []identitymanagement.User {
	var wrapperUsers []identitymanagement.User
	for _, user := range users {
		wrapperUsers = append(wrapperUsers, FromGRPCUser(user))
	}
	return wrapperUsers
}
