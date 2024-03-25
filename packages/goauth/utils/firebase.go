package authutils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"firebase.google.com/go/auth"
)

const exchangeIDTokenTemplate = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=%s"

// Init firebase auth token for user with uid.
// Return id token and auth token
func LoadFirebaseAuthForUser(
	client *auth.Client,
	uid string,
	webAPIKey string,
) (string, *auth.Token, error) {
	customToken, _ := client.CustomToken(context.TODO(), uid)

	requestBody := map[string]any{
		"token":             customToken,
		"returnSecureToken": true,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	res, err := http.Post(
		fmt.Sprintf(exchangeIDTokenTemplate, webAPIKey),
		"application/json",
		bytes.NewBuffer(requestBodyBytes),
	)
	if err != nil {
		return "", nil, err
	}

	resBodyBytes, _ := io.ReadAll(res.Body)
	var resBody map[string]any
	err = json.Unmarshal(resBodyBytes, &resBody)
	if err != nil {
		return "", nil, err
	}

	idToken := resBody["idToken"].(string)
	authToken, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return "", nil, err
	}

	return idToken, authToken, nil
}

// Init firebase auth token for user with uid (with cache file)
// Return id token and auth token
func LoadFirebaseAuthForUserWithCache(
	client *auth.Client,
	uid string,
	webAPIKey string,
	cacheFile string,
) (string, *auth.Token, error) {
	if cacheFile == "" {
		cacheFile = "auth.json"
	}
	idToken, authToken, err := LoadFirebaseCredentials(client, cacheFile)
	if err != nil || authToken.UID != uid {
		idToken, authToken, err = LoadFirebaseAuthForUser(client, uid, webAPIKey)
		err := StoreFirebaseCredentials(idToken, cacheFile)
		if err != nil {
			return "", nil, err
		}
	}

	return idToken, authToken, err
}

// Load idToken from json file {"idToken": "..."}, verify it and generate actual firebase auth token
func LoadFirebaseCredentials(client *auth.Client, filename string) (string, *auth.Token, error) {
	dataBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", nil, err
	}

	var data map[string]string
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return "", nil, err
	}

	idToken := data["idToken"]
	token, err := client.VerifyIDToken(context.TODO(), idToken)
	if err != nil {
		return "", nil, err
	}

	return idToken, token, nil
}

// Store idToken to json file {"idToken": "..."}
func StoreFirebaseCredentials(idToken string, filename string) error {
	data := map[string]string{
		"idToken": idToken,
	}
	dataBytes, _ := json.MarshalIndent(data, "", "    ")

	return os.WriteFile(filename, dataBytes, 0o777)
}
