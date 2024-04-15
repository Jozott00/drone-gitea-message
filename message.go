package main

import (
	"code.gitea.io/sdk/gitea"
)

type messageClient struct {
	Client  *gitea.Client
	Owner   string
	Repo    string
	Index   int64
	Title   string
	Message string
}

func (mc *messageClient) sendMessage() (*gitea.Comment, *gitea.Response, error) {
	opt := gitea.CreateIssueCommentOption{
		Body: mc.Message,
	}

	return mc.Client.CreateIssueComment(
		mc.Owner,
		mc.Repo,
		mc.Index,
		opt,
	)
}
