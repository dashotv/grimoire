package grimoire

import (
	"fmt"

	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueryBuilder[T mgm.Model] struct {
	store  *Store[T]
	values []bson.M
	limit  int64
	skip   int64
	sort   bson.D
}

func (q *QueryBuilder[T]) String() string {
	return fmt.Sprintf("QueryBuilder[T] %#v", q.values)
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

func (q *QueryBuilder[T]) options() *options.FindOptions {
	o := &options.FindOptions{}
	if q.limit > 0 {
		o.SetLimit(q.limit)
	}
	o.SetSkip(q.skip)
	o.SetSort(q.sort)
	return o
}

// Run executes the query and returns a list of objects.
func (q *QueryBuilder[T]) Run() ([]T, error) {
	result := make([]T, 0)
	filter := bson.M{}
	if len(q.values) > 0 {
		filter["$and"] = q.values
	}
	err := q.store.Collection.SimpleFind(&result, filter, q.options())
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Count executes the query and returns the number of objects.
func (q *QueryBuilder[T]) Count() (int64, error) {
	filter := bson.M{}
	if len(q.values) > 0 {
		filter["$and"] = q.values
	}
	return q.store.Collection.CountDocuments(mgm.Ctx(), filter)
}

// DeleteMany executes the query and deletes the objects.
func (q *QueryBuilder[T]) DeleteMany() (int64, error) {
	filter := bson.M{}
	if len(q.values) > 0 {
		filter["$and"] = q.values
	}
	n, err := q.store.Collection.DeleteMany(mgm.Ctx(), filter)
	if err != nil {
		return 0, err
	}
	return n.DeletedCount, nil
}

// Where adds a where clause to the query.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	Where("name", "value")
func (q *QueryBuilder[T]) Where(field string, value interface{}) *QueryBuilder[T] {
	q.values = append(q.values, bson.M{field: bson.M{operator.Eq: value}})
	return q
}

// In adds an in clause to the query.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	In("name", []string{"foo", "bar"})
func (q *QueryBuilder[T]) In(field string, value interface{}) *QueryBuilder[T] {
	q.values = append(q.values, bson.M{field: bson.M{operator.In: value}})
	return q
}

// NotIn adds a not in clause to the query.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	NotIn("name", []string{"foo", "bar"})
func (q *QueryBuilder[T]) NotIn(field string, value interface{}) *QueryBuilder[T] {
	q.values = append(q.values, bson.M{field: bson.M{operator.Nin: value}})
	return q
}

// NotEqual adds a not equal clause to the query.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	NotEqual("name", "value")
func (q *QueryBuilder[T]) NotEqual(field string, value interface{}) *QueryBuilder[T] {
	q.values = append(q.values, bson.M{field: bson.M{operator.Ne: value}})
	return q
}

// LessThan adds a less than clause to the query.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	LessThan("name", 10)
func (q *QueryBuilder[T]) LessThan(field string, value interface{}) *QueryBuilder[T] {
	q.values = append(q.values, bson.M{field: bson.M{operator.Lt: value}})
	return q
}

// LessThanEqual adds a less than or equal clause to the query.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	LessThanEqual("name", 10)
func (q *QueryBuilder[T]) LessThanEqual(field string, value interface{}) *QueryBuilder[T] {
	q.values = append(q.values, bson.M{field: bson.M{operator.Lte: value}})
	return q
}

// GreaterThan adds a greater than clause to the query.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	GreaterThan("name", 10)
func (q *QueryBuilder[T]) GreaterThan(field string, value interface{}) *QueryBuilder[T] {
	q.values = append(q.values, bson.M{field: bson.M{operator.Gt: value}})
	return q
}

// GreaterThanEqual adds a greater than or equal clause to the query.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	GreaterThanEqual("name", 10)
func (q *QueryBuilder[T]) GreaterThanEqual(field string, value interface{}) *QueryBuilder[T] {
	q.values = append(q.values, bson.M{field: bson.M{operator.Gte: value}})
	return q
}
