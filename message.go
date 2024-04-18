package main

import (
	"code.gitea.io/sdk/gitea"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type messageClient struct {
	Client           *gitea.Client
	Owner            string
	Repo             string
	Index            int64
	Message          string
	DeleteIdentifier string
}

func (mc *messageClient) sendMessage() (*gitea.Comment, *gitea.Response, error) {
	var body = mc.Message
	if mc.DeleteIdentifier != "" {
		body = mc._getDeleteIdentifierMd() + "\n" + body
	}

	opt := gitea.CreateIssueCommentOption{
		Body: body,
	}

	return mc.Client.CreateIssueComment(
		mc.Owner,
		mc.Repo,
		mc.Index,
		opt,
	)
}

// Deletes all comments in the PR that include the DeleteIdentifier.
// If the DeleteIdentifier is "", the search will not be performed.
//
// Returns the number of deleted comments or an error.
func (mc *messageClient) deletePreviousMessages() (int, error) {
	if mc.DeleteIdentifier == "" {
		log.Info("No DeleteIdentifier specified... skipping comment deletion")
		return 0, nil
	}

	identifier := mc._getDeleteIdentifierMd()
	log.WithField("deleteIdentifier", mc.DeleteIdentifier).Info("Start deletion of PR comments")
	matchingComments, err := mc._searchCommentsForIdentifier(identifier)
	log.Info("Found comments with identifier: ", len(matchingComments))

	if err != nil {
		return 0, fmt.Errorf("failed to search for comments. %s", err)
	}

	for i, comment := range matchingComments {
		_, err := mc.Client.DeleteIssueComment(mc.Owner, mc.Repo, comment.ID)
		if err != nil {
			return i, fmt.Errorf("error deleting comment: %w", err)
		}
		log.WithField("id", comment.ID).Info("Deleted comment")
	}

	log.WithField("numberOfComments", len(matchingComments)).Info("Deletion of old comment(s) completed")
	return len(matchingComments), nil
}

// searchCommentsForIdentifier searches for comments containing a specific identifier within a pull request
func (mc *messageClient) _searchCommentsForIdentifier(identifier string) ([]*gitea.Comment, error) {
	// Fetch all comments on the specified pull request
	comments, _, err := mc.Client.ListIssueComments(mc.Owner, mc.Repo, mc.Index, gitea.ListIssueCommentOptions{})
	if err != nil {
		return nil, fmt.Errorf("error fetching comments: %w", err)
	}

	// Filter comments to find those that contain the identifier
	var matchingComments []*gitea.Comment
	for _, comment := range comments {
		if strings.Contains(comment.Body, identifier) {
			matchingComments = append(matchingComments, comment)
		}
	}

	return matchingComments, nil
}

func (mc *messageClient) _getDeleteIdentifierMd() string {
	return "<!-- delete-identifier=\"" + mc.DeleteIdentifier + "\" -->"
}
