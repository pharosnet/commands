package commands

type CommandBus interface {
	Subscribe(subject string, handler CommandHandler)
	UnSubscribe(subject string)
	dispatch(subject string, cmd Command, cb CommandCallback) error
}

type CommandHandler interface {
	Handle(cmd Command, err error) (interface{}, error)
}

