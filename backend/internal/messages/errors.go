package messages

import "errors"

var (
	ErrMessageNotFound     = errors.New("message not found")
	ErrNotSender           = errors.New("only sender can delete their messages")
	ErrInvalidContentType  = errors.New("invalid content type")
	ErrContentTooLarge     = errors.New("content exceeds maximum size")
	ErrNotChannelMember    = errors.New("not a member of this channel")
	ErrInvalidEncryption   = errors.New("invalid encryption metadata")
)
