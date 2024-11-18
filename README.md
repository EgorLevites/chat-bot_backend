# Chat-Bot Backend

This repository contains the backend implementation for a chat-bot application, developed in Go. 
The backend establishes WebSocket connections with clients and integrates with the Gemini API to generate responses.

## Project Structure

- `gemini/gemini.go`: Handles interactions with the Gemini API, including client creation and response generation.
- `handlers/websocket.go`: Manages WebSocket connections, facilitating communication between clients and the Gemini API.
- `models/message.go`: Defines the `Message` struct used for message exchange.
- `main.go`: The entry point of the application, setting up routes and initiating the server.
- `.env`: Contains environment variables, notably the `API_KEY` for the Gemini API.

## Prerequisites

- [Go](https://golang.org/dl/) installed on your system.
- [Docker](https://www.docker.com/get-started) installed for containerization.
- A valid API key for the Gemini API.

## Setup Instructions

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/yourusername/chat-bot-backend.git
   cd chat-bot-backend
   ```

2. **Configure Environment Variables**:

   Create a `.env` file in the root directory with the following content:

   ```
   API_KEY=your_gemini_api_key_here
   ```

   Alternatively, set the `API_KEY` environment variable in your system.

## Running the Application

### Using Docker

1. **Build the Docker Image**:

   ```bash
   docker build -t chat-bot-backend .
   ```

2. **Run the Docker Container**:

   ```bash
   docker run -p 8080:8080 --env-file .env chat-bot-backend
   ```

   This command maps port 8080 of the container to port 8080 on your host machine and provides the necessary environment variables.

### Running Locally

1. **Install Dependencies**:

   ```bash
   go mod download
   ```

2. **Build the Application**:

   ```bash
   go build -o main .
   ```

3. **Run the Application**:

   ```bash
   ./main
   ```

   The server will start on port 8080.

## API Endpoints

- `/`: Serves static files from the `../chat-bot-frontend` directory.
- `/ws`: Handles WebSocket connections for real-time communication.

## Environment Variables

- `API_KEY`: Your Gemini API key. Ensure this is set in the `.env` file or your system's environment variables.

## Dependencies

- [Google Generative AI Go Client](https://pkg.go.dev/github.com/google/generative-ai-go/genai)
- [Gorilla WebSocket](https://pkg.go.dev/github.com/gorilla/websocket)
- [godotenv](https://pkg.go.dev/github.com/joho/godotenv)

## Notes

- Ensure the `API_KEY` is kept confidential and not shared publicly.
- The application includes CORS middleware to allow requests from any origin. Adjust this in `main.go` as needed for production environments.

For further information or assistance, please refer to the official documentation of the respective dependencies.
