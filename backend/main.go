package main

import (
    "encoding/json"
    "log"
    "net/http"
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
var videos []Video

func main() {
    // Initialize router
    router := mux.NewRouter()

    // Define endpoints
    router.HandleFunc("/videos", getVideos).Methods("GET")
    router.HandleFunc("/videos", createVideo).Methods("POST")

    // Start server
    log.Println("Server starting on :8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
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
