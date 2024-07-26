package grimoire

import (
	"reflect"
	"strings"

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

// CreateIndexes creates indexes on the collection
// descriptor is a string of index specs separated by semicolons
// each spec is a comma separated list of fields, with an optional direction
func CreateIndexes[T mgm.Model](s *Store[T], o T, descriptor string) {
	if descriptor == "" {
		return
	}

	specs := strings.Split(descriptor, ";")
	for _, spec := range specs {
		d := bson.D{}
		fields := strings.Split(spec, ",")
		for _, field := range fields {
			parts := strings.Split(field, ":")
			if len(parts) > 1 {
				if parts[1] == "desc" || parts[1] == "-1" {
					d = append(d, bson.E{Key: parts[0], Value: -1})
				} else if parts[1] == "text" {
					d = append(d, bson.E{Key: parts[0], Value: "text"})
				}
			}
		}
		s.Collection.Indexes().CreateOne(mgm.Ctx(), mongo.IndexModel{Keys: d})
	}
}

// Indexes creates indexes on the collection based on struct tags
// deprecated: use CreateIndexesFromTags
func Indexes[T mgm.Model](s *Store[T], o T) {
	CreateIndexesFromTags(s, o)
}

// Indexes creates indexes on the collection based on struct tags
func CreateIndexesFromTags[T mgm.Model](s *Store[T], o T) {
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if v, ok := field.Tag.Lookup("grimoire"); ok {
			vals := strings.Split(v, ",")
			if vals[0] == "index" {
				dir := 1
				if len(vals) > 1 {
					if vals[1] == "desc" {
						dir = -1
					}
				}
				name := strings.ToLower(field.Name) // default to field name
				if v, ok := field.Tag.Lookup("bson"); ok {
					vals := strings.Split(v, ",")
					if len(vals) > 0 {
						name = vals[0] // use bson tag if available
					}
				}
				s.Collection.Indexes().CreateOne(mgm.Ctx(), mongo.IndexModel{Keys: bson.M{name: dir}})
			}
		}
	}
}

// New creates a new store object
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

// SetQueryDefaults sets defaults used for all queries
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
