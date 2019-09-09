package subscribing

//ExchangeConfiguration allows some customization as to how the
//exchange is created on first use.
type ExchangeConfiguration struct {
	AutoDelete bool
	Durable bool
	Internal bool
}