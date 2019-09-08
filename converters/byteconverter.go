package converters

import "reflect"

//ByteConverter is a generic interface that takes a
//byte array and converts it to a known type.
type ByteConverter interface {
	ToBytes(source interface{}) ([]byte, error)
	FromBytes(source []byte, objType reflect.Type) (interface{}, error)
}
