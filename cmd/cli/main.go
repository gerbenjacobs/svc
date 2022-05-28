// This executable Go file acts as a CLI tool for use with the UserService
// it allows you to create a JSON file out of your Users table
// Rationale: This shows you how you can have multiple tools within your /cmd/ folder,
// it also shows you the benefit of separating your storage layer from your code
// as we will re-use it.
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/gerbenjacobs/svc"
	"github.com/gerbenjacobs/svc/internal"
	"github.com/gerbenjacobs/svc/storage"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// load our config and db settings
	// our storage layer requires a *sql.DB object
	db, err := internal.NewDB(internal.NewConfig())
	if err != nil {
		log.Fatalf("failed to set up database: %v", err)
	}

	// create our user storage layer
	userRepo := storage.NewUserRepository(db)

	// select all users
	users := userRepo.AllUsers(context.Background())

	// print output to stdout
	log.Printf("Found %d users\n", len(users))
	for _, u := range users {
		log.Printf("%s\n", u) // %s calls the .String() method on our *app.User objects
	}

	// create slice of outputUser
	// this is a custom struct with custom marshaller that removes Token information
	var outputUsers []OutputUser
	for _, u := range users {
		outputUsers = append(outputUsers, OutputUser{User: u})
	}

	// marshall our data to JSON, using human-readable settings
	b, err := json.MarshalIndent(outputUsers, "", " ")
	if err != nil {
		log.Fatalf("failed to marshal output users: %v", err)
	}

	// write data to output.json
	// os.WriteFile creates and truncates the file
	err = os.WriteFile("output.json", b, 0644)
	if err != nil {
		log.Fatalf("failed to write output.json: %v", err)
	}
	log.Println("Finished writing to output.json")
}

// OutputUser is a local struct that encapsulates the user
// In combination with MarshalJSON we can exclude Token from being used in our output
type OutputUser struct {
	*svc.User
}

func (o OutputUser) MarshalJSON() ([]byte, error) {
	type ou OutputUser // prevent recursion
	u := ou(o)
	u.Token = "" // clear the Token from the struct
	return json.Marshal(u)
}
