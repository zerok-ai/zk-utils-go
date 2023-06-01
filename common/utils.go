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
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type StructUtils interface {
	ToString(iInstance interface{}) *string
	ToReader(iString string) *strings.Reader
	ToJsonReader(iInstance interface{}) *strings.Reader
	ToJsonString(iInstance interface{}) *string
	FromJsonString(iString string, iType reflect.Type) interface{}
}

type structUtils struct {
}

func NewStructUtils() StructUtils {
	return &structUtils{}
}

func (structUtils structUtils) ToString(iInstance interface{}) *string {
	if iInstance == nil {
		return nil
	}
	return ToPtr[string](fmt.Sprint(iInstance))
}

func (structUtils structUtils) ToReader(iString string) *strings.Reader {
	iReader := strings.NewReader(iString)
	return iReader
}

func (structUtils structUtils) ToJsonReader(iInstance interface{}) *strings.Reader {
	if iInstance == nil {
		return nil
	}
	iReader := strings.NewReader(*structUtils.ToJsonString(iInstance))
	return iReader
}

func (structUtils structUtils) ToJsonString(iInstance interface{}) *string {
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

func (structUtils structUtils) FromJsonString(iString string, iType reflect.Type) interface{} {
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
type CryptoUtils struct {
}

func (cryptoUtils CryptoUtils) ToSha256(input string) [sha256.Size]byte {
	return sha256.Sum256([]byte(input))
}

func (cryptoUtils CryptoUtils) ToSha256String(prefix string, input string, suffix string) string {
	bytes := cryptoUtils.ToSha256(input)
	return prefix + hex.EncodeToString(bytes[:]) + suffix
}

// General Utils
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

func (cryptoUtils CryptoUtils) GenerateRandomToken(length int) (string, error) {

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
		zkError := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZK_ERROR_NOT_FOUND, nil)
		return &zkError
	case nil:
		return nil
	default:
		zkError := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZK_ERROR_INTERNAL_SERVER, nil)
		zkLogger.Debug(logTag, "unable to scan rows", err)
		return &zkError
	}
}

func RollbackTransaction(tx *sql.Tx, logTag string) (bool, *zkErrors.ZkError) {
	dbErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZK_ERROR_DB_ERROR, nil)

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
