package commands

import (
	"sync"
	"errors"
)

type commandProcessEntry struct {
	cmd Command
	handler CommandHandler
	callback CommandCallback
}

type CommandProcessor interface {
	Process(cmd Command, handler CommandHandler, callback CommandCallback) error
	Stop()
}

type defaultCommandProcessor struct {
	ch chan *commandProcessEntry
	running bool
	mutex *sync.RWMutex
}

func NewDefaultCommandProcessor(bufferSize int) CommandProcessor {
	p := &defaultCommandProcessor{
		ch: make(chan *commandProcessEntry, bufferSize),
		running: true,
		mutex: new(sync.RWMutex),
	}
	go p.run()
	return p
}

func (p *defaultCommandProcessor) Process(cmd Command, handler CommandHandler, callback CommandCallback) error {
	if cmd == nil {
		return errors.New("command processor process failed, cmd is nil.")
	}
	if handler == nil {
		return errors.New("command processor process failed, handler is nil.")
	}
	if callback == nil {
		return errors.New("command processor process failed, callback is nil.")
	}
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	if !p.running {
		return errors.New("command processor process failed, stopped.")
	}
	p.ch <- &commandProcessEntry{cmd:cmd, handler:handler, callback:callback}
	return nil
}

func (p *defaultCommandProcessor) Stop()  {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	close(p.ch)
	p.running = false
}

func (p *defaultCommandProcessor) run() {
	for p.running {
		entry, ok := <- p.ch
		if !ok {
			break
		}
		resultData, err := entry.handler.Handle(entry.cmd, nil)
		if err != nil {
			entry.callback.On(nil, err)
			continue
		}
		entry.callback.On(newDefaultCommandResult(resultData), nil)
	}
}