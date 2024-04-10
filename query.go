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

// Run executes the query and returns a list of objects.
func (q *QueryBuilder[T]) Batch(size int64, f func(results []T) error) error {
	filter := bson.M{}
	if len(q.values) > 0 {
		filter["$and"] = q.values
	}

	total, err := q.Count()
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
		result := make([]T, 0)
		q.Skip(int(i))
		q.Limit(int(size))
		err := q.store.Collection.SimpleFind(&result, filter, q.options())
		if err != nil {
			return err
		}
		err = f(result)
		if err != nil {
			return err
		}
	}

	return nil
}

// Raw executes the raw bson.M query and returns a list of objects.
// NOTE: This does not use the query builder values.
func (q *QueryBuilder[T]) Raw(query bson.M) ([]T, error) {
	result := make([]T, 0)
	err := q.store.Collection.SimpleFind(&result, query, q.options())
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

// Exists adds an exists clause to the query to check if a field exists.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	Exists("name")
func (q *QueryBuilder[T]) Exists(field string) *QueryBuilder[T] {
	q.values = append(q.values, bson.M{field: bson.M{operator.Exists: true}})
	return q
}

// NotExists adds an exists clause to the query to check if a field does not exist.
// NOTE: field should be a valid BSON field.
//
// Example:
//
//	NotExists("name")
func (q *QueryBuilder[T]) NotExists(field string) *QueryBuilder[T] {
	q.values = append(q.values, bson.M{field: bson.M{operator.Exists: false}})
	return q
}

// Or adds an or clause to the query. This is used when the or clause compares different fields. If you need to
// compare the same field, use the In or NotIn functions.
// NOTE: f should be a function that accepts a querybuilder.
//
// Example:
//
//	Or(func(q *QueryBuilder[T]) {
//		return q.Where("field1", "value").Where("field2", "value2")
//	})
func (q *QueryBuilder[T]) Or(f func(q *QueryBuilder[T])) *QueryBuilder[T] {
	ss := &Store[T]{
		Client:     q.store.Client,
		Database:   q.store.Database,
		Collection: q.store.Collection,
	}
	qq := ss.Query()
	f(qq)
	q.values = append(q.values, bson.M{operator.Or: qq.values})
	return q
}

// ComplexOr adds an or clause to the query using two separate query builders. This is used
// when the or clause requires two queries that are structurally different.
// NOTE: f should be a function that accepts two query builders.
//
// Example:
//
//	ComplexOr(func(qq *QueryBuilder[T], qr *QueryBuilder[T])) *QueryBuilder[T] {
//		qq.Where("name", "value")
//		qr.Where("type", "value2")
//	})
func (q *QueryBuilder[T]) ComplexOr(f func(qq *QueryBuilder[T], qr *QueryBuilder[T])) *QueryBuilder[T] {
	ss := &Store[T]{
		Client:     q.store.Client,
		Database:   q.store.Database,
		Collection: q.store.Collection,
	}
	qq := ss.Query()
	qr := ss.Query()
	f(qq, qr)
	q.values = append(q.values, bson.M{operator.Or: bson.A{bson.M{operator.And: qq.values}, bson.M{operator.And: qr.values}}})
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
		q.values = append(q.values, bson.M{field: value})
	}
	return q
}
