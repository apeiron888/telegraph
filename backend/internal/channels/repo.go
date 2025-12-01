package channels

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChannelRepo interface {
	Create(ctx context.Context, c *Channel) error
	GetByID(ctx context.Context, id uuid.UUID) (*Channel, error)
	GetUserChannels(ctx context.Context, userID uuid.UUID) ([]*Channel, error)
	AddMember(ctx context.Context, channelID, userID uuid.UUID) error
	RemoveMember(ctx context.Context, channelID, userID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, channelID, userID uuid.UUID, role string) error
	IsMember(ctx context.Context, channelID, userID uuid.UUID) (bool, error)
	Update(ctx context.Context, c *Channel) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type mongoChannelRepo struct {
	collection *mongo.Collection
}

func NewMongoChannelRepo(db *mongo.Database) ChannelRepo {
	return &mongoChannelRepo{
		collection: db.Collection("channels"),
	}
}

func (r *mongoChannelRepo) Create(ctx context.Context, c *Channel) error {
	c.ID = uuid.New()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	// Ensure members are correctly structured if passed (though service usually handles this)
	// But here we just insert what's given. The service must construct the ChannelMember list.
	_, err := r.collection.InsertOne(ctx, c)
	return err
}

func (r *mongoChannelRepo) GetByID(ctx context.Context, id uuid.UUID) (*Channel, error) {
	var channel Channel
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&channel)
	if err == mongo.ErrNoDocuments {
		return nil, ErrChannelNotFound
	}
	return &channel, err
}

func (r *mongoChannelRepo) GetUserChannels(ctx context.Context, userID uuid.UUID) ([]*Channel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"members.user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var channels []*Channel
	if err := cursor.All(ctx, &channels); err != nil {
		return nil, err
	}
	return channels, nil
}

func (r *mongoChannelRepo) AddMember(ctx context.Context, channelID, userID uuid.UUID) error {
	member := ChannelMember{
		UserID:   userID,
		Role:     ChannelRoleMember,
		JoinedAt: time.Now(),
	}
	update := bson.M{"$addToSet": bson.M{"members": member}, "$set": bson.M{"updated_at": time.Now()}}
	_, err := r.collection.UpdateOne(ctx, bson.M{"id": channelID}, update)
	return err
}

func (r *mongoChannelRepo) RemoveMember(ctx context.Context, channelID, userID uuid.UUID) error {
	update := bson.M{"$pull": bson.M{"members": bson.M{"user_id": userID}}, "$set": bson.M{"updated_at": time.Now()}}
	_, err := r.collection.UpdateOne(ctx, bson.M{"id": channelID}, update)
	return err
}

func (r *mongoChannelRepo) IsMember(ctx context.Context, channelID, userID uuid.UUID) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"id": channelID, "members.user_id": userID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *mongoChannelRepo) Update(ctx context.Context, c *Channel) error {
	c.UpdatedAt = time.Now()
	
	update := bson.M{"$set": bson.M{
		"name":           c.Name,
		"description":    c.Description,
		"permissions":    c.Permissions,
		"security_label": c.SecurityLabel,
		"updated_at":     c.UpdatedAt,
	}}
	
	_, err := r.collection.UpdateOne(ctx, bson.M{"id": c.ID}, update)
	return err
}

func (r *mongoChannelRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}


func (r *mongoChannelRepo) UpdateMemberRole(ctx context.Context, channelID, userID uuid.UUID, role string) error {
	filter := bson.M{"id": channelID, "members.user_id": userID}
	update := bson.M{"$set": bson.M{"members.$.role": role, "updated_at": time.Now()}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
