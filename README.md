# IntelliMerge

> ⚠️ **Note**: This project is currently in active development. Features and APIs may change without notice.

IntelliMerge is an intelligent GitLab merge request automation tool that uses LLM-powered agents to review, analyze, and manage merge requests automatically.

## Features

- 🤖 **AI-Powered Reviews**: Automated code review using LLM agents with optimized prompts for best results
- 🔌 **LLM Flexibility**: Use any LLM provider you want - OpenAI, Anthropic, local models, or any compatible API
- 📝 **Slash Commands**: Use `/summary` command in GitLab merge requests to automatically fetch Jira task summaries based on code changes
- 🔗 **GitLab Integration**: Seamless integration with GitLab via webhooks
- ⚡ **Webhook Support**: Real-time processing of merge request events
- 🎯 **Configurable**: Environment-based configuration for easy deployment

## How It Works

### `/summary` Command

The `/summary` command allows you to automatically generate summaries of Jira tasks based on the code changes in your merge request. Simply comment `/summary` in a GitLab merge request, and IntelliMerge will:

1. Analyze the code changes in the merge request
2. Fetch the associated Jira task
3. Generate an intelligent summary of the task based on the actual changes made
4. Post the summary as a comment in the merge request

This helps keep your Jira tasks up-to-date with the actual implementation, ensuring accurate documentation and better project tracking.

## Prerequisites

- Go 1.25.3 or later
- GitLab instance (GitLab.com or self-hosted)
- GitLab bot token with at least Maintainer role
- LLM API access (configured in the LLM module)
- Jira integration (for `/summary` command functionality)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.