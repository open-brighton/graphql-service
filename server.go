package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-gonic/gin"
	"github.com/openbrighton/graphql-service/graph"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aws/aws-lambda-go/events"
	adapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	// Server setup:
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})

	h.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/v1/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// TODO: figure out why the /v1 is needed here
	r.POST("/v1/graphql", graphqlHandler())
	r.GET("/v1/graphql", playgroundHandler())

	return r
}

func isRunningInLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}

func startLambda() {
	r := SetupRouter()
	ginLambda := adapter.NewV2(r)
	lambda.Start(func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		resp, err := ginLambda.ProxyWithContext(ctx, event)
		if err != nil {
			log.Printf("Lambda error: %v", err)
		}
		return resp, err
	})
}

func startLocal() {
	r := SetupRouter()
	log.Println("Starting local server on http://localhost:8080/v1/graphql")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start local server:", err)
	}
}

func main() {
	if isRunningInLambda() {
		log.Println("Starting lambda server")
		startLambda()
	} else {
		startLocal()
	}
}
