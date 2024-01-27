package datamongo

import (
	"context"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/gateways/authentication"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client   *mongo.Client
	ctx      context.Context
	database *mongo.Database
)

func InitDB(uri string, databaseName string) {
	ctx = context.TODO()
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	c, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	client = c
	database = client.Database(databaseName)
	createCollections(ctx)
}

func CleanDB() {
	database.Collection("users").Drop(ctx)
	database.Collection("messages").Drop(ctx)
	database.Collection("contactRequests").Drop(ctx)
	database.Collection("rooms").Drop(ctx)
	database.Collection("counters").Drop(ctx)

	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func createCollections(context context.Context) {
	passwordManager := authentication.NewPasswordManager()
	password, _ := passwordManager.EncryptPassword("12345678")
	database.Collection("users").InsertMany(
		context,
		[]any{
			models.User{Id: 0, Username: "Juan", Email: "juan@gmail.com", Password: password, Status: true, Code: 0},
			models.User{Id: 1, Username: "Pedro", Email: "pedro@gmail.com", Password: password, Status: true, Code: 0},
			models.User{Id: 2, Username: "Pepe", Email: "pepe@gmail.com", Password: password, Status: true, Code: 0},
			models.User{Id: 3, Username: "FooBar", Email: "foobar@gmail.com", Password: password, Status: true, Code: 0},
		},
	)
	database.Collection("rooms").InsertMany(
		context,
		[]any{
			models.Room{Id: 0, Name: "CHAT 2", Users: []models.UserRoom{
				{UserId: 0, Role: "admin"},
				{UserId: 1, Role: "user"},
				{UserId: 2, Role: "user"},
			}, Type: models.RoomType(models.Group)},
		},
	)
	database.Collection("counters").InsertMany(
		context,
		[]any{
			map[string]any{"_id": "contactrequestid", "seq": 0},
			map[string]any{"_id": "roomid", "seq": 1},
			map[string]any{"_id": "userid", "seq": 3},
		},
	)
}
