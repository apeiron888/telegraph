package audit

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EventType represents types of audit events
type EventType string

const (
	EventLogin          EventType = "login"
	EventLoginFailed    EventType = "login_failed"
	EventMFAVerified    EventType = "mfa_verified"
	EventMFAFailed      EventType = "mfa_failed"
	EventLogout         EventType = "logout"
	EventAccessDenied   EventType = "access_denied"
	EventChannelCreated EventType = "channel_created"
	EventChannelDeleted EventType = "channel_deleted"
	EventMessageSent    EventType = "message_sent"
	EventMessageDeleted EventType = "message_deleted"
)

// AuditLog represents a single audit event
type AuditLog struct {
	ID        uuid.UUID  `bson:"_id"`
	UserID    *uuid.UUID `bson:"user_id,omitempty"`
	Action    EventType  `bson:"action"`
	Resource  string     `bson:"resource,omitempty"`
	IPAddress string     `bson:"ip_address,omitempty"`
	Result    string     `bson:"result"`
	Details   string     `bson:"details,omitempty"`
	Timestamp time.Time  `bson:"timestamp"`
}

// Logger provides audit logging functionality
type Logger struct {
	collection *mongo.Collection
	fileLogger *log.Logger
	file       *os.File
}

func NewLogger(db *mongo.Database) *Logger {
	file, err := os.OpenFile("audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open audit log file: %v", err)
		return &Logger{
			collection: db.Collection("audit_logs"),
		}
	}

	return &Logger{
		collection: db.Collection("audit_logs"),
		fileLogger: log.New(file, "", 0),
		file:       file,
	}
}

// Log records an audit event
func (l *Logger) Log(ctx context.Context, event AuditLog) error {
	event.ID = uuid.New()
	event.Timestamp = time.Now()

	// Log to file
	if l.fileLogger != nil {
		logEntry := fmt.Sprintf("[%s] Action=%s UserID=%v Resource=%s Result=%s IP=%s Details=%s",
			event.Timestamp.Format(time.RFC3339),
			event.Action,
			event.UserID,
			event.Resource,
			event.Result,
			event.IPAddress,
			event.Details,
		)
		l.fileLogger.Println(logEntry)
	}

	// Log to MongoDB
	_, err := l.collection.InsertOne(ctx, event)
	return err
}

// GetUserLogs retrieves audit logs for a specific user
func (l *Logger) GetUserLogs(ctx context.Context, userID uuid.UUID, limit int) ([]AuditLog, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := l.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []AuditLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}
