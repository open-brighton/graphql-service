package graph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/openbrighton/graphql-service/graph/model"
)

const anthropicAPIURL = "https://api.anthropic.com/v1/messages"
const anthropicVersion = "2023-06-01"
const claudeModel = "claude-3-haiku-20240307"

const systemPrompt = `You are a helpful assistant for OpenBrighton, a community organisation based in Brighton, UK. 
OpenBrighton runs events, workshops, and initiatives that bring together the local tech, creative, and civic communities.
Answer questions about OpenBrighton, local events, and community activities in a friendly and informative way.
If you don't know something specific about OpenBrighton, be honest and suggest the user check the website or get in touch directly.`

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type anthropicResponse struct {
	Content []anthropicContentBlock `json:"content"`
	Error   *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func CallClaude(ctx context.Context, messages []*model.ChatMessageInput) (string, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY is not set")
	}

	anthropicMessages := make([]anthropicMessage, 0, len(messages))
	for _, m := range messages {
		if m == nil {
			continue
		}
		anthropicMessages = append(anthropicMessages, anthropicMessage{
			Role:    string(m.Role),
			Content: m.Content,
		})
	}

	reqBody := anthropicRequest{
		Model:     claudeModel,
		MaxTokens: 1024,
		System:    systemPrompt,
		Messages:  anthropicMessages,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, anthropicAPIURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Anthropic API: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var anthropicResp anthropicResponse
	if err := json.Unmarshal(respBytes, &anthropicResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if anthropicResp.Error != nil {
		return "", fmt.Errorf("Anthropic API error (%s): %s", anthropicResp.Error.Type, anthropicResp.Error.Message)
	}

	if len(anthropicResp.Content) == 0 {
		return "", fmt.Errorf("empty response from Anthropic API")
	}

	return anthropicResp.Content[0].Text, nil
}
