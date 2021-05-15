package main

import (
	"email_podcast/email_client"
	"fmt"
)

func main() {
	/*
		helloHandler := func(w http.ResponseWriter, req *http.Request) {
			io.WriteString(w, "Hello, world!\n")
		}

		http.HandleFunc("/auth", helloHandler)
		log.Println("Listing for requests at http://localhost:8000/auth")
		log.Fatal(http.ListenAndServe(":8000", nil))
	*/
	//user := "me"
	user := "arvid.b901@gmail.com"
	client := email_client.MakeClient(user)

	//Todo: Query for emails that are unread
	emails := client.QueryEmails("from:(Medium Daily Digest noreply@medium.com)")
	fmt.Printf("%+v\n", emails)
}
