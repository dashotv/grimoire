package grimoire

import (
	"github.com/chenmingyong0423/go-mongox/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Document struct {
	mongox.Model `bson:",inline"`
}

type Model interface {
	GetID() bson.ObjectID
}

func (d Document) GetID() bson.ObjectID {
	return d.ID
}

//
// func (d *Document) PrepareID(id interface{}) (interface{}, error) {
// 	return d.IDField.PrepareID(id)
// }
//
// func (d *Document) GetID() interface{} {
// 	return d.IDField.GetID()
// }
//
// func (d *Document) SetID(id interface{}) {
// 	d.IDField.SetID(id)
// }
