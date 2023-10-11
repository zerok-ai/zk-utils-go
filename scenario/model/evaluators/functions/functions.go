package functions

import (
	"fmt"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
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
	serviceIPStore *stores.LocalCacheHSetStore
	attrStore      *stores.ExecutorAttrStore
}

func NewFunctionFactory(serviceIPStore *stores.LocalCacheHSetStore, attrStore *stores.ExecutorAttrStore) *FunctionFactory {
	return &FunctionFactory{serviceIPStore: serviceIPStore, attrStore: attrStore}
}

func (ff FunctionFactory) GetFunction(name string, args []string, attrStoreKey *cache.AttribStoreKey) *Function {

	newArgs := make([]string, 0)
	for _, arg := range args {
		newArg := ff.attrStore.GetAttributeFromStore(*attrStoreKey, arg)
		if newArg != nil {
			newArgs = append(newArgs, *newArg)
		} else {
			newArgs = append(newArgs, arg)
		}
	}

	var fn Function
	switch name {
	case JsonExtract:
		fn = ExtractJson{name: name, args: args, attrStore: ff.attrStore, attrStoreKey: attrStoreKey}
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

func (ff FunctionFactory) GetPathAndFunctions(input string, attrStoreKey *cache.AttribStoreKey) []Function {

	defer func() {
		if r := recover(); r != nil {
			zkLogger.ErrorF(LoggerTag, "In GetPathAndFunctions %s\n: Recovered from panic: %v", input, r)
		}
	}()

	// Define regular expressions for path, function name, and function parameters
	// Define the regular expression patternForInput.
	patternForInput := `([^.#]+(?:\.[^.#]+)*|#[^.)]+\(.*?\))`

	// Define a regular expression pattern to match a function.
	patternForFunction := `#(\w+)\(([^)]*)\)`

	// Compile the regular expression.
	compiledRegexFullMatch := regexp.MustCompile(patternForInput)
	compiledRegexForFunction := regexp.MustCompile(patternForFunction)

	// Find all matches.
	matches := compiledRegexFullMatch.FindAllString(input, -1)

	// create the fucntions
	functions := make([]Function, 0)
	for _, match := range matches {
		var fn *Function
		if strings.HasPrefix(match, "#") {
			functionMatch := compiledRegexForFunction.FindStringSubmatch(input)

			if len(functionMatch) > 0 {
				// functionMatch[0] is the full functionMatch, functionMatch[1] is the function name, functionMatch[2] is the arguments.
				functionName := functionMatch[1]
				arguments := functionMatch[2]

				// Split arguments into a list.
				args := strings.Split(arguments, ", ")

				fn = ff.GetFunction(functionName, args, attrStoreKey)
			}
		} else {
			fn = ff.GetFunction(JsonExtract, []string{match}, attrStoreKey)
		}
		if fn != nil {
			functions = append(functions, *fn)
		}
	}
	return functions
}

func (ff FunctionFactory) GetPathAndFunctions1(input string, attrStoreKey *cache.AttribStoreKey) (path string, functions []Function) {

	// Define regular expressions for path, function name, and function parameters
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
		// get the name
		name := match[1]
		fmt.Print("\n---fn name:" + name)

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
		fmt.Printf("\tnumber of args: %d\n", len(args))

		// append the function to the list
		fn := ff.GetFunction(name, args, attrStoreKey)
		if fn != nil {
			functions = append(functions, *fn)
		}
	}
	return path, functions
}
