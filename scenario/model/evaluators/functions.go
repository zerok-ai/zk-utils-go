package evaluators

import (
	"encoding/json"
	"fmt"
	"github.com/jmespath/go-jmespath"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"regexp"
	"strings"
)

type Function struct {
	Name string
	Args []string
}

func GetPathAndFunctions(input string) (path string, functions []Function) {

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
		fmt.Printf("\tnumber of Args: %d", len(args))

		// append the function to the list
		functions = append(functions, Function{name, args})
	}
	return path, functions
}

func GetValueFromStore(path string, store map[string]interface{}) (interface{}, bool) {

	var err error
	var ok bool
	var functions []Function
	var valueAtObject interface{}

	valueAtObject = store
	path, functions = GetPathAndFunctions(path)

	// handle path
	if len(path) > 0 {
		valueAtObject, err = jmespath.Search(path, valueAtObject)
		if err != nil {
			zkLogger.ErrorF(LoggerTag, "Error evaluating jmespath at path:%s for store %v\n%v", path, store, err)
			return valueAtObject, false
		}
	}

	// handle functions
	for _, fn := range functions {
		if valueAtObject == nil {
			return valueAtObject, false
		}

		if fn.Name == "jsonExtract" {
			valueAtObject, ok = jsonExtract(fn, valueAtObject)
		} else if fn.Name == "toUpperCase" || fn.Name == "toLowerCase" {
			valueAtObject, ok = stringFunction(fn, valueAtObject)
		}
		if !ok {
			return valueAtObject, false
		}
	}

	return valueAtObject, true
}

func stringFunction(fn Function, valueAtObject interface{}) (value string, ok bool) {
	var stringVal string
	stringVal, ok = valueAtObject.(string)
	if ok {
		if fn.Name == "toUpperCase" {
			return strings.ToUpper(stringVal), true
		} else if fn.Name == "toLowerCase" {
			return strings.ToLower(stringVal), true
		}
	}
	return "", false
}

func jsonExtract(fn Function, valueAtObject interface{}) (value string, ok bool) {

	// check if valueAtObject is a string
	var err error
	var stringVal string

	// if valueAtObject is a string, convert it to json, else directly read the json
	var jsonObject interface{}
	stringVal, ok = valueAtObject.(string)
	if ok {
		// convert string to json
		err = json.Unmarshal([]byte(stringVal), &jsonObject)
		if err != nil {
			zkLogger.Error(LoggerTag, "Error marshalling string:", stringVal)
			return "", false
		}
	} else {
		jsonObject = valueAtObject
	}
	path := fn.Args[0]
	valueAtObject, err = jmespath.Search(path, jsonObject)

	if err != nil {
		zkLogger.ErrorF(LoggerTag, "Error evaluating jmespath at path:%s for store %v", path, jsonObject)
		return "", false
	}
	return valueAtObject.(string), true
}
