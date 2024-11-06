package webserver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wojcikp/go-weather-go/weather-hub/internal"
)

type ScoresServer struct {
	scoresInfo []internal.ScoreInfo
}

func NewScoresServer() *ScoresServer {
	return &ScoresServer{}
}

func (s *ScoresServer) SetScoresInfo(scores []internal.ScoreInfo) {
	s.scoresInfo = scores
}

func (s ScoresServer) RunWeatherScoresServer() {
	http.HandleFunc("/scores", s.scoresHandler)

	log.Println("Server is running on http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

func (s ScoresServer) scoresHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.URL.Path != "/scores" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(s.scoresInfo); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
