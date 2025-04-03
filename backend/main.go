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

// Video struct (unchanged)
type Video struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Duration    int    `json:"duration"`
    UploadDate  string `json:"uploadDate"`
}

// In-memory store (unchanged)
var videos []Video

// Secret key for signing JWTs (in production, use an env variable)
var jwtKey = []byte("my_secret_key")

// Claims struct for JWT
type Claims struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func main() {
    router := mux.NewRouter()

    // Public endpoint
    router.HandleFunc("/login", login).Methods("POST")

    // Protected endpoints with auth middleware
    router.HandleFunc("/videos", getVideos).Methods("GET").Handler(authMiddleware(http.HandlerFunc(getVideos)))
    router.HandleFunc("/videos", createVideo).Methods("POST").Handler(authMiddleware(http.HandlerFunc(createVideo)))

    log.Println("Server starting on :8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
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

// Handlers (unchanged from Step 2)
func getVideos(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(videos)
}

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
