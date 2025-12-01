package users

import (
	"time"

	"github.com/google/uuid"
)

// Address represents user's physical address
type Address struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Street  string `json:"street"`
}

type User struct {
	ID           uuid.UUID              `json:"id" bson:"id"`
	Username     string                 `json:"username" bson:"username"`
	Email        string                 `json:"email" bson:"email"`
	Phone        string                 `json:"phone,omitempty" bson:"phone,omitempty"`
	PasswordHash string                 `json:"-" bson:"password_hash"`
	Bio          string                 `json:"bio" bson:"bio"`
	BirthDate    *time.Time             `json:"birth_date,omitempty" bson:"birth_date,omitempty"`
	
	// Address fields
	Country string `json:"country,omitempty" bson:"country,omitempty"`
	City    string `json:"city,omitempty" bson:"city,omitempty"`
	Street  string `json:"street,omitempty" bson:"street,omitempty"`
	
	// Account details
	AccountType    string     `json:"account_type" bson:"account_type"`     // "basic" or "premium"
	AccountStart   *time.Time `json:"account_start,omitempty" bson:"account_start,omitempty"`
	RenewalPeriod  int        `json:"renewal_period,omitempty" bson:"renewal_period,omitempty"` // days
	
	// Access Control fields
	Role           string                 `json:"role" bson:"role"`           // RBAC: "member", "moderator", "admin"
	SecurityLabel  string                 `json:"security_label" bson:"security_label"` // MAC: "public", "internal", "confidential"
	Attributes     map[string]interface{} `json:"attributes" bson:"attributes"`     // ABAC: custom attributes

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
