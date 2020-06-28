package freshdesk

import (
	"encoding/json"
	"net/http"
	"time"
)

type ServiceRequestManager interface {
	Create(CreateTicket) (Ticket, error)
	View(int64) (Ticket, error)
}

type serviceRequestManager struct {
	client *ApiClient
}

func newServiceRequestManager(client *ApiClient) serviceRequestManager {
	return serviceRequestManager{
		client,
	}
}

type ServiceRequest struct {
	ID             int64       `json:"id"`
	CreatedAt      *time.Time  `json:"created_at"`
	UpdatedAt      *time.Time  `json:"updated_at"`
	Quantity       int         `json:"quantity"`
	Stage          int         `json:"stage"`
	Loaned         bool        `json:"loaned"`
	CostPerRequest float32     `json:"cost_per_request"`
	Remarks        string      `json:"remarks"`
	DeliveryTime   int         `json:"delivery_time"`
	IsParent       bool        `json:"is_parent"`
	ServiceItemID  int64       `json:"service_item_id"`
	CustomFields   interface{} `json:"custom_fields"`
}

type RespServiceRequests struct {
	ServiceRequests []ServiceRequest `json:"requested_items,omitempty"`
}

type CreateServiceRequest struct {
	Quantity     int         `json:"quantity"`
	RequestedFor string      `json:"requested_for"`
	Email        string      `json:"email"`
	CustomFields interface{} `json:"custom_fields"`
}

type SRStage int

const (
	SRStageRequested SRStage = 1 + iota
	SRStageDelivered
	SRStageCancelled
	SRStageFulfilled
	SRStagePartiallyFulfilled
)

func (s SRStage) Value() int {
	return int(s)
}

func (manager serviceRequestManager) Create(ticket CreateTicket) (Ticket, error) {
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

func (manager serviceRequestManager) View(id int64) (Ticket, error) {
	output := RespTicket{}
	_, err := manager.client.get(endpoints.servicerequest.view(id), &output)
	if err != nil {
		return Ticket{}, err
	}

	return output.Ticket, nil
}
