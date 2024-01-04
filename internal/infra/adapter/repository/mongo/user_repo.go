package mongo

import (
	"context"
	"errors"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	mongoRepo  *MongoRepo
	collection *mongo.Collection
	ctx        context.Context
}

func (repo *UserRepo) GetUserByEmail(user *models.User) (*models.User, error) {
	userDB := models.User{}
	filter := bson.D{{Key: "email", Value: user.Email}}
	result := repo.collection.FindOne(context.TODO(), filter)
	err := result.Decode(&userDB)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("no user found")
	}
	return &userDB, nil
}

func (ur *UserRepo) Register(user *models.User) error {
	userId := ur.mongoRepo.FindNextId(ur.ctx, "userid")
	user.Id = userId
	_, err := ur.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	} else {
		return nil
	}
}

var _ ports.UserRepo = (*UserRepo)(nil)

func NewUserRepo(mongoRepo *MongoRepo, ctx context.Context) *UserRepo {
	collection := mongoRepo.database.Collection("users")

	return &UserRepo{
		collection: collection,
		mongoRepo:  mongoRepo,
		ctx:        ctx,
	}
}
