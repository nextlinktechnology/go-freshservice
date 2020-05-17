package freshdesk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nextlinktechnology/go-freshservice/querybuilder"
	"github.com/nextlinktechnology/mgm/v3"
)

type RequesterManager interface {
	All() (RequesterSlice, error)
	Create(*Requester) (*Requester, error)
	Search(querybuilder.Query) (RequesterResults, error)
	Update(int64, *Requester) (*Requester, error)
}

type requesterManager struct {
	client *ApiClient
}

func newrequesterManager(client *ApiClient) requesterManager {
	return requesterManager{
		client,
	}
}

type RequesterResults struct {
	next    string
	Results RequesterSlice
	client  *ApiClient
}

type Requester struct {
	mgm.DefaultModel                          `bson:",inline"`
	ID                                        int64                  `bson:"id" json:"id,omitempty"`
	FirstName                                 string                 `bson:"first_name" json:"first_name,omitempty"`
	LastName                                  string                 `bson:"last_name" json:"last_name,omitempty"`
	JobTitle                                  string                 `bson:"job_title" json:"job_title,omitempty"`
	PrimaryEmail                              string                 `bson:"primary_email" json:"primary_email,omitempty"`
	SecondaryEmails                           []string               `bson:"secondary_emails" json:"secondary_emails,omitempty"`
	WorkPhoneNumber                           string                 `bson:"work_phone_number" json:"work_phone_number,omitempty"`
	MobilePhoneNumber                         string                 `bson:"mobile_phone_number" json:"mobile_phone_number,omitempty"`
	DepartmentIDs                             []int64                `bson:"department_ids" json:"department_ids,omitempty"`
	CanSeeAllTicketsFromAssociatedDepartments bool                   `bson:"can_see_all_tickets_from_associated_departments" json:"can_see_all_tickets_from_associated_departments,omitempty"`
	ReportingManagerID                        int64                  `bson:"reporting_manager_id" json:"reporting_manager_id,omitempty"`
	Address                                   string                 `bson:"address" json:"address,omitempty"`
	TimeZone                                  string                 `bson:"time_zone" json:"time_zone,omitempty"`
	TimeFormat                                string                 `bson:"time_format" json:"time_format,omitempty"`
	Language                                  string                 `bson:"language" json:"language,omitempty"`
	LocationID                                int64                  `bson:"location_id" json:"location_id,omitempty"`
	BackgroundInformation                     string                 `bson:"background_information" json:"background_information,omitempty"`
	CustomFields                              map[string]interface{} `bson:"custom_fields" json:"custom_fields,omitempty"`
	Active                                    bool                   `bson:"active" json:"active,omitempty"`
	HasLoggedIn                               bool                   `bson:"has_logged_in" json:"has_logged_in,omitempty"`
	CreatedAt                                 *time.Time             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt                                 *time.Time             `bson:"updated_at" json:"updated_at,omitempty"`
}

type RequesterSlice []Requester

func (s RequesterSlice) Len() int { return len(s) }

func (s RequesterSlice) Less(i, j int) bool { return s[i].ID < s[j].ID }

func (s RequesterSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s RequesterSlice) Print() {
	for _, requester := range s {
		fmt.Println(requester.FirstName)
	}
}

func (manager requesterManager) All() (RequesterSlice, error) {
	output := RequesterSlice{}
	headers, err := manager.client.get(endpoints.requesters.all, &output)
	if err != nil {
		return RequesterSlice{}, err
	}
	for {
		nextLink := manager.client.getNextLink(headers)
		if nextLink == "" {
			break
		}
		nextSlice := RequesterSlice{}
		headers, err = manager.client.get(nextLink, &nextSlice)
		if err != nil {
			return RequesterSlice{}, err
		}
		output = append(output, nextSlice...)
	}
	return output, nil
}

func (manager requesterManager) Search(query querybuilder.Query) (RequesterResults, error) {
	output := struct {
		Slice RequesterSlice `json:"results,omitempty"`
	}{}
	headers, err := manager.client.get(endpoints.requesters.search(query.URLSafe()), &output)
	if err != nil {
		return RequesterResults{}, err
	}
	return RequesterResults{
		Results: output.Slice,
		client:  manager.client,
		next:    manager.client.getNextLink(headers),
	}, nil
}

func (manager requesterManager) Create(requester *Requester) (*Requester, error) {
	output := &Requester{}
	jsonb, err := json.Marshal(requester)
	if err != nil {
		return output, err
	}
	err = manager.client.postJSON(endpoints.requesters.create, jsonb, &output, http.StatusCreated)
	if err != nil {
		return &Requester{}, err
	}
	return output, nil
}

func (manager requesterManager) Update(id int64, requester *Requester) (*Requester, error) {
	output := &Requester{}
	jsonb, err := json.Marshal(requester)
	if err != nil {
		return output, err
	}
	err = manager.client.put(endpoints.requesters.update(id), jsonb, &output, http.StatusOK)
	if err != nil {
		return &Requester{}, err
	}
	return output, nil
}
