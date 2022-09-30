package grimoire

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store[T mgm.Model] struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mgm.Collection
}

func New[T mgm.Model](URI, database, collection string) (*Store[T], error) {
	c, err := newClient(URI)
	if err != nil {
		return nil, err
	}

	db := c.Database(database)
	col := mgm.NewCollection(db, collection)

	s := &Store[T]{
		Client:     c,
		Database:   db,
		Collection: col,
	}
	return s, nil
}

func (s *Store[T]) FindByID(id primitive.ObjectID, out T) error {
	err := s.Collection.FindByID(id, out)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store[T]) Find(id string, out T) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.FindByID(oid, out)
}

func (s *Store[T]) Save(o *Document) error {
	// TODO: if id is nil create otherwise, call update
	return s.Collection.Create(o)
}

func (s *Store[T]) Update(o *Document) error {
	return s.Collection.Update(o)
}

func (s *Store[T]) Delete(o *Document) error {
	return s.Collection.Delete(o)
}

func (s *Store[T]) Count(query bson.M) (int64, error) {
	return s.Collection.CountDocuments(mgm.Ctx(), query)
}

func (s *Store[T]) Query() *QueryBuilder[T] {
	values := make(bson.M)
	return &QueryBuilder[T]{
		store:  s,
		values: values,
		limit:  25,
		skip:   0,
		sort:   bson.D{},
	}
}
