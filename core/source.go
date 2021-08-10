package core

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	SourceTypeBook   = "book"
	SourceTypeUrl    = "url"
	SourceTypeCustom = "custom"
)

type sourceType = string

type Source struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Type  sourceType         `bson:"type,omitempty"`
	Title string             `bson:"title,omitempty"`
}

type SourceRepository interface {
	GetSources(ctx context.Context) ([]Source, error)
	GetSourceByID(ctx context.Context, id string) (Source, error)
	AddSource(ctx context.Context, source *Source) error
	UpdateSource(ctx context.Context, source Source) error
	DeleteSource(ctx context.Context, source Source) error
}

type SourceTypeNotValidError struct {
	SourceType sourceType
}

func (err *SourceTypeNotValidError) Error() string {
	return fmt.Sprintf("source type '%s* not valid", err.SourceType)
}

type SourceNotFoundError struct {
	ID string
}

func (err *SourceNotFoundError) Error() string {
	return fmt.Sprintf("source with id '%s' not found", err.ID)
}

type SourceIDNotValidError struct {
	ID string
}

func (err *SourceIDNotValidError) Error() string {
	return fmt.Sprintf("source id '%s' not valid", err.ID)
}
