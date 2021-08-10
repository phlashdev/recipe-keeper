package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/phlashdev/recipe-keeper-api/api"
	mongodb "github.com/phlashdev/recipe-keeper-api/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	DatabaseName         = "recipe-keeper"
	RecipeCollectionName = "recipes"
	SourceCollectionName = "sources"
)

const (
	MongoDbConStrEnv = "RECIPEKEEPER_MONGODB_CONSTR"
)

func main() {
	connectionString := os.Getenv(MongoDbConStrEnv)
	if len(connectionString) == 0 {
		log.Fatal(fmt.Sprintf("Environment variable %q is not set", MongoDbConStrEnv))
	}

	dbClient, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = dbClient.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer dbClient.Disconnect(ctx)

	err = dbClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	recipesCollection := dbClient.Database(DatabaseName).Collection(RecipeCollectionName)
	recipeRepository := mongodb.NewMongoRecipeRepository(recipesCollection)

	sourcesCollection := dbClient.Database(DatabaseName).Collection(SourceCollectionName)
	sourceRepository := mongodb.NewMongoSourceRepository(sourcesCollection)

	router := mux.NewRouter()

	recipesSubrouter := router.PathPrefix("/api/recipes").Subrouter()
	recipesSubrouter.Handle("/{id}", api.NewGetRecipeHandler(recipeRepository)).Methods(http.MethodGet)
	recipesSubrouter.Handle("/{id}", api.NewUpdateRecipeHandler(recipeRepository)).Methods(http.MethodPut)
	recipesSubrouter.Handle("/{id}", api.NewDeleteRecipeHandler(recipeRepository)).Methods(http.MethodDelete)
	recipesSubrouter.Handle("/", api.NewGetRecipesHandler(recipeRepository)).Methods(http.MethodGet)
	recipesSubrouter.Handle("", api.NewGetRecipesHandler(recipeRepository)).Methods(http.MethodGet)
	recipesSubrouter.Handle("", api.NewAddRecipeHandler(recipeRepository)).Methods(http.MethodPost)

	sourcesSubrouter := router.PathPrefix("/api/sources").Subrouter()
	sourcesSubrouter.Handle("/{id}", api.NewGetSourceHandler(sourceRepository)).Methods(http.MethodGet)
	sourcesSubrouter.Handle("/{id}", api.NewUpdateSourceHandler(sourceRepository)).Methods(http.MethodPut)
	sourcesSubrouter.Handle("/{id}", api.NewDeleteSourceHandler(sourceRepository)).Methods(http.MethodDelete)
	sourcesSubrouter.Handle("/", api.NewGetSourcesHandler(sourceRepository)).Methods(http.MethodGet)
	sourcesSubrouter.Handle("", api.NewGetSourcesHandler(sourceRepository)).Methods(http.MethodGet)
	sourcesSubrouter.Handle("", api.NewAddSourceHandler(sourceRepository)).Methods(http.MethodPost)

	log.Print("Starting web server")
	log.Fatal(http.ListenAndServe(":5000", router))
}
