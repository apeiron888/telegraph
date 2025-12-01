package messages

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepo interface {
	Create(ctx context.Context, m *Message) error
	GetByID(ctx context.Context, id uuid.UUID) (*Message, error)
	GetByChannelID(ctx context.Context, channelID uuid.UUID, limit, offset int) ([]*Message, error)
	Update(ctx context.Context, m *Message) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountAfter(ctx context.Context, channelID uuid.UUID, after time.Time) (int64, error)
}

type mongoMessageRepo struct {
	collection *mongo.Collection
}

func NewMongoMessageRepo(db *mongo.Database) MessageRepo {
	return &mongoMessageRepo{
		collection: db.Collection("messages"),
	}
}

func (r *mongoMessageRepo) Create(ctx context.Context, m *Message) error {
	m.ID = uuid.New()
	m.Timestamp = time.Now()
	m.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, m)
	return err
}

func (r *mongoMessageRepo) GetByID(ctx context.Context, id uuid.UUID) (*Message, error) {
	var message Message
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&message)
	if err == mongo.ErrNoDocuments {
		return nil, ErrMessageNotFound
	}
	return &message, err
}

func (r *mongoMessageRepo) GetByChannelID(ctx context.Context, channelID uuid.UUID, limit, offset int) ([]*Message, error) {
	filter := bson.M{
		"channel_id": channelID,
		"deleted":    false,
	}
	
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *mongoMessageRepo) Update(ctx context.Context, m *Message) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"id": m.ID}, m)
	return err
}

func (r *mongoMessageRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	update := bson.M{"$set": bson.M{"deleted": true}}
	_, err := r.collection.UpdateOne(ctx, bson.M{"id": id}, update)
	return err
}

func (r *mongoMessageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r *mongoMessageRepo) CountAfter(ctx context.Context, channelID uuid.UUID, after time.Time) (int64, error) {
	filter := bson.M{
		"channel_id": channelID,
		"timestamp":  bson.M{"$gt": after},
		"deleted":    false,
	}
	return r.collection.CountDocuments(ctx, filter)
}
