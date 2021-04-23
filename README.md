# embedqueries

This shows a use case for the Go 1.16+ 'embed' feature, which allows for static files to be embedded into the built binary.

In this repo, there is a sample app that issues queries to the SpaceX GraphQL API. Each query is a file that will be embedded at compile time and hence available at runtime. This example also makes use of Go's 'templating' feature, substituting query parameters at execution time.

With the queries in standalone files as opposed to strings in a Go file, it is much more developer friendly to write and visualise the queries, allowing IDE plugins for formatting etc.

A similar approach could of course be used for SQL queries.

For more information on the embed feature, see https://golang.org/pkg/embed/

# Notes on the SpaceX GraphQL API:
 - A playground is available here: https://api.spacex.land/graphql/
 
# This repo:
 * /queries dir:
   * Contains the sample GraphQL queries
   * The files contain various Go Template syntax for variables
 * store.go
   * Embeds the files from the /queries sub dir. The embedding occurs at compile time. Note the ```//go:embed queries/*``` annotation
   * the newStore func loads the files (at runtime) and parses into Go Templates
 * client.go
   * Injects variables into the parsed templates
   * Posts the requests to the GraphQL endpoint
 * main.go
   * Simply runs some sample API requests
  