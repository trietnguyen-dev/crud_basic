package database

import (
	"fmt"
	"math/big"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

var (
	bigInt = reflect.TypeOf(&big.Int{})

	bigIntMongoRegistry = bson.NewRegistryBuilder().
				RegisterTypeEncoder(bigInt, bsoncodec.ValueEncoderFunc(bigIntEncodeValue)).
				RegisterTypeDecoder(bigInt, bsoncodec.ValueDecoderFunc(bigIntDecodeValue)).
				Build()
)

// https://gist.github.com/SupaHam/3afe982dc75039356723600ccc91ff77

func bigIntEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != bigInt {
		return bsoncodec.ValueEncoderError{Name: "bigIntEncodeValue", Types: []reflect.Type{bigInt}, Received: val}
	}

	b := val.Interface().(*big.Int)
	if b == nil {
		return vw.WriteNull()
	}

	return vw.WriteString(b.String())
}

func bigIntDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != bigInt {
		return bsoncodec.ValueDecoderError{Name: "uuidDecodeValue", Types: []reflect.Type{bigInt}, Received: val}
	}

	var data string
	var err error

	switch vrType := vr.Type(); vrType {
	case bsontype.String:
		data, err = vr.ReadString()
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a BigInt", vrType)
	}

	if err != nil {
		return err
	}

	if data != "" {
		bigInt, _ := new(big.Int).SetString(data, 10)
		val.Set(reflect.ValueOf(bigInt))
	}

	return nil
}
