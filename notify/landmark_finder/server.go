package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sashabaranov/go-openai"

	"github.com/joho/godotenv"
)

type RouteFinderRequest struct {
	CurrentLocation string `json:"current_location"`
	Destination     string `json:"destination"`
	// TODO: Interests       string `json:"interests"`
}

func GoogleRouteFinder(currentLocation string, destination string) (string, error) {

	return "", nil
}

// returns the audio in a base64 encoded string for .wav format
func assistantAudioRequest(nearbyInformationPrompt string) (string, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: nearbyInformationPrompt,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("Error: %v\n", err)
	}
	return resp.Choices[0].Message.Content, nil
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// Add console log to show server is starting
	log.Println("Starting server...")

	// Handle root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})

	// Handle notification endpoint
	http.HandleFunc("/notification", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		log.Println("Received request to /notification")
		// Get query parameters with defaults
		currentLocation := r.URL.Query().Get("current_location")
		destination := r.URL.Query().Get("destination")
		nearbyInformationPrompt, err := GoogleRouteFinder(currentLocation, destination) // This is where the prompt is being made based on the Maps API response
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: the nearbyInformation should be in an audio format that gets sent back to the user in the response.
		nearbyInformation, err := assistantAudioRequest(nearbyInformationPrompt)

		w.Write([]byte(nearbyInformation))
	})

	port := "8080"
	log.Printf("Server running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
