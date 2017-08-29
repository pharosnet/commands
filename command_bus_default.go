package commands

import "errors"

type defaultCommandBus struct {
	handlers CommandHandlers
	processor CommandProcessor
}

func NewDefaultCommandBus(chs CommandHandlers, processor CommandProcessor) CommandBus {
	bus := new(defaultCommandBus)
	bus.handlers = chs
	bus.processor = processor
	return bus
}

func (bus *defaultCommandBus) Subscribe(subject string, handler CommandHandler) {
	bus.handlers.Set(subject, handler)
}

func (bus *defaultCommandBus) UnSubscribe(subject string) {
	bus.handlers.Del(subject)
}

func (bus *defaultCommandBus) dispatch(subject string, cmd Command, cb CommandCallback) error  {
	handler, has := bus.handlers.Get(subject)
	if !has {
		return errors.New("command bus dispatch failed, can not find command handler named " + subject)
	}
	return bus.processor.Process(cmd, handler, cb)
}


