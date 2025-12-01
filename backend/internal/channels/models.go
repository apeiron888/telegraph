package channels

import (
	"time"

	"github.com/google/uuid"
)

// ChannelType defines the type of channel
type ChannelType string

const (
	ChannelTypePrivate   ChannelType = "private"  // 1-on-1 direct message
	ChannelTypeGroup     ChannelType = "group"    // Multi-user group chat
	ChannelTypeBroadcast ChannelType = "channel"  // Broadcast channel (one-to-many)
)

const (
	ChannelRoleOwner  = "owner"
	ChannelRoleAdmin  = "admin"
	ChannelRoleMember = "member"
)

type ChannelMember struct {
	UserID            uuid.UUID  `json:"user_id" bson:"user_id"`
	Role              string     `json:"role" bson:"role"`
	JoinedAt          time.Time  `json:"joined_at" bson:"joined_at"`
	LastReadMessageID *uuid.UUID `json:"last_read_message_id,omitempty" bson:"last_read_message_id,omitempty"`
}

// Channel represents a messaging channel/group/private chat
type Channel struct {
	ID            uuid.UUID              `json:"id" bson:"id"`
	Type          ChannelType            `json:"type" bson:"type"`
	Name          string                 `json:"name,omitempty" bson:"name,omitempty"` // Optional for private chats
	Description   string                 `json:"description,omitempty" bson:"description,omitempty"`
	OwnerID       uuid.UUID              `json:"owner_id" bson:"owner_id"`
	Members       []ChannelMember        `json:"members" bson:"members"` // Stored as array of objects
	Permissions   map[string]interface{} `json:"permissions,omitempty" bson:"permissions,omitempty"` // ABAC policies
	SecurityLabel string                 `json:"security_label" bson:"security_label"` // MAC classification
	CreatedAt     time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" bson:"updated_at"`
}

// CreateChannelRequest is the payload for creating a new channel
type CreateChannelRequest struct {
	Type          ChannelType            `json:"type"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Members       []uuid.UUID            `json:"members"`       // Initial members (default to 'member' role)
	Permissions   map[string]interface{} `json:"permissions"`   // Optional ABAC policies
	SecurityLabel string                 `json:"security_label"` // Optional, defaults to owner's label
}

// AddMemberRequest is the payload for adding a member
type AddMemberRequest struct {
	UserID uuid.UUID `json:"user_id,omitempty"`
	Email  string    `json:"email,omitempty"`
	Phone  string    `json:"phone,omitempty"`
}
