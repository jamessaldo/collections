package service

import (
	"auth/infrastructure/worker"
	"fmt"
	"reflect"
)

type Message interface{}

type MessageBus struct {
	commandHandlers map[reflect.Type]func(uow *UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error
	uow             *UnitOfWork
	mailer          worker.WorkerInterface
}

func NewMessageBus(commandHandlers map[reflect.Type]func(uow *UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error, uow *UnitOfWork, mailer worker.WorkerInterface) *MessageBus {
	return &MessageBus{
		commandHandlers: commandHandlers,
		uow:             uow,
		mailer:          mailer,
	}
}

func (mb *MessageBus) Handle(message Message) error {
	err := mb.handleCommand(message)
	if err != nil {
		return err
	}
	return nil
}

func (mb *MessageBus) handleCommand(message Message) error {
	handler, ok := mb.commandHandlers[reflect.TypeOf(message)]
	if !ok {
		return fmt.Errorf("no handler for %v", reflect.TypeOf(message))
	}
	err := handler(mb.uow, mb.mailer, message)
	if err != nil {
		return err
	}
	return nil
}
