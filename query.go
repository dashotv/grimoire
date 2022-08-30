package grimoire

import (
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueryBuilder[OUTPUT mgm.Model] struct {
	store  *Store[OUTPUT]
	values bson.M
	limit  int64
	skip   int64
	sort   bson.D
}

func (q *QueryBuilder[OUTPUT]) addSort(field string, value int) *QueryBuilder[OUTPUT] {
	q.sort = append(q.sort, bson.E{Key: field, Value: value})
	return q
}

func (q *QueryBuilder[OUTPUT]) Asc(field string) *QueryBuilder[OUTPUT] {
	return q.addSort(field, 1)
}

func (q *QueryBuilder[OUTPUT]) Desc(field string) *QueryBuilder[OUTPUT] {
	return q.addSort(field, -1)
}

func (q *QueryBuilder[OUTPUT]) Limit(limit int) *QueryBuilder[OUTPUT] {
	q.limit = int64(limit)
	return q
}

func (q *QueryBuilder[OUTPUT]) Skip(skip int) *QueryBuilder[OUTPUT] {
	q.skip = int64(skip)
	return q
}

func (q *QueryBuilder[OUTPUT]) options() *options.FindOptions {
	o := &options.FindOptions{}
	o.SetLimit(q.limit)
	o.SetSkip(q.skip)
	o.SetSort(q.sort)
	return o
}

func (q *QueryBuilder[OUTPUT]) Run() ([]OUTPUT, error) {
	result := make([]OUTPUT, 0)
	err := q.store.Collection.SimpleFind(&result, q.values, q.options())
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (q *QueryBuilder[OUTPUT]) Where(key string, value interface{}) *QueryBuilder[OUTPUT] {
	q.values[key] = bson.M{operator.Eq: value}
	return q
}

func (q *QueryBuilder[OUTPUT]) In(key string, value interface{}) *QueryBuilder[OUTPUT] {
	q.values[key] = bson.M{operator.In: value}
	return q
}

func (q *QueryBuilder[OUTPUT]) NotIn(key string, value interface{}) *QueryBuilder[OUTPUT] {
	q.values[key] = bson.M{operator.Nin: value}
	return q
}

func (q *QueryBuilder[OUTPUT]) NotEqual(key string, value interface{}) *QueryBuilder[OUTPUT] {
	q.values[key] = bson.M{operator.Ne: value}
	return q
}

func (q *QueryBuilder[OUTPUT]) LessThan(key string, value interface{}) *QueryBuilder[OUTPUT] {
	q.values[key] = bson.M{operator.Lt: value}
	return q
}

func (q *QueryBuilder[OUTPUT]) LessThanEqual(key string, value interface{}) *QueryBuilder[OUTPUT] {
	q.values[key] = bson.M{operator.Lte: value}
	return q
}

func (q *QueryBuilder[OUTPUT]) GreaterThan(key string, value interface{}) *QueryBuilder[OUTPUT] {
	q.values[key] = bson.M{operator.Gt: value}
	return q
}

func (q *QueryBuilder[OUTPUT]) GreaterThanEqual(key string, value interface{}) *QueryBuilder[OUTPUT] {
	q.values[key] = bson.M{operator.Gte: value}
	return q
}
