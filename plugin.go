package main

import (
	"code.gitea.io/sdk/gitea"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type (
	Repo struct {
		Owner string
		Name  string
	}
	Pr struct {
		Index int64
	}
	Build struct {
		Event string
	}
	Config struct {
		APIKey      string
		MessageFile string
		BaseURl     string
		Title       string
	}
	Plugin struct {
		Repo   Repo
		Pr     Pr
		Build  Build
		Config Config
	}
)

func (p Plugin) Exec() error {

	if p.Build.Event != "pull_request" {
		return fmt.Errorf("only pull request may trigger this plugin. (was %s)", p.Build.Event)
	}

	if p.Config.APIKey == "" {
		return fmt.Errorf("you must provide an API key")
	}

	if p.Config.BaseURl == "" {
		return fmt.Errorf("you must provide the repo's base url")
	}

	if !strings.HasSuffix(p.Config.BaseURl, "/") {
		p.Config.BaseURl = p.Config.BaseURl + "/"
	}

	if p.Pr.Index == 0 {
		return fmt.Errorf("pull request number is not set")
	}

	glob, err := filepath.Glob(p.Config.MessageFile)
	if err != nil {
		return fmt.Errorf("failed to glob %s. %s", p.Config.MessageFile, err)
	}

	content, err := os.ReadFile(glob[0])
	if err != nil {
		return fmt.Errorf("failed to read the file %s. %s", glob[0], err)
	}

	httpClient := &http.Client{}
	client, err := gitea.NewClient(p.Config.BaseURl, gitea.SetToken(p.Config.APIKey), gitea.SetHTTPClient(httpClient))

	mc := messageClient{
		Client:  client,
		Owner:   p.Repo.Owner,
		Repo:    p.Repo.Name,
		Index:   p.Pr.Index,
		Title:   p.Config.Title,
		Message: string(content),
	}

	_, _, err = mc.sendMessage()

	return err
}
