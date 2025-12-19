package gemini

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	drawdomain "backend/internal/domain/draw"
	"backend/internal/port/llm"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	defaultModelName   = "gemini-2.5-flash"
	maxFormattedLength = 400
	minFormattedLength = 12
)

var rejectionKeywords = []string{"kill", "suicide", "die"}

var newGeminiClient = genai.NewClient

// contentGenerator は gemini.GenerativeModel をテストしやすい形に抽象化したもの。
type contentGenerator interface {
	GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error)
}

// Formatter は Gemini モデルを利用した整形ロジックを提供する。
type Formatter struct {
	generator contentGenerator
	closeFn   func() error
	modelName string
}

// NewFormatter は API キーとモデル名から Formatter を生成する。
func NewFormatter(ctx context.Context, apiKey, modelName string, extraOpts ...option.ClientOption) (*Formatter, error) {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, fmt.Errorf("gemini formatter: API キーが設定されていません")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	opts := append([]option.ClientOption{option.WithAPIKey(apiKey)}, extraOpts...)
	client, err := newGeminiClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", llm.ErrFormatterUnavailable, err)
	}

	resolvedModel := resolveModelName(modelName)
	model := client.GenerativeModel(resolvedModel)
	configured := configureModel(model)

	return &Formatter{
		generator: configured,
		closeFn:   makeCloseFn(client),
		modelName: resolvedModel,
	}, nil
}

// Close は内部で保持している Gemini クライアントをクローズする。
func (f *Formatter) Close() error {
	if f == nil || f.closeFn == nil {
		return nil
	}
	return f.closeFn()
}

// Format は闇投稿を Gemini で整形し、pending 状態の結果を返す。
func (f *Formatter) Format(ctx context.Context, req *llm.FormatRequest) (*llm.FormatResult, error) {
	if err := validateFormatRequest(req); err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if f.generator == nil {
		return nil, fmt.Errorf("%w: gemini formatter: 生成器が初期化されていません", llm.ErrFormatterUnavailable)
	}

	prompt := buildPrompt(string(req.DarkContent))
	resp, err := f.generator.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", llm.ErrFormatterUnavailable, err)
	}

	text, err := extractFirstText(resp)
	if err != nil {
		return nil, err
	}

	return &llm.FormatResult{
		DarkPostID:       req.DarkPostID,
		FormattedContent: drawdomain.FormattedContent(text),
		Status:           drawdomain.StatusPending,
	}, nil
}

// Validate は Gemini の整形結果を検証し、公開可能かどうかを判定する。
func (f *Formatter) Validate(ctx context.Context, result *llm.FormatResult) (*llm.FormatResult, error) {
	if result == nil || result.DarkPostID == "" {
		return nil, llm.ErrInvalidFormat
	}

	trimmed := strings.TrimSpace(string(result.FormattedContent))
	if trimmed == "" {
		result.Status = drawdomain.StatusRejected
		result.ValidationReason = "整形結果が空です"
		return result, llm.ErrInvalidFormat
	}

	if reason, rejected := shouldReject(trimmed); rejected {
		result.Status = drawdomain.StatusRejected
		result.ValidationReason = reason
		return result, llm.ErrContentRejected
	}

	result.Status = drawdomain.StatusVerified
	result.FormattedContent = drawdomain.FormattedContent(trimmed)
	result.ValidationReason = ""

	return result, nil
}

func validateFormatRequest(req *llm.FormatRequest) error {
	if req == nil || req.DarkPostID == "" {
		return llm.ErrInvalidFormat
	}
	if strings.TrimSpace(string(req.DarkContent)) == "" {
		return llm.ErrInvalidFormat
	}
	return nil
}

func resolveModelName(name string) string {
	if strings.TrimSpace(name) == "" {
		return defaultModelName
	}
	return name
}

func buildPrompt(content string) string {
	template := `
あなたは匿名の悩み相談を受け取り、投稿者を肯定しながら穏やかで前向きな 200 字以内の日本語メッセージに整形する編集者です。
- です・ます調で丁寧に書く
- URL や顔文字、箇条書きは禁止
- 余計な前置きは書かず、すぐ本文を書き始める

原文:
%s
`
	return fmt.Sprintf(strings.TrimSpace(template), strings.TrimSpace(content))
}

func extractFirstText(resp *genai.GenerateContentResponse) (string, error) {
	if resp == nil {
		return "", llm.ErrInvalidFormat
	}
	for _, candidate := range resp.Candidates {
		if candidate == nil || candidate.Content == nil {
			continue
		}
		for _, part := range candidate.Content.Parts {
			if part == nil {
				continue
			}
			if text, ok := part.(genai.Text); ok {
				trimmed := strings.TrimSpace(string(text))
				if trimmed != "" {
					return trimmed, nil
				}
			}
		}
	}
	return "", llm.ErrInvalidFormat
}

func shouldReject(text string) (string, bool) {
	length := utf8.RuneCountInString(text)
	if length < minFormattedLength {
		return "整形結果が短すぎます", true
	}
	if length > maxFormattedLength {
		return "整形結果が長すぎます", true
	}

	lower := strings.ToLower(text)
	for _, keyword := range rejectionKeywords {
		if strings.Contains(lower, keyword) {
			return fmt.Sprintf("不適切な語句(%s)が含まれています", keyword), true
		}
	}
	if strings.Contains(lower, "http://") || strings.Contains(lower, "https://") {
		return "URL は含めないでください", true
	}
	return "", false
}

func configureModel(model *genai.GenerativeModel) contentGenerator {
	if model == nil {
		return nil
	}
	model.SetCandidateCount(1)
	model.SetMaxOutputTokens(512)
	model.SetTemperature(0.4)
	return model
}

func makeCloseFn(client *genai.Client) func() error {
	if client == nil {
		return nil
	}
	return client.Close
}
