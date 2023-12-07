package grimoire

import (
	"fmt"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestStore_Create(t *testing.T) {
	s, err := New[*Download]("mongodb://localhost:27017", "seer_development", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	o := &Download{
		MediumId:  primitive.NewObjectID(),
		Auto:      true,
		Multi:     false,
		Force:     false,
		Url:       "https://example.com",
		ReleaseId: "1234567890",
		Thash:     "1234567890",
	}

	err = s.Save(o)
	assert.NoError(t, err, "save")
	assert.NotNil(t, o.ID, "id")

	createdId = o.ID
}

func TestStore_CreateWithTransaction(t *testing.T) {
	s, err := New[*Download]("mongodb://localhost:27017", "seer_development", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	o := &Download{
		MediumId:  primitive.NewObjectID(),
		Auto:      true,
		Multi:     false,
		Force:     false,
		Url:       "https://example.com",
		ReleaseId: "1234567890",
		Thash:     "1234567890",
	}

	err = s.Save(o)
	assert.NoError(t, err, "save")
	assert.NotNil(t, o.ID, "id")
}

func TestStore_Get(t *testing.T) {
	s, err := New[*Download]("mongodb://localhost:27017", "seer_development", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	assert.False(t, createdId.IsZero(), "created id")

	o, err := s.Get(createdId, &Download{})
	assert.NoError(t, err)
	assert.NotNil(t, o)

	fmt.Printf("%# v\n", pretty.Formatter(o))
}

func TestStore_Find(t *testing.T) {
	s, err := New[*Download]("mongodb://localhost:27017", "seer_development", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	assert.False(t, createdId.IsZero(), "created id")

	o := &Download{}
	err = s.Find(createdId.Hex(), o)
	assert.NoError(t, err)
	assert.NotNil(t, o)

	fmt.Printf("%# v\n", pretty.Formatter(o))
}

func TestStore_Update(t *testing.T) {
	s, err := New[*Download]("mongodb://localhost:27017", "seer_development", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	o := &Download{}
	err = s.Find(createdId.Hex(), o)
	assert.NoError(t, err)
	assert.NotNil(t, o)
	//fmt.Printf("%# v\n", pretty.Formatter(o))

	o.Status = "searching"
	err = s.Update(o)
	assert.NoError(t, err)

	o2 := &Download{}
	err = s.Find(createdId.Hex(), o2)
	assert.NoError(t, err)
	assert.NotNil(t, o)

	assert.Equal(t, "searching", o.Status, "status should match")
}

func TestStore_SaveUpdate(t *testing.T) {
	s, err := New[*Download]("mongodb://localhost:27017", "seer_development", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	o := &Download{}
	err = s.Find(createdId.Hex(), o)
	assert.NoError(t, err)
	assert.NotNil(t, o)
	//fmt.Printf("%# v\n", pretty.Formatter(o))

	o.Status = "searching"
	err = s.Update(o)
	assert.NoError(t, err)

	o2 := &Download{}
	err = s.Find(createdId.Hex(), o2)
	assert.NoError(t, err)
	assert.NotNil(t, o)

	assert.Equal(t, "searching", o.Status, "status should match")
}

func TestStore_Delete(t *testing.T) {
	s, err := New[*Download]("mongodb://localhost:27017", "seer_development", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	assert.False(t, createdId.IsZero(), "created id")

	d := &Download{}
	d.ID = createdId
	err = s.Delete(d)
	assert.NoError(t, err)
}

func TestStore_CountQuery(t *testing.T) {
	s, err := New[*Download]("mongodb://localhost:27017", "seer_development", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	q, err := s.Query().Where("status", "done").Count()
	assert.NoError(t, err)
	c, err := s.Count(bson.M{"status": "done"})
	assert.NoError(t, err)
	assert.Equal(t, c, q, "download count")
}

func TestStore_CountDownloads(t *testing.T) {
	s, err := New[*Download]("mongodb://localhost:27017", "seer_development", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	count, err := s.Count(bson.M{})
	assert.NoError(t, err)
	assert.Equal(t, int64(TOTAL_DOWNLOADS), count, "download count")
}

func TestStore_CountSeries(t *testing.T) {
	s, err := New[*Medium]("mongodb://localhost:27017", "seer_development", "media")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	count, err := s.Count(bson.M{"_type": "Series"})
	assert.NoError(t, err)
	assert.Equal(t, int64(TOTAL_SERIES), count, "series count")
}
