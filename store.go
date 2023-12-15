package grimoire

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store[T mgm.Model] struct {
	Client        *mongo.Client
	Database      *mongo.Database
	Collection    *mgm.Collection
	queryDefaults []bson.M
}

func New[T mgm.Model](URI, database, collection string) (*Store[T], error) {
	c, err := newClient(URI)
	if err != nil {
		return nil, err
	}

	db := c.Database(database)
	col := mgm.NewCollection(db, collection)

	s := &Store[T]{
		Client:        c,
		Database:      db,
		Collection:    col,
		queryDefaults: []bson.M{},
	}
	return s, nil
}

func (s *Store[T]) SetQueryDefaults(values []bson.M) {
	s.queryDefaults = append(s.queryDefaults, values...)
}

func (s *Store[T]) GetByID(id primitive.ObjectID, out T) (T, error) {
	err := s.Collection.FindByID(id, out)
	return out, err
}

func (s *Store[T]) Get(id string, out T) (T, error) {
	oid, err := idFromHex(id)
	if err != nil {
		return out, err
	}
	return s.GetByID(oid, out)
}

func (s *Store[T]) FindByID(id primitive.ObjectID, out T) error {
	err := s.Collection.FindByID(id, out)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store[T]) Find(id string, out T) error {
	oid, err := idFromHex(id)
	if err != nil {
		return err
	}
	return s.FindByID(oid, out)
}

func idFromHex(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return oid, nil
}

func (s *Store[T]) Save(o T) error {
	if o.GetID().(primitive.ObjectID).IsZero() {
		return s.Collection.Create(o)
	}
	return s.Collection.Update(o)
}

func (s *Store[T]) CreateWithTransaction(o T) error {
	return mgm.TransactionWithClient(mgm.Ctx(), s.Client, func(session mongo.Session, ctx mongo.SessionContext) error {
		err := s.Collection.CreateWithCtx(ctx, o)
		if err != nil {
			return err
		}
		return session.CommitTransaction(ctx)
	})
}

func (s *Store[T]) Update(o T) error {
	return s.Collection.Update(o)
}

func (s *Store[T]) Delete(o T) error {
	return s.Collection.Delete(o)
}

func (s *Store[T]) Count(query bson.M) (int64, error) {
	return s.Collection.CountDocuments(mgm.Ctx(), query)
}

func (s *Store[T]) Query() *QueryBuilder[T] {
	values := make([]bson.M, 0)
	if len(s.queryDefaults) > 0 {
		values = append(values, s.queryDefaults...)
	}
	return &QueryBuilder[T]{
		store:  s,
		values: values,
		limit:  25,
		skip:   0,
		sort:   bson.D{},
	}
}
