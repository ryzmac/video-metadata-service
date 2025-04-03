package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/golang-jwt/jwt/v5"
    "github.com/gorilla/mux"
)

func TestGetVideos(t *testing.T) {
    // Reset videos
    videos = []Video{}

    // Set up router
    router := mux.NewRouter()
    router.HandleFunc("/videos", getVideos).Methods("GET").Handler(authMiddleware(http.HandlerFunc(getVideos)))

    // Create request with token
    req, err := http.NewRequest("GET", "/videos", nil)
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Authorization", "Bearer "+generateTestToken(t))

    // Record response
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Expected status 200, got %v, body: %v", status, rr.Body.String())
    }

    expected := "[]"
    if rr.Body.String() != expected {
        t.Errorf("Expected body %v, got %v", expected, rr.Body.String())
    }
}

func TestCreateVideo(t *testing.T) {
    videos = []Video{}

    router := mux.NewRouter()
    router.HandleFunc("/videos", createVideo).Methods("POST").Handler(authMiddleware(http.HandlerFunc(createVideo)))

    video := Video{ID: "1", Title: "Test", Description: "Test video", Duration: 60, UploadDate: "2025-04-05"}
    body, _ := json.Marshal(video)
    req, err := http.NewRequest("POST", "/videos", bytes.NewBuffer(body))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Authorization", "Bearer "+generateTestToken(t))
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Expected status 200, got %v, body: %v", status, rr.Body.String())
    }

    var returnedVideo Video
    if err := json.Unmarshal(rr.Body.Bytes(), &returnedVideo); err != nil {
        t.Errorf("Failed to unmarshal response: %v", err)
    }
    if returnedVideo != video {
        t.Errorf("Expected video %v, got %v", returnedVideo, video)
    }
}

// Helper to generate a valid test token
func generateTestToken(t *testing.T) string {
    claims := &Claims{
        Username: "demo",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        t.Fatal(err)
    }
    return tokenString
}
