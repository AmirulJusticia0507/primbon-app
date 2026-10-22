package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const griphubChatURL = "https://griphubrouter.web.id/v1/chat/completions"

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

// AIChat forwards a short user prompt to the optional OpenAI-compatible Griphub endpoint.
func AIChat(c *fiber.Ctx, apiKey, model string) error {
	apiKey = strings.TrimSpace(apiKey)
	model = strings.TrimSpace(model)
	if apiKey == "" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Fitur AI belum dikonfigurasi. Atur AI_API_KEY dan AI_MODEL."})
	}
	if model == "" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Model AI belum dikonfigurasi. Atur AI_MODEL."})
	}

	var input struct {
		Message string `json:"message"`
	}
	if err := c.BodyParser(&input); err != nil || strings.TrimSpace(input.Message) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Isi pesan terlebih dahulu."})
	}
	input.Message = strings.TrimSpace(input.Message)
	if len(input.Message) > 2000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Pesan terlalu panjang (maksimal 2000 karakter)."})
	}

	reply, err := requestGriphubChat(c.UserContext(), apiKey, model, input.Message)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Layanan AI tidak dapat dihubungi saat ini. Silakan coba lagi nanti."})
	}
	return c.JSON(fiber.Map{"reply": reply})
}

func requestGriphubChat(parent context.Context, apiKey, model, message string) (string, error) {
	ctx, cancel := context.WithTimeout(parent, 30*time.Second)
	defer cancel()

	payload, err := json.Marshal(chatCompletionRequest{
		Model:    model,
		Messages: []chatMessage{{Role: "user", Content: message}},
	})
	if err != nil {
		return "", fmt.Errorf("encode chat request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, griphubChatURL, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create chat request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send chat request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("chat endpoint returned HTTP %d", resp.StatusCode)
	}

	var result chatCompletionResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&result); err != nil {
		return "", fmt.Errorf("decode chat response: %w", err)
	}
	if len(result.Choices) == 0 || strings.TrimSpace(result.Choices[0].Message.Content) == "" {
		return "", errors.New("chat endpoint returned no message")
	}
	return result.Choices[0].Message.Content, nil
}
