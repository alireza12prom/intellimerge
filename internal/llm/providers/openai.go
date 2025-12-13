package providers

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"text/template"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
)

//go:embed templates/*.txt
var templateFS embed.FS

type OpenAIProvider struct {
	model    string
	client   openai.Client
	template *template.Template
}

type TemplateData struct {
	JiraTaskSummary     string
	JiraTaskDescription string
	CodeChanges         string
}

func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}

	client := openai.NewClient(opts...)

	// Load template
	tmpl, err := templateFS.ReadFile("templates/merge_request_summary.txt")
	if err != nil {
		panic(fmt.Sprintf("Failed to load template: %v", err))
	}

	t, err := template.New("merge_request_summary").Parse(string(tmpl))
	if err != nil {
		panic(fmt.Sprintf("Failed to parse template: %v", err))
	}

	modelName := model
	if modelName == "" {
		modelName = "gpt-4o"
	}

	return &OpenAIProvider{
		model:    modelName,
		client:   client,
		template: t,
	}
}

func (p *OpenAIProvider) GenerateMergeRequestSummary(jiraTaskSummary, jiraTaskDescription, codeChanges string) (string, error) {
	// Prepare template data
	data := TemplateData{
		JiraTaskSummary:     jiraTaskSummary,
		JiraTaskDescription: jiraTaskDescription,
		CodeChanges:         codeChanges,
	}

	// Execute template
	var promptBuf bytes.Buffer
	if err := p.template.Execute(&promptBuf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	prompt := promptBuf.String()

	// Split prompt into system and user messages
	lines := bytes.Split([]byte(prompt), []byte("\n\n"))
	if len(lines) < 2 {
		return "", fmt.Errorf("template must contain system and user prompts separated by blank line")
	}

	systemPrompt := string(lines[0])
	userPrompt := string(bytes.Join(lines[1:], []byte("\n\n")))

	return p.chat(systemPrompt, userPrompt)
}

func (p *OpenAIProvider) chat(systemPrompt, userPrompt string) (string, error) {
	ctx := context.Background()

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
		openai.UserMessage(userPrompt),
	}

	req := openai.ChatCompletionNewParams{
		Model:       p.model,
		Messages:    messages,
		MaxTokens:   param.NewOpt(int64(2000)),
		Temperature: param.NewOpt(0.7),
	}

	completion, err := p.client.Chat.Completions.New(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	content := completion.Choices[0].Message.Content
	if content == "" {
		return "", fmt.Errorf("empty content in response")
	}

	return content, nil
}
