package grimoire

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store[OUTPUT mgm.Model] struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mgm.Collection
}

func NewStore[OUTPUT mgm.Model](uri, database, collection string) (*Store[OUTPUT], error) {
	c, err := mgm.NewClient(CustomClientOptions(uri))
	if err != nil {
		return nil, err
	}

	db := c.Database(database)
	col := mgm.NewCollection(db, collection)

	s := &Store[OUTPUT]{
		Client:     c,
		Database:   db,
		Collection: col,
	}
	return s, nil
}

func (s *Store[OUTPUT]) FindByID(id primitive.ObjectID, output OUTPUT) error {
	err := s.Collection.FindByID(id, output)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store[OUTPUT]) Find(id string, output OUTPUT) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.FindByID(oid, output)
}

func (s *Store[OUTPUT]) Save(o *Document) error {
	// TODO: if id is nil create otherwise, call update
	return s.Collection.Create(o)
}

func (s *Store[OUTPUT]) Update(o *Document) error {
	return s.Collection.Update(o)
}

func (s *Store[OUTPUT]) Delete(o *Document) error {
	return s.Collection.Delete(o)
}

func (s *Store[OUTPUT]) Query() *QueryBuilder[OUTPUT] {
	values := make(bson.M)
	return &QueryBuilder[OUTPUT]{
		store:  s,
		values: values,
		limit:  25,
		skip:   0,
		sort:   bson.D{},
	}
}
