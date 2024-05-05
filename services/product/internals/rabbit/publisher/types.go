package publisher

type ProductEventType int

const (
	ProductCreated ProductEventType = iota
	ProductUpdated
)
