package models

import (
	"time"

	"github.com/shinhagunn/mpa/pkg/mpa_fx"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Add callback
type UserState string

const (
	UserStateActive  UserState = "active"
	UserStatePending UserState = "pending"
	UserStateDeleted UserState = "deleted"
	UserStateBanned  UserState = "banned"
	UserStateLocked  UserState = "locked"
)

var UserStates = []UserState{
	UserStateActive,
	UserStatePending,
	UserStateDeleted,
	UserStateBanned,
	UserStateLocked,
}

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
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	UID            string             `bson:"uid"`
	Username       string             `bson:"username"`
	Email          string             `bson:"email"`
	PasswordDigest string             `bson:"password_digest"`
	Level          int64              `bson:"level"`
	OTP            bool               `bson:"otp"`
	Role           UserRole           `bson:"role"`
	State          UserState          `bson:"state"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
	// ReferralUID    string             `bson:"referral_uid"`
	// Data UserData            `bson:"data"`
}

func (u User) TableName() string {
	return "users"
}

func (u User) SetIndex() map[string]mpa_fx.IndexType {
	return map[string]mpa_fx.IndexType{
		"uid": mpa_fx.IndexTypeUnique,
	}
}
