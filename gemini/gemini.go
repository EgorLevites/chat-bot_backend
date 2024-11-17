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

// Load environment variables from the .env file
// func init() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatalf("Error loading .env file: %v", err)
// 	}
// }
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}
}

// Create a client for Gemini API
func GetGeminiClient() (*genai.Client, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API_KEY not set in environment variables")
	}

	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating Gemini client: %v", err)
	}
	return client, nil
}

// GenerateResponse sends the message to Gemini API and returns a response.
func GenerateResponse(message string) (string, error) {
    client, err := GetGeminiClient()
    if err != nil {
        return "", fmt.Errorf("failed to get Gemini client: %w", err)
    }
    defer client.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    model := client.GenerativeModel("gemini-1.5-flash-latest")
    resp, err := model.GenerateContent(ctx, genai.Text(message))

    if err != nil {
        log.Printf("Primary model error: %v. Switching to fallback model.", err)
        model = client.GenerativeModel("gemini-1.5-pro")
        resp, err = model.GenerateContent(ctx, genai.Text(message))
        if err != nil {
            return "", fmt.Errorf("error generating response from fallback model: %w", err)
        }
    }

    // Extract the response text
	// Check the type and extract the content
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
    	rawContent := fmt.Sprintf("%v", resp.Candidates[0].Content) // Convert the content to a string

    	// Remove unnecessary characters using a regular expression
    	cleanContent := cleanResponse(rawContent)
    	return cleanContent, nil
	}



    return "", fmt.Errorf("no response text found in the generated content")
}

// cleanResponse removes unnecessary characters from the response
func cleanResponse(rawContent string) string {
    // Regular expression to remove unnecessary characters
    re := regexp.MustCompile(`[&{}\[\]]+|model`)
    cleanContent := re.ReplaceAllString(rawContent, "")

    // Remove extra spaces
    return strings.TrimSpace(cleanContent)
}


