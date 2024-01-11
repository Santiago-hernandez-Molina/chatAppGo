package test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ContactTest struct {
	Name   string
	Args   map[string]int
	Status int
}

var ContactTestCase []ContactTest = []ContactTest{
	{
		"Add contact",
		map[string]int{
			"toUserId": 1,
		},
		200,
	},
	{
		"Add contact",
		map[string]int{
			"toUserId": 2,
		},
		200,
	},
	{
		"Add contact",
		map[string]int{
			"toUserId": 2,
		},
		404,
	},
	{
		"Add contact",
		map[string]int{
			"toUserId": 10,
		},
		404,
	},
}

func TestContact(t *testing.T) {
	t.Run("Send Request", func(t *testing.T) {
		for _, request := range ContactTestCase {
			w := MakeRequest[map[string]int]("POST", "/v1/contact/new", request.Args, true, LoginUser)
			assert.Equal(t, request.Status, w.Code)
		}
	})
	t.Run("Get sended requests test", func(t *testing.T) {
		w := MakeRequest[map[string]string]("POST", "/v1/contact/find/sended", nil, true, LoginUser)
		var response []map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 2, len(response))
	})
	t.Run("Get received requests test", func(t *testing.T) {
		w := MakeRequest[map[string]string]("POST", "/v1/contact/find/received", nil, true, LoginUser2)
		var response []map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, len(response))
	})
	t.Run("Accept Request", func(t *testing.T) {
		w := MakeRequest[map[string]string]("POST", "/v1/contact/accept/0", nil, true, LoginUser2)
		assert.Equal(t, 202, w.Code)
	})
	t.Run("View Contact Room", func(t *testing.T) {
		w := MakeRequest[map[string]string]("GET", "/v1/room/show/1", nil, true, LoginUser2)
		assert.Equal(t, 200, w.Code)
		var response map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Juan", response["name"])
	})
	t.Run("View Contact Room 2", func(t *testing.T) {
		w := MakeRequest[map[string]string]("GET", "/v1/room/show/1", nil, true, LoginUser)
		assert.Equal(t, 200, w.Code)
		var response map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Pedro", response["name"])
	})
}
