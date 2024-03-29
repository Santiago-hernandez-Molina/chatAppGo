package mongo

import (
	"context"
	"errors"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/exceptions"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	mongoRepo  *MongoRepo
	collection *mongo.Collection
	ctx        context.Context
}

func (repo *UserRepo) GetUsersCount(filter string) (int, error) {
	regexFilter := primitive.Regex{Pattern: filter, Options: "i"} // "i" para hacer la búsqueda sin distinción entre mayúsculas y minúsculas
	filterQuery := bson.M{"username": regexFilter}

	count, err := repo.collection.CountDocuments(repo.ctx, filterQuery)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserRepo) GetUsersByUsername(userId int, filter string, size int, offset int) (*models.PaginatedModel[models.UserContact], error) {
	regexFilter := primitive.Regex{Pattern: filter, Options: "i"}
	filterQuery := bson.M{"username": regexFilter}
	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": bson.M{"$ne": userId}},
		},
		{
			"$lookup": bson.M{
				"from": "contactRequests",
				"let":  bson.M{"userid": "$_id"},
				"pipeline": []bson.M{{
					"$match": bson.M{"$expr": bson.M{
						"$or": bson.A{
							bson.D{{Key: "$and", Value: bson.A{
								bson.M{"$eq": bson.A{"$fromuserid", userId}},
								bson.M{"$eq": bson.A{"$touserid", "$$userid"}},
							}}},
							bson.D{{Key: "$and", Value: bson.A{
								bson.M{"$eq": bson.A{"$fromuserid", "$$userid"}},
								bson.M{"$eq": bson.A{"$touserid", userId}},
							}}},
						},
					}},
				}},
				"as": "joined",
			},
		},
		{
			"$match": bson.M{"$and": bson.A{filterQuery, bson.M{"joined": []any{}}}},
		},
		{
			"$project": bson.M{
				"_id":      1,
				"username": 1,
			},
		},
		{
			"$skip": offset,
		},
		{
			"$limit": size,
		},
	}

	users := []models.UserContact{}

	cursor, err := repo.collection.Aggregate(repo.ctx, pipeline)
	defer cursor.Close(repo.ctx)
	if err != nil {
		return nil, err
	}

	err = cursor.All(repo.ctx, &users)
	if err != nil {
		return nil, err
	}
	paginatedModel := models.PaginatedModel[models.UserContact]{
		Offset: offset,
		Limit:  size,
		Data:   users,
	}

	return &paginatedModel, nil
}

func (repo *UserRepo) GetUserById(userId int) (*models.User, error) {
	filter := bson.D{{Key: "_id", Value: userId}}
	result := repo.collection.FindOne(repo.ctx, filter)
	user := models.User{}
	err := result.Decode(&user)
	if err == nil {
		return &user, nil
	}
	if err == mongo.ErrNoDocuments {
		return nil, &exceptions.UserNotFound{}
	}
	return nil, err
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

func (repo *UserRepo) DeleteUserByEmailAndStatus(email string, status bool) error {
	filter := bson.D{
		{Key: "email", Value: email},
		{Key: "status", Value: status},
	}
	_, err := repo.collection.DeleteOne(repo.ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	return nil
}

func (repo *UserRepo) DeleteUser(userId int) error {
	filter := bson.D{{Key: "_id", Value: userId}}
	result, err := repo.collection.DeleteOne(repo.ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("User not found")
	}
	return nil
}

func (repo *UserRepo) GetUserByEmail(user *models.User) (*models.User, error) {
	userDB := models.User{}
	filter := bson.D{{Key: "email", Value: user.Email}}
	result := repo.collection.FindOne(repo.ctx, filter)
	err := result.Decode(&userDB)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no user found")
		}
		return nil, err
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
