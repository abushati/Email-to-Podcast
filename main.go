package main

import (
	"email_podcast/database"
	"email_podcast/email_client"
	"io"
	"log"
	"net/http"
)

func main() {

	//user := "me"
	user := "arvid.b901@gmail.com"
	client := email_client.MakeClient(user)
	//client.SaveEmails()
	//Todo: Query for emails that are unread
	//client.QueryEmails("from:(Medium Daily Digest noreply@medium.com)")
	//client.SaveEmails()

	//fmt.Printf("%+v\n", client.QueriedEmails)

	client.MarkEmailRead("2314")
	//fmt.Printf("%+v\n", client.QueriedEmails[0])

	db := database.OpenConnection()
	sqlStatement := `INSERT INTO person (name, nickname) VALUES ('12', '12')`
	_, err := db.Exec(sqlStatement)
	print(err.Error())
}

func getToken() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}

	http.HandleFunc("/auth", helloHandler)
	log.Println("Listing for requests at http://localhost:8000/auth")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
