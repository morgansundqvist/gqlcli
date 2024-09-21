package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/blues/jsonata-go"
	"github.com/morgansundqvist/gqlcli/pkg/executioncontext"
)

func LoadGraphQLQuery(path string) (string, error) {
	query, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(query), nil
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

type Config struct {
	GraphQLFile string            `json:"graphql_file"`
	Variables   map[string]string `json:"variables"`
	Headers     map[string]string `json:"headers"`
	Output      map[string]string `json:"output"`
	GraphQLURL  string            `json:"graphql_url"`
}

func PromptForMissingVariables(vars map[string]interface{}) (map[string]interface{}, error) {
	reader := bufio.NewReader(os.Stdin)
	for key, val := range vars {
		if val == nil || val == "" {
			fmt.Printf("Enter value for %s (or type 'quit' to exit): ", key)
			input, err := reader.ReadString('\n')
			if err != nil {
				return nil, err
			}
			input = strings.TrimSpace(input)
			if strings.ToLower(input) == "quit" {
				os.Exit(0)
			}
			vars[key] = input
		}
	}
	return vars, nil
}

var funcMap = template.FuncMap{
	"required": func(value interface{}, name string) (interface{}, error) {
		if value == nil || value == "" {
			return nil, fmt.Errorf("missing required variable: %s", name)
		}
		return value, nil
	},
}

func ParseVariables(vars map[string]string, ctx *executioncontext.ExecutionContext) (map[string]interface{}, error) {
	parsedVars := make(map[string]interface{})
	for key, val := range vars {
		var tmpl *template.Template
		var err error

		// Check if val is a string to process template

		tmpl, err = template.New("").Funcs(funcMap).Parse(val)
		if err != nil {
			return nil, fmt.Errorf("error parsing template for key '%s': %v", key, err)
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, ctx.Vars)
		if err != nil {
			return nil, fmt.Errorf("error executing template for key '%s': %v", key, err)
		}
		parsedVars[key] = buf.String()

	}

	//iterate over paresdVars and check environment variables prefixed with GQLCLI_ and replace the value if the value is ""
	for key, val := range parsedVars {
		if val == "" {
			envVar := os.Getenv("GQLCLI_" + key)
			if envVar != "" {
				parsedVars[key] = envVar
			}
		}
	}
	return parsedVars, nil
}

func StoreOutputInContext(outputConfig map[string]string, data map[string]interface{}, ctx *executioncontext.ExecutionContext) error {
	for key, tmplStr := range outputConfig {
		tmpl, err := template.New("").Parse(tmplStr)
		if err != nil {
			return fmt.Errorf("error parsing output template for key '%s': %v", key, err)
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, data)
		if err != nil {
			return fmt.Errorf("error executing output template for key '%s': %v", key, err)
		}
		ctx.Set(key, buf.String())
	}
	return nil
}

func StoreOutputInContextJSONATA(outputConfig map[string]string, data map[string]interface{}, ctx *executioncontext.ExecutionContext) error {
	for key, tmplStr := range outputConfig {
		// Create expression.
		e := jsonata.MustCompile(tmplStr)

		// Evaluate.
		res, err := e.Eval(data)
		if err != nil {
			log.Fatal(err)
		}

		ctx.Set(key, res)
	}
	return nil
}

func ParseHeaders(headers map[string]string, ctx *executioncontext.ExecutionContext) (map[string]string, error) {
	parsedHeaders := make(map[string]string)
	for key, val := range headers {
		tmpl, err := template.New("").Funcs(funcMap).Parse(val)
		if err != nil {
			return nil, fmt.Errorf("error parsing header template for key '%s': %v", key, err)
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, ctx.Vars)
		if err != nil {
			return nil, fmt.Errorf("error executing header template for key '%s': %v", key, err)
		}
		parsedHeaders[key] = buf.String()
	}
	return parsedHeaders, nil
}

func PrintContext(ctx *executioncontext.ExecutionContext, outputFields []string) {
	println("")
	if len(outputFields) == 0 {
		for key, val := range ctx.Vars {
			fmt.Printf("%s: %v\n", key, val)
		}
	} else {
		for _, field := range outputFields {
			val := ctx.Get(field)
			fmt.Printf("%s: %v\n", field, val)
		}
	}
}

type InputVariables struct {
	OutputFields []string
	InputFiles   []string
	DoTiming     bool
}

func ParseInput() (*InputVariables, error) {
	var outputFields []string
	var inputFiles []string
	var doTiming = false

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-o" && i+1 < len(args) {

			outputFields = append(outputFields, args[i+1])
			i++
		} else if arg == "-t" {
			//do timing
			//set the doTiming flag to true
			doTiming = true
		} else {

			inputFiles = append(inputFiles, arg)
		}
	}

	if len(inputFiles) == 0 {
		return nil, fmt.Errorf("no input files provided")
	}

	returnVariable := &InputVariables{
		OutputFields: outputFields,
		InputFiles:   inputFiles,
		DoTiming:     doTiming,
	}

	return returnVariable, nil
}
