package grimoire

import (
	"fmt"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
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

type Medium struct {
	Document `bson:",inline"` // include mgm.DefaultModel
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Type         string           `json:"type" bson:"_type"`
	Kind         primitive.Symbol `json:"kind" bson:"kind"`
	Source       string           `json:"source" bson:"source"`
	SourceId     string           `json:"source_id" bson:"source_id"`
	Title        string           `json:"title" bson:"title"`
	Description  string           `json:"description" bson:"description"`
	Slug         string           `json:"slug" bson:"slug"`
	Text         []string         `json:"text" bson:"text"`
	Display      string           `json:"display" bson:"display"`
	Directory    string           `json:"directory" bson:"directory"`
	Search       string           `json:"search" bson:"search"`
	SearchParams struct {
		Type       string `json:"type" bson:"type"`
		Verified   bool   `json:"verified" bson:"verified"`
		Group      string `json:"group" bson:"group"`
		Author     string `json:"author" bson:"author"`
		Resolution int    `json:"resolution" bson:"resolution"`
		Source     string `json:"source" bson:"source"`
		Uncensored bool   `json:"uncensored" bson:"uncensored"`
		Bluray     bool   `json:"bluray" bson:"bluray"`
	} `json:"search_params" bson:"search_params"`
	Active      bool      `json:"active" bson:"active"`
	Downloaded  bool      `json:"downloaded" bson:"downloaded"`
	Completed   bool      `json:"completed" bson:"completed"`
	Skipped     bool      `json:"skipped" bson:"skipped"`
	Watched     bool      `json:"watched" bson:"watched"`
	Broken      bool      `json:"broken" bson:"broken"`
	ReleaseDate time.Time `json:"release_date" bson:"release_date"`
	Paths       []struct {
		Id        primitive.ObjectID `json:"id" bson:"_id"`
		Type      primitive.Symbol   `json:"type" bson:"type"`
		Remote    string             `json:"remote" bson:"remote"`
		Local     string             `json:"local" bson:"local"`
		Extension string             `json:"extension" bson:"extension"`
		Size      int                `json:"size" bson:"size"`
		UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	} `json:"paths" bson:"paths"`
	Cover      string `json:"cover" bson:"cover"`
	Background string `json:"background" bson:"background"`
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

func TestStore_Save(t *testing.T) {
	s, err := New[*Download]("mongodb://localhost:27017", "seer_development", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	o := &Download{}
	err = s.Find("62f661903359bbbe05a5bb2c", o)
	assert.NoError(t, err)
	assert.NotNil(t, o)
	//fmt.Printf("%# v\n", pretty.Formatter(o))

	o.Status = "searching"
	err = s.Update(o)
	assert.NoError(t, err)

	o2 := &Download{}
	err = s.Find("62f661903359bbbe05a5bb2c", o2)
	assert.NoError(t, err)
	assert.NotNil(t, o)

	assert.Equal(t, "searching", o.Status, "status should match")
}

func TestStore_CountDownloads(t *testing.T) {
	s, err := New[*Download]("mongodb://localhost:27017", "seer_development", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	count, err := s.Count(bson.M{})
	assert.NoError(t, err)
	assert.Equal(t, int64(795), count, "download count")
}

func TestStore_CountSeries(t *testing.T) {
	s, err := New[*Medium]("mongodb://localhost:27017", "seer_development", "media")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	count, err := s.Count(bson.M{"_type": "Series"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1293), count, "series count")
}
