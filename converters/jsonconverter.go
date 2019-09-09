package converters

import (
	"encoding/json"
	"reflect"
)

//JSONConverter is a wrapper around the builtin json marshalling methods.
type JSONConverter struct{}

//ToBytes is a simple wrapper over the json.Marshal builtin. By default, it
//uses UTF-8 encoding for the json bytes.
func (converter *JSONConverter) ToBytes(source interface{}) ([]byte, error) {
	var bytes []byte

	bytes, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

//FromBytes is a simple wrapper on top of json.Unmarshal. It's assumed that
//the bytes to convert are encoded in UTF-8.
func (converter *JSONConverter) FromBytes(source []byte, eventType reflect.Type) (interface{}, error) {
	obj := reflect.New(eventType).Interface()

	err := json.Unmarshal(source, &obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
