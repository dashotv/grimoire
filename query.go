package grimoire

import (
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueryBuilder[T mgm.Model] struct {
	store  *Store[T]
	values bson.M
	limit  int64
	skip   int64
	sort   bson.D
}

func (q *QueryBuilder[T]) addSort(field string, value int) *QueryBuilder[T] {
	q.sort = append(q.sort, bson.E{Key: field, Value: value})
	return q
}

func (q *QueryBuilder[T]) Asc(field string) *QueryBuilder[T] {
	return q.addSort(field, 1)
}

func (q *QueryBuilder[T]) Desc(field string) *QueryBuilder[T] {
	return q.addSort(field, -1)
}

func (q *QueryBuilder[T]) Limit(limit int) *QueryBuilder[T] {
	q.limit = int64(limit)
	return q
}

func (q *QueryBuilder[T]) Skip(skip int) *QueryBuilder[T] {
	q.skip = int64(skip)
	return q
}

func (q *QueryBuilder[T]) options() *options.FindOptions {
	o := &options.FindOptions{}
	o.SetLimit(q.limit)
	o.SetSkip(q.skip)
	o.SetSort(q.sort)
	return o
}

func (q *QueryBuilder[T]) Run() ([]T, error) {
	result := make([]T, 0)
	err := q.store.Collection.SimpleFind(&result, q.values, q.options())
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (q *QueryBuilder[T]) Where(key string, value interface{}) *QueryBuilder[T] {
	q.values[key] = bson.M{operator.Eq: value}
	return q
}

func (q *QueryBuilder[T]) In(key string, value interface{}) *QueryBuilder[T] {
	q.values[key] = bson.M{operator.In: value}
	return q
}

func (q *QueryBuilder[T]) NotIn(key string, value interface{}) *QueryBuilder[T] {
	q.values[key] = bson.M{operator.Nin: value}
	return q
}

func (q *QueryBuilder[T]) NotEqual(key string, value interface{}) *QueryBuilder[T] {
	q.values[key] = bson.M{operator.Ne: value}
	return q
}

func (q *QueryBuilder[T]) LessThan(key string, value interface{}) *QueryBuilder[T] {
	q.values[key] = bson.M{operator.Lt: value}
	return q
}

func (q *QueryBuilder[T]) LessThanEqual(key string, value interface{}) *QueryBuilder[T] {
	q.values[key] = bson.M{operator.Lte: value}
	return q
}

func (q *QueryBuilder[T]) GreaterThan(key string, value interface{}) *QueryBuilder[T] {
	q.values[key] = bson.M{operator.Gt: value}
	return q
}

func (q *QueryBuilder[T]) GreaterThanEqual(key string, value interface{}) *QueryBuilder[T] {
	q.values[key] = bson.M{operator.Gte: value}
	return q
}
