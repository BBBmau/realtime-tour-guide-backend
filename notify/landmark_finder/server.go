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

	// Geocode origin
	origResults, err := c.Geocode(context.Background(), &maps.GeocodingRequest{
		Address: currentLocation,
	})
	if err != nil || len(origResults) == 0 {
		return "", fmt.Errorf("failed to geocode origin: %v", err)
	}

	// Geocode destination
	destResults, err := c.Geocode(context.Background(), &maps.GeocodingRequest{
		Address: destination,
	})
	if err != nil || len(destResults) == 0 {
		return "", fmt.Errorf("failed to geocode destination: %v", err)
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

	// Build prompt with route information and coordinates
	prompt := fmt.Sprintf("As a driver going from %s (%.6f, %.6f) to %s (%.6f, %.6f), here are important points along your route:\n",
		currentLocation, origResults[0].Geometry.Location.Lat, origResults[0].Geometry.Location.Lng,
		destination, destResults[0].Geometry.Location.Lat, destResults[0].Geometry.Location.Lng)

	fmt.Println(prompt)
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
					Role:    "system",
					Content: "You are currently on roadtrip helping the driver find unique places, minimum 3, along their route that will provide the driver an unforgettable experience. It should be based on where they are from the given directions, no more than 10 miles away. You should also provide the address of the landmark and the distance from the current location. It should also be straight to the point while also being excited about the places you are finding.",
				},
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
		nearbyInformation, err := assistantAudioRequest(nearbyInformationPrompt)
		fmt.Println(nearbyInformation)
		w.Write([]byte(nearbyInformation))
	})

	port := "8080"
	log.Printf("Server running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
