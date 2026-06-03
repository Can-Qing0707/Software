package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"training_eval_system/config"
)

type LLMService struct {
	client *http.Client
}

func NewLLMService() *LLMService {
	return &LLMService{
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (s *LLMService) Chat(messages []ChatMessage) (string, error) {
	cfg := config.AppConfig.LLM
	reqBody := ChatRequest{
		Model:       cfg.Model,
		Messages:    messages,
		MaxTokens:   cfg.MaxTokens,
		Temperature: cfg.Temperature,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	apiURL := strings.TrimRight(cfg.APIURL, "/") + "/chat/completions"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if chatResp.Error.Message != "" {
		return "", fmt.Errorf("llm error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from llm")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func (s *LLMService) AnalyzeContent(content, requirements string) (string, error) {
	prompt := fmt.Sprintf(`你是一个软件实训教学评价助手。请对学生的实训提交内容进行分析。

任务要求：
%s

学生提交内容：
%s

请从以下几个方面进行分析：
1. 步骤完整性：学生是否完成了所有要求的步骤
2. 逻辑漏洞：是否存在逻辑错误或遗漏
3. 要求匹配度：提交内容与任务要求的匹配程度

请以JSON格式返回分析结果。`, requirements, content)

	messages := []ChatMessage{
		{Role: "system", Content: "你是一个严谨的软件工程实训评价助手。"},
		{Role: "user", Content: prompt},
	}
	return s.Chat(messages)
}

func (s *LLMService) ScoreContent(content, requirements string, indicators []string) (string, error) {
	indicatorStr := ""
	for _, ind := range indicators {
		indicatorStr += "- " + ind + "\n"
	}

	prompt := fmt.Sprintf(`你是一个软件实训教学评分助手。请对学生的实训提交内容进行评分。

任务要求：
%s

学生提交内容：
%s

评价指标：
%s

请按每个指标给出0-100的分数和评语。严格按以下JSON数组格式返回，不要包含任何其他文字：
[
  {"indicator_name": "指标名称", "score": 85.0, "comment": "评语"},
  ...
]`, requirements, content, indicatorStr)

	messages := []ChatMessage{
		{Role: "system", Content: "你是一个严谨的软件工程实训评分助手。"},
		{Role: "user", Content: prompt},
	}
	return s.Chat(messages)
}
