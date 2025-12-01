package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



type MFACode struct {
	UserID    uuid.UUID `bson:"_id"`
	CodeHash  string    `bson:"code_hash"`
	ExpiresAt time.Time `bson:"expires_at"`
}

type mongoMFACodeRepo struct {
	collection *mongo.Collection
}

func NewMFACodeRepo(db *mongo.Database) MFARepo {
	return &mongoMFACodeRepo{
		collection: db.Collection("mfa_codes"),
	}
}

func (r *mongoMFACodeRepo) Store(ctx context.Context, userID uuid.UUID, codeHash string, expiresAt time.Time) error {
	code := MFACode{
		UserID:    userID,
		CodeHash:  codeHash,
		ExpiresAt: expiresAt,
	}
	
	// Upsert - replace if exists
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": code}
	opts := options.Update().SetUpsert(true)
	
	_, err := r.collection.UpdateOne(ctx, filter, update,opts)
	return err
}

func (r *mongoMFACodeRepo) Find(ctx context.Context, userID uuid.UUID) (string, time.Time, error) {
	var code MFACode
	err := r.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&code)
	if err == mongo.ErrNoDocuments {
		return "", time.Time{}, nil
	}
	if err != nil {
		return "", time.Time{}, err
	}
	return code.CodeHash, code.ExpiresAt, nil
}

func (r *mongoMFACodeRepo) Delete(ctx context.Context, userID uuid.UUID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": userID})
	return err
}
