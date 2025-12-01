package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RefreshTokenRepo interface {
	Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	GetByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	Revoke(ctx context.Context, tokenHash string) error
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
}

type RefreshToken struct {
	ID        uuid.UUID `bson:"_id"`
	UserID    uuid.UUID `bson:"user_id"`
	TokenHash string    `bson:"token_hash"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
	Revoked   bool      `bson:"revoked"`
}

type mongoRefreshTokenRepo struct {
	collection *mongo.Collection
}

func NewRefreshTokenRepo(db *mongo.Database) RefreshTokenRepo {
	return &mongoRefreshTokenRepo{
		collection: db.Collection("refresh_tokens"),
	}
}

func (r *mongoRefreshTokenRepo) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	token := RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		Revoked:   false,
	}
	_, err := r.collection.InsertOne(ctx, token)
	return err
}

func (r *mongoRefreshTokenRepo) GetByHash(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	var token RefreshToken
	err := r.collection.FindOne(ctx, bson.M{
		"token_hash": tokenHash,
		"revoked":    false,
	}).Decode(&token)
	
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &token, err
}

func (r *mongoRefreshTokenRepo) Revoke(ctx context.Context, tokenHash string) error {
	update := bson.M{"$set": bson.M{"revoked": true}}
	_, err := r.collection.UpdateOne(ctx, bson.M{"token_hash": tokenHash}, update)
	return err
}

func (r *mongoRefreshTokenRepo) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	update := bson.M{"$set": bson.M{"revoked": true}}
	_, err := r.collection.UpdateMany(ctx, bson.M{"user_id": userID}, update)
	return err
}
