package grimoire

import (
	"errors"
	"reflect"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//type Client struct {
//	client *mongo.Client
//}
//
//func New(URI string) (*Client, error) {
//	err := mgm.SetDefaultConfig(nil, "golem", CustomClientOptions("mongodb://localhost:27017/"))
//	if err != nil {
//		panic(err)
//	}
//
//	c, err := mgm.NewClient(CustomClientOptions(uri))
//	if err != nil {
//		return nil, err
//	}
//
//	return &Client{client: c}, nil
//}

func init() {
	// TODO: this shouldn't be necessary
	err := mgm.SetDefaultConfig(nil, "golem", CustomClientOptions("mongodb://localhost:27017/"))
	if err != nil {
		panic(err)
	}
}

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

// https://stackoverflow.com/questions/58984435/how-to-ignore-nulls-while-unmarshalling-a-mongodb-document
type nullawareDecoder struct {
	defDecoder bsoncodec.ValueDecoder
	zeroValue  reflect.Value
}

func (d *nullawareDecoder) DecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	// add symbol => string decoding here too
	if vr.Type() != bsontype.Null {
		return d.defDecoder.DecodeValue(dctx, vr, val)
	}

	if !val.CanSet() {
		return errors.New("value not settable")
	}

	if err := vr.ReadNull(); err != nil {
		return err
	}

	// Set the zero value of val's type:
	val.Set(d.zeroValue)
	return nil
}

func CustomClientOptions(URI string) *options.ClientOptions {
	customValues := []interface{}{
		"",          // string
		int(0),      // int
		int32(0),    // int32
		time.Time{}, // time
	}

	rb := bson.NewRegistryBuilder()
	for _, v := range customValues {
		t := reflect.TypeOf(v)
		defDecoder, err := bson.DefaultRegistry.LookupDecoder(t)
		if err != nil {
			panic(err)
		}
		rb.RegisterDecoder(t, &nullawareDecoder{defDecoder, reflect.Zero(t)})
	}

	return options.Client().ApplyURI(URI).SetRegistry(rb.Build())
}
