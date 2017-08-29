package commands

import (
	"reflect"
	"errors"
)

type CommandResult interface {
	Scan(v interface{}) error
	Get() interface{}
	Set(v interface{})
}

type defaultCommandResult struct {
	v interface{}
}

func newDefaultCommandResult(v interface{}) CommandResult {
	return &defaultCommandResult{v: v}
}

func (f *defaultCommandResult) Scan(v interface{}) error {
	rv := reflect.TypeOf(v)
	if rv.Kind() != reflect.Ptr  {
		return errors.New("The argument to Scan must be a non-nil pointer.")
	}
	rsv := reflect.ValueOf(f.v)
	rvv := reflect.ValueOf(v)
	if rsv.Kind() == reflect.Ptr {
		rvv.Elem().Set(rsv.Elem())
		return nil
	}
	rvv.Elem().Set(rsv)
	return nil
}

func (f *defaultCommandResult) Get() interface{} {
	return f.v
}

func (f *defaultCommandResult) Set(v interface{}) {
	f.v = v
}