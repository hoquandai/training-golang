package main

import (
	"net/http"
	"fmt"
	"os"
	"encoding/json"
)

const clientID = "d42039fc6522975a6d58"
const clientSecret = "514bd3fad8a4a9ba659c72de6b6d129f5b992470"

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	httpClient := http.Client{}

	http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		code := r.FormValue("code")

		reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
		req, err := http.NewRequest(http.MethodPost, reqURL, nil)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		req.Header.Set("Accept", "application/json")

		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		var t OauthAccessResponse
		if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", "/welcome.html?access_token="+t.AccessToken)
		w.WriteHeader(http.StatusFound)
	})

	http.ListenAndServe(":8080", nil)
}

type OauthAccessResponse struct {
	AccessToken string `json:"access_token"`
}
