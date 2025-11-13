package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"push-notification-microservice/internal/models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type FCMSender struct {
	projectID   string
	httpClient  *http.Client
	tokenSource oauth2.TokenSource
}

func NewFCMSender(serviceAccountPath string) (*FCMSender, error) {
	// Read service account file
	jsonKey, err := os.ReadFile(serviceAccountPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read service account: %w", err)
	}

	var serviceAccount struct {
		ProjectID string `json:"project_id"`
	}
	if err := json.Unmarshal(jsonKey, &serviceAccount); err != nil {
		return nil, fmt.Errorf("failed to parse service account: %w", err)
	}

	// Create credentials and token source
	creds, err := google.CredentialsFromJSON(
		context.Background(),
		jsonKey,
		"https://www.googleapis.com/auth/firebase.messaging",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create credentials: %w", err)
	}

	return &FCMSender{
		projectID:   serviceAccount.ProjectID,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		tokenSource: creds.TokenSource,
	}, nil
}

func (fs *FCMSender) Send(notification *models.PushNotification) error {
	ctx := context.Background()

	// Get access token
	token, err := fs.getAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// FCM v1 API endpoint
	url := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", fs.projectID)

	// Create message payload - using type-safe approach
	message := fs.buildMessage(notification)

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Send request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := fs.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("FCM returned %d: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("Successfully sent notification to %s\n", notification.PushToken)
	return nil
}

func (fs *FCMSender) getAccessToken() (string, error) {
	token, err := fs.tokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	if !token.Valid() {
		return "", fmt.Errorf("token is invalid")
	}

	return token.AccessToken, nil
}

func (fs *FCMSender) buildMessage(notification *models.PushNotification) map[string]interface{} {
	message := map[string]interface{}{
		"message": map[string]interface{}{
			"token": notification.PushToken,
			"notification": map[string]interface{}{
				"title": notification.Data.Title,
				"body":  notification.Data.Body,
			},
		},
	}

	return message
}

func (fs *FCMSender) SendToMultiple(notifications []*models.PushNotification) []error {
	var errors []error

	for _, notification := range notifications {
		if err := fs.Send(notification); err != nil {
			errors = append(errors, fmt.Errorf("failed to send to %s: %w", notification.PushToken, err))
		}
	}

	return errors
}
