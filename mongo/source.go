package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/phlashdev/recipe-keeper-api/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSourceRepository struct {
	sourcesCollection *mongo.Collection
}

func NewMongoSourceRepository(sourcesCollection *mongo.Collection) *MongoSourceRepository {
	return &MongoSourceRepository{
		sourcesCollection: sourcesCollection,
	}
}

func (repo *MongoSourceRepository) GetSources(ctx context.Context) ([]core.Source, error) {
	var sources []core.Source
	cursor, err := repo.sourcesCollection.Find(ctx, bson.M{})
	if err != nil {
		return []core.Source{}, fmt.Errorf("error while executing query: %v", err)
	}

	if err = cursor.All(ctx, &sources); err != nil {
		return []core.Source{}, fmt.Errorf("error while iterating cursor: %v", err)
	}

	// cursor.All returns nil if collection is empty
	if sources == nil {
		return []core.Source{}, nil
	}

	return sources, nil
}

func (repo *MongoSourceRepository) GetSourceByID(ctx context.Context, id string) (core.Source, error) {
	var source core.Source

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Printf("not a valid object id: %s", id)
		return core.Source{}, &core.SourceIDNotValidError{
			ID: id,
		}
	}

	filter := bson.M{"_id": objectID}
	if err := repo.sourcesCollection.FindOne(ctx, filter).Decode(&source); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.Source{}, &core.SourceNotFoundError{
				ID: id,
			}
		}
		return core.Source{}, fmt.Errorf("error while executing query: %v", err)
	}

	return source, nil
}

func (repo *MongoSourceRepository) AddSource(ctx context.Context, source *core.Source) error {
	sourceTypeIsValid := source.Type == core.SourceTypeBook || source.Type == core.SourceTypeUrl || source.Type == core.SourceTypeCustom
	if !sourceTypeIsValid {
		return &core.SourceTypeNotValidError{
			SourceType: source.Type,
		}
	}

	source.ID = primitive.NewObjectID()

	_, err := repo.sourcesCollection.InsertOne(ctx, source)
	if err != nil {
		return fmt.Errorf("error while executing insert: %v", err)
	}

	return nil
}

func (repo *MongoSourceRepository) UpdateSource(ctx context.Context, source core.Source) error {
	filter := bson.M{"_id": source.ID}
	_, err := repo.sourcesCollection.ReplaceOne(ctx, filter, source)
	if err != nil {
		return fmt.Errorf("error while executing update: %v", err)
	}

	return nil
}

func (repo *MongoSourceRepository) DeleteSource(ctx context.Context, source core.Source) error {
	filter := bson.M{"_id": source.ID}
	_, err := repo.sourcesCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error while executing delete: %v", err)
	}

	return nil
}
