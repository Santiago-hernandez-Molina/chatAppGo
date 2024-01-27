package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

type ModelTest struct {
	Name   string
	Args   map[string]any
	Status int
}

var LoginTestCase []ModelTest = []ModelTest{
	{
		"User log-in successfully",
		map[string]any{
			"email":    "juan@gmail.com",
			"password": "12345678",
		},
		200,
	},
	{
		"User log-in Failed",
		map[string]any{
			"password": "12345678",
		},
		401,
	},
	{
		"User log-in Failed",
		map[string]any{
			"email":    "juan@gmail.com",
			"password": "1234",
		},
		401,
	},
}

var GetUsersTestCase []ModelTest = []ModelTest{
	{
		"Get 1 user",
		map[string]any{
			"username": "pedro",
			"limit":    "1",
			"offset":   "0",
			"users":    1,
		},
		200,
	},
	{
		"Get 3 users",
		map[string]any{
			"username": "",
			"limit":    "3",
			"offset":   "0",
			"users":    3,
		},
		200,
	},
	{
		"Get code 400",
		map[string]any{
			"username": "",
			"limit":    "100",
			"offset":   "100",
			"users":    0,
		},
		400,
	},
	{
		"Get more than 0 users",
		map[string]any{
			"username": "",
			"limit":    "-1",
			"offset":   "-1",
			"users":    0,
		},
		400,
	},
	{
		"Get more than 0 users",
		map[string]any{
			"username": "",
			"limit":    "0",
			"offset":   "-1",
			"users":    0,
		},
		400,
	},
	{
		"Get more than 0 users",
		map[string]any{
			"username": "",
			"limit":    "-1",
			"offset":   "0",
			"users":    0,
		},
		400,
	},
}

func TestLogin(t *testing.T) {
	for _, user := range LoginTestCase {
		w := MakeRequest("POST", "/login", user.Args, false, nil)
		assert.Equal(t, user.Status, w.Code)
	}
}

func TestGetUsers(t *testing.T) {
	for _, user := range GetUsersTestCase {
		t.Run((user.Name), func(t *testing.T) {
			w := MakeRequest("GET",
				fmt.Sprintf(
					"/v1/user/find?username=%s&limit=%s&offset=%s",
					user.Args["username"],
					user.Args["limit"],
					user.Args["offset"],
				),
				map[string]any{}, true, LoginUser,
			)
			var response models.PaginatedModel[models.UserContact]
			json.Unmarshal(w.Body.Bytes(), &response)

			usersCount := response.Data
			assert.Equal(t, user.Args["users"], len(usersCount))
			assert.Equal(t, user.Status, w.Code)
		})
	}
}
