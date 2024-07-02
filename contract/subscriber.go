package contract

// Subscriber subscriber
type Subscriber interface {
	// Subscribe subscribe
	Subscribe(dispatcher Bus)
}
