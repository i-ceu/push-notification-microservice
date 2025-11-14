package config

import (
	"encoding/base64"
	"fmt"
	"os"
)

func GetFirebaseCredentials() (string, error) {
	// if path := os.Getenv("FCM_SERVICE_ACCOUNT"); path != "" {
	// 	return path, nil
	// }
	if base64Key := os.Getenv("FCM_SERVICE_ACCOUNT_BASE64"); base64Key != "" {
		decoded, err := base64.StdEncoding.DecodeString(base64Key)
		if err != nil {
			return "", fmt.Errorf("failed to decode FCM_SERVICE_ACCOUNT_BASE64: %w", err)
		}

		tmpFile := "tmp/firebase-credentials.json"
		os.MkdirAll("tmp", 0755)
		err = os.WriteFile(tmpFile, decoded, 0600)
		if err != nil {
			return "", fmt.Errorf("failed to write credentials file: %w", err)
		}
		return tmpFile, nil
	}

	return "", fmt.Errorf("no Firebase credentials found")
}
