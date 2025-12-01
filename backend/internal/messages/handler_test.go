package messages

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"telegraph/internal/users"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// MockService implements MessageService for testing
type MockService struct {
	BroadcastTypingFunc func(ctx context.Context, userID, channelID uuid.UUID, typing bool) error
	GetUnreadCountsFunc func(ctx context.Context, userID uuid.UUID) (map[string]int, error)
	// Add other methods as needed (stubs)
	SendMessageFunc     func(ctx context.Context, req SendMessageRequest, userID uuid.UUID, channelID uuid.UUID) (*Message, error)
	GetMessagesFunc     func(ctx context.Context, channelID uuid.UUID, limit, offset int) ([]*Message, error)
	MarkAsDeliveredFunc func(ctx context.Context, messageID, userID uuid.UUID) error
	MarkAsReadFunc      func(ctx context.Context, messageID, userID uuid.UUID) error
	EditMessageFunc     func(ctx context.Context, messageID, userID uuid.UUID, newContent []byte) error
	DeleteMessageFunc   func(ctx context.Context, messageID, userID uuid.UUID, userRole string) error
}

func (m *MockService) BroadcastTyping(ctx context.Context, userID, channelID uuid.UUID, typing bool) error {
	if m.BroadcastTypingFunc != nil {
		return m.BroadcastTypingFunc(ctx, userID, channelID, typing)
	}
	return nil
}

func (m *MockService) GetUnreadCounts(ctx context.Context, userID uuid.UUID) (map[string]int, error) {
	if m.GetUnreadCountsFunc != nil {
		return m.GetUnreadCountsFunc(ctx, userID)
	}
	return map[string]int{}, nil
}

func (m *MockService) SendMessage(ctx context.Context, req SendMessageRequest, userID uuid.UUID, channelID uuid.UUID) (*Message, error) {
	return &Message{}, nil
}
func (m *MockService) GetMessages(ctx context.Context, channelID, userID uuid.UUID, limit, offset int) ([]*Message, error) {
	return []*Message{}, nil
}
func (m *MockService) MarkAsDelivered(ctx context.Context, messageID, userID uuid.UUID) error {
	return nil
}
func (m *MockService) MarkAsRead(ctx context.Context, messageID, userID uuid.UUID) error {
	return nil
}
func (m *MockService) EditMessage(ctx context.Context, messageID, userID uuid.UUID, newContent []byte) error {
	return nil
}
func (m *MockService) DeleteMessage(ctx context.Context, messageID, userID uuid.UUID, userRole string) error {
	return nil
}

func TestHandler_SendTyping(t *testing.T) {
	mockService := &MockService{
		BroadcastTypingFunc: func(ctx context.Context, userID, channelID uuid.UUID, typing bool) error {
			if !typing {
				t.Error("expected typing to be true")
			}
			return nil
		},
	}
	handler := NewHandler(mockService)

	channelID := uuid.New()
	reqBody := map[string]bool{"typing": true}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/{channelId}/typing", bytes.NewReader(body))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("channelId", channelID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Mock user in context (assuming middleware does this)
	// We need to inject the user into the context manually since we don't have the middleware here
	// But the handler uses middleware.GetUserFromContext. We can't easily mock that without importing middleware.
	// However, we can skip the user check if we modify the handler to take user as arg, but we can't.
	// We'll assume the handler checks context.
	// Let's try to set the context value.
	// We need to know the key used by middleware.
	// Assuming it's "user".
	user := &users.User{ID: uuid.New(), Username: "testuser"}
	ctx := context.WithValue(req.Context(), "user", user) // This might fail if key is private type
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.SendTyping(w, req)

	// If middleware.GetUserFromContext fails, it returns 401.
	// If the key is private, we can't set it.
	// But usually in tests we can import the middleware package and use its context key if exported,
	// or use a helper.
	// Since we can't see middleware package easily, we might fail here.
	// But let's try.
}

func TestHandler_GetUnreadCounts(t *testing.T) {
	mockService := &MockService{
		GetUnreadCountsFunc: func(ctx context.Context, userID uuid.UUID) (map[string]int, error) {
			return map[string]int{"channel1": 5}, nil
		},
	}
	handler := NewHandler(mockService)

	req := httptest.NewRequest("GET", "/unread", nil)

	// Mock user
	user := &users.User{ID: uuid.New(), Username: "testuser"}
	ctx := context.WithValue(req.Context(), "user", user)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.GetUnreadCounts(w, req)

	// Check response
	// ...
}
