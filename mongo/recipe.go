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

type MongoRecipeRepository struct {
	recipesCollection *mongo.Collection
}

func NewMongoRecipeRepository(recipesCollection *mongo.Collection) *MongoRecipeRepository {
	return &MongoRecipeRepository{
		recipesCollection: recipesCollection,
	}
}

func (repo *MongoRecipeRepository) GetRecipes(ctx context.Context) ([]core.Recipe, error) {
	var recipes []core.Recipe
	cursor, err := repo.recipesCollection.Find(ctx, bson.M{})
	if err != nil {
		return []core.Recipe{}, fmt.Errorf("error while executing query: %v", err)
	}

	if err = cursor.All(ctx, &recipes); err != nil {
		return []core.Recipe{}, fmt.Errorf("error while iterating cursor: %v", err)
	}

	// cursor.All returns nil if collection is empty
	if recipes == nil {
		return []core.Recipe{}, nil
	}

	return recipes, nil
}

func (repo *MongoRecipeRepository) GetRecipeByID(ctx context.Context, id string) (core.Recipe, error) {
	var recipe core.Recipe

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Printf("not a valid object id: %s", id)
		return core.Recipe{}, &core.RecipeIDNotValidError{
			ID: id,
		}
	}

	filter := bson.M{"_id": objectID}
	if err := repo.recipesCollection.FindOne(ctx, filter).Decode(&recipe); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.Recipe{}, &core.RecipeNotFoundError{
				ID: id,
			}
		}
		return core.Recipe{}, fmt.Errorf("error while executing query: %v", err)
	}

	return recipe, nil
}

func (repo *MongoRecipeRepository) AddRecipe(ctx context.Context, recipe *core.Recipe) error {
	recipe.ID = primitive.NewObjectID()

	_, err := repo.recipesCollection.InsertOne(ctx, recipe)
	if err != nil {
		return fmt.Errorf("error while executing insert: %v", err)
	}

	return nil
}

func (repo *MongoRecipeRepository) UpdateRecipe(ctx context.Context, recipe core.Recipe) error {
	filter := bson.M{"_id": recipe.ID}
	_, err := repo.recipesCollection.ReplaceOne(ctx, filter, recipe)
	if err != nil {
		return fmt.Errorf("error while executing update: %v", err)
	}

	return nil
}

func (repo *MongoRecipeRepository) DeleteRecipe(ctx context.Context, recipe core.Recipe) error {
	filter := bson.M{"_id": recipe.ID}
	_, err := repo.recipesCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error while executing delete: %v", err)
	}

	return nil
}
