package main

import (
	"io"
	"net/http"
	"os"
)

func MessengerRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("hub.verify_token") != "" {
		verifyToken(r, w)
		return
	}
}

func verifyToken(w http.ResponseWriter, r *http.Request) {
	fbVerificationToken = os.Getenv("FB_VERIFICATION_TOKEN")
	if r.FormValue("hub.verify_token") == fbVerificationToken() {
		io.WriteString(w, r.FormValue("hub.challenge"))
		return
	} else {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, "incorrect verification token")
		return
	}

}
