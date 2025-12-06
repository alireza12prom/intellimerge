# IntelliMerge

> ⚠️ **Note**: This project is currently in active development. Features and APIs may change without notice.

IntelliMerge is an intelligent GitLab merge request automation tool that uses LLM-powered agents to review, analyze, and manage merge requests automatically.

## Features

- 🤖 **AI-Powered Reviews**: Automated code review using LLM agents with optimized prompts for best results
- 🔌 **LLM Flexibility**: Use any LLM provider you want - OpenAI, Anthropic, local models, or any compatible API
- 🔗 **GitLab Integration**: Seamless integration with GitLab via webhooks
- ⚡ **Webhook Support**: Real-time processing of merge request events
- 🎯 **Configurable**: Environment-based configuration for easy deployment

## Prerequisites

- Go 1.25.3 or later
- GitLab instance (GitLab.com or self-hosted)
- GitLab bot token with at least Maintainer role
- LLM API access (configured in the LLM module)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.