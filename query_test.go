package grimoire

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const TOTAL_DOWNLOADS = 745
const TOTAL_SERIES = 1533

var createdId primitive.ObjectID

type Fake struct {
	Document `bson:",inline"` // include mgm.DefaultModel
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Name string `json:"name" bson:"name" grimoire:"index"`
	Age  int    `json:"age" bson:"age" grimoire:"index,desc"`
}

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

func TestStore_QueryMedium(t *testing.T) {
	s, err := New[*Medium]("mongodb://localhost:27017", "seer_development", "media")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	tomorrow := time.Now().Add(time.Hour * 48)
	yesterday := time.Now().Add(-time.Hour * 48)

	q := s.Query()
	list, err := q.
		Where("_type", "Episode").
		LessThan("release_date", tomorrow).
		GreaterThan("release_date", yesterday).
		Asc("release_date").
		Run()
	assert.NoError(t, err)
	assert.NotNil(t, list)

	fmt.Printf("between: %s and %s\n", yesterday.Format("2006-01-02T15:04:05.000Z"), tomorrow.Format("2006-01-02T15:04:05.000Z"))
	fmt.Printf("## weird off by 1 issue, but it works\n")
	for _, e := range list {
		assert.LessOrEqual(t, e.ReleaseDate, tomorrow)
		assert.GreaterOrEqual(t, e.ReleaseDate, yesterday)
		fmt.Printf("episode: %s: %s\n", e.ID.Hex(), e.ReleaseDate.Format("2006-01-02T15:04:05.000Z"))
	}
}
func TestStore_QueryEmpty(t *testing.T) {
	s, err := New[*Medium]("mongodb://localhost:27017", "seer_development", "media")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	q := s.Query()
	list, err := q.
		Asc("release_date").
		Run()
	assert.NoError(t, err)
	assert.NotNil(t, list)
}

func TestStore_QueryOr(t *testing.T) {
	s, err := New[*Medium]("mongodb://localhost:27017", "seer_development", "media")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	q := s.Query().Or(func(qq *QueryBuilder[*Medium]) {
		qq.Where("_type", "Series").Where("_type", "Movie")
	})
	list, err := q.Asc("release_date").Run()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	for _, e := range list {
		fmt.Printf("medium: %s %s\n", e.Type, e.Title)
	}
}

func TestStore_ComplexOr(t *testing.T) {
	s, err := New[*Medium]("mongodb://localhost:27017", "seer_development", "media")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	q := s.Query().ComplexOr(func(qq *QueryBuilder[*Medium], qr *QueryBuilder[*Medium]) {
		qq.Where("_type", "Movie").Where("kind", "movies3d").Where("title", "Up")
		qr.Where("_type", "Series").Where("kind", "donghua").Where("title", "The Great Ruler")
	})
	list, err := q.Asc("release_date").Run()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	for _, e := range list {
		fmt.Printf("medium: %s %s %s\n", e.Type, e.Kind, e.Title)
	}
}

func TestStore_QueryLimit(t *testing.T) {
	s, err := New[*Download]("mongodb://localhost:27017", "seer_development", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	q := s.Query()
	list, err := q.
		Asc("release_date").
		Limit(1).
		Run()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))

	q = s.Query()
	list, err = q.
		Limit(-1).
		Run()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Greater(t, len(list), 25, "should be more than 25")
}
