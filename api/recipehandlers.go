package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/phlashdev/recipe-keeper-api/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type recipeModelBase struct {
	Title            string   `json:"title"`
	SourceID         string   `json:"sourceId"`
	SourceAnnotation string   `json:"sourceAnnotation"`
	Category         string   `json:"category"`
	Allergens        []string `json:"allergens"`
}

type recipeModel struct {
	ID string `json:"id"`
	recipeModelBase
}

type recipeForCreationModel struct {
	recipeModelBase
}

type recipeForUpdateModel struct {
	recipeModelBase
}

type GetRecipesHandler struct {
	recipeRepository core.RecipeRepository
}

func NewGetRecipesHandler(recipeRepository core.RecipeRepository) *GetRecipesHandler {
	return &GetRecipesHandler{
		recipeRepository: recipeRepository,
	}
}

func (handler *GetRecipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	recipes, err := handler.recipeRepository.GetRecipes(ctx)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var recipeModels = make([]recipeModel, 0, len(recipes))
	for _, recipe := range recipes {
		sourceID := ""
		if !recipe.Source.IsZero() {
			sourceID = recipe.Source.Hex()
		}

		recipeModels = append(recipeModels, recipeModel{
			ID: recipe.ID.Hex(),
			recipeModelBase: recipeModelBase{
				Title:            recipe.Title,
				SourceID:         sourceID,
				SourceAnnotation: recipe.SourceAnnotation,
				Category:         recipe.Category,
				Allergens:        recipe.Allergens,
			},
		})
	}

	jsonRecipes, err := json.Marshal(recipeModels)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonRecipes)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type GetRecipeHandler struct {
	recipeRepository core.RecipeRepository
}

func NewGetRecipeHandler(recipeRepository core.RecipeRepository) *GetRecipeHandler {
	return &GetRecipeHandler{
		recipeRepository: recipeRepository,
	}
}

func (handler *GetRecipeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]
	recipe, err := handler.recipeRepository.GetRecipeByID(ctx, id)
	if err != nil {
		fmt.Println(err)

		var recipeNotFoundErr *core.RecipeNotFoundError
		if errors.As(err, &recipeNotFoundErr) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var idNotValidErr *core.RecipeIDNotValidError
		if errors.As(err, &idNotValidErr) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sourceID := ""
	if !recipe.Source.IsZero() {
		sourceID = recipe.Source.Hex()
	}

	recipeModel := recipeModel{
		ID: recipe.ID.Hex(),
		recipeModelBase: recipeModelBase{
			Title:            recipe.Title,
			SourceID:         sourceID,
			SourceAnnotation: recipe.SourceAnnotation,
			Category:         recipe.Category,
			Allergens:        recipe.Allergens,
		},
	}

	jsonRecipe, err := json.Marshal(recipeModel)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonRecipe)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type AddRecipeHandler struct {
	recipeRepository core.RecipeRepository
}

func NewAddRecipeHandler(recipeRepository core.RecipeRepository) *AddRecipeHandler {
	return &AddRecipeHandler{
		recipeRepository: recipeRepository,
	}
}

func (handler *AddRecipeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var recipeForCreation recipeForCreationModel
	err := json.NewDecoder(r.Body).Decode(&recipeForCreation)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sourceID, err := primitive.ObjectIDFromHex(recipeForCreation.SourceID)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recipe := core.Recipe{
		Title:            recipeForCreation.Title,
		Source:           sourceID,
		SourceAnnotation: recipeForCreation.SourceAnnotation,
		Category:         recipeForCreation.Category,
		Allergens:        recipeForCreation.Allergens,
	}

	err = handler.recipeRepository.AddRecipe(ctx, &recipe)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type UpdateRecipeHandler struct {
	recipeRepository core.RecipeRepository
}

func NewUpdateRecipeHandler(recipeRepository core.RecipeRepository) *UpdateRecipeHandler {
	return &UpdateRecipeHandler{
		recipeRepository: recipeRepository,
	}
}

func (handler *UpdateRecipeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]
	recipe, err := handler.recipeRepository.GetRecipeByID(ctx, id)
	if err != nil {
		fmt.Println(err)

		var e *core.RecipeNotFoundError
		if errors.As(err, &e) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var recipeForUpdate recipeForUpdateModel
	err = json.NewDecoder(r.Body).Decode(&recipeForUpdate)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sourceID, err := primitive.ObjectIDFromHex(recipeForUpdate.SourceID)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recipe.Title = recipeForUpdate.Title
	recipe.Source = sourceID
	recipe.SourceAnnotation = recipeForUpdate.SourceAnnotation
	recipe.Category = recipeForUpdate.Category
	recipe.Allergens = recipeForUpdate.Allergens

	err = handler.recipeRepository.UpdateRecipe(ctx, recipe)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type DeleteRecipeHandler struct {
	recipeRepository core.RecipeRepository
}

func NewDeleteRecipeHandler(recipeRepository core.RecipeRepository) *DeleteRecipeHandler {
	return &DeleteRecipeHandler{
		recipeRepository: recipeRepository,
	}
}

func (handler *DeleteRecipeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]
	recipe, err := handler.recipeRepository.GetRecipeByID(ctx, id)
	if err != nil {
		fmt.Println(err)

		var e *core.RecipeNotFoundError
		if errors.As(err, &e) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = handler.recipeRepository.DeleteRecipe(ctx, recipe)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
