package publishing

//Publisher is an interface that simply allows placing an 
//event onto the bus based on an exchange/topic combination.
//All of the details of serialization and other settings is 
// implementation specific. 
type Publisher interface {
	Publish(exchange string, topic string, event interface{}) error
}