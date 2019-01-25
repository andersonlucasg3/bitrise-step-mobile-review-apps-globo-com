package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

var appNameKey = "app_name"
var artifactURLKey = "artifact_url"
var branchKey = "branch"
var commitHashKey = "commit_hash"
var endPointAuthTokenKey = "SERVICE_AUTH_TOKEN"
var endPoint = "https://fake.globo.com/apps/mobile/"

type AppInfo struct {
	App         string `json:"app"`
	ArtifactURL string `json:"artifact_url"`
	Branch      string `json:"branch"`
	CommitHash  string `json:"commit_hash"`
}

func failedValidation(condition bool, inputVarName string, inputVarKey string) {
	if condition {
		fmt.Printf("%s not provided, please fulfill %s`s input variable.", inputVarName, inputVarKey)
		os.Exit(1)
	} else {
		fmt.Printf("Input variable %s provided correctly.", inputVarKey)
		fmt.Println()
	}
}

func validateAuthToken(token string) {
	failedValidation(len(token) < 10, "Auth token", endPointAuthTokenKey)
}

func validateAppName(name string) {
	failedValidation(len(name) == 0, "App name", appNameKey)
}

func validateArtifactURL(artifact string) {
	_, err := url.Parse(artifact)
	failedValidation(err != nil, "Artifact URL", artifactURLKey)
}

func validateBranch(branch string) {
	failedValidation(len(branch) == 0, "Branch name", branchKey)
}

func validateHash(hash string) {
	failedValidation(len(hash) < 7, "Commit hash", commitHashKey)
}

func validateAll(appName string, artifact string, branch string, hash string) {
	fmt.Printf("Input variable '%s': %s", appNameKey, appName)
	fmt.Println()
	fmt.Printf("Input variable '%s': %s", artifactURLKey, artifact)
	fmt.Println()
	fmt.Printf("Input variable '%s': %s", branchKey, branch)
	fmt.Println()
	fmt.Printf("Input variable '%s': %s", commitHashKey, hash)
	fmt.Println()

	validateAppName(appName)
	validateArtifactURL(artifact)
	validateBranch(branch)
	validateHash(hash)
}

func encodeJSON(object interface{}) *bytes.Buffer {
	b := &bytes.Buffer{}
	err := json.NewEncoder(b).Encode(object)
	if err != nil {
		fmt.Println("Failed to Encode Json Object: ", object)
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	return b
}

func makeRequest(byteArray *bytes.Buffer, authToken string) {
	req, err := http.NewRequest(http.MethodPost, endPoint, byteArray)
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", authToken))
	if err != nil {
		fmt.Println("Failed to create new request with EndPoint: ", endPoint)
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Failed to execute request with EndPoint: ", endPoint)
		fmt.Println("Error: ", err)
		os.Exit(1)
	} else if response.StatusCode != 201 {
		fmt.Println("Request data: ", *req)
		fmt.Println("Response Status Code: ", response.StatusCode)
		fmt.Println("Response string: ", response.Status)
		os.Exit(1)
	}
}

func assertSuccess() {
	fmt.Println("Mobile Review created with success.")
	os.Exit(0)
}

func main() {
	appName := os.Getenv(appNameKey)
	artifactURL := os.Getenv(artifactURLKey)
	branch := os.Getenv(branchKey)
	commitHash := os.Getenv(commitHashKey)
	authToken := os.Getenv(endPointAuthTokenKey)

	validateAll(appName, artifactURL, branch, commitHash)
	validateAuthToken(authToken)

	jsonObject := AppInfo{
		App:         appName,
		ArtifactURL: artifactURL,
		Branch:      branch,
		CommitHash:  commitHash,
	}

	b := encodeJSON(jsonObject)
	makeRequest(b, authToken)

	assertSuccess()
}
