package grimoire

import (
	"context"
	"fmt"
	"time"

	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type QueryBuilder[T Model] struct {
	store *Store[T]
	query *query.Builder
	limit int64
	skip  int64
	sort  bson.D
}

func (q *QueryBuilder[T]) String() string {
	return fmt.Sprintf("QueryBuilder[T] %#v", q.Build())
}

// Run executes the query and returns a list of objects.
func (q *QueryBuilder[T]) Run() ([]*T, error) {
	return q.find()
}

func (q *QueryBuilder[T]) find() ([]*T, error) {
	return q.findWithContext(context.Background())
}

func (q *QueryBuilder[T]) findWithContext(ctx context.Context) ([]*T, error) {
	f := q.store.Collection.Finder().Filter(q.query.Build())
	if q.limit > 0 {
		f.Limit(q.limit)
	}
	if q.skip > 0 {
		f.Skip(q.skip)
	}
	if len(q.sort) > 0 {
		f.Sort(q.sort)
	}
	return f.Find(ctx)
}

// Batch executes the query and yields 'size' objects at a time.
func (q *QueryBuilder[T]) Batch(size int64, f func(results []*T) error) error {
	ctx, timeout := context.WithTimeout(context.Background(), 120*time.Second)
	defer timeout()

	total, err := q.CountWithContext(ctx)
	if err != nil {
		return err
	}
	if total <= size {
		q.Skip(0)
		q.Limit(int(size))
		list, err := q.Run()
		if err != nil {
			return err
		}
		return f(list)
	}

	for i := int64(0); i < total; i += size {
		q.Skip(int(i))
		q.Limit(int(size))
		list, err := q.find()
		if err != nil {
			return err
		}
		err = f(list)
		if err != nil {
			return err
		}
	}

	return nil
}

// Batch executes the query in batches of 'batchSize' and yields one object at a time
func (q *QueryBuilder[T]) Each(batchSize int64, f func(result *T) error) error {
	return q.Batch(batchSize, func(results []*T) error {
		for _, result := range results {
			if err := f(result); err != nil {
				return err
			}
		}
		return nil
	})
}

// First executes the query and returns the first object.
func (q *QueryBuilder[T]) First() (*T, error) {
	list, err := q.Limit(1).Run()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("no results found")
	}
	return list[0], nil
}

// Build returns the result of the underlying mongox query builder
func (q *QueryBuilder[T]) Build() bson.D {
	return q.query.Build()
}

// Raw executes the raw bson.M query and returns a list of objects.
// NOTE: This does not use the query builder values.
func (q *QueryBuilder[T]) Raw(query bson.M) ([]*T, error) {
	return q.store.Collection.Finder().Filter(query).Find(context.Background())
}

// Count executes the query and returns the number of objects.
func (q *QueryBuilder[T]) Count() (int64, error) {
	return q.CountWithContext(context.Background())
}

// CountWithContext executes the query and returns the number of objects.
func (q *QueryBuilder[T]) CountWithContext(ctx context.Context) (int64, error) {
	return q.store.Collection.Finder().Filter(q.query.Build()).Count(ctx)
}

// DeleteMany executes the query and deletes the objects.
func (q *QueryBuilder[T]) DeleteMany() (int64, error) {
	n, err := q.store.Collection.Deleter().Filter(q.query.Build()).DeleteMany(context.Background())
	if err != nil {
		return 0, err
	}
	return n.DeletedCount, nil
}

func (q *QueryBuilder[T]) addSort(field string, value int) *QueryBuilder[T] {
	q.sort = append(q.sort, bson.E{Key: field, Value: value})
	return q
}

// Asc adds an ascending sort to the query.
// NOTE: field should be a valid BSON field.
//
// Examples:
//
//	Asc("name")
//	Asc("name").Asc("age")
func (q *QueryBuilder[T]) Asc(field string) *QueryBuilder[T] {
	return q.addSort(field, 1)
}

// Desc adds a descending sort to the query.
// NOTE: field should be a valid BSON field.
//
// Examples:
//
//	Desc("name")
//	Desc("name").Desc("age")
func (q *QueryBuilder[T]) Desc(field string) *QueryBuilder[T] {
	return q.addSort(field, -1)
}

// Limit sets the limit of the query.
//
// Examples:
//
//	Limit(10)
func (q *QueryBuilder[T]) Limit(limit int) *QueryBuilder[T] {
	q.limit = int64(limit)
	return q
}

// Skip sets the how many objects to skip of the query.
// Examples:
//
//	Skip(10)
func (q *QueryBuilder[T]) Skip(skip int) *QueryBuilder[T] {
	q.skip = int64(skip)
	return q
}

// func (q *QueryBuilder[T]) options() *options.FindOptions {
// 	o := &options.FindOptions{}
// 	if q.limit > 0 {
// 		o.SetLimit(q.limit)
// 	}
// 	o.SetSkip(q.skip)
// 	o.SetSort(q.sort)
// 	return o
// }

// Where adds a where clause to the query.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	Where("name", "value")
func (q *QueryBuilder[T]) Where(field string, value interface{}) *QueryBuilder[T] {
	q.query.Eq(field, value)
	return q
}

// In adds an in clause to the query.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	In("name", []string{"foo", "bar"})
func (q *QueryBuilder[T]) In(field string, value ...interface{}) *QueryBuilder[T] {
	q.query.In(field, value...)
	return q
}

// NotIn adds a not in clause to the query.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	NotIn("name", []string{"foo", "bar"})
func (q *QueryBuilder[T]) NotIn(field string, value interface{}) *QueryBuilder[T] {
	q.query.Nin(field, value)
	return q
}

// NotEqual adds a not equal clause to the query.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	NotEqual("name", "value")
func (q *QueryBuilder[T]) NotEqual(field string, value interface{}) *QueryBuilder[T] {
	q.query.Ne(field, value)
	return q
}

// LessThan adds a less than clause to the query.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	LessThan("name", 10)
func (q *QueryBuilder[T]) LessThan(field string, value interface{}) *QueryBuilder[T] {
	q.query.Lt(field, value)
	return q
}

// LessThanEqual adds a less than or equal clause to the query.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	LessThanEqual("name", 10)
func (q *QueryBuilder[T]) LessThanEqual(field string, value interface{}) *QueryBuilder[T] {
	q.query.Lte(field, value)
	return q
}

// GreaterThan adds a greater than clause to the query.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	GreaterThan("name", 10)
func (q *QueryBuilder[T]) GreaterThan(field string, value interface{}) *QueryBuilder[T] {
	q.query.Gt(field, value)
	return q
}

// GreaterThanEqual adds a greater than or equal clause to the query.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	GreaterThanEqual("name", 10)
func (q *QueryBuilder[T]) GreaterThanEqual(field string, value interface{}) *QueryBuilder[T] {
	q.query.Gte(field, value)
	return q
}

// Exists adds an exists clause to the query to check if a field exists.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	Exists("name")
func (q *QueryBuilder[T]) Exists(field string) *QueryBuilder[T] {
	q.query.Exists(field, true)
	return q
}

// NotExists adds an exists clause to the query to check if a field does not exist.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	NotExists("name")
func (q *QueryBuilder[T]) NotExists(field string) *QueryBuilder[T] {
	q.query.Exists(field, false)
	return q
}

// Or adds an or clause to the query. This is used when the or clause compares different fields. If you need to
// compare the same field, use the In or NotIn functions.
//
// Example:
//
// Or(
//
//	query.Eq("title", "The Incredibles"),
//	query.Eq("_type", "Series"),
//
// )
//
// Or(
//
//	q.Where("_type", "Movie").Where("kind", "movies3d").Where("title", "Up").Build(),
//	q.Where("_type", "Series").Where("kind", "donghua").Where("title", "The Great Ruler").Build(),
//
// )
func (q *QueryBuilder[T]) Or(conditions ...any) *QueryBuilder[T] {
	q.query.Or(conditions...)
	return q
}

// If adds a field and value to the query if the condition is true.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	If(true, "name", "value")
func (q *QueryBuilder[T]) If(cond bool, field string, value interface{}) *QueryBuilder[T] {
	if cond {
		q.query.Eq(field, value)
	}
	return q
}
