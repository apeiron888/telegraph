package messages

import (
	"time"

	"github.com/google/uuid"
)

// ContentType defines the type of message content
type ContentType string

const (
	ContentTypeText     ContentType = "text"
	ContentTypeImage    ContentType = "image"
	ContentTypeAudio    ContentType = "audio"
	ContentTypeVideo    ContentType = "video"
	ContentTypeDocument ContentType = "document"
)

// MessageStatus represents delivery/read status
type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
)

// FileAttachment represents an attached file
type FileAttachment struct {
	FileID      string `json:"file_id" bson:"file_id"`
	FileName    string `json:"file_name" bson:"file_name"`
	ContentType string `json:"content_type" bson:"content_type"`
	Size        int64  `json:"size" bson:"size"`
	URL         string `json:"url" bson:"url"`
}

// Message represents an encrypted message in a channel
type Message struct {
	ID             uuid.UUID              `json:"id" bson:"id"`
	SenderID       uuid.UUID              `json:"sender_id" bson:"sender_id"`
	ChannelID      uuid.UUID              `json:"channel_id" bson:"channel_id"`
	Content        []byte                 `json:"content" bson:"content"` // Encrypted blob
	ContentType    ContentType            `json:"content_type" bson:"content_type"`
	EncryptionMeta map[string]interface{} `json:"encryption_meta" bson:"encryption_meta"` // IV, algorithm info
	
	// Attachments
	Attachments []FileAttachment `json:"attachments,omitempty" bson:"attachments,omitempty"`
	
	// Status tracking
	Status        MessageStatus            `json:"status" bson:"status"`
	DeliveredTo   []uuid.UUID              `json:"delivered_to,omitempty" bson:"delivered_to,omitempty"`
	ReadBy        []uuid.UUID              `json:"read_by,omitempty" bson:"read_by,omitempty"`
	
	// Reply/Forward
	ReplyTo       *uuid.UUID `json:"reply_to,omitempty" bson:"reply_to,omitempty"`
	ForwardedFrom *uuid.UUID `json:"forwarded_from,omitempty" bson:"forwarded_from,omitempty"`
	
	// Editing
	Edited    bool       `json:"edited" bson:"edited"`
	EditedAt  *time.Time `json:"edited_at,omitempty" bson:"edited_at,omitempty"`
	
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Deleted   bool      `json:"deleted" bson:"deleted"` // Soft delete flag
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// SendMessageRequest is the payload for sending a message
type SendMessageRequest struct {
	Content        []byte                 `json:"content"` // Already encrypted by client
	ContentType    ContentType            `json:"content_type"`
	EncryptionMeta map[string]interface{} `json:"encryption_meta"`
	Attachments    []FileAttachment       `json:"attachments,omitempty"`
	ReplyTo        *uuid.UUID             `json:"reply_to,omitempty"`
	ForwardedFrom  *uuid.UUID             `json:"forwarded_from,omitempty"`
}
