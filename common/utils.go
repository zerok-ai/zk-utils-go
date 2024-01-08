package common

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/lib/pq"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
	"io"
	"math"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var LogTag = "common_utils"

func ToString(iInstance interface{}) *string {
	if iInstance == nil {
		return nil
	}
	return ToPtr[string](fmt.Sprint(iInstance))
}

func ToReader(iString string) *strings.Reader {
	iReader := strings.NewReader(iString)
	return iReader
}

func ToJsonReader(iInstance interface{}) *strings.Reader {
	if iInstance == nil {
		return nil
	}
	iReader := strings.NewReader(*ToJsonString(iInstance))
	return iReader
}

func ToJsonString(iInstance interface{}) *string {
	if iInstance == nil {
		return nil
	}
	bytes, err := json.Marshal(iInstance)
	if err != nil {
		return nil
	} else {
		iString := string(bytes)
		return &iString

	}
}

func FromJsonString(iString string, iType reflect.Type) interface{} {
	if iType.Kind() == reflect.Ptr {
		iType = iType.Elem()
	}
	iTypeInterface := reflect.New(iType).Interface()
	iReader := strings.NewReader(iString)
	decoder := json.NewDecoder(iReader)
	err := decoder.Decode(iTypeInterface)
	if err != nil {
		//TODO:Refactor
	}
	return iTypeInterface
}

func ToFloat32(input string) (float32, error) {
	str := "3.14"
	f, err := strconv.ParseFloat(str, 32)
	if err != nil {
		zkLogger.ErrorF(LogTag, "Error while converting string %s to float32, err %v.\n", str, err)
		return float32(0.0), err
	}
	return float32(f), nil
}

// String Utils

func ToSha256(input string) [sha256.Size]byte {
	return sha256.Sum256([]byte(input))
}

func ToSha256String(prefix string, input string, suffix string) string {
	bytes := ToSha256(input)
	return prefix + hex.EncodeToString(bytes[:]) + suffix
}

// General Utils

func GetFloatFromString(k string, b int) (float64, error) {
	return strconv.ParseFloat(k, b)
}

func GetIntegerFromString(k string) (int, error) {
	return strconv.Atoi(k)
}

func ToPtr[T any](arg T) *T {
	return &arg
}

func PtrTo[T any](arg *T) T {
	if arg == nil {
		panic(errors.New("PtrTo - Passed pointer is nil"))
	}
	return *arg
}

func GenerateRandomToken(length int) (string, error) {

	tokenLength := length

	tokenBytes := make([]byte, tokenLength)

	_, err := rand.Read(tokenBytes)
	if err != nil {
		zkLogger.Error(LogTag, "Error while creating random bytes of length %v, err %v.\n", length, err)
		return "", err
	}

	token := base64.StdEncoding.EncodeToString(tokenBytes)
	return token, nil
}

func RollbackTransaction(tx *sql.Tx, logTag string) (bool, *zkErrors.ZkError) {
	dbErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)

	if err := tx.Rollback(); err != nil {
		zkLogger.Debug(logTag, "unable to rollback transaction, "+err.Error())
	}
	return false, &dbErr
}

func CommitTransaction(tx *sql.Tx, logTag string) (bool, *zkErrors.ZkError) {
	if err := tx.Commit(); err != nil {
		zkLogger.Debug(logTag, "unable to commit transaction, "+err.Error())
		return RollbackTransaction(tx, logTag)
	}
	return true, nil
}

func CurrentTime() time.Time {
	return time.Now()
}

func Generate256SHA(params ...string) string {
	currentDate := CurrentTime().Format("2006#01#02")
	salt := "ydqnk@93765"
	currentTokenRaw := currentDate + salt
	for _, v := range params {
		currentTokenRaw += v
	}

	currentTokenBytes := []byte(currentTokenRaw)
	currentTokenSha256Hash := sha256.Sum256(currentTokenBytes)
	return hex.EncodeToString(currentTokenSha256Hash[:])
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func IsEmpty(v string) bool {
	return len(v) == 0
}

// Round Rounds to nearest like 12.3456 -> 12.35
func Round(val float64, precision int) float64 {
	return math.Round(val*(math.Pow10(precision))) / math.Pow10(precision)
}

func SetResponseInCtxAndReturn[T any](ctx iris.Context, resp *T, zkError *zkErrors.ZkError) {
	if zkError != nil {
		zkHttpResponse := zkHttp.ZkHttpResponseBuilder[T]{}.WithZkErrorType(zkError.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	z := &zkHttp.ZkHttpResponseBuilder[T]{}
	zkHttpResponse := z.WithStatus(iris.StatusOK).Data(resp).Build()
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
	return
}

func ExtractRegexStringFromString(target, jsonExtractPattern string) string {

	// Compile the regular expressions
	jsonExtractRegex := regexp.MustCompile(jsonExtractPattern)

	// Extract response_payload and message using regular expressions
	message := jsonExtractRegex.FindStringSubmatch(target)

	if len(message) > 1 {
		return message[1]
	}
	return ""
}

func GetBytesFromFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		zkLogger.Error(LogTag, "Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		zkLogger.Error(LogTag, "Error reading file:", err)
		return nil
	}

	return content
}

func GetStmtForCopyIn(tx *sql.Tx, tableName string, columns []string) (*sql.Stmt, error) {
	stmt, err := tx.Prepare(pq.CopyIn(tableName, columns...))
	if err != nil {
		zkLogger.Error(LogTag, "Error preparing insert statement:", err)
		return nil, err
	}
	return stmt, nil
}

func GetStmtRawQuery(tx *sql.Tx, stmt string) (*sql.Stmt, error) {
	preparedStmt, err := tx.Prepare(stmt)
	if err != nil {
		zkLogger.Error(LogTag, "Error preparing insert statement:", err)
		return nil, err
	}
	return preparedStmt, nil
}

func DeepCopy[T any](input *T) (*T, error) {

	byteArray, err := json.Marshal(input)
	if err != nil {
		zkLogger.Error(LogTag, "Error in Deep copy[1]:", err)
		return nil, err
	}

	var newObject T
	err = json.Unmarshal(byteArray, &newObject)
	if err != nil {
		zkLogger.Error(LogTag, "Error in Deep copy[2]:", err)
		return nil, err
	}

	return &newObject, nil
}

type Cleanup func() error

func BlockUntilChannelClosed(cleanup Cleanup) {

	// Create a channel to receive signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Press Ctrl+C (SIGINT) or send SIGTERM to exit.")

	// Block until a signal is received
	<-sig

	fmt.Println("Received signal. Shutting down gracefully...")
	// Additional cleanup and shutdown logic can be added here

	// Cleanup
	if cleanup != nil {
		err := cleanup()
		if err != nil {
			zkLogger.Error(LogTag, "Got error during cleanup:", err)
			return
		}
	}

	// Exit the application
	os.Exit(0)
}
