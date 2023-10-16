package functions

import (
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/scenario/model/evaluators/cache"
	"github.com/zerok-ai/zk-utils-go/storage/redis/stores"
	"regexp"
	"strings"
)

type Function interface {
	Execute(valueAtObject interface{}) (value interface{}, ok bool)
	GetName() string
}

type FunctionFactory struct {
	podDetailsStore *stores.LocalCacheHSetStore
	attrStore       *stores.ExecutorAttrStore
}

func NewFunctionFactory(podDetailsStore *stores.LocalCacheHSetStore, attrStore *stores.ExecutorAttrStore) *FunctionFactory {
	return &FunctionFactory{podDetailsStore: podDetailsStore, attrStore: attrStore}
}

func (ff FunctionFactory) GetFunction(name string, args []string, attrStoreKey *cache.AttribStoreKey) *Function {

	defer func() {
		if r := recover(); r != nil {
			zkLogger.ErrorF(LoggerTag, "In GetFunction name:%s args:%v attrStoreKey:%v", name, args, attrStoreKey)
			zkLogger.ErrorF(LoggerTag, "In GetFunction: Recovered from panic: %v", r)
		}
	}()

	var fn Function
	switch name {
	case JsonExtract:
		fn = ExtractJson{name: name, args: args, attrStore: ff.attrStore, attrStoreKey: attrStoreKey, ff: &ff}
	case getWorkloadFromIP:
		fn = ExtractWorkLoadFromIP{name: name, args: args, attrStore: ff.attrStore, attrStoreKey: attrStoreKey, ff: &ff, podDetailsStore: ff.podDetailsStore}
	case toLowerCase:
		fn = LowerCase{name, args}
	case toUpperCase:
		fn = UpperCase{name, args}
	default:
		fn = NoNameFunction{name: NoName, args: args, attrStore: ff.attrStore, attrStoreKey: attrStoreKey, ff: &ff}
	}
	return &fn
}

func (ff FunctionFactory) HandleStringForFunctions(input string, attrStoreKey *cache.AttribStoreKey) []Function {
	return ff.GetPathAndFunctionsInternal(input, attrStoreKey, true)
}

func (ff FunctionFactory) GetPathAndFunctions(input string, attrStoreKey *cache.AttribStoreKey) []Function {
	return ff.GetPathAndFunctionsInternal(input, attrStoreKey, true)
}

func (ff FunctionFactory) GetPathAndFunctionsInternal(input string, attrStoreKey *cache.AttribStoreKey, allowNoNameFn bool) []Function {

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

	// create the functions
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

				if !allowNoNameFn && functionName == NoName {
					continue
				}
				fn = ff.GetFunction(functionName, args, attrStoreKey)
			}
		} else if allowNoNameFn {
			fn = ff.GetFunction(NoName, []string{match}, attrStoreKey)
		}

		if fn != nil {
			functions = append(functions, *fn)
		}
	}
	return functions
}

func (ff FunctionFactory) EvaluateString(inputPath string, store map[string]interface{}, attrStoreKey *cache.AttribStoreKey) (interface{}, bool) {
	return getValueFromStoreInternal(inputPath, store, &ff, attrStoreKey, true)
}

func getValueFromStoreInternal(inputPath string, store map[string]interface{}, ff *FunctionFactory, attrStoreKey *cache.AttribStoreKey, allowNoNameFn bool) (interface{}, bool) {

	defer func() {
		if r := recover(); r != nil {
			zkLogger.ErrorF(LoggerTag, "In EvaluateString: inputPath: %s \nstore:  %v \nattrStoreKey:%v", inputPath, store, attrStoreKey, r)
			zkLogger.ErrorF(LoggerTag, "In EvaluateString: Recovered from panic: %v", r)
		}
	}()

	var ok bool
	var valueAtObject interface{}
	var newValueAtObject interface{}

	valueAtObject = store
	functionArr := ff.GetPathAndFunctionsInternal(inputPath, attrStoreKey, allowNoNameFn)
	if len(functionArr) == 0 {
		return valueAtObject, false
	}

	// handle functionArr
	for _, fn := range functionArr {
		if valueAtObject == nil {
			return valueAtObject, false
		}
		newValueAtObject, ok = fn.Execute(valueAtObject)
		if !ok {
			return valueAtObject, false
		}
		valueAtObject = newValueAtObject
	}

	return valueAtObject, true
}
