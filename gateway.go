package commands


type Gateway interface {
	Send(subject string, cmd Command) (*Future, error)
	SendWithCallback(subject string, cmd Command, callback CommandCallback) error
}

type basicGateway struct {
	commandBus CommandBus
}

func NewGateway(bus CommandBus) Gateway {
	return &basicGateway{commandBus:bus}
}

func (gateway *basicGateway) Send(subject string, cmd Command) (*Future, error) {
	future := new(Future)
	future.promiseChan = make(chan promise, 1)
	callback := &futureCommandCallback{
		future:future,
	}
	err := gateway.commandBus.dispatch(subject, cmd, callback)
	future.closeChan()
	if err != nil {
		return nil, err
	}
	return future, nil
}


func (gateway *basicGateway) SendWithCallback(subject string, cmd Command, callback CommandCallback) error {
	return gateway.commandBus.dispatch(subject, cmd, callback)
}

