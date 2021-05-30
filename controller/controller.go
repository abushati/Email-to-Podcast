package controller

import "email_podcast/email_client"

func GetEmailsToConvert(user string, query string) {
	
	emailClient := email_client.MakeClient(user)
	emailClient.QueryEmails(query)
	
	print(len(emailClient.QueriedEmails))
}

func convertEmailToMP3() {}

func saveEmailInfoToDB() {}

func saveMP3FileToDB() {}

func MailLoop() {
	user := "arvid.b901@gmail.com"
	//Todo: This would need to be the json file
	emailQueries := []string{"from:(Medium Daily Digest noreply@medium.com)label:unread",
		"from:(robinhood snacks <noreply@robinhood.com>) label:unread "}
	for _, query := emailQueries{
		GetEmailsToConvert(user)
	}
	
}
