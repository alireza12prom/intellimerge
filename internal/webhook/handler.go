package webhook

import (
	"crypto/hmac"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/alireza12prom/intellimerge/internal/commands"
	"github.com/alireza12prom/intellimerge/internal/config"
	"github.com/alireza12prom/intellimerge/internal/gitlab"
	"github.com/alireza12prom/intellimerge/internal/jira"
	"github.com/alireza12prom/intellimerge/internal/llm"
	"github.com/alireza12prom/intellimerge/internal/models"
)

type Handler struct {
	config *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		config: cfg,
	}
}

func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify webhook secret
	if !h.verifyWebhookSecret(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var webhook models.GitLabCommentEvent
	if err := json.Unmarshal(body, &webhook); err != nil {
		http.Error(w, "Failed to parse webhook", http.StatusBadRequest)
		return
	}

	if webhook.ObjectKind == "note" && webhook.ObjectAttributes.NoteableType == "MergeRequest" {
		note := strings.TrimSpace(webhook.ObjectAttributes.Note)

		if strings.HasPrefix(note, "/summary") {
			gitlabClient := gitlab.NewClient(h.config.GitLabURL, h.config.GitLabToken)
			jiraClient := jira.NewClient(h.config.JiraURL, h.config.JiraEmail, h.config.JiraAPIToken)
			llmClient := llm.NewClient("openai", h.config.LLMAPIKey)
			summaryCmd := commands.NewSummaryCommand(&webhook, gitlabClient, jiraClient, llmClient)

			go func() {
				if err := summaryCmd.Execute(); err != nil {
					fmt.Printf("Error handling summary command: %v\n", err)
				}
			}()

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *Handler) verifyWebhookSecret(r *http.Request) bool {
	token := r.Header.Get("X-Gitlab-Token")
	if token == "" {
		return false
	}
	return hmac.Equal([]byte(token), []byte(h.config.WebhookSecret))
}
