package commands

import (
	"fmt"
	"strings"

	"github.com/alireza12prom/intellimerge/internal/gitlab"
	"github.com/alireza12prom/intellimerge/internal/jira"
	"github.com/alireza12prom/intellimerge/internal/llm"
	"github.com/alireza12prom/intellimerge/internal/models"
)

type SummaryCommand struct {
	webhook      *models.GitLabCommentEvent
	gitlabClient *gitlab.Client
	jiraClient   *jira.Client
	llmClient    llm.Provider
}

func NewSummaryCommand(
	webhook *models.GitLabCommentEvent,
	gitlabClient *gitlab.Client,
	jiraClient *jira.Client,
	llmClient llm.Provider,
) *SummaryCommand {
	return &SummaryCommand{
		webhook:      webhook,
		gitlabClient: gitlabClient,
		jiraClient:   jiraClient,
		llmClient:    llmClient,
	}
}

func (c *SummaryCommand) Execute() error {
	projectID := c.webhook.Project.ID
	mergeRequestIID := c.webhook.ObjectAttributes.MergeRequest.IID
	sourceBranch := c.webhook.ObjectAttributes.MergeRequest.SourceBranch

	jiraKey := jira.ExtractJiraKey(sourceBranch)
	issue, err := c.jiraClient.GetIssue(jiraKey)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Failed to fetch Jira issue %s: %v", jiraKey, err)
		return c.gitlabClient.CreateMergeRequestNote(projectID, mergeRequestIID, errorMsg)
	}

	diffs, err := c.gitlabClient.GetMergeRequestDiffs(projectID, mergeRequestIID)
	if err != nil {
		return fmt.Errorf("failed to get merge request diffs: %w", err)
	}

	summary, err := c.llmClient.GenerateMergeRequestSummary(
		issue.Fields.Summary,
		issue.Fields.Description,
		gitlab.FormatDiffsAsString(diffs),
	)
	if err != nil {
		summary = "خطایی در تولید خلاصه رخ داده است. لطفاً دوباره تلاش کنید."
	}

	summaryComment := c.formatSummaryComment(jiraKey, issue, summary, diffs)
	return c.gitlabClient.CreateMergeRequestNote(projectID, mergeRequestIID, summaryComment)
}

func (c *SummaryCommand) formatSummaryComment(jiraKey string, issue *jira.Issue, summary string, diffs []gitlab.Diff) string {
	var sb strings.Builder
	sb.WriteString("## 📌 خلاصه درخواست:\n\n")

	// Metadata section
	sb.WriteString("### 📊 اطلاعات کلی\n\n")
	sb.WriteString(fmt.Sprintf("**وضعیت تسک:** %s\n", issue.Fields.Status.Name))
	sb.WriteString(fmt.Sprintf("**لینک تسک:** [%s](%s)\n", jiraKey, jiraKey))

	// Calculate files changed
	filesChanged := len(diffs)
	sb.WriteString(fmt.Sprintf("**تعداد فایل‌های تغییر یافته:** %d\n", filesChanged))

	// Calculate volume changed
	volumeChanged := c.calculateVolumeChanged(diffs)
	sb.WriteString(fmt.Sprintf("**حجم تغییرات:** %s\n", volumeChanged))

	sb.WriteString("\n---\n\n")
	sb.WriteString("### 📋 تغییرات\n\n")
	sb.WriteString(summary)
	sb.WriteString("\n\n---\n\n")
	sb.WriteString("*✨ تولید شده توسط IntelliMerge*")
	return sb.String()
}

func (c *SummaryCommand) calculateVolumeChanged(diffs []gitlab.Diff) string {
	totalLines := 0
	for _, diff := range diffs {
		lines := strings.Count(diff.Diff, "\n")
		totalLines += lines
	}

	switch {
	case totalLines < 100:
		return "🟢 کوچک (Small)"
	case totalLines < 500:
		return "🟡 متوسط (Medium)"
	default:
		return "🔴 بزرگ (Big)"
	}
}
