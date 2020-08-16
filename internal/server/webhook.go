package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/go-chi/render"
)

type event struct {
	ForUserID         string          `json:"for_user_id"`
	TweetCreateEvents []twitter.Tweet `json:"tweet_create_events"`
}

func (s *server) webhook(w http.ResponseWriter, r *http.Request) {
	var params event

	body, err := parse(r.Body, &params)
	if err != nil {
		log.Println(fmt.Sprintf("[ERROR] failed to parse body: %v", err))
		return
	}

	if s.debug {
		log.Println(fmt.Sprintf("[DEBUG] %s", string(body)))
	}

	err = s.app.Handle(params.ForUserID, params.TweetCreateEvents)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())

		return
	}

	render.PlainText(w, r, "ok")
}
