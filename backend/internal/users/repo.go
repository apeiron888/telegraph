package users

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo interface {
	Create(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByEmailOrPhone(ctx context.Context, identifier string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, query string) ([]*User, error)
}

type mongoUserRepo struct {
	collection *mongo.Collection
}

func NewMongoUserRepo(db *mongo.Database) UserRepo {
	return &mongoUserRepo{
		collection: db.Collection("users"),
	}
}

func (r *mongoUserRepo) Create(ctx context.Context, u *User) error {
	u.ID = uuid.New()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, u)
	return err
}

func (r *mongoUserRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepo) Update(ctx context.Context, u *User) error {
	u.UpdatedAt = time.Now()
	
	update := bson.M{
		"$set": bson.M{
			"username":       u.Username,
			"bio":            u.Bio,
			"birth_date":     u.BirthDate,
			"country":        u.Country,
			"city":           u.City,
			"street":         u.Street,
			"account_type":   u.AccountType,
			"security_label": u.SecurityLabel,
			"attributes":     u.Attributes,
			"updated_at":     u.UpdatedAt,
		},
	}
	
	_, err := r.collection.UpdateOne(ctx, bson.M{"id": u.ID}, update)
	return err
}

func (r *mongoUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r *mongoUserRepo) GetByEmailOrPhone(ctx context.Context, identifier string) (*User, error) {
	var user User
	// Try email first
	filter := bson.M{"$or": []bson.M{
		{"email": identifier},
		{"phone": identifier},
	}}
	
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepo) Search(ctx context.Context, query string) ([]*User, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"username": bson.M{"$regex": query, "$options": "i"}},
			{"email": bson.M{"$regex": query, "$options": "i"}},
			{"phone": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

