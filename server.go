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
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/graphql", graphqlHandler())
	r.GET("/graphql", playgroundHandler())

	r.POST("/", graphqlHandler())
	r.GET("/", playgroundHandler())

	return r
}

func isRunningInLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}

func startLambda() {
	r := SetupRouter()
	ginLambda := adapter.New(r)
	lambda.Start(func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		resp, err := ginLambda.ProxyWithContext(ctx, event)
		if err != nil {
			log.Printf("Lambda error: %v", err)
		}
		return resp, err
	})
}

// func startLocal() {
// 	r := SetupRouter()
// 	log.Println("Starting local server on http://localhost:8080")
// 	if err := r.Run(":8080"); err != nil {
// 		log.Fatal("Failed to start local server:", err)
// 	}
// }

func main() {
	log.Println("Dumping environment variables:")
	for _, env := range os.Environ() {
		log.Println(env)
	}

	log.Println("Running Lambda", isRunningInLambda())
	// if isRunningInLambda() {
	startLambda()
	// } else {
	// 	startLocal()
	// }
}
