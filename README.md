# drone-gitea-messenger

Drone plugin to send message as comments to Gitea pull requests. 
All releases are available on [Docker Hub](https://hub.docker.com/r/jozott/drone-gitea-message). 

## Drone YAML Usage

The plugin support sending text as well as some file's content.
By setting the `delete_identifier`, all older comments with the same identifier will be
deleted before sending the message.

### Example of sending a text
```yaml
steps:
  - name: send text to pr
    image: jozott/drone-gitea-message:v0.2.0
    settings:
      api_key:
        from_secret: gitea_token
      base_url: http://gitea.example.com
      message_text: "Hello world"
```

### Example of sending content of file
```yaml
steps:
  - name: send file content to pr
    image: jozott/drone-gitea-message:v0.2.0
    settings:
      api_key:
        from_secret: gitea_token
      base_url: http://gitea.example.com
      message_file: path/to/hello_world.md
```

### Example with `delete_identifier`
If the `delete_identifier` is set, the plugin will delete all existing PR comments 
that contain the delete identifier. This is handy if sending status updates that make
older status updates obsolete.
```yaml
steps:
  - name: send test report to pr
    image: jozott/drone-gitea-message:v0.2.0
    settings:
      api_key:
        from_secret: gitea_token
      base_url: http://gitea.example.com
      message_file: build/test/report.md
      delete_identifier: test-report-delete-id
```

## Build
To build the binary execute the following commands
```bash
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/linux/amd64/drone-gitea-message
```

## Run
To get all available commands run
```
./drone-gitea-message --help
```

