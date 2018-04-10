package main

import (
	"fmt"
	"github.com/ashwanthkumar/slack-go-webhook"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/go-playground/webhooks.v3"
	"gopkg.in/go-playground/webhooks.v3/gitlab"
	"strconv"
)

var (
	port       = kingpin.Flag("port", "Bind to this port").Short('p').Default("3016").Int()
	endpoint   = kingpin.Flag("endpoint", "Webhook endpoint to receive events").Default("/webhook").Short('e').String()
	webhookUrl = kingpin.Flag("webhookUrl", "Slack webhook to send events").Required().Short('s').String()
	gitlabUrl  = kingpin.Flag("gitlabUrl", "Your Gitlab URL").Default("https://www.gitlab.com/").String()
)

func main() {
	kingpin.Version("0.0.2")
	kingpin.Parse()

	hook := gitlab.New(&gitlab.Config{})
	hook.RegisterEvents(HandlePipeline, gitlab.PipelineEvents)

	err := webhooks.Run(hook, ":"+strconv.Itoa(*port), *endpoint)
	if err != nil {
		fmt.Println(err)
	}
}

// HandleRelease handles GitHub release events
func HandlePipeline(payload interface{}, header webhooks.Header) {

	fmt.Println("Handling Pipeline event")

	pl := payload.(gitlab.PipelineEventPayload)
	if pl.ObjectAttributes.Status == "failed" {
		attachment := slack.Attachment{}
		attachment.AddField(slack.Field{Title: "Failed commit", Value: fmt.Sprintf("<%s|link>", pl.Commit.URL)}).AddField(slack.Field{Title: "Author", Value: pl.Commit.Author.Name})
		slackPayload := slack.Payload{
			Text:        fmt.Sprintf("Hi! Project: *%s* failed to build! Here's a link to the <%s/%s/pipelines/%d|failed pipeline>",  pl.Project.PathWithNamespace, *gitlabUrl, pl.Project.PathWithNamespace, pl.ObjectAttributes.ID),
			Username:    "Gitlab",
			Attachments: []slack.Attachment{attachment},
		}
		err := slack.Send(*webhookUrl, "", slackPayload)
		if len(err) > 0 {
			fmt.Printf("error: %s\n", err)
		}
	}

}
