package grimoire

import (
	"context"
	"reflect"
	"strings"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store[T Model] struct {
	Client        *mongox.Client
	Database      *mongox.Database
	Collection    *mongox.Collection[T]
	queryDefaults []bson.M
}

// New creates a new store object
func New[T Model](URI, database, collection string) (*Store[T], error) {
	c, err := mongo.Connect(options.Client().ApplyURI(URI))
	if err != nil {
		return nil, err
	}

	return NewWithClient[T](c, database, collection)
}

// New creates a new store object using an existing mongo client
func NewWithClient[T Model](c *mongo.Client, database, collection string) (*Store[T], error) {
	// ensure we're connected
	if err := c.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	client := mongox.NewClient(c, &mongox.Config{})
	db := client.NewDatabase(database)
	col := mongox.NewCollection[T](db, collection)

	s := &Store[T]{
		Client:        client,
		Database:      db,
		Collection:    col,
		queryDefaults: []bson.M{},
	}
	return s, nil
}

// CreateIndexes creates indexes on the collection
// descriptor is a string of index specs separated by semicolons
// each spec is a comma separated list of fields, with an optional direction
func CreateIndexes[T Model](s *Store[T], descriptor string) {
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
		// s.Collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: d})
		s.Collection.Collection().Indexes().CreateMany(context.Background(), []mongo.IndexModel{{Keys: d}})
	}
}

// Indexes creates indexes on the collection based on struct tags
// deprecated: use CreateIndexesFromTags
func Indexes[T Model](s *Store[T], o *T) {
	CreateIndexesFromTags(s, o)
}

// Indexes creates indexes on the collection based on struct tags
func CreateIndexesFromTags[T Model](s *Store[T], o *T) {
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
				s.Collection.Collection().Indexes().CreateMany(context.Background(), []mongo.IndexModel{{Keys: bson.D{{Key: name, Value: dir}}}})
			}
		}
	}
}

// SetQueryDefaults sets defaults used for all queries
func (s *Store[T]) SetQueryDefaults(values []bson.M) {
	s.queryDefaults = append(s.queryDefaults, values...)
}

func (s *Store[T]) GetByID(id bson.ObjectID) (*T, error) {
	return s.Collection.Finder().Filter(query.Id(id)).FindOne(context.Background())
}

func (s *Store[T]) Get(id string) (*T, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.Collection.Finder().Filter(query.Id(oid)).FindOne(context.Background())
}

func (s *Store[T]) FindByID(id bson.ObjectID) (*T, error) {
	return s.Collection.Finder().Filter(query.Id(id)).FindOne(context.Background())
}

func (s *Store[T]) Find(id string) (*T, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.FindByID(oid)
}

func (s *Store[T]) Save(o *T) error {
	if (*o).GetID().IsZero() {
		_, err := s.Collection.Creator().InsertOne(context.Background(), o)
		return err
	}
	return s.Update(o)
}

func (s *Store[T]) Update(o *T) error {
	_, err := s.Collection.Updater().Filter(query.Id((*o).GetID())).Updates(bson.M{"$set": (*o)}).UpdateOne(context.Background())
	return err
}

func (s *Store[T]) Delete(o *T) error {
	_, err := s.Collection.Deleter().Filter(query.Id((*o).GetID())).DeleteOne(context.Background())
	return err
}

func (s *Store[T]) Count(query bson.M) (int64, error) {
	return s.Collection.Finder().Filter(query).Count(context.Background())
}

func (s *Store[T]) Query() *QueryBuilder[T] {
	q := query.NewBuilder()

	if len(s.queryDefaults) > 0 {
		for _, m := range s.queryDefaults {
			for k, v := range m {
				q.Eq(k, v)
			}
		}
	}

	return &QueryBuilder[T]{
		store: s,
		query: q,
		limit: 25,
		skip:  0,
		sort:  bson.D{},
	}
}
