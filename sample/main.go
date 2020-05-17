package main

import (
	"log"
	"os"

	freshservice "github.com/nextlinktechnology/go-freshservice"
)

func main() {
	logger := log.New(os.Stdout, "[logger] ", 0)

	client := freshservice.Init("domain", "apiKey", &freshservice.ClientOptions{Logger: logger}) // Or use freshdesk.EmptyOptions()

	departments, err := client.Departments.All()
	if err != nil {
		panic(err)
	}
	departments.Print()

	tickets, err := client.Tickets.All()
	if err != nil {
		panic(err)
	}
	tickets.Results.Print()

	ticket, err := client.Tickets.Create(freshservice.CreateTicket{
		Subject:     "Ticket Subject",
		Description: "Ticket description.",
		Email:       "identifier@domain.tld",
		Status:      freshservice.StatusOpen.Value(),
		Priority:    freshservice.PriorityLow.Value(),
	})
	if err != nil {
		panic(err)
	}
	ticket.Print()
}
