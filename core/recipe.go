package core

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Recipe struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Title            string             `bson:"title,omitempty"`
	Source           primitive.ObjectID `bson:"source,omitempty"`
	SourceAnnotation string             `bson:"sourceAnnotation,omitempty"`
	Category         string             `bson:"category,omitempty"`
	Allergens        []string           `bson:"allergens,omitempty"`
}

type RecipeRepository interface {
	GetRecipes(ctx context.Context) ([]Recipe, error)
	GetRecipeByID(ctx context.Context, id string) (Recipe, error)
	AddRecipe(ctx context.Context, recipe *Recipe) error
	UpdateRecipe(ctx context.Context, recipe Recipe) error
	DeleteRecipe(ctx context.Context, recipe Recipe) error
}

type RecipeNotFoundError struct {
	ID string
}

func (err *RecipeNotFoundError) Error() string {
	return fmt.Sprintf("recipe with id '%s' not found", err.ID)
}

type RecipeIDNotValidError struct {
	ID string
}

func (err *RecipeIDNotValidError) Error() string {
	return fmt.Sprintf("recipe id '%s' not valid", err.ID)
}
