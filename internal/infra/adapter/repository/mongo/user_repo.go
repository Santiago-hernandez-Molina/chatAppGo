package mongo

import (
	"context"
	"errors"
	"log"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/exceptions"
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

func (repo *UserRepo) ActivateAccount(code int, email string) error {
	filter := bson.D{
		{Key: "email", Value: email},
		{Key: "code", Value: code},
	}
	update := bson.D{{
		Key:   "$set",
		Value: bson.D{{Key: "status", Value: true}},
	}}
	result, err := repo.collection.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return err
	}
    if result.MatchedCount == 0 {
        return errors.New("Cannot found the user")
    }

	return nil
}

func (repo *UserRepo) DeleteInactiveUser(email string) error {
	filter := bson.D{
		{Key: "email", Value: email},
		{Key: "status", Value: false},
	}
	result, err := repo.collection.DeleteOne(repo.ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	log.Println(result)
	log.Println(result.DeletedCount)
	return nil
}

func (*UserRepo) DeleteUser(userId int) error {
	panic("unimplemented")
}

func (repo *UserRepo) GetUserByEmail(user *models.User) (*models.User, error) {
	userDB := models.User{}
	filter := bson.D{{Key: "email", Value: user.Email}}
	result := repo.collection.FindOne(repo.ctx, filter)
	err := result.Decode(&userDB)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("no user found")
	}
	return &userDB, nil
}

func (ur *UserRepo) Register(user *models.User) error {
	userId := ur.mongoRepo.FindNextId(ur.ctx, "userid")
	user.Id = userId
	_, err := ur.collection.InsertOne(ur.ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return &exceptions.DuplicatedUser{}
		}
		return err
	}
	return nil
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
