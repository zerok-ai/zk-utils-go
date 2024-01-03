package zkClient

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type ClusterKeyData struct {
	TokenString string `json:"token"`
	ExpiresAt   int64  `json:"expiresAt"`
}

func DecodeToken(base64Str string) (*ClusterKeyData, error) {
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64 string: %w", err)
	}

	var tokenData ClusterKeyData
	err = json.Unmarshal(data, &tokenData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return &tokenData, nil
}
