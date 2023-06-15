package test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/crypto"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"reflect"
	"sort"
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

	output := common.ToString(inputStruct)
	assert.Equal(t, inputStructFmtString, *output)
}

func Test_StructUtils_NilInput_ToString_Success(t *testing.T) {
	output := common.ToString(nil)
	assert.Nil(t, output)
}

func Test_StructUtils_EmptyInput_ToString_Success(t *testing.T) {
	output := common.ToString(inputStruct{})
	assert.Equal(t, "{ { 0}}", *output)
}

func Test_StructUtils_NormalInput_ToReader_Success(t *testing.T) {
	inputStructFmtString := "{HelloFrom-inputStruct {HelloFrom-inputStructInner 55}}"

	output := common.ToReader(inputStructFmtString)
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

	output := common.ToJsonReader(inputStruct)
	buffer := make([]byte, len(inputStructJsonString))
	_, err := output.Read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, inputStructJsonString, string(buffer))
}

func Test_StructUtils_NilInput_ToJsonReader_Success(t *testing.T) {
	output := common.ToJsonReader(nil)
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

	output := common.ToJsonString(inputStruct)
	assert.Equal(t, inputStructJsonString, *output)
}

func Test_StructUtils_NilInput_ToJsonString_Success(t *testing.T) {
	output := common.ToJsonString(nil)
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

	output := common.FromJsonString(inputStructJsonString, reflect.TypeOf(inputStruct))
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

	output := common.FromJsonString(inputStructJsonString, reflect.TypeOf(&inputStruct))
	isDeepEqual := reflect.DeepEqual(output, &inputStruct)
	assert.Equal(t, true, isDeepEqual)
}

func Test_CryptoUtils_NormalInput_ToSha256_Success(t *testing.T) {
	input := "input-string"
	expectedSha256 := sha256.Sum256([]byte(input))
	output := common.ToSha256(input)
	assert.Equal(t, expectedSha256, output)
}

func Test_CryptoUtils_NormalInput_ToSha256String_Success(t *testing.T) {
	input := "input-string"
	prefix := "prefix"
	suffix := "suffix"
	sha256Bytes := sha256.Sum256([]byte(input))
	expectedOutput := prefix + hex.EncodeToString(sha256Bytes[:]) + suffix
	output := common.ToSha256String(prefix, input, suffix)
	assert.Equal(t, expectedOutput, output)
}

// General Utils Tests

func Test_ToPtr_Success(t *testing.T) {
	var input = "hello"
	output := common.ToPtr[string](input)
	assert.Equal(t, &input, output)
}

func Test_PtrTo_Success(t *testing.T) {
	var input = "hello"
	output := common.PtrTo[string](&input)
	assert.Equal(t, input, output)
}

func Test_PtrTo_Failure(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	common.PtrTo[string](nil)
}

func Test_GetIntegerFromString_Success(t *testing.T) {
	input := "55"
	expectedOutput := 55
	output, error := common.GetIntegerFromString(input)
	assert.Nil(t, error)
	assert.Equal(t, expectedOutput, output)
}

func Test_GetIntegerFromString_Failure(t *testing.T) {
	input := "55t"
	_, error := common.GetIntegerFromString(input)
	assert.NotNil(t, error)
}

func TestCalculateHash(t *testing.T) {
	input := "example string"
	expected := "0d355fb4-a093-520b-9d7e-f1f615e3fa10"
	result := crypto.CalculateHashNewSHA2(input)

	if result.String() != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

// TODO: discuss with mudit which one to keep, this or below one
func TestSortedScenariosHash1(t *testing.T) {
	unsortedWorkloadJS := string(common.GetBytesFromFile("files/unsortedWorkloadJs.json"))

	var wUnsorted model.Workload
	errUnsorted := json.Unmarshal([]byte(unsortedWorkloadJS), &wUnsorted)
	assert.NoError(t, errUnsorted)
	sort.Sort(wUnsorted.Rule.Rules)

	sortedWorkloadJS := string(common.GetBytesFromFile("files/sortedWorkloadJs.json"))

	var wSorted model.Workload
	errSorted := json.Unmarshal([]byte(sortedWorkloadJS), &wSorted)
	assert.NoError(t, errSorted)
	sort.Sort(wSorted.Rule.Rules)

	assert.Equal(t, model.WorkLoadUUID(wUnsorted), model.WorkLoadUUID(wSorted))
}

// TODO: discuss with mudit which one to keep, this or above one
func TestSortedScenariosHash2(t *testing.T) {
	unsortedWorkloadJS := string(common.GetBytesFromFile("files/unsortedWorkloadJs.json"))
	var wUnsorted model.Workload
	errUnsorted := json.Unmarshal([]byte(unsortedWorkloadJS), &wUnsorted)
	assert.NoError(t, errUnsorted)
	wUnsorted.Rule.Rules.Sort()
	sortedWorkloadJS := string(common.GetBytesFromFile("files/sortedWorkloadJs.json"))
	var wSorted model.Workload
	errSorted := json.Unmarshal([]byte(sortedWorkloadJS), &wSorted)
	assert.NoError(t, errSorted)
	sort.Sort(wSorted.Rule.Rules)
	a, b := model.WorkLoadUUID(wUnsorted), model.WorkLoadUUID(wSorted)

	assert.Equal(t, a, b)
}
