package llm

type Provider interface {
	GenerateMergeRequestSummary(jiraTaskSummary, jiraTaskDescription, codeChanges string) (string, error)
}
