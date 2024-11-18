package gemini

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"

	"strings"
	"regexp"
)

// Load environment variables from the .env file during initialization
func init() {
	// Attempt to load variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}
}

// Create a client for Gemini API
func GetGeminiClient() (*genai.Client, error) {
	// Retrieve the API key from environment variables
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API_KEY not set in environment variables")
	}

	// Initialize the Gemini API client using the API key
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating Gemini client: %v", err)
	}
	return client, nil
}

// GenerateResponse sends the message to Gemini API and returns a response.
func GenerateResponse(message string) (string, error) {
	// Create the Gemini client
	client, err := GetGeminiClient()
	if err != nil {
		return "", fmt.Errorf("failed to get Gemini client: %w", err)
	}
	defer client.Close() // Ensure the client is properly closed after usage

	// Create a context with a timeout to limit the API call duration
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel() // Cancel the context to release resources

	// Use the primary model for generating content
	model := client.GenerativeModel("gemini-1.5-flash-latest")
	resp, err := model.GenerateContent(ctx, genai.Text(message))

	if err != nil {
		// If the primary model fails, log the error and switch to the fallback model
		log.Printf("Primary model error: %v. Switching to fallback model.", err)
		model = client.GenerativeModel("gemini-1.5-pro")
		resp, err = model.GenerateContent(ctx, genai.Text(message))
		if err != nil {
			return "", fmt.Errorf("error generating response from fallback model: %w", err)
		}
	}

	// Extract the response text from the generated candidates
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		rawContent := fmt.Sprintf("%v", resp.Candidates[0].Content) // Convert the content to a string

		// Clean the raw response content by removing unnecessary characters
		cleanContent := cleanResponse(rawContent)
		return cleanContent, nil
	}

	return "", fmt.Errorf("no response text found in the generated content")
}

// cleanResponse removes unnecessary characters from the response
func cleanResponse(rawContent string) string {
	// Regular expression to remove unwanted characters like '&', '{', '}', '[', ']', and the word 'model'
	re := regexp.MustCompile(`[&{}\[\]]+|model`)
	cleanContent := re.ReplaceAllString(rawContent, "") // Replace matches with an empty string

	// Trim leading and trailing whitespace from the cleaned content
	return strings.TrimSpace(cleanContent)
}
