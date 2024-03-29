package mongo

import (
	"context"
	"errors"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/exceptions"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RoomRepo struct {
	ctx        context.Context
	collection *mongo.Collection
	mongoRepo  *MongoRepo
}

func (repo *RoomRepo) AddUserToRoom(userId int, roomId int) error {
	filter := bson.D{{Key: "_id", Value: roomId}}
	update := bson.D{{
		Key: "$push",
		Value: bson.D{
			{
				Key: "users", Value: bson.D{
					{Key: "userid", Value: userId},
					{Key: "role", Value: "user"},
				},
			},
		},
	}}

	_, err := repo.collection.UpdateOne(repo.ctx, filter, update)
	if err == nil {
		return nil
	}
	if err == mongo.ErrNoDocuments {
		return errors.New("Room Not Found")
	}
	return &exceptions.AccesDataException{}
}

var _ ports.RoomRepo = (*RoomRepo)(nil)

func (repo *RoomRepo) GetRoomById(roomId int) (*models.Room, error) {
	filter := bson.D{{Key: "_id", Value: roomId}}
	result := repo.collection.FindOne(repo.ctx, filter)
	var room models.Room
	err := result.Decode(&room)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}

	return &room, nil
}

func (repo *RoomRepo) NewRoom(room *models.Room) error {
	room.Id = repo.mongoRepo.FindNextId(repo.ctx, "roomid")
	_, err := repo.collection.InsertOne(repo.ctx, room)
	return err
}

func (repo *RoomRepo) GetRoomsByUserId(userId int) ([]models.Room, error) {
	filter := bson.D{{Key: "users.userid", Value: userId}}
	rooms := []models.Room{}
	cursor, err := repo.collection.Find(repo.ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(repo.ctx)
	err = cursor.All(repo.ctx, &rooms)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (repo *RoomRepo) GetUserRoom(userId int, roomId int) (*models.UserRoom, error) {
	filter := bson.M{
		"_id":          roomId,
		"users.userid": userId,
	}
	projection := bson.M{
		"_id":   0,
		"users": bson.M{"$elemMatch": bson.M{"userid": userId}},
	}
	var room models.Room
	err := repo.collection.FindOne(
		repo.ctx,
		filter,
		options.FindOne().SetProjection(projection),
	).Decode(&room)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &exceptions.UserNotFound{}
		}
		return nil, err
	}

	return &room.Users[0], nil
}

func NewRoomRepo(mongoRepo *MongoRepo, ctx context.Context) *RoomRepo {
	collection := mongoRepo.database.Collection("rooms")

	return &RoomRepo{
		collection: collection,
		mongoRepo:  mongoRepo,
		ctx:        ctx,
	}
}
