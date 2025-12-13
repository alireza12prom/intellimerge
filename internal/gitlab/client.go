package gitlab

import (
	"fmt"
	"strings"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Client struct {
	client *gitlab.Client
}

type MergeRequest struct {
	ID           int
	IID          int
	ProjectID    int
	Title        string
	Description  string
	SourceBranch string
	TargetBranch string
}

type Diff struct {
	Diff        string
	NewPath     string
	OldPath     string
	NewFile     bool
	DeletedFile bool
}

func NewClient(url, token string) *Client {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		panic(fmt.Sprintf("Failed to create GitLab client: %v", err))
	}

	return &Client{
		client: client,
	}
}

func (c *Client) GetMergeRequest(projectID, mergeRequestIID int) (*MergeRequest, error) {
	mr, _, err := c.client.MergeRequests.GetMergeRequest(int64(projectID), int64(mergeRequestIID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get merge request: %w", err)
	}

	return &MergeRequest{
		ID:           int(mr.ID),
		IID:          int(mr.IID),
		ProjectID:    int(mr.ProjectID),
		Title:        mr.Title,
		Description:  mr.Description,
		SourceBranch: mr.SourceBranch,
		TargetBranch: mr.TargetBranch,
	}, nil
}

func (c *Client) GetMergeRequestDiffs(projectID, mergeRequestIID int) ([]Diff, error) {
	diffs, _, err := c.client.MergeRequests.ListMergeRequestDiffs(int64(projectID), int64(mergeRequestIID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get merge request diffs: %w", err)
	}

	result := make([]Diff, 0, len(diffs))
	for _, diff := range diffs {
		result = append(result, Diff{
			Diff:        diff.Diff,
			NewPath:     diff.NewPath,
			OldPath:     diff.OldPath,
			NewFile:     diff.NewFile,
			DeletedFile: diff.DeletedFile,
		})
	}

	return result, nil
}

func (c *Client) CreateMergeRequestNote(projectID, mergeRequestIID int, body string) error {
	opt := &gitlab.CreateMergeRequestNoteOptions{
		Body: &body,
	}

	_, _, err := c.client.Notes.CreateMergeRequestNote(int64(projectID), int64(mergeRequestIID), opt)
	if err != nil {
		return fmt.Errorf("failed to create merge request note: %w", err)
	}

	return nil
}

func FormatDiffsAsString(diffs []Diff) string {
	var sb strings.Builder
	for _, diff := range diffs {
		if diff.NewFile {
			sb.WriteString(fmt.Sprintf("**New file:** %s\n", diff.NewPath))
		} else if diff.DeletedFile {
			sb.WriteString(fmt.Sprintf("**Deleted file:** %s\n", diff.OldPath))
		} else {
			sb.WriteString(fmt.Sprintf("**Modified file:** %s\n", diff.NewPath))
		}
		sb.WriteString("```diff\n")
		sb.WriteString(diff.Diff)
		sb.WriteString("\n```\n\n")
	}
	return sb.String()
}
