package messages

import (
	"context"
	"fmt"
	"time"

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
	MarkAsDelivered(ctx context.Context, messageID, userID uuid.UUID) error
	MarkAsRead(ctx context.Context, messageID, userID uuid.UUID) error
	EditMessage(ctx context.Context, messageID, userID uuid.UUID, newContent []byte) error
	BroadcastTyping(ctx context.Context, userID, channelID uuid.UUID, typing bool) error
	GetUnreadCounts(ctx context.Context, userID uuid.UUID) (map[string]int, error)
}

type messageService struct {
	repo        MessageRepo
	channelRepo channels.ChannelRepo
	audit       *audit.Logger
	hub         Hub
}

// Hub interface for WebSocket broadcasting
type Hub interface {
	SendToUser(userID string, message interface{})
	BroadcastTyping(userID, channelID string, typing bool)
}

func NewMessageService(repo MessageRepo, channelRepo channels.ChannelRepo, audit *audit.Logger, hub Hub) MessageService {
	return &messageService{repo: repo, channelRepo: channelRepo, audit: audit, hub: hub}
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
		Attachments:    req.Attachments,
		ReplyTo:        req.ReplyTo,
		ForwardedFrom:  req.ForwardedFrom,
		Status:         MessageStatusSent,
		Deleted:        false,
		Edited:         false,
	}

	if err := s.repo.Create(ctx, message); err != nil {
		return nil, err
	}

	// Broadcast message to all channel members via WebSocket
	channel, _ := s.channelRepo.GetByID(ctx, channelID)
	if channel != nil && s.hub != nil {
		wsMessage := map[string]interface{}{
			"type":       "MESSAGE_NEW",
			"channel_id": channelID.String(),
			"message":    message,
		}
		
		for _, member := range channel.Members {
			s.hub.SendToUser(member.UserID.String(), wsMessage)
		}
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

func (s *messageService) MarkAsDelivered(ctx context.Context, messageID, userID uuid.UUID) error {
	// Logic to mark as delivered (e.g., update message status if needed, or just notify sender)
	// For now, we'll just broadcast the update
	message, err := s.repo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}

	// Broadcast read receipt
	if s.hub != nil {
		s.hub.SendToUser(message.SenderID.String(), map[string]interface{}{
			"type":       "MESSAGE_DELIVERED",
			"message_id": messageID,
			"user_id":    userID,
			"channel_id": message.ChannelID,
		})
	}
	return nil
}

func (s *messageService) MarkAsRead(ctx context.Context, messageID, userID uuid.UUID) error {
	message, err := s.repo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}

	// Update last read pointer in channel member
	if err := s.channelRepo.UpdateLastRead(ctx, message.ChannelID, userID, messageID); err != nil {
		return err
	}

	// Broadcast read receipt
	if s.hub != nil {
		s.hub.SendToUser(message.SenderID.String(), map[string]interface{}{
			"type":       "MESSAGE_READ",
			"message_id": messageID,
			"user_id":    userID,
			"channel_id": message.ChannelID,
		})
	}
	return nil
}

func (s *messageService) EditMessage(ctx context.Context, messageID, userID uuid.UUID, newContent []byte) error {
	message, err := s.repo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}

	if message.SenderID != userID {
		return ErrNotSender
	}

	message.Content = newContent
	message.Edited = true
	now := time.Now()
	message.EditedAt = &now

	if err := s.repo.Update(ctx, message); err != nil {
		return err
	}

	// Broadcast edit
	if s.hub != nil {
		// We need to broadcast to the channel, but Hub.SendToUser is 1:1.
		// Ideally Hub should have BroadcastToChannel.
		// For now, we rely on the handler or assume Hub handles it?
		// Actually, let's just use the existing mechanism if possible.
		// The current Hub implementation in `hub.go` has `BroadcastTyping` which iterates all clients.
		// We should probably add `BroadcastToChannel` to Hub, but for now let's stick to the plan.
		// Wait, `SendMessage` iterates channel members. We should do the same here.
		channel, _ := s.channelRepo.GetByID(ctx, message.ChannelID)
		if channel != nil {
			wsMessage := map[string]interface{}{
				"type":       "MESSAGE_UPDATED",
				"channel_id": message.ChannelID.String(),
				"message":    message,
			}
			for _, member := range channel.Members {
				s.hub.SendToUser(member.UserID.String(), wsMessage)
			}
		}
	}
	return nil
}

func (s *messageService) BroadcastTyping(ctx context.Context, userID, channelID uuid.UUID, typing bool) error {
	if s.hub != nil {
		s.hub.BroadcastTyping(userID.String(), channelID.String(), typing)
	}
	return nil
}

func (s *messageService) GetUnreadCounts(ctx context.Context, userID uuid.UUID) (map[string]int, error) {
	channels, err := s.channelRepo.GetUserChannels(ctx, userID)
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int)
	for _, ch := range channels {
		var lastReadTime time.Time
		// Find user's last read message
		for _, m := range ch.Members {
			if m.UserID == userID {
				if m.LastReadMessageID != nil {
					msg, err := s.repo.GetByID(ctx, *m.LastReadMessageID)
					if err == nil {
						lastReadTime = msg.Timestamp
					}
				} else {
					// If never read, count from joined_at
					lastReadTime = m.JoinedAt
				}
				break
			}
		}

		count, err := s.repo.CountAfter(ctx, ch.ID, lastReadTime)
		if err == nil {
			counts[ch.ID.String()] = int(count)
		}
	}
	return counts, nil
}
