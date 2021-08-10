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
)

type sourceModelBase struct {
	Title      string `json:"title"`
	SourceType string `json:"type"`
}

type sourceModel struct {
	ID string `json:"id"`
	sourceModelBase
}

type sourceForCreationModel struct {
	sourceModelBase
}

type sourceForUpdateModel struct {
	sourceModelBase
}

type GetSourcesHandler struct {
	sourceRepository core.SourceRepository
}

func NewGetSourcesHandler(sourceRepository core.SourceRepository) *GetSourcesHandler {
	return &GetSourcesHandler{
		sourceRepository: sourceRepository,
	}
}

func (handler *GetSourcesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sources, err := handler.sourceRepository.GetSources(ctx)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var sourceModels = make([]sourceModel, 0, len(sources))
	for _, source := range sources {
		sourceModels = append(sourceModels, sourceModel{
			ID: source.ID.Hex(),
			sourceModelBase: sourceModelBase{
				Title:      source.Title,
				SourceType: source.Type,
			},
		})
	}

	jsonRecipes, err := json.Marshal(sourceModels)
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

type GetSourceHandler struct {
	sourceRepository core.SourceRepository
}

func NewGetSourceHandler(sourceRepository core.SourceRepository) *GetSourceHandler {
	return &GetSourceHandler{
		sourceRepository: sourceRepository,
	}
}

func (handler *GetSourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]
	source, err := handler.sourceRepository.GetSourceByID(ctx, id)
	if err != nil {
		fmt.Println(err)

		var e *core.SourceNotFoundError
		if errors.As(err, &e) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	sourceModel := sourceModel{
		ID: source.ID.Hex(),
		sourceModelBase: sourceModelBase{
			Title:      source.Title,
			SourceType: source.Type,
		},
	}

	jsonRecipe, err := json.Marshal(sourceModel)
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

type AddSourceHandler struct {
	sourceRepository core.SourceRepository
}

func NewAddSourceHandler(sourceRepository core.SourceRepository) *AddSourceHandler {
	return &AddSourceHandler{
		sourceRepository: sourceRepository,
	}
}

func (handler *AddSourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var sourceForCreation sourceForCreationModel
	err := json.NewDecoder(r.Body).Decode(&sourceForCreation)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	source := core.Source{
		Title: sourceForCreation.Title,
		Type:  sourceForCreation.SourceType,
	}
	err = handler.sourceRepository.AddSource(ctx, &source)
	if err != nil {
		log.Print(err)

		var e *core.SourceTypeNotValidError
		if errors.As(err, &e) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type UpdateSourceHandler struct {
	sourceRepository core.SourceRepository
}

func NewUpdateSourceHandler(sourceRepository core.SourceRepository) *UpdateSourceHandler {
	return &UpdateSourceHandler{
		sourceRepository: sourceRepository,
	}
}

func (handler *UpdateSourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]
	source, err := handler.sourceRepository.GetSourceByID(ctx, id)
	if err != nil {
		fmt.Println(err)

		var e *core.SourceNotFoundError
		if errors.As(err, &e) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var sourceForUpdate sourceForUpdateModel
	err = json.NewDecoder(r.Body).Decode(&sourceForUpdate)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	source.Title = sourceForUpdate.Title
	source.Type = sourceForUpdate.SourceType

	err = handler.sourceRepository.UpdateSource(ctx, source)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type DeleteSourceHandler struct {
	sourceRepository core.SourceRepository
}

func NewDeleteSourceHandler(sourceRepository core.SourceRepository) *DeleteSourceHandler {
	return &DeleteSourceHandler{
		sourceRepository: sourceRepository,
	}
}

func (handler *DeleteSourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]
	source, err := handler.sourceRepository.GetSourceByID(ctx, id)
	if err != nil {
		fmt.Println(err)

		var e *core.SourceNotFoundError
		if errors.As(err, &e) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = handler.sourceRepository.DeleteSource(ctx, source)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
