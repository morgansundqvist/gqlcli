package main

import (
	"fmt"
	"log"
	"time"

	"github.com/morgansundqvist/gqlcli/pkg/executioncontext"
	"github.com/morgansundqvist/gqlcli/pkg/executor"
	"github.com/morgansundqvist/gqlcli/pkg/utils"
)

func main() {
	ctx := executioncontext.NewExecutionContext()

	// Parse input
	inputVariables, err := utils.ParseInput()
	if err != nil {
		//print usage
		fmt.Println("Usage: gqlcli [-o output_field] input_file1 [input_file2 ...]")
		return
	}

	//setup variables for starttime to print out the time taken to execute the queries
	startTime := time.Now()

	// Execute GraphQL operations for each input file
	for _, inputFile := range inputVariables.InputFiles {
		config, err := utils.LoadConfig(inputFile)
		if err != nil {
			log.Fatalf("Failed to load config %s: %v", inputFile, err)
		}

		query, err := utils.LoadGraphQLQuery(config.GraphQLFile)
		if err != nil {
			log.Fatalf("Failed to load GraphQL file %s: %v", config.GraphQLFile, err)
		}

		// Parse variables
		vars, err := utils.ParseVariables(config.Variables, ctx)
		if err != nil {
			log.Fatalf("Failed to parse variables: %v", err)
		}

		// Prompt for missing variables
		vars, err = utils.PromptForMissingVariables(vars)
		if err != nil {
			log.Fatalf("Failed during user prompt: %v", err)
		}

		// Parse headers
		headers, err := utils.ParseHeaders(config.Headers, ctx)
		if err != nil {
			log.Fatalf("Failed to parse headers: %v", err)
		}

		// Execute GraphQL operation
		responseData, err := executor.ExecuteGraphQL(query, vars, headers, config)
		if err != nil {
			log.Fatalf("GraphQL execution failed: %v", err)
		}

		// Store output in context
		//err = utils.StoreOutputInContext(config.Output, responseData, ctx)
		err = utils.StoreOutputInContextJSONATA(config.Output, responseData, ctx)
		if err != nil {
			log.Fatalf("Failed to store output in context: %v", err)
		}
	}

	// Print context
	utils.PrintContext(ctx, inputVariables.OutputFields)
	endTime := time.Now()
	if inputVariables.DoTiming {
		fmt.Println("Time taken to execute the queries: ", endTime.Sub(startTime))
	}
}
