package freshdesk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nextlinktechnology/go-freshservice/querybuilder"
	"github.com/nextlinktechnology/mgm/v3"
)

type TicketManager interface {
	All() (TicketResults, error)
	Create(CreateTicket) (Ticket, error)
	View(int64) (Ticket, error)
	Search(querybuilder.Query) (TicketResults, error)
	Reply(int64, CreateConversation) (Conversation, error)
	Conversations(int64) (ConversationSlice, error)
	UpdatedSinceAll(string) (TicketResults, error)
}

type ticketManager struct {
	client *ApiClient
}

type TicketResults struct {
	next    string
	Results TicketSlice
	client  *ApiClient
}

func newTicketManager(client *ApiClient) ticketManager {
	return ticketManager{
		client,
	}
}

type Ticket struct {
	mgm.DefaultModel       `bson:",inline" json:"-"`
	Attachments            []interface{}          `bson:"attachments" json:"attachments"`
	CCEmails               []string               `bson:"cc_emails" json:"cc_emails"`
	DepartmentID           int64                  `bson:"department_id" json:"department_id"`
	CustomFields           map[string]interface{} `bson:"custom_fields" json:"custom_fields"`
	Deleted                bool                   `bson:"deleted" json:"deleted"`
	Description            string                 `bson:"description" json:"description"`
	DescriptionText        string                 `bson:"description_text" json:"description_text"`
	DueBy                  *time.Time             `bson:"due_by" json:"due_by"`
	Email                  string                 `bson:"email" json:"email"`
	EmailConfigID          int64                  `bson:"email_config_id" json:"email_config_id"`
	FirstResponseDueBy     *time.Time             `bson:"fr_due_by" json:"fr_due_by"`
	FirstResponseEscalated bool                   `bson:"fr_escalated" json:"fr_escalated"`
	FwdEmails              []string               `bson:"fwd_emails" json:"fwd_emails"`
	GroupID                int64                  `bson:"group_id" json:"group_id"`
	ID                     int64                  `bson:"id" json:"id"`
	IsEscalated            bool                   `bson:"is_escalated" json:"is_escalated"`
	Name                   string                 `bson:"name" json:"name"`
	Phone                  string                 `bson:"phone" json:"phone"`
	Priority               int                    `bson:"priority" json:"priority"`
	Category               int                    `bson:"category" json:"category"`
	SubCategory            []string               `bson:"sub_category" json:"sub_category"`
	ItemCategory           []string               `bson:"item_category" json:"item_category"`
	ReplyCCEmails          []string               `bson:"reply_cc_emails" json:"reply_cc_emails"`
	RequesterID            int64                  `bson:"requester_id" json:"requester_id"`
	ResponderID            int64                  `bson:"responder_id" json:"responder_id"`
	Source                 int                    `bson:"source" json:"source"`
	Spam                   bool                   `bson:"spam" json:"spam"`
	Status                 int                    `bson:"status" json:"status"`
	Subject                string                 `bson:"subjecte" json:"subject"`
	Tags                   []string               `bson:"tags" json:"tags"`
	ToEmails               []string               `bson:"to_emails" json:"to_emails"`
	Type                   string                 `bson:"type" json:"type"`
	CreatedAt              *time.Time             `bson:"created_at" json:"created_at"`
	UpdatedAt              *time.Time             `bson:"updated_at" json:"updated_at"`
	Urgency                string                 `bson:"urgency" json:"urgency"`
	Impact                 int64                  `bson:"impact" json:"impact"`
	Conversations          []Conversation         `bson:"-" json:"conversations"`
}

type RespTickets struct {
	Tickets []Ticket `json:"tickets,omitempty"`
}

type RespTicket struct {
	Ticket Ticket `json:"ticket,omitempty"`
}

type CreateTicket struct {
	Name               string                 `json:"name,omitempty"`
	RequesterID        int                    `json:"requester_id,omitempty"`
	Email              string                 `json:"email,omitempty"`
	Phone              string                 `json:"phone,omitempty"`
	Subject            string                 `json:"subject,omitempty"`
	Type               string                 `json:"type,omitempty"`
	Status             int                    `json:"status,omitempty"`
	Priority           int                    `json:"priority,omitempty"`
	Description        string                 `json:"description,omitempty"`
	ResponderID        int                    `json:"responder_id,omitempty"`
	Attachments        []interface{}          `json:"attachments,omitempty"`
	CCEmails           []string               `json:"cc_emails,omitempty"`
	CustomFields       map[string]interface{} `json:"custom_fields,omitempty"`
	DueBy              *time.Time             `json:"due_by,omitempty"`
	EmailConfigID      int                    `json:"email_config_id,omitempty"`
	FirstResponseDueBy *time.Time             `json:"fr_due_by,omitempty"`
	GroupID            int                    `json:"group_id,omitempty"`
	Source             int                    `json:"source,omitempty"`
	Tags               []string               `json:"tags,omitempty"`
	DepartmentID       int64                  `json:"department_id,omitempty"`
	Category           string                 `json:"category,omitempty"`
	SubCategory        []string               `json:"sub_category,omitempty"`
	ItemCategory       []string               `json:"item_category,omitempty"`
	Assets             string                 `json:"assets,omitempty"`
	Urgency            string                 `json:"urgency,omitempty"`
	Impact             int64                  `json:"impact,omitempty"`
}

type Conversation struct {
	mgm.DefaultModel `bson:",inline" json:"-"`
	Attachments      []interface{} `json:"attachments"`
	Body             string        `bson:"body" json:"body"`
	BodyText         string        `bson:"body_text" json:"body_text"`
	ID               int64         `bson:"id" json:"id"`
	Incoming         bool          `bson:"incoming" json:"incoming"`
	ToEmails         []string      `bson:"to_emails" json:"to_emails"`
	Private          bool          `bson:"private" json:"private"`
	Source           int           `bson:"source" json:"source"`
	SupportEmail     string        `bson:"support_email" json:"support_email"`
	TicketID         int64         `bson:"ticket_id" json:"ticket_id"`
	UserID           int64         `bson:"user_id" json:"user_id"`
	CreatedAt        *time.Time    `bson:"created_at" json:"created_at"`
	UpdatedAt        *time.Time    `bson:"updated_at" json:"updated_at"`
}

type RespConversations struct {
	Conversations []Conversation `json:"conversations,omitempty"`
}

type RespConversation struct {
	Conversation Conversation `json:"conversation,omitempty"`
}

type CreateConversation struct {
	Body        string        `json:"body,omitempty"`
	FromEmail   string        `json:"from_email,omitempty"`
	Attachments []interface{} `json:"attachments,omitempty"`
	UserID      int           `json:"user_id,omitempty"`
	CCEmails    []string      `json:"cc_emails,omitempty"`
	BCCEmails   []string      `json:"bcc_emails,omitempty"`
}

type Source int
type Status int
type Priority int

const (
	SourceEmail Source = 1 + iota
	SourcePortal
	SourcePhone
	SourceChat
	SourceFeedbackWidget
	SourceYammer
	SourceAWSCloudwatch
	SourcePagerduty
	SourceWalkup
	SourceSlack
)

const (
	StatusOpen Status = 2 + iota
	StatusPending
	StatusResolved
	StatusClosed
)

const (
	PriorityLow Priority = 1 + iota
	PriorityMedium
	PriorityHigh
	PriorityUrgent
)

func (s Source) Value() int {
	return int(s)
}

func (s Status) Value() int {
	return int(s)
}

func (p Priority) Value() int {
	return int(p)
}

func (t Ticket) Print() {
	jsonb, _ := json.MarshalIndent(t, "", "    ")
	fmt.Println(string(jsonb))
}
func (r Conversation) Print() {
	jsonb, _ := json.MarshalIndent(r, "", "    ")
	fmt.Println(string(jsonb))
}

type TicketSlice []Ticket

func (s TicketSlice) Len() int { return len(s) }

func (s TicketSlice) Less(i, j int) bool { return s[i].ID < s[j].ID }

func (s TicketSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s TicketSlice) Print() {
	for _, ticket := range s {
		fmt.Println(ticket.Subject)
	}
}

type ConversationSlice []Conversation

func (s ConversationSlice) Len() int { return len(s) }

func (s ConversationSlice) Less(i, j int) bool { return s[i].CreatedAt.Unix() > s[j].CreatedAt.Unix() }

func (s ConversationSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ConversationSlice) Print() {
	for _, ticket := range s {
		fmt.Println(ticket.BodyText)
	}
}

func (manager ticketManager) All() (TicketResults, error) {
	resp := RespTickets{}
	output := TicketSlice{}
	headers, err := manager.client.get(endpoints.tickets.all, &resp)
	if err != nil {
		return TicketResults{}, err
	}
	output = append(output, resp.Tickets...)

	return TicketResults{
		Results: output,
		client:  manager.client,
		next:    manager.client.getNextLink(headers),
	}, nil
}

func (manager ticketManager) UpdatedSinceAll(timeString string) (TicketResults, error) {
	resp := RespTickets{}
	output := TicketSlice{}
	headers, err := manager.client.get(endpoints.tickets.updatedSinceAll(timeString), &resp)
	if err != nil {
		return TicketResults{}, err
	}
	output = append(output, resp.Tickets...)

	return TicketResults{
		Results: output,
		client:  manager.client,
		next:    manager.client.getNextLink(headers),
	}, nil
}

func (manager ticketManager) Create(ticket CreateTicket) (Ticket, error) {
	output := RespTicket{}
	jsonb, err := json.Marshal(ticket)
	if err != nil {
		return Ticket{}, err
	}
	err = manager.client.postJSON(endpoints.tickets.create, jsonb, &output, http.StatusCreated)
	if err != nil {
		return Ticket{}, err
	}
	return output.Ticket, nil
}

func (manager ticketManager) View(id int64) (Ticket, error) {
	output := RespTicket{}
	_, err := manager.client.get(endpoints.tickets.view(id), &output)
	if err != nil {
		return Ticket{}, err
	}

	return output.Ticket, nil
}

func (manager ticketManager) Conversations(id int64) (ConversationSlice, error) {
	resp := RespConversations{}
	output := ConversationSlice{}
	_, err := manager.client.get(endpoints.tickets.conversations(id), &resp)
	if err != nil {
		return ConversationSlice{}, err
	}
	output = append(output, resp.Conversations...)
	return output, nil
}

func (manager ticketManager) Reply(id int64, reply CreateConversation) (Conversation, error) {
	output := RespConversation{}
	jsonb, err := json.Marshal(reply)
	if err != nil {
		return Conversation{}, err
	}
	err = manager.client.postJSON(endpoints.tickets.reply(id), jsonb, &output, http.StatusCreated)
	if err != nil {
		return Conversation{}, err
	}
	return output.Conversation, nil
}

func (manager ticketManager) Search(query querybuilder.Query) (TicketResults, error) {
	resp := RespTickets{}
	output := TicketSlice{}
	headers, err := manager.client.get(endpoints.tickets.search(query.URLSafe()), &resp)
	if err != nil {
		return TicketResults{}, err
	}
	output = append(output, resp.Tickets...)

	return TicketResults{
		Results: output,
		client:  manager.client,
		next:    manager.client.getNextLink(headers),
	}, nil
}

func (results TicketResults) Next() (TicketResults, error) {
	if results.next == "" {
		return TicketResults{}, errors.New("no more tickets")
	}
	resp := RespTickets{}
	output := TicketSlice{}
	headers, err := results.client.get(results.next, &resp)
	if err != nil {
		return TicketResults{}, err
	}
	output = append(output, resp.Tickets...)

	return TicketResults{
		Results: output,
		client:  results.client,
		next:    results.client.getNextLink(headers),
	}, nil
}

func (results *TicketResults) FilterTags(tags ...string) *TicketResults {
	filtered := TicketSlice{}
	for _, ticket := range results.Results {
		_filterFlag := false
		for _, ticketTag := range ticket.Tags {
			for _, filterTag := range tags {
				if ticketTag == filterTag {
					_filterFlag = true
					break
				}
			}
		}
		if _filterFlag {
			continue
		}
		filtered = append(filtered, ticket)
	}
	results.Results = filtered
	return results
}

func (results *TicketResults) FilterTypes(filterTypes ...string) *TicketResults {
	filtered := TicketSlice{}
	for _, ticket := range results.Results {
		_filterFlag := false
		for _, filterType := range filterTypes {
			if ticket.Type == filterType {
				_filterFlag = true
				break
			}
		}
		if _filterFlag {
			continue
		}
		filtered = append(filtered, ticket)
	}
	results.Results = filtered
	return results
}
