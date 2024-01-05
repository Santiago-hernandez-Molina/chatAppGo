package mongo

import (
	"context"
	"time"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/adapter/repository/mongo/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MongoTimeOut = 200
)

type MongoRepo struct {
	client   *mongo.Client
	database *mongo.Database
}

func (repo *MongoRepo) FindNextId(ctx context.Context, counterName string) int {
	counters := repo.database.Collection("counters")
	filter := bson.D{{Key: "_id", Value: counterName}}
	update := bson.D{{
		Key:   "$inc",
		Value: bson.D{{Key: "seq", Value: 1}},
	}}
	counter := data.Counter{}
	result := counters.FindOneAndUpdate(ctx, filter, update)
    result.Decode(&counter)
	return counter.Seq
}

func NewRepo(conn string) (*MongoRepo, error) {
	ctx, cancelF := context.WithTimeout(
		context.Background(),
		MongoTimeOut*time.Second,
	)
	defer cancelF()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		conn,
	))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	database := client.Database("chatApp")

	return &MongoRepo{
		client:   client,
		database: database,
	}, nil
}