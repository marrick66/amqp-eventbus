package subscribing

//QueueConfig allows some customization of queue settings
//for all subscriptions.
type QueueConfiguration struct {
	Name string
	Durable bool
	AutoDelete bool
	Exclusive bool
}