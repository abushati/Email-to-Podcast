package email_client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "./email_client/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

type Email_client struct {
	srv  *gmail.Service
	user string
}

func MakeClient(user string) *Email_client {
	b, err := ioutil.ReadFile("./email_client/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	eml := &Email_client{srv, user}

	return eml
}

type Email struct {
	id   string
	body string
}

type queriedEmails struct {
	emails []Email
}

//Should return a list of emails.
// Todo: change this into a list of emails. email struct with metadata.
func (c *Email_client) QueryEmails(queryString string) queriedEmails {
	queriedEmails := queriedEmails{}
	messages_servicer := gmail.NewUsersMessagesService(c.srv)

	emailsCall := c.srv.Users.Messages.List(c.user)
	emailsCall.Q(queryString)
	matchedMessagesRes, _ := emailsCall.Do()
	matchedMessages := matchedMessagesRes.Messages
	//print(&matchedMessages)
	v := matchedMessages
	for _, email := range v {
		body, _ := messages_servicer.Get(c.user, email.Id).Format("full").Do()
		//fmt.Printf("%+v\n", body.Payload.Parts[0].Body.Data)
		emailByte, _ := base64.URLEncoding.DecodeString(body.Payload.Parts[0].Body.Data)
		emailString := string(emailByte)

		email := Email{email.Id, emailString}
		queriedEmails.emails = append(queriedEmails.emails, email)
	}
	return queriedEmails
}
