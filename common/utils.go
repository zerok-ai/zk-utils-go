package zkcommon

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
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
	"io"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

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
	bytes, error := json.Marshal(iInstance)
	if error != nil {
		//TODO:Refactor
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
	error := decoder.Decode(iTypeInterface)
	if error != nil {
		//TODO:Refactor
	}
	return iTypeInterface
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
		fmt.Printf("Error while creating random bytes of length %v, err %v.\n", length, err)
		return "", err
	}

	token := base64.StdEncoding.EncodeToString(tokenBytes)

	fmt.Println("Random token:", token)

	return token, nil
}

func CheckSqlError(err error, logTag string) *zkErrors.ZkError {
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		zkError := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorNotFound, nil)
		return &zkError
	case nil:
		return nil
	default:
		zkError := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorInternalServer, nil)
		zkLogger.Debug(logTag, "unable to scan rows", err)
		return &zkError
	}
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
		z := &zkHttp.ZkHttpResponseBuilder[T]{}
		z.WithZkErrorType(zkError.Error).Build()
		ctx.StatusCode(zkError.Error.Status)
		return
	}

	z := &zkHttp.ZkHttpResponseBuilder[T]{}
	zkHttpResponse := z.WithStatus(iris.StatusOK).Data(resp).Build()
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
	return
}

func GetBytesFromFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	return content

}
