package commands

type CommandCallback interface {
	On(result CommandResult, err error)
}

type futureCommandCallback struct {
	future *Future
}

func (cb *futureCommandCallback) On(result CommandResult, err error)  {
	defer cb.future.closeChan()
	cb.future.setPromise(result, err)
}