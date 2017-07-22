package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func MessengerRequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Message recieved")
	if r.URL.Query().Get("hub.verify_token") != "" {
		verifyToken(w, r)
		return
	}
	sendMessage(r)
}

func verifyToken(w http.ResponseWriter, r *http.Request) {
	fbVerificationToken := os.Getenv("FB_VERIFICATION_TOKEN")
	if r.FormValue("hub.verify_token") == fbVerificationToken {
		io.WriteString(w, r.FormValue("hub.challenge"))
		return
	} else {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, "incorrect verification token")
		return
	}

}

func sendMessage(r *http.Request) {
	FBUserID, err := getFBUserID(r)
	if err != nil {
		log.Println("error getting FB UserID: ", err)
		return
	}

	resp := MessengerResponse{
		FBUser: FBUser{
			ID: FBUserID,
		},
		Message: FBMessage{
			Text: "What's up?",
		},
	}
	err = sendToMessenger(resp)
	if err != nil {
		log.Println("error sending response to messenger: ", err)
	}
	return
}

func getFBUserID(r *http.Request) (string, error) {
	var req FBWebhookMsg
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, &req)
	if err != nil {
		return "", err
	}
	return req.Entry[0].Messaging[0].Sender.ID, nil
}

func sendToMessenger(payload MessengerResponse) error {
	buf, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://graph.facebook.com/v2.8/me/messages?access_token=%s", os.Getenv("FB_PAGE_TOKEN"))
	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewReader(buf))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("error response from Messenger: ", string(body)))
	}
	return nil
}

type FBWebhookMsg struct {
	Entry []struct {
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Message struct {
				Text string `json:"text"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}

type FBUser struct {
	ID string `json:"id"`
}

type FBMessage struct {
	Text string `json:"text"`
}

type MessengerResponse struct {
	FBUser  FBUser    `json:"recipient"`
	Message FBMessage `json:"message"`
}
