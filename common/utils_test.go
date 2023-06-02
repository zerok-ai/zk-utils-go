package zkcommon

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type inputStructInner struct {
	StringKey  string `json:"stringKey"`
	IntegerKey int    `json:"integerKey"`
}

type inputStruct struct {
	StringKey string           `json:"stringKey"`
	StructKey inputStructInner `json:"structKey"`
}

type StructUtilsTestSuite struct {
	inputStruct           inputStruct
	inputStructFmtString  string
	inputStructJsonString string
}

func (structUtilsTestSuite *StructUtilsTestSuite) SetupTest() {
	structUtilsTestSuite.inputStructFmtString = "{HelloFrom-inputStruct {HelloFrom-inputStructInner 55}}"
	structUtilsTestSuite.inputStructJsonString = "{\"stringKey\":\"HelloFrom-inputStruct\",\"structKey\":{\"stringKey\":\"HelloFrom-inputStructInner\",\"integerKey\":55}}"
	structUtilsTestSuite.inputStruct = inputStruct{
		StringKey: "HelloFrom-inputStruct",
		StructKey: inputStructInner{
			StringKey:  "HelloFrom-inputStructInner",
			IntegerKey: 55,
		},
	}
}

// Test_<StructName>_<Case-Description>_<MethodName>_<Success/Failure>
func Test_StructUtils_NormalInput_ToString_Success(t *testing.T) {
	inputStructFmtString := "{HelloFrom-inputStruct {HelloFrom-inputStructInner 55}}"
	inputStruct := inputStruct{
		StringKey: "HelloFrom-inputStruct",
		StructKey: inputStructInner{
			StringKey:  "HelloFrom-inputStructInner",
			IntegerKey: 55,
		},
	}

	output := ToString(inputStruct)
	assert.Equal(t, inputStructFmtString, *output)
}

func Test_StructUtils_NilInput_ToString_Success(t *testing.T) {
	output := ToString(nil)
	assert.Nil(t, output)
}

func Test_StructUtils_EmptyInput_ToString_Success(t *testing.T) {
	output := ToString(inputStruct{})
	assert.Equal(t, "{ { 0}}", *output)
}

func Test_StructUtils_NormalInput_ToReader_Success(t *testing.T) {
	inputStructFmtString := "{HelloFrom-inputStruct {HelloFrom-inputStructInner 55}}"

	output := ToReader(inputStructFmtString)
	buffer := make([]byte, len(inputStructFmtString))
	_, err := output.Read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, inputStructFmtString, string(buffer))
}

func Test_StructUtils_NormalInput_ToJsonReader_Success(t *testing.T) {
	inputStructJsonString := "{\"stringKey\":\"HelloFrom-inputStruct\",\"structKey\":{\"stringKey\":\"HelloFrom-inputStructInner\",\"integerKey\":55}}"
	inputStruct := inputStruct{
		StringKey: "HelloFrom-inputStruct",
		StructKey: inputStructInner{
			StringKey:  "HelloFrom-inputStructInner",
			IntegerKey: 55,
		},
	}

	output := ToJsonReader(inputStruct)
	buffer := make([]byte, len(inputStructJsonString))
	_, err := output.Read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, inputStructJsonString, string(buffer))
}

func Test_StructUtils_NilInput_ToJsonReader_Success(t *testing.T) {
	output := ToJsonReader(nil)
	assert.Nil(t, output)
}

func Test_StructUtils_NormalInput_ToJsonString_Success(t *testing.T) {
	inputStructJsonString := "{\"stringKey\":\"HelloFrom-inputStruct\",\"structKey\":{\"stringKey\":\"HelloFrom-inputStructInner\",\"integerKey\":55}}"
	inputStruct := inputStruct{
		StringKey: "HelloFrom-inputStruct",
		StructKey: inputStructInner{
			StringKey:  "HelloFrom-inputStructInner",
			IntegerKey: 55,
		},
	}

	output := ToJsonString(inputStruct)
	assert.Equal(t, inputStructJsonString, *output)
}

func Test_StructUtils_NilInput_ToJsonString_Success(t *testing.T) {
	output := ToJsonString(nil)
	assert.Nil(t, output)
}

func Test_StructUtils_NormalInput_FromString_Success(t *testing.T) {
	inputStructJsonString := "{\"stringKey\":\"HelloFrom-inputStruct\",\"structKey\":{\"stringKey\":\"HelloFrom-inputStructInner\",\"integerKey\":55}}"
	inputStruct := inputStruct{
		StringKey: "HelloFrom-inputStruct",
		StructKey: inputStructInner{
			StringKey:  "HelloFrom-inputStructInner",
			IntegerKey: 55,
		},
	}

	output := FromJsonString(inputStructJsonString, reflect.TypeOf(inputStruct))
	isDeepEqual := reflect.DeepEqual(output, &inputStruct)
	assert.Equal(t, true, isDeepEqual)
}

func Test_StructUtils_PtrInput_FromString_Success(t *testing.T) {
	inputStructJsonString := "{\"stringKey\":\"HelloFrom-inputStruct\",\"structKey\":{\"stringKey\":\"HelloFrom-inputStructInner\",\"integerKey\":55}}"
	inputStruct := inputStruct{
		StringKey: "HelloFrom-inputStruct",
		StructKey: inputStructInner{
			StringKey:  "HelloFrom-inputStructInner",
			IntegerKey: 55,
		},
	}

	output := FromJsonString(inputStructJsonString, reflect.TypeOf(&inputStruct))
	isDeepEqual := reflect.DeepEqual(output, &inputStruct)
	assert.Equal(t, true, isDeepEqual)
}

func Test_CryptoUtils_NormalInput_ToSha256_Success(t *testing.T) {
	input := "input-string"
	expectedSha256 := sha256.Sum256([]byte(input))
	output := ToSha256(input)
	assert.Equal(t, expectedSha256, output)
}

func Test_CryptoUtils_NormalInput_ToSha256String_Success(t *testing.T) {
	input := "input-string"
	prefix := "prefix"
	suffix := "suffix"
	sha256Bytes := sha256.Sum256([]byte(input))
	expectedOutput := prefix + hex.EncodeToString(sha256Bytes[:]) + suffix
	output := ToSha256String(prefix, input, suffix)
	assert.Equal(t, expectedOutput, output)
}

// General Utils Tests

func Test_ToPtr_Success(t *testing.T) {
	var input = "hello"
	output := ToPtr[string](input)
	assert.Equal(t, &input, output)
}

func Test_PtrTo_Success(t *testing.T) {
	var input = "hello"
	output := PtrTo[string](&input)
	assert.Equal(t, input, output)
}

func Test_PtrTo_Failure(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	PtrTo[string](nil)
}

func Test_GetIntegerFromString_Success(t *testing.T) {
	input := "55"
	expectedOutput := 55
	output, error := GetIntegerFromString(input)
	assert.Nil(t, error)
	assert.Equal(t, expectedOutput, output)
}

func Test_GetIntegerFromString_Failure(t *testing.T) {
	input := "55t"
	_, error := GetIntegerFromString(input)
	assert.NotNil(t, error)
}
