package converters_test

import (
	"reflect"
	"testing"

	"github.com/marrick66/amqp-eventbus/converters"
)

type testObj struct {
	Index int
	Value string
}

func TestSimpleStructToAndFrom(t *testing.T) {
	obj := testObj{
		Index: 3,
		Value: "Test"}

	converter := &converters.JSONConverter{}

	var bytes []byte
	var newObj interface{}
	var convObj *testObj
	var err error
	var ok bool

	if bytes, err = converter.ToBytes(obj); err != nil {
		t.Errorf("ToBytes: Can't convert object: %v. Error: %v", obj, err)
		return
	}

	if newObj, err = converter.FromBytes(bytes, reflect.TypeOf(obj)); err != nil {
		t.Errorf("FromBytes: Can't convert bytes to object %v", err)
		return
	}

	if convObj, ok = newObj.(*testObj); !ok {
		t.Errorf("FromBytes: Unexpected conversion type.")
		return
	}

	if convObj.Index != obj.Index || convObj.Value != obj.Value {
		t.Errorf("Converted objects don't match. Original: %v,  Converted: %v", obj, convObj)
		return
	}

	return
}

func TestInvalidMessageBytes(t *testing.T) {
	garbled := "AAGS(BEGW@##GA"

	converter := converters.JSONConverter{}

	var bytes []byte
	var err error
	var obj testObj

	if bytes, err = converter.ToBytes(garbled); err != nil {
		t.Errorf("ToBytes: Unable to convert garbled %s", garbled)
		return
	}

	if _, err = converter.FromBytes(bytes, reflect.TypeOf(obj)); err == nil {
		t.Errorf("FromBytes: Expected to fail, but did not.")
		return
	}
}
