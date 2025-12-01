package messages

import (
	"time"

	"github.com/google/uuid"
)

// ContentType defines the type of message content
type ContentType string

const (
	ContentTypeText  ContentType = "text"
	ContentTypeImage ContentType = "image"
	ContentTypeAudio ContentType = "audio"
	ContentTypeVideo ContentType = "video"
)

// Message represents an encrypted message in a channel
type Message struct {
	ID             uuid.UUID              `json:"id" bson:"id"`
	SenderID       uuid.UUID              `json:"sender_id" bson:"sender_id"`
	ChannelID      uuid.UUID              `json:"channel_id" bson:"channel_id"`
	Content        []byte                 `json:"content" bson:"content"` // Encrypted blob
	ContentType    ContentType            `json:"content_type" bson:"content_type"`
	EncryptionMeta map[string]interface{} `json:"encryption_meta" bson:"encryption_meta"` // IV, algorithm info
	Timestamp      time.Time              `json:"timestamp" bson:"timestamp"`
	Deleted        bool                   `json:"deleted" bson:"deleted"` // Soft delete flag
	CreatedAt      time.Time              `json:"created_at" bson:"created_at"`
}

// SendMessageRequest is the payload for sending a message
type SendMessageRequest struct {
	Content        []byte                 `json:"content"` // Already encrypted by client
	ContentType    ContentType            `json:"content_type"`
	EncryptionMeta map[string]interface{} `json:"encryption_meta"`
}
