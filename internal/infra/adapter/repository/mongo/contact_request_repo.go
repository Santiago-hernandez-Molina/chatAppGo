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

func (repo *ContactRequestRepo) GetRequestByToUserId(requestId int, userId int) (*models.ContactRequest, error) {
	filter := bson.D{
		{Key: "_id", Value: requestId},
		{Key: "touserid", Value: userId},
	}
	request := models.ContactRequest{}
	result := repo.collection.FindOne(repo.ctx, filter)
	err := result.Decode(&request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (*ContactRequestRepo) GetReceivedRequests(userid int) ([]models.ContactRequest, error) {
	panic("unimplemented")
}

func (*ContactRequestRepo) GetSendedRequests(userid int) ([]models.ContactRequest, error) {
	panic("unimplemented")
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
