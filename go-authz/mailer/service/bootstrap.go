package service

import (
	"context"

	"gorm.io/gorm"
)

type key int

const (
	ContextDBkey key = iota
)

func Bootstrap(db *gorm.DB) func() context.Context {
	return func() context.Context {
		ctx := context.WithValue(context.Background(), ContextDBkey, db)
		return ctx
	}
}
