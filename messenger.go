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

type MessengerResponse struct {
	FBUser  FBUser    `json:"recipient"`
	Message FBMessage `json:"message"`
}

type FBUser struct {
	ID string `json:"id"`
}

type FBMessage struct {
	Text        string         `json:"text,omitempty"`
	Attachment  *FBAttachment  `json:"attachment,omitempty"`
	Attachments []FBAttachment `json:"attachments,omitempty"`
}

type FBAttachment struct {
	Type    string    `json:"type,omitempty"`
	Payload FBPayload `json:"payload,omitempty"`
}

type FBPayload struct {
	TemplateType string             `json:"template_type,omitempty"`
	Coordinates  *FBCoordinates     `json:"coordinates,omitempty"`
	Elements     []FBPayloadElement `json:"elements,omitempty"`
}

type FBCoordinates struct {
	Lat  float64 `json:"lat,omitempty"`
	Long float64 `json:"long,omitempty"`
}

type FBPayloadElement struct {
	Title         string          `json:"title,omitempty"`
	ImageUrl      string          `json:"image_url,omitempty"`
	DefaultAction FBDefaultAction `json:"default_action,omitempty"`
}

type FBDefaultAction struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type FBWebhookMsg struct {
	Entry []struct {
		Messaging []struct {
			Sender  FBUser    `json:"sender"`
			Message FBMessage `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}

type ThreadSetting struct {
	SettingType   string         `json:"setting_type"`
	ThreadState   string         `json:"thread_state"`
	CallToActions []CallToAction `json:"call_to_actions"`
}

type CallToAction struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Payload string `json:"payload"`
	Url     string `json:"url"`
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
	googleRecommendations, err := location.GetPlacesFromGoogle(client)
	if err != nil {
		log.Fatal("error accessing db: ", err)
		return
	}
	curatedRecommendations, err := GetPlacesFromDB(DB, *location)
	if err != nil {
		fmt.Println(err)
	}

	if len(curatedRecommendations) != 0 {
		sendText(FBUserID, "I've been researching this area! I recommend...")
		sendPlaces(curatedRecommendations, client, FBUserID)
		return
	}
	sendText(FBUserID, "I don't have any recommendations in this area, but this is what turns up on Google...")
	sendPlaces(googleRecommendations, client, FBUserID)
}

func sendPlaces(places []Place, client GooglePlacesClient, FBUserID string) {
	for _, place := range places {
		err := place.GetDetails(client)
		if err != nil {
			fmt.Println(err)
			return
		}
		sendLocation(FBUserID, place)
		sendText(FBUserID, fmt.Sprintf("%v\n%s", convertToStars(place.Rating), place.Website))
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
	message := FBMessage{
		Text: text,
	}
	err := sendToMessenger(user, message)
	if err != nil {
		log.Println("error sending text response to messenger: ", err)
	}
	return
}

func sendLocation(user string, p Place) {
	attachment := FBAttachment{
		Type: "template",
		Payload: FBPayload{
			TemplateType: "generic",
			Elements: []FBPayloadElement{
				{
					Title: p.Name,
					DefaultAction: FBDefaultAction{
						Type: "web_url",
						Url:  p.LinkMapUrl(),
					},
					ImageUrl: p.StaticMapUrl(),
				},
			},
		},
	}

	message := FBMessage{
		Attachment: &attachment,
	}
	err := sendToMessenger(user, message)
	if err != nil {
		log.Println("error sending location response to messenger: ", err)
	}
	return
}

func sendToMessenger(user string, message FBMessage) error {
	payload := MessengerResponse{
		FBUser: FBUser{
			ID: user,
		},
		Message: message,
	}
	buf, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://graph.facebook.com/v2.8/me/messages?access_token=%s", os.Getenv("FB_PAGE_TOKEN"))
	return post(url, buf)
}

func setPersistentMenu() error {
	threadSetting := ThreadSetting{
		SettingType: "call_to_actions",
		ThreadState: "existing_thread",
		CallToActions: []CallToAction{
			CallToAction{
				Type:  "web_url",
				Title: "Make a recommendation",
				Url:   "https://nicolaa.typeform.com/to/noVPUi",
			},
		},
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/thread_settings?access_token=%s", os.Getenv("FB_PAGE_ID"), os.Getenv("FB_PAGE_TOKEN"))

	payload, err := json.Marshal(threadSetting)
	if err != nil {
		return err
	}
	return post(url, payload)
}

func post(url string, payload []byte) error {
	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewReader(payload))
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
