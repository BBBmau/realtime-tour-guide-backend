package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sashabaranov/go-openai"
	"googlemaps.github.io/maps"

	"github.com/joho/godotenv"
)

type RouteFinderRequest struct {
	CurrentLocation string `json:"current_location"`
	Destination     string `json:"destination"`
	// TODO: Interests       string `json:"interests"`
}

func GoogleRouteFinder(request *http.Request) (string, error) {
	currentLocation := request.URL.Query().Get("current_location")
	destination := request.URL.Query().Get("destination")

	// Create Google Maps client
	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_API_KEY")))
	if err != nil {
		return "", fmt.Errorf("failed to create maps client: %v", err)
	}

	// Create directions request
	r := &maps.DirectionsRequest{
		Origin:      currentLocation,
		Destination: destination,
	}

	// Get route
	routes, _, err := c.Directions(context.Background(), r)
	if err != nil {
		return "", fmt.Errorf("failed to get directions: %v", err)
	}

	if len(routes) == 0 {
		return "", fmt.Errorf("no routes found")
	}

	// Build prompt with route information
	prompt := fmt.Sprintf("As a driver going from %s to %s, here are important points along your route:\n",
		currentLocation, destination)

	// Add route information
	for _, leg := range routes[0].Legs {
		prompt += fmt.Sprintf("\nTotal distance: %s\nEstimated duration: %s\n",
			leg.Distance.HumanReadable, leg.Duration.String())

		// Add important steps
		for _, step := range leg.Steps {
			if step.HTMLInstructions != "" {
				prompt += fmt.Sprintf("- %s\n", step.HTMLInstructions)
			}
		}
	}

	return prompt, nil
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
		nearbyInformationPrompt, err := GoogleRouteFinder(r) // This is where the prompt is being made based on the Maps API response
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: the nearbyInformation should be in an audio format that gets sent back to the user in the response.
		// nearbyInformation, err := assistantAudioRequest(nearbyInformationPrompt)

		w.Write([]byte(nearbyInformationPrompt))
	})

	port := "8080"
	log.Printf("Server running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
