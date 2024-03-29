package mongo

import (
	"context"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContactRequestRepo struct {
	mongoRepo  *MongoRepo
	collection *mongo.Collection
	ctx        context.Context
}

func (repo *ContactRequestRepo) UpdateRequestStatus(accepted bool, requestId int) error {
	filter := bson.D{
		{Key: "_id", Value: requestId},
	}
	update := bson.D{{
		Key: "$set",
		Value: bson.D{
			{
				Key: "accepted", Value: true,
			},
		},
	}}

	_, err := repo.collection.UpdateOne(repo.ctx, filter, update)
	return err
}

func (repo *ContactRequestRepo) GetRequestById(requestId int) (*models.ContactRequest, error) {
	filter := bson.D{
		{Key: "_id", Value: requestId},
	}
	request := models.ContactRequest{}
	result := repo.collection.FindOne(repo.ctx, filter)
	err := result.Decode(&request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (repo *ContactRequestRepo) GetRequestByToUserId(userId int, fromUserId int) (*models.ContactRequest, error) {
	filter := bson.D{
		{Key: "touserid", Value: userId},
		{Key: "fromuserid", Value: fromUserId},
	}
	request := models.ContactRequest{}
	result := repo.collection.FindOne(repo.ctx, filter)
	err := result.Decode(&request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (repo *ContactRequestRepo) GetReceivedRequests(userid int) ([]models.ContactRequestWithUser, error) {
	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "fromuserid"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}},
	}
	matchStage := bson.D{
		{
			Key: "$match",
			Value: bson.M{
				"touserid": userid,
				"accepted": false,
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

	requests := []models.ContactRequestWithUser{}
	result, err := repo.collection.Aggregate(repo.ctx, mongo.Pipeline{
		lookupStage, matchStage, fieldStage,
	})
	if err != nil {
		return nil, err
	}
	err = result.All(repo.ctx, &requests)
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (repo *ContactRequestRepo) GetSendedRequests(userid int) ([]models.ContactRequestWithUser, error) {
	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "touserid"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}},
	}
	matchStage := bson.D{
		{
			Key: "$match",
			Value: bson.M{
				"fromuserid": userid,
				"accepted":   false,
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

	requests := []models.ContactRequestWithUser{}
	result, err := repo.collection.Aggregate(repo.ctx, mongo.Pipeline{
		lookupStage, matchStage, fieldStage,
	})
	err = result.All(repo.ctx, &requests)
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (repo *ContactRequestRepo) SaveRequest(request *models.ContactRequest) error {
	request.Id = repo.mongoRepo.FindNextId(repo.ctx, "contactrequestid")
	_, err := repo.collection.InsertOne(repo.ctx, request)
	if err != nil {
		return err
	}
	return nil
}

var _ ports.ContactRequestRepo = (*ContactRequestRepo)(nil)

func NewContactRequestRepo(
	mongoRepo *MongoRepo,
	ctx context.Context,
) *ContactRequestRepo {
	collection := mongoRepo.database.Collection("contactRequests")
	return &ContactRequestRepo{
		mongoRepo:  mongoRepo,
		collection: collection,
		ctx:        ctx,
	}
}
