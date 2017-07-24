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

const errNoLocation = "no location sent"

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

type FBWebhookMsg struct {
	Entry []struct {
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Message struct {
				Attachments []struct {
					Type    string `json:"type"`
					Payload struct {
						Coordinates struct {
							Lat  float64 `json:"lat"`
							Long float64 `json:"long"`
						} `json:"coordinates"`
					} `json:"payload"`
				} `json:"attachments"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}

func MessengerRequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Message recieved")
	if r.URL.Query().Get("hub.verify_token") != "" {
		verifyToken(w, r)
		return
	}
	FBUserID, location, err := getUserDetails(r)
	if err != nil {
		if err.Error() == errNoLocation {
			sendText(FBUserID, "Send your location to get some delicious recommendations!")
			return
		}
		log.Println("error getting FB User details: ", err)
		return
	}

	client := NewGooglePlacesClient(Config{})
	places, err := location.GetPlaces(client)
	if err != nil {
		fmt.Println(err)
	}

	for _, place := range places {
		err := place.GetWebsite(client)
		if err != nil {
			fmt.Println(err)
		}
		sendText(FBUserID, fmt.Sprintf("%s, %v\n %s", place.Name, convertToStars(place.Rating), place.Website))
	}
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

func getUserDetails(r *http.Request) (string, *Location, error) {
	var req FBWebhookMsg
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", nil, err
	}
	err = json.Unmarshal(body, &req)
	if err != nil {
		return "", nil, err
	}

	userID := req.Entry[0].Messaging[0].Sender.ID

	message := req.Entry[0].Messaging[0].Message
	if message.Attachments == nil {
		return userID, nil, errors.New(errNoLocation)
	}

	if message.Attachments[0].Type != "location" {
		return userID, nil, errors.New(errNoLocation)
	}

	lat := message.Attachments[0].Payload.Coordinates.Lat
	long := message.Attachments[0].Payload.Coordinates.Long

	location, err := NewLocation(lat, long)
	return userID, location, err
}

func sendText(user, text string) {
	resp := MessengerResponse{
		FBUser: FBUser{
			ID: user,
		},
		Message: FBMessage{
			Text: text,
		},
	}
	err := sendToMessenger(resp)
	if err != nil {
		log.Println("error sending response to messenger: ", err)
	}
	return
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

func convertToStars(rating float64) string {
	var stars string
	for i := 0; i < int(rating); i++ {
		stars += "★"
	}
	if float64(int64(rating)) != rating {
		stars += " ½"
	}
	return stars
}
