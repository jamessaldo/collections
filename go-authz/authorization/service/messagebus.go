package service

import (
	"authorization/infrastructure/worker"
	"context"
	"fmt"
	"reflect"
)

type Message interface{}

type MessageBus struct {
	commandHandlers map[reflect.Type]func(ctx context.Context, uow *UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error
	UoW             *UnitOfWork
	mailer          worker.WorkerInterface
}

func NewMessageBus(commandHandlers map[reflect.Type]func(ctx context.Context, uow *UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error, uow *UnitOfWork, mailer worker.WorkerInterface) *MessageBus {
	return &MessageBus{
		commandHandlers: commandHandlers,
		UoW:             uow,
		mailer:          mailer,
	}
}

func (mb *MessageBus) Handle(ctx context.Context, message Message) error {
	handler, ok := mb.commandHandlers[reflect.TypeOf(message)]
	if !ok {
		return fmt.Errorf("no handler for %v", reflect.TypeOf(message))
	}
	err := handler(ctx, mb.UoW, mb.mailer, message)
	if err != nil {
		return err
	}
	return nil
}
