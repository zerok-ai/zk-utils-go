package functions

import (
	"fmt"
	"github.com/zerok-ai/zk-utils-go/storage/redis/stores"
	"regexp"
	"strings"
)

type Function interface {
	Execute(valueAtObject interface{}) (value interface{}, ok bool)
}

type BlankFunction struct {
}

func (fn BlankFunction) Execute(valueAtObject interface{}) (value interface{}, ok bool) {
	return nil, false
}

type FunctionFactory struct {
	serviceIPStore stores.LocalCacheHSetStore
}

func NewFunctionFactory(serviceIPStore stores.LocalCacheHSetStore) *FunctionFactory {
	return &FunctionFactory{serviceIPStore}
}

func (ff FunctionFactory) GetFunction(name string, args []string) *Function {
	var fn Function
	switch name {
	case JsonExtract:
		fn = ExtractJson{name, args}
	case getWorkloadFromIP:
		fn = ExtractWorkLoadFromIP{name, args, ff.serviceIPStore}
	case toLowerCase:
		fn = LowerCase{name, args}
	case toUpperCase:
		fn = UpperCase{name, args}
	default:
		fn = BlankFunction{}
	}
	return &fn
}

func (ff FunctionFactory) GetPathAndFunctions(input string) (path string, functions []Function) {

	// Define regular expressions for path, function Name, and function parameters
	pathRegex := regexp.MustCompile(`^([^#]*)`)

	// Extract path
	pathMatches := pathRegex.FindStringSubmatch(input)
	path = input
	if len(pathMatches) > 1 {
		path = pathMatches[1]
		fmt.Println("Path:", path)
	}

	// Regular expression pattern to match function calls
	pattern := `#(\w+)\(([^)]*)\)`

	// Compile the regular expression pattern
	regex := regexp.MustCompile(pattern)

	// Find all matches of the pattern in the input string
	matches := regex.FindAllStringSubmatch(input, -1)

	// Iterate over the matches and print the extracted function calls
	functions = make([]Function, 0)
	for _, match := range matches {
		// get the Name
		name := match[1]
		fmt.Print("\n---fn Name:" + name)

		// get the params and trim spaces from each substring
		params := strings.Split(match[2], ",")

		args := make([]string, 0)
		for _, s := range params {
			temp := strings.TrimSpace(s)
			if len(temp) > 0 {
				args = append(args, temp)
				fmt.Print("\t" + temp)
			}
		}
		fmt.Printf("\tnumber of Args: %d\n", len(args))

		// append the function to the list
		fn := ff.GetFunction(name, args)
		if fn != nil {
			functions = append(functions, *fn)
		}
	}
	return path, functions
}
