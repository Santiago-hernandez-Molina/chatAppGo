package mongo

import (
	"context"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ ports.MessageRepo = (*MessageRepo)(nil)

type MessageRepo struct {
	mongoRepo  *MongoRepo
	collection *mongo.Collection
	ctx        context.Context
}

func (mr *MessageRepo) GetMessagesByRoomId(roomId int) ([]models.MessageUser, error) {
	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "userid"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}},
	}
	matchStage := bson.D{
		{
			Key: "$match",
			Value: bson.M{
				"roomid": roomId,
			},
		},
	}
	fieldStage := bson.D{
		{
			Key: "$addFields",
			Value: bson.D{
				{
					Key: "user",
					Value: bson.D{
						{Key: "$arrayElemAt", Value: bson.A{"$user", 0}},
					},
				},
			},
		},
	}
	projectStage := bson.D{
		{
			Key: "$project",
			Value: bson.D{
				{
					Key: "user.password",
                    Value: 0,
				},
			},
		},
	}

	cursor, err := mr.collection.Aggregate(mr.ctx, mongo.Pipeline{
		lookupStage, matchStage, fieldStage, projectStage,
	})
	if err != nil {
		return nil, err
	}
	messages := []models.MessageUser{}
	defer cursor.Close(mr.ctx)
	cursor.All(mr.ctx, &messages)
	return messages, nil
}

func (mr *MessageRepo) SaveMessage(message *models.Message) error {
	messsageId := mr.mongoRepo.FindNextId(mr.ctx, "messageid")
	message.Id = messsageId
	_, err := mr.collection.InsertOne(mr.ctx, message)
	if err != nil {
		return err
	}
	return nil
}

func NewMessageRepo(mongoRepo *MongoRepo, ctx context.Context) *MessageRepo {
	collection := mongoRepo.database.Collection("messages")

	return &MessageRepo{
		mongoRepo:  mongoRepo,
		collection: collection,
		ctx:        ctx,
	}
}
