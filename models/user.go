package models

import (
	"time"

	"github.com/shinhagunn/mpa/pkg/mpa_fx"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole string

const (
	UserRoleMember     UserRole = "member"
	UserRoleAdmin      UserRole = "admin"
	UserRoleSuperAdmin UserRole = "superadmin"
)

var UserRoles = []UserRole{
	UserRoleMember,
	UserRoleAdmin,
	UserRoleSuperAdmin,
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UID       string             `bson:"uid"`
	Email     string             `bson:"email"`
	Role      UserRole           `bson:"role"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (u User) TableName() string {
	return "users"
}

func (u User) SetIndex() map[string]mpa_fx.IndexType {
	return map[string]mpa_fx.IndexType{
		"uid": mpa_fx.IndexTypeUnique,
	}
}
