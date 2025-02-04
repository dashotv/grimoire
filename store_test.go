package grimoire

import (
	"fmt"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestStore_Create(t *testing.T) {
	s, err := New[Download]("mongodb://localhost:27017", "seer_test", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	o := &Download{
		MediumId:  bson.NewObjectID(),
		Auto:      true,
		Multi:     false,
		Force:     false,
		Url:       "https://example.com",
		ReleaseId: "1234567890",
		Thash:     "1234567890",
	}

	err = s.Save(o)
	assert.NoError(t, err, "save")
	assert.NotNil(t, o.GetID(), "id")

	createdId = o.GetID()
	fmt.Printf("created: %s\n", createdId)
}

func TestStore_CreateWithClient(t *testing.T) {
	c, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	require.NotNil(t, c)

	s, err := NewWithClient[Download](c, "seer_test", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	o := &Download{
		MediumId:  bson.NewObjectID(),
		Auto:      true,
		Multi:     false,
		Force:     false,
		Url:       "https://example.com",
		ReleaseId: "1234567890",
		Thash:     "1234567890",
	}

	err = s.Save(o)
	assert.NoError(t, err, "save")
	assert.NotNil(t, o.GetID(), "id")

	createdId = o.GetID()
	fmt.Printf("created: %s\n", createdId)
}

func TestStore_Get(t *testing.T) {
	s, err := New[Download]("mongodb://localhost:27017", "seer_test", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	assert.False(t, createdId.IsZero(), "created id")

	o, err := s.Get(createdId.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, o)

	o2, err := s.GetByID(createdId)
	assert.NoError(t, err)
	assert.NotNil(t, o2)

	fmt.Printf("%# v\n", pretty.Formatter(o))
}

func TestStore_Find(t *testing.T) {
	s, err := New[Download]("mongodb://localhost:27017", "seer_test", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	assert.False(t, createdId.IsZero(), "created id")

	o, err := s.Find(createdId.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, o)

	fmt.Printf("%# v\n", pretty.Formatter(o))
}

func TestStore_Update(t *testing.T) {
	s, err := New[Download]("mongodb://localhost:27017", "seer_test", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	o, err := s.Find(createdId.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, o)
	//fmt.Printf("%# v\n", pretty.Formatter(o))

	o.Status = "searching"
	err = s.Update(o)
	assert.NoError(t, err)

	o2, err := s.Find(createdId.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, o2)

	assert.Equal(t, "searching", o2.Status, "status should match")
}

func TestStore_SaveUpdate(t *testing.T) {
	s, err := New[Download]("mongodb://localhost:27017", "seer_test", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	o, err := s.Find(createdId.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, o)
	//fmt.Printf("%# v\n", pretty.Formatter(o))

	o.Status = "searching"
	err = s.Update(o)
	assert.NoError(t, err)

	o2, err := s.Find(createdId.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, o2)

	assert.Equal(t, "searching", o2.Status, "status should match")
}

func TestStore_Delete(t *testing.T) {
	s, err := New[Download]("mongodb://localhost:27017", "seer_test", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	assert.False(t, createdId.IsZero(), "created id")

	d := &Download{}
	d.ID = createdId
	err = s.Delete(d)
	assert.NoError(t, err)
}

func TestStore_CountQuery(t *testing.T) {
	s, err := New[Download]("mongodb://localhost:27017", "seer_test", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	_, err = s.Query().DeleteMany()
	assert.NoError(t, err)

	list := []*Download{
		{Status: "searching", Thash: "1234567890"},
		{Status: "loading", Thash: "1234567890"},
		{Status: "managing", Thash: "1234567890"},
		{Status: "downloading", Thash: "1234567890"},
		{Status: "done", Thash: "1234567890"},
	}
	for _, d := range list {
		s.Save(d)
	}

	q, err := s.Query().Where("status", "done").Count()
	assert.NoError(t, err)
	c, err := s.Count(bson.M{"status": "done"})
	assert.NoError(t, err)
	assert.Equal(t, c, q, "download count")

	_, err = s.Query().DeleteMany()
	assert.NoError(t, err)
}

func TestStore_CountDownloads(t *testing.T) {
	s, err := New[Download]("mongodb://localhost:27017", "seer_test", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	_, err = s.Query().DeleteMany()
	assert.NoError(t, err)

	list := []*Download{
		{Status: "searching", Thash: "1234567890"},
		{Status: "loading", Thash: "1234567890"},
		{Status: "managing", Thash: "1234567890"},
		{Status: "downloading", Thash: "1234567890"},
		{Status: "done", Thash: "1234567890"},
	}

	for _, d := range list {
		s.Save(d)
	}

	count, err := s.Count(bson.M{})
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count, "download count")

	count, err = s.Count(bson.M{"status": "done"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count, "download done count")

	count, err = s.Query().Where("status", "managing").Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count, "download managing count")

	_, err = s.Query().DeleteMany()
	assert.NoError(t, err)
}

func TestStore_QueryDefaults(t *testing.T) {
	s, err := New[Download]("mongodb://localhost:27017", "seer_test", "downloads")
	assert.NoError(t, err)
	assert.NotNil(t, s)
	s.SetQueryDefaults([]bson.M{{"status": "done"}})

	_, err = s.Query().DeleteMany()
	assert.NoError(t, err)

	list := []*Download{
		{Status: "searching", Thash: "1234567890"},
		{Status: "loading", Thash: "1234567890"},
		{Status: "managing", Thash: "1234567890"},
		{Status: "downloading", Thash: "1234567890"},
		{Status: "done", Thash: "1234567890"},
	}

	for _, d := range list {
		s.Save(d)
	}

	q := s.Query().GreaterThan("created_at", time.Now().Add(-30*time.Hour*24)).Limit(100).Desc("created_at")
	list, err = q.Run()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Len(t, list, 1)

	for _, e := range list {
		assert.Equal(t, "done", e.Status, "status should match")
	}

	_, err = s.Query().DeleteMany()
	assert.NoError(t, err)
}

func TestStore_Index(t *testing.T) {
	s, err := New[Fake]("mongodb://localhost:27017", "seer_test", "fakes")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	f := &Fake{Name: "blarg"}
	err = s.Save(f)
	assert.NoError(t, err)

	CreateIndexes(s, "created_at")
	CreateIndexes(s, "name:text")
	CreateIndexesFromTags(s, Fake{})
}
