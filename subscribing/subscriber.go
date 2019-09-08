package subscribing

import "reflect"

//Subscriber handles events sent to it by the event
//feeds it's subscribed to.
type Subscriber interface {
	Handle(event interface{}) error
	EventType() reflect.Type
}
