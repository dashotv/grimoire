package grimoire

import "github.com/kamva/mgm/v3"

type Document struct {
	mgm.DefaultModel `bson:",inline"`
}

func (d *Document) PrepareID(id interface{}) (interface{}, error) {
	return d.IDField.PrepareID(id)
}

func (d *Document) GetID() interface{} {
	return d.IDField.GetID()
}

func (d *Document) SetID(id interface{}) {
	d.IDField.SetID(id)
}
