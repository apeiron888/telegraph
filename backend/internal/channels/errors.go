package channels

import "errors"

var (
	ErrChannelNotFound     = errors.New("channel not found")
	ErrNotChannelMember    = errors.New("not a channel member")
	ErrNotChannelOwner     = errors.New("not the channel owner")
	ErrInvalidChannelType  = errors.New("invalid channel type")
	ErrBroadcastRestricted = errors.New("only admins can create broadcast channels")
	ErrAlreadyMember       = errors.New("user is already a member")
)
