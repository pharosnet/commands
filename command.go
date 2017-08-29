package commands

import "sync"

type Command interface {
	GetIdentifier() string
	FindMeta(key string) (value interface{}, has bool)
}

type BasicCommand struct {
	Identifier string
	Meta map[string] interface{}
}

func (c BasicCommand) GetIdentifier() string {
	return c.Identifier
}

func (c *BasicCommand) FindMeta(key string) (interface{}, bool) {
	if c.Meta == nil || len(c.Meta) == 0 {
		return nil, false
	}
	value, has := c.Meta[key]
	return value, has
}

type CommandHandlers interface {
	Get(subject string) (CommandHandler, bool)
	Set(subject string, handler CommandHandler)
	Del(subject string)

}

type defaultCommandHandlers struct {
	handlers map[string]CommandHandler
	mutex *sync.RWMutex
}

func NewDefaultCommandHandlers() CommandHandlers {
	return &defaultCommandHandlers{
		handlers:make(map[string]CommandHandler),
		mutex: new(sync.RWMutex),
	}
}

func (ch *defaultCommandHandlers) Get(subject string) (CommandHandler, bool) {
	ch.mutex.RLock()
	defer ch.mutex.RUnlock()
	handler, has := ch.handlers[subject]
	if !has {
		return nil, has
	}
	return handler, has
}

func (ch *defaultCommandHandlers) Del(subject string) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	delete(ch.handlers, subject)
}

func (ch *defaultCommandHandlers) Set(subject string, handler CommandHandler) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	ch.handlers[subject] = handler
}
