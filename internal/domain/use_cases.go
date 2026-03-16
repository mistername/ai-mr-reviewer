package domain

import (
	"context"
)

//go:generate go tool mockery

type MRProviderPort interface {
	GetMergeRequestChanges(ctx context.Context) ([]Diff, error)
	GetExistingComments(ctx context.Context) (map[string][]string, error)
	AddMergeRequestDiscussion(ctx context.Context, file string, line int, note string) error
	DeleteBotCommentsExceptResolved(ctx context.Context) error
}

type AIProviderPort interface {
	ReviewCode(ctx context.Context, diff string) (string, error)
}
