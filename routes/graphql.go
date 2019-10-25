package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"munchserver/gqlfields"
	"net/http"

	"github.com/graphql-go/graphql"
)

func GetGraphQLHandler(w http.ResponseWriter, r *http.Request) {

	// Root query
	rootFields := gqlfields.RootFields
	rootQuery := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "root",
			Fields: rootFields,
		},
	)

	// Mutation query
	mutationFields := gqlfields.MutationFields
	mutationQuery := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "mutation",
			Fields: mutationFields,
		},
	)

	// Defines query schemas
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: mutationQuery,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Get query and apply
	body, err := ioutil.ReadAll(r.Body)

	// Send result
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: string(body),
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}
