package messages

import (
	"context"
	"fmt"

	"telegraph/internal/acl"
	"telegraph/internal/audit"
	"telegraph/internal/channels"

	"github.com/google/uuid"
)

const MaxContentSize = 10 * 1024 * 1024 // 10MB

type MessageService interface {
	SendMessage(ctx context.Context, req SendMessageRequest, senderID, channelID uuid.UUID) (*Message, error)
	GetMessages(ctx context.Context, channelID, userID uuid.UUID, limit, offset int) ([]*Message, error)
	DeleteMessage(ctx context.Context, messageID, userID uuid.UUID, userRole string) error
}

type messageService struct {
	repo        MessageRepo
	channelRepo channels.ChannelRepo
	audit       *audit.Logger
}

func NewMessageService(repo MessageRepo, channelRepo channels.ChannelRepo, audit *audit.Logger) MessageService {
	return &messageService{repo: repo, channelRepo: channelRepo, audit: audit}
}

func (s *messageService) SendMessage(ctx context.Context, req SendMessageRequest, senderID, channelID uuid.UUID) (*Message, error) {
	// Validate content type
	if req.ContentType != ContentTypeText && req.ContentType != ContentTypeImage && 
	   req.ContentType != ContentTypeAudio && req.ContentType != ContentTypeVideo {
		return nil, ErrInvalidContentType
	}

	// Validate content size
	if len(req.Content) > MaxContentSize {
		return nil, ErrContentTooLarge
	}

	// Verify sender is member of channel
	isMember, err := s.channelRepo.IsMember(ctx, channelID, senderID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrNotChannelMember
	}

	// Validate encryption metadata exists
	if req.EncryptionMeta == nil || len(req.EncryptionMeta) == 0 {
		return nil, ErrInvalidEncryption
	}

	message := &Message{
		SenderID:       senderID,
		ChannelID:      channelID,
		Content:        req.Content,
		ContentType:    req.ContentType,
		EncryptionMeta: req.EncryptionMeta,
		Deleted:        false,
	}

	if err := s.repo.Create(ctx, message); err != nil {
		return nil, err
	}

	// Audit Log
	s.audit.Log(ctx, audit.AuditLog{
		UserID:   &senderID,
		Action:   audit.EventMessageSent,
		Resource: message.ID.String(),
		Result:   "success",
		Details:  fmt.Sprintf("Sent message to channel %s", channelID),
	})

	return message, nil
}

func (s *messageService) GetMessages(ctx context.Context, channelID, userID uuid.UUID, limit, offset int) ([]*Message, error) {
	// Verify user is member of channel
	isMember, err := s.channelRepo.IsMember(ctx, channelID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrNotChannelMember
	}

	// Apply default pagination limits
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.GetByChannelID(ctx, channelID, limit, offset)
}

func (s *messageService) DeleteMessage(ctx context.Context, messageID, userID uuid.UUID, userRole string) error {
	message, err := s.repo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}

	// Users can delete their own messages
	if message.SenderID == userID {
		return s.repo.SoftDelete(ctx, messageID)
	}

	// Moderators and admins can delete any message
	if acl.HasPermission(userRole, acl.PermissionDeleteAnyMessage) {
		return s.repo.SoftDelete(ctx, messageID)
	}

	return fmt.Errorf("insufficient permissions to delete message")
}
