package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/gorilla/mux"
)

// Video represents the metadata structure
type Video struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Duration    int    `json:"duration"` // in seconds
    UploadDate  string `json:"uploadDate"`
}

// In-memory store (we’ll swap this for a database later if you want)
var videos = []Video{}  // Changed from var videos []Video

// Secret key for signing JWTs (in production, use an env variable)
var jwtKey = []byte("my_secret_key")

// Claims struct for JWT
type Claims struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func main() {
    router := mux.NewRouter()

    // Define routes
    router.HandleFunc("/login", login).Methods("POST")
    router.HandleFunc("/videos", getVideos).Methods("GET").Handler(authMiddleware(http.HandlerFunc(getVideos)))
    router.HandleFunc("/videos", createVideo).Methods("POST").Handler(authMiddleware(http.HandlerFunc(createVideo)))

    // Apply CORS middleware to all routes
    corsHandler := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
            w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            next.ServeHTTP(w, r)
        })
    }

    // Start server with CORS-wrapped router
    log.Println("Server starting on :8080...")
    log.Fatal(http.ListenAndServe(":8080", corsHandler(router)))
}

// Login handler to issue a JWT
func login(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var creds struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Dummy check (in real apps, verify against a database)
    if creds.Username != "demo" || creds.Password != "password" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Create JWT
    expirationTime := time.Now().Add(1 * time.Hour)
    claims := &Claims{
        Username: creds.Username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// Middleware to validate JWT
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Token is valid; proceed to the next handler
        next.ServeHTTP(w, r)
    })
}

// Handler to return all videos
func getVideos(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(videos)
}

// Handler to create a new video
func createVideo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var video Video
    err := json.NewDecoder(r.Body).Decode(&video)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    videos = append(videos, video)
    json.NewEncoder(w).Encode(video)
}
