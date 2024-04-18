package main

import (
	"github.com/sirupsen/logrus"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var (
	version = "unkown"
)

func main() {
	app := cli.NewApp()
	app.Name = "gitea-message plugin"
	app.Usage = "gitea-message plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api-key",
			Usage:  "api key to access gitea api",
			EnvVar: "PLUGIN_API_KEY,GITEA_MESSAGE_API_KEY,GITEA_TOKEN",
		},
		cli.StringSliceFlag{
			Name:   "message-text",
			Usage:  "text of message content. If this is set, message-file is going to be ignored.",
			EnvVar: "PLUGIN_MESSAGE_TEXT,GITEA_MESSAGE_FILE",
		},
		cli.StringSliceFlag{
			Name:   "message-file",
			Usage:  "file with content for message",
			EnvVar: "PLUGIN_MESSAGE_FILE,GITEA_MESSAGE_FILE",
		},
		cli.StringFlag{
			Name:   "base-url",
			Usage:  "url of the gitea instance",
			EnvVar: "PLUGIN_BASE_URL,GITEA_MESSAGE_BASE_URL",
		},
		cli.StringFlag{
			Name:   "title",
			Value:  "",
			Usage:  "string for the title shown in the gitea pr comment",
			EnvVar: "PLUGIN_TITLE,GITEA_MESSAGE_TITLE",
		},
		cli.StringFlag{
			Name:  "delete-identifier",
			Value: "",
			Usage: "string that is used as identifier for deletion upon a new message. " +
				"So every previous comment in the PR that has this identifier will be deleted before the message is sent",
			EnvVar: "PLUGIN_DELETE_IDENTIFIER",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "build.event",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.Int64Flag{
			Name:   "pull.request",
			Usage:  "pull request",
			EnvVar: "DRONE_PULL_REQUEST",
		},

		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Pr: Pr{
			Index: c.Int64("pull.request"),
		},
		Build: Build{
			Event: c.String("build.event"),
		},
		Config: Config{
			APIKey:           c.String("api-key"),
			MessageText:      c.String("message-text"),
			MessageFile:      c.String("message-file"),
			BaseURl:          c.String("base-url"),
			Title:            c.String("title"),
			DeleteIdentifier: c.String("delete-identifier"),
		},
	}

	return plugin.Exec()
}
