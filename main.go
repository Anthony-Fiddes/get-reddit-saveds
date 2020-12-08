package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	tokenAPIURL       = "https://www.reddit.com/api/v1/access_token"
	savedsAPIURL      = "https://oauth.reddit.com/user/%s/saved"
	userAgentTemplate = "desktop:saveds-downloader:v0.0.1 (by /u/%s)"
	settings          = "settings.json"
	output            = "output.json"
	limit             = "1000"
)

var (
	client    = &http.Client{Timeout: time.Second * 10}
	userAgent = ""
)

type apiLogin struct {
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
}

func must(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func requestToken(config apiLogin) string {
	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("username", config.Username)
	form.Set("password", config.Password)
	formStr := form.Encode()
	// ? how are these errors handled in prod?
	request, err := http.NewRequest("POST", tokenAPIURL, strings.NewReader(formStr))
	must(err)
	request.SetBasicAuth(config.ClientID, config.ClientSecret)
	request.Header.Set("User-Agent", userAgent)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", fmt.Sprintf("%d", len(formStr)))
	response, err := client.Do(request)
	must(err)
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Panicf(`Error "%s" received from Reddit`, response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	must(err)
	result := struct {
		Error       string
		AccessToken string `json:"access_token"`
	}{}
	err = json.Unmarshal(body, &result)
	must(err)
	if result.Error != "" {
		serverErr := fmt.Errorf(`Server returned error: "%s"`, result.Error)
		log.Print(serverErr)
		log.Print("Please check your login details.")
		panic(serverErr)
	}
	return result.AccessToken
}

func main() {
	data, err := ioutil.ReadFile(settings)
	login := apiLogin{}
	if err == nil {
		jsonErr := json.Unmarshal(data, &login)
		if jsonErr != nil {
			log.Print(jsonErr)
			log.Fatal("Your settings.json file may be malformatted")
		}
	} else if !os.IsNotExist(err) {
		log.Printf(`Was unable to read configuration at "%s", proceeding to manual entry.`, settings)
	}
	setFieldIfEmpty := func(field *string, name string) {
		if *field == "" {
			fmt.Printf("%s: ", name)
			fmt.Scanln(field)
		}
	}
	setFieldIfEmpty(&login.Username, "Username")
	setFieldIfEmpty(&login.Password, "Password")
	setFieldIfEmpty(&login.ClientID, "Client ID")
	setFieldIfEmpty(&login.ClientSecret, "Client Secret")
	userAgent = fmt.Sprintf(userAgentTemplate, login.Username)

	token := requestToken(login)
	fmt.Printf("Token: %s\n", token)
	savedsGet, err := url.Parse(fmt.Sprintf(savedsAPIURL, login.Username))
	must(err)
	params := savedsGet.Query()
	params.Set("limit", limit)
	savedsGet.RawQuery = params.Encode()
	request, err := http.NewRequest(
		"GET",
		savedsGet.String(),
		nil,
	)
	must(err)
	request.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	request.Header.Set("User-Agent", userAgent)
	response, err := client.Do(request)
	must(err)
	defer response.Body.Close()
	outFile, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0744)
	must(err)
	_, err = io.Copy(outFile, response.Body)
	must(err)
	fmt.Printf("Done. Check %s to find (at most) %s of your most recent saved posts.\n", output, limit)
}
