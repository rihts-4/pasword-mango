package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type User struct {
	Name  string `firestore:"name"`
	Email string `firestore:"email"`
	Age   int    `firestore:"age"`
}

func main() {
	ctx := context.Background()

	config := &firebase.Config{
		ProjectID: "pass-mang0",
	}
	sa := option.WithCredentialsFile("adminkey.json")
	app, err := firebase.NewApp(ctx, config, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// Create
	_, _, err = client.Collection("users").Add(ctx, User{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   28,
	})
	if err != nil {
		log.Fatalf("Failed adding user: %v", err)
	}

	// Read all
	iter := client.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}

		var user User
		doc.DataTo(&user)
		fmt.Printf("User: %+v\n", user)
	}
}
