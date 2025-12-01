package channels

import (
	"context"
	"fmt"
	"time"

	"telegraph/internal/acl"
	"telegraph/internal/audit"
	"telegraph/internal/users"

	"github.com/google/uuid"
)

type ChannelService interface {
	CreateChannel(ctx context.Context, req CreateChannelRequest, creatorID uuid.UUID, creatorRole string) (*Channel, error)
	GetChannel(ctx context.Context, channelID uuid.UUID) (*Channel, error)
	GetUserChannels(ctx context.Context, userID uuid.UUID) ([]*Channel, error)
	AddMember(ctx context.Context, channelID, requestorID, newMemberID uuid.UUID) error
	RemoveMember(ctx context.Context, channelID, requestorID, memberID uuid.UUID) error
	DeleteChannel(ctx context.Context, channelID, requestorID uuid.UUID) error
	IsMember(ctx context.Context, channelID, userID uuid.UUID) (bool, error)
	PromoteToAdmin(ctx context.Context, channelID, requestorID, memberID uuid.UUID) error
	DemoteAdmin(ctx context.Context, channelID, requestorID, memberID uuid.UUID) error
}

type channelService struct {
	repo     ChannelRepo
	userRepo users.UserRepo
	audit    *audit.Logger
}

func NewChannelService(repo ChannelRepo, userRepo users.UserRepo, audit *audit.Logger) ChannelService {
	return &channelService{repo: repo, userRepo: userRepo, audit: audit}
}

func (s *channelService) CreateChannel(ctx context.Context, req CreateChannelRequest, creatorID uuid.UUID, creatorRole string) (*Channel, error) {
	// Validate channel type
	if req.Type != ChannelTypePrivate && req.Type != ChannelTypeGroup && req.Type != ChannelTypeBroadcast {
		return nil, ErrInvalidChannelType
	}

	// Validate Name for Group/Broadcast
	if (req.Type == ChannelTypeGroup || req.Type == ChannelTypeBroadcast) && req.Name == "" {
		return nil, fmt.Errorf("channel name is required for groups and broadcasts")
	}

	// Only admins can create broadcast channels (RBAC enforcement)
	if req.Type == ChannelTypeBroadcast && !acl.HasPermission(creatorRole, acl.PermissionBroadcast) {
		return nil, ErrBroadcastRestricted
	}

	// Get creator to inherit security label if not specified
	creator, err := s.userRepo.GetByID(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	// Account Limit Check (Basic users max 5 owned channels)
	if creator.AccountType == "basic" && req.Type != ChannelTypePrivate {
		userChannels, err := s.repo.GetUserChannels(ctx, creatorID)
		if err != nil {
			return nil, err
		}
		ownedCount := 0
		for _, c := range userChannels {
			if c.OwnerID == creatorID && c.Type != ChannelTypePrivate {
				ownedCount++
			}
		}
		if ownedCount >= 5 {
			return nil, fmt.Errorf("basic account limit reached: upgrade to premium to create more groups")
		}
	}

	// Default security label to creator's label
	securityLabel := req.SecurityLabel
	if securityLabel == "" {
		securityLabel = creator.SecurityLabel
	}

	// Validate security label
	if err := acl.ValidateLabel(securityLabel); err != nil {
		return nil, err
	}

	// Creator must have clearance for the channel's security label
	if !acl.CanAccessResource(creator.SecurityLabel, securityLabel) {
		return nil, fmt.Errorf("insufficient clearance to create channel with label: %s", securityLabel)
	}

	// Initialize members with creator as Owner
	members := []ChannelMember{{
		UserID:   creatorID,
		Role:     ChannelRoleOwner,
		JoinedAt: time.Now(),
	}}

	if req.Members != nil {
		// Add requested members (avoiding duplicates)
		memberMap := map[uuid.UUID]bool{creatorID: true}
		for _, m := range req.Members {
			if !memberMap[m] {
				members = append(members, ChannelMember{
					UserID:   m,
					Role:     ChannelRoleMember, // Default role
					JoinedAt: time.Now(),
				})
				memberMap[m] = true
			}
		}
	}

	channel := &Channel{
		Type:          req.Type,
		Name:          req.Name,
		Description:   req.Description,
		OwnerID:       creatorID,
		Members:       members,
		Permissions:   req.Permissions,
		SecurityLabel: securityLabel,
	}

	if err := s.repo.Create(ctx, channel); err != nil {
		return nil, err
	}

	// Audit Log
	s.audit.Log(ctx, audit.AuditLog{
		UserID:   &creatorID,
		Action:   audit.EventChannelCreated,
		Resource: channel.ID.String(),
		Result:   "success",
		Details:  fmt.Sprintf("Created channel '%s' (%s)", channel.Name, channel.Type),
	})

	return channel, nil
}

func (s *channelService) GetChannel(ctx context.Context, channelID uuid.UUID) (*Channel, error) {
	return s.repo.GetByID(ctx, channelID)
}

func (s *channelService) GetUserChannels(ctx context.Context, userID uuid.UUID) ([]*Channel, error) {
	return s.repo.GetUserChannels(ctx, userID)
}

func (s *channelService) AddMember(ctx context.Context, channelID, requestorID, newMemberID uuid.UUID) error {
	// Get channel
	channel, err := s.repo.GetByID(ctx, channelID)
	if err != nil {
		return err
	}

	// Check if requestor is owner or admin
	isAuthorized := false
	if channel.OwnerID == requestorID {
		isAuthorized = true
	} else {
		for _, m := range channel.Members {
			if m.UserID == requestorID && m.Role == ChannelRoleAdmin {
				isAuthorized = true
				break
			}
		}
	}

	if !isAuthorized {
		// Fallback to system permission check if not channel admin
		requestor, err := s.userRepo.GetByID(ctx, requestorID)
		if err != nil {
			return err
		}
		if !acl.HasPermission(requestor.Role, acl.PermissionManageMembers) {
			return ErrNotChannelOwner
		}
	}

	// Check if new member has required clearance for channel
	newMember, err := s.userRepo.GetByID(ctx, newMemberID)
	if err != nil {
		return err
	}

	if !acl.CanAccessResource(newMember.SecurityLabel, channel.SecurityLabel) {
		return fmt.Errorf("user lacks clearance for this channel")
	}

	return s.repo.AddMember(ctx, channelID, newMemberID)
}

func (s *channelService) RemoveMember(ctx context.Context, channelID, requestorID, memberID uuid.UUID) error {
	channel, err := s.repo.GetByID(ctx, channelID)
	if err != nil {
		return err
	}

	// Owner can remove anyone. Admins can remove members but not other admins/owner.
	isAuthorized := false
	requestorRole := ""
	
	if channel.OwnerID == requestorID {
		isAuthorized = true
		requestorRole = ChannelRoleOwner
	} else {
		for _, m := range channel.Members {
			if m.UserID == requestorID {
				requestorRole = m.Role
				if m.Role == ChannelRoleAdmin {
					isAuthorized = true
				}
				break
			}
		}
	}

	if !isAuthorized {
		// System admin override
		requestor, err := s.userRepo.GetByID(ctx, requestorID)
		if err != nil {
			return err
		}
		if acl.HasPermission(requestor.Role, acl.PermissionManageMembers) {
			isAuthorized = true
			requestorRole = ChannelRoleOwner // Treat as owner
		}
	}

	if !isAuthorized {
		return ErrNotChannelOwner
	}

	// Check target role
	targetRole := ChannelRoleMember
	for _, m := range channel.Members {
		if m.UserID == memberID {
			targetRole = m.Role
			break
		}
	}

	// Admin cannot remove another Admin or Owner
	if requestorRole == ChannelRoleAdmin && (targetRole == ChannelRoleAdmin || targetRole == ChannelRoleOwner) {
		return fmt.Errorf("admins cannot remove other admins or the owner")
	}

	return s.repo.RemoveMember(ctx, channelID, memberID)
}

func (s *channelService) DeleteChannel(ctx context.Context, channelID, requestorID uuid.UUID) error {
	channel, err := s.repo.GetByID(ctx, channelID)
	if err != nil {
		return err
	}

	// Only owner or admin can delete
	if channel.OwnerID != requestorID {
		requestor, err := s.userRepo.GetByID(ctx, requestorID)
		if err != nil {
			return err
		}
		if !acl.HasPermission(requestor.Role, acl.PermissionDeleteChannel) {
			return ErrNotChannelOwner
		}
	}

	return s.repo.Delete(ctx, channelID)
}

func (s *channelService) IsMember(ctx context.Context, channelID, userID uuid.UUID) (bool, error) {
	return s.repo.IsMember(ctx, channelID, userID)
}



func (s *channelService) PromoteToAdmin(ctx context.Context, channelID, requestorID, memberID uuid.UUID) error {
	channel, err := s.repo.GetByID(ctx, channelID)
	if err != nil {
		return err
	}

	// Only Owner can promote
	if channel.OwnerID != requestorID {
		return ErrNotChannelOwner
	}

	// Verify member exists
	isMember := false
	for _, m := range channel.Members {
		if m.UserID == memberID {
			isMember = true
			break
		}
	}
	if !isMember {
		return fmt.Errorf("user is not a member of this channel")
	}

	return s.repo.UpdateMemberRole(ctx, channelID, memberID, ChannelRoleAdmin)
}

func (s *channelService) DemoteAdmin(ctx context.Context, channelID, requestorID, memberID uuid.UUID) error {
	channel, err := s.repo.GetByID(ctx, channelID)
	if err != nil {
		return err
	}

	// Only Owner can demote
	if channel.OwnerID != requestorID {
		return ErrNotChannelOwner
	}

	// Verify member exists
	isMember := false
	for _, m := range channel.Members {
		if m.UserID == memberID {
			isMember = true
			break
		}
	}
	if !isMember {
		return fmt.Errorf("user is not a member of this channel")
	}

	return s.repo.UpdateMemberRole(ctx, channelID, memberID, ChannelRoleMember)
}
