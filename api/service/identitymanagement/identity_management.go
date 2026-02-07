package identitymanagement

import (
	"context"
)

type IdentityManagement interface {
	GetGroup(ctx context.Context, req *GetGroupRequest) (*GetGroupResponse, error)
	GetAllGroups(ctx context.Context, req *GetAllGroupsRequest) (*GetAllGroupsResponse, error)
	GetUsersForGroup(ctx context.Context, req *GetUsersForGroupRequest) (*GetUsersForGroupResponse, error)
	GetGroupsForUser(ctx context.Context, req *GetGroupsForUserRequest) (*GetGroupsForUserResponse, error)
}

type AuthContext struct {
	// V1 Fields
	Data map[string]string
}

type GetGroupRequest struct {
	// V1 Fields
	GroupName   string
	AuthContext AuthContext
}

type GetGroupResponse struct {
	// V1 Fields
	Group Group
}

type GetAllGroupsRequest struct {
	// V1 Fields
	AuthContext AuthContext
}

type GetAllGroupsResponse struct {
	// V1 Fields
	Groups []Group
}

type GetUsersForGroupRequest struct {
	// V1 Fields
	GroupID     string
	AuthContext AuthContext
}

type GetUsersForGroupResponse struct {
	// V1 Fields
	Users []User
}

type User struct {
	// V1 Fields
	ID    string
	Name  string
	Email string
}

type GetGroupsForUserRequest struct {
	// V1 Fields
	UserID      string
	AuthContext AuthContext
}

type GetGroupsForUserResponse struct {
	// V1 Fields
	Groups []Group
}

type Group struct {
	// V1 Fields
	ID   string
	Name string
}
