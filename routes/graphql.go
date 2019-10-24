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
	rootFields := gqlfields.GetRootQueryFields()
	rootQuery := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "root",
			Fields: rootFields,
		},
	)

	// Mutation query
	mutationFields := gqlfields.MutationTypes
	mutationQuery := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: mutationFields,
		},
	)

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: mutationQuery,
	})
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(r.Body)

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: string(body),
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}
