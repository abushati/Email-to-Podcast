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
	srv            *gmail.Service
	user           string
	messageService *gmail.UsersMessagesService
	QueriedEmails  []Email
}

type Email struct {
	Id   string
	Body string
}

func getMessageService(c *Email_client) *gmail.UsersMessagesService {
	if c.messageService == nil {
		messages_servicer := gmail.NewUsersMessagesService(c.srv)
		return messages_servicer
	} else {
		messages_servicer := c.messageService
		return messages_servicer
	}
}

func MakeClient(user string) *Email_client {
	b, err := ioutil.ReadFile("./email_client/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	// If modifying these scopes, delete your previously saved token.json.
	scope := []string{gmail.GmailReadonlyScope, gmail.GmailModifyScope}
	config, err := google.ConfigFromJSON(b, scope...)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	eml := &Email_client{srv, user, nil, []Email{}}

	return eml
}

//Should return a list of emails.
// Todo: change this into a list of emails. email struct with metadata.
func (c *Email_client) QueryEmails(queryString string) {
	var queriedEmails []Email

	messages_servicer := getMessageService(c)

	emailsCall := c.srv.Users.Messages.List(c.user)
	emailsCall.Q(queryString)
	matchedMessagesRes, _ := emailsCall.Do()
	matchedMessages := matchedMessagesRes.Messages
	//print(&matchedMessages)
	v := matchedMessages
	for _, email := range v {
		pulledEmail, _ := messages_servicer.Get(c.user, email.Id).Format("full").Do()
		//fmt.Printf("%+v\n", body.Payload.Parts[0].Body.Data)
		emailByte, _ := base64.URLEncoding.DecodeString(pulledEmail.Payload.Parts[0].Body.Data)
		emailString := string(emailByte)

		email := Email{email.Id, emailString}
		queriedEmails = append(queriedEmails, email)
	}
	c.QueriedEmails = append(c.QueriedEmails, queriedEmails...)

}

//Todo: perhaps private method
func (c *Email_client) SaveEmails() {
	if len(c.QueriedEmails) != 0 {
		print(" emails saved")

	} else {
		print("no emails to save")
	}
}

func (c *Email_client) MarkEmailRead(emailID string) {
	modifyRequest := &gmail.ModifyMessageRequest{RemoveLabelIds: []string{"UNREAD"}}
	messages_servicer := getMessageService(c)
	res, err := messages_servicer.Modify(c.user, emailID, modifyRequest).Do()
	//removeLabelId of "UNREAD
	if err != nil {
		print("here")
		println(res)
	}

}

/*
func (c *Email_client) MarkEmailRead(e Email) {
	modifyRequest := &gmail.ModifyMessageRequest{RemoveLabelIds: []string{"UNREAD"}}
	messages_servicer := getMessageService(c)
	messages_servicer.Modify(c.user, "179742756ef03ead", modifyRequest)
	//removeLabelId of "UNREAD

}
*/
