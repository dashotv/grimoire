package grimoire

import (
	"fmt"
	"testing"
	"time"

	"github.com/kr/pretty"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Download struct {
	Document `bson:",inline"` // include mgm.DefaultModel
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	MediumId   primitive.ObjectID `json:"medium_id" bson:"medium_id"`
	Auto       bool               `json:"auto" bson:"auto"`
	Multi      bool               `json:"multi" bson:"multi"`
	Force      bool               `json:"force" bson:"force"`
	Url        string             `json:"url" bson:"url"`
	ReleaseId  string             `json:"release_id" bson:"tdo_id"`
	Thash      string             `json:"thash" bson:"thash"`
	Timestamps struct {
		Found      time.Time `json:"found" bson:"found"`
		Loaded     time.Time `json:"loaded" bson:"loaded"`
		Downloaded time.Time `json:"downloaded" bson:"downloaded"`
		Completed  time.Time `json:"completed" bson:"completed"`
		Deleted    time.Time `json:"deleted" bson:"deleted"`
	} `json:"timestamps" bson:"timestamps"`
	Selected string `json:"selected" bson:"selected"`
	Status   string `json:"status" bson:"status"`
	Files    []struct {
		Id       primitive.ObjectID `json:"id" bson:"_id"`
		MediumId primitive.ObjectID `json:"medium_id" bson:"medium_id"`
		Num      int                `json:"num" bson:"num"`
	} `json:"download_files" bson:"download_files"`
}

func TestStore_Query(t *testing.T) {
	s, err := New[*Download]("mongodb://localhost:27017", "seer_development", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	q := s.Query()
	list, err := q.In("status", []string{"searching", "loading", "managing", "downloading", "reviewing"}).Run()
	assert.NoError(t, err)
	assert.NotNil(t, list)

	//fmt.Printf("%# v\n", pretty.Formatter(list))
	for _, e := range list {
		fmt.Printf("download: %s\n", e.ID.Hex())
	}
}

func TestStore_Find(t *testing.T) {
	s, err := New[*Download]("mongodb://localhost:27017", "seer_development", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	o := &Download{}
	err = s.Find("62f661903359bbbe05a5bb2c", o)
	assert.NoError(t, err)
	assert.NotNil(t, o)

	fmt.Printf("%# v\n", pretty.Formatter(o))
}
