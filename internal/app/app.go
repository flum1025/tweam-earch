package app

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/flum1025/tweam-earch/internal/config"
	"github.com/slack-go/slack"
)

type App struct {
	config *config.Config
}

func NewApp(config *config.Config) (*App, error) {
	return &App{
		config: config,
	}, nil
}

type tweet twitter.Tweet

func (t tweet) Match(rules config.Rules) bool {
	for _, rule := range rules {
		r := regexp.MustCompile(rule.Text)

		if r.MatchString(t.Text) && !t.Retweeted && t.InReplyToUserIDStr == "" {
			return true
		}
	}

	return false
}

func (t tweet) Attachment(config config.Slack) slack.Attachment {
	return slack.Attachment{
		Fallback:   fmt.Sprintf("%s by %s", t.Text, t.User.ScreenName),
		Color:      "#4169e1",
		AuthorName: fmt.Sprintf("%s(%s)", t.User.Name, t.User.ScreenName),
		AuthorLink: t.User.URL,
		AuthorIcon: t.User.ProfileImageURLHttps,
		Text:       t.Text,
		Fields: []slack.AttachmentField{
			{
				Value: fmt.Sprintf("<https://twitter.com/%s/status/%d|Tweet>", t.User.ScreenName, t.ID),
			},
		},
		Footer: fmt.Sprintf("%s %s", config.Icon, config.SourceUser),
		Ts:     json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
}

func (a *App) Handle(userID string, tweets []twitter.Tweet) error {
	account := a.config.Accounts.Find(userID)
	if account == nil {
		log.Println(fmt.Sprintf("not target user: %v", userID))
		return nil
	}

	client := slack.New(account.Slack.Token)
	attachments := make([]slack.Attachment, 0)

	for _, t := range tweets {
		ts := tweet(t)
		if ts.Match(account.Rules) {
			attachments = append(attachments, ts.Attachment(account.Slack))
		}
	}

	if len(attachments) > 0 {
		_, _, err := client.PostMessage(
			account.Slack.Channel,
			slack.MsgOptionAttachments(attachments...),
			slack.MsgOptionIconEmoji(account.Slack.UserIcon),
		)

		if err != nil {
			return fmt.Errorf("post message: %w", err)
		}
	}

	return nil
}
