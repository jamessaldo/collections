package command

import "reflect"

// Command is an interface that all commands must implement
type Command struct{}

func (c Command) CheckType() reflect.Type {
	return reflect.TypeOf(c)
}
