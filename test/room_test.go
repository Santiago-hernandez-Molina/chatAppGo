package test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type RoomTest struct {
	Name   string
	user   map[string]any
	count  int
	Status int
}

var RoomTestCase []RoomTest = []RoomTest{
	{
		"GetAllRooms",
		LoginUser,
		2,
		200,
	},
	{
		"GetAllRooms (With length 0)",
		LoginUser3,
		0,
		200,
	},
}

func TestRoom(t *testing.T) {
	t.Run(("View all Rooms"), func(t *testing.T) {
		for _, room := range RoomTestCase {
			w := MakeRequest("GET", "/v1/room/find", map[string]any{}, true, room.user)
			var response []map[string]any
			json.Unmarshal(w.Body.Bytes(), &response)

			assert.Equal(t, room.Status, w.Code)
			assert.Equal(t, room.count, len(response))
		}
	})
	t.Run("View Contact Room first user", func(t *testing.T) {
		w := MakeRequest[map[string]any]("GET", "/v1/room/show/1", nil, true, LoginUser2)
		var response map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "Juan", response["name"])
	})
	t.Run("View Contact Room second user", func(t *testing.T) {
		w := MakeRequest[map[string]any]("GET", "/v1/room/show/1", nil, true, LoginUser)
		var response map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "Pedro", response["name"])
	})
}
