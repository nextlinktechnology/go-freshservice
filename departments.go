package freshdesk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nextlinktechnology/mgm/v3"
)

type DepartmentManager interface {
	All() (DepartmentSlice, error)
	Create(CreateDepartment) (Department, error)
	Update(int64, CreateDepartment) (Department, error)
}

type departmentManager struct {
	client *ApiClient
}

func newDepartmentManager(client *ApiClient) departmentManager {
	return departmentManager{
		client,
	}
}

type Department struct {
	mgm.DefaultModel `bson:",inline"`
	ID               int64                  `bson:"id" json:"id"`
	Name             string                 `bson:"name" json:"name,omitempty"`
	Description      string                 `bson:"description" json:"description,omitempty"`
	HeadUserID       string                 `bson:"head_user_id" json:"head_user_id,omitempty"`
	PrimeUserID      string                 `bson:"prime_user_id" json:"prime_user_id,omitempty"`
	Domains          []string               `bson:"domains" json:"domains,omitempty"`
	CustomFields     map[string]interface{} `bson:"custom_fields" json:"custom_fields,omitempty"`
	CreatedAt        *time.Time             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt        *time.Time             `bson:"updated_at" json:"updated_at,omitempty"`
}

type CreateDepartment struct {
	ID           int64                  `json:"id,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	HeadUserID   string                 `json:"head_user_id,omitempty"`
	PrimeUserID  string                 `json:"prime_user_id,omitempty"`
	Domains      []string               `json:"domains,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
	CreatedAt    *time.Time             `json:"created_at,omitempty"`
	UpdatedAt    *time.Time             `json:"updated_at,omitempty"`
}

type RespDepartment struct {
	Departments []Department `json:"departments,omitempty"`
}

type DepartmentSlice []Department

func (c DepartmentSlice) Len() int {
	return len(c)
}

func (c DepartmentSlice) Less(i, j int) bool {
	return c[i].ID < c[j].ID
}

func (c DepartmentSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c DepartmentSlice) Print() {
	for _, department := range c {
		fmt.Println(department.Name)
	}
}

func (manager departmentManager) All() (DepartmentSlice, error) {
	resp := RespDepartment{}
	output := DepartmentSlice{}
	headers, err := manager.client.get(endpoints.departments.all, &resp)
	if err != nil {
		return DepartmentSlice{}, err
	}
	output = append(output, resp.Departments...)

	for {
		nextLink := manager.client.getNextLink(headers)
		if nextLink == "" {
			break
		}
		nextResp := RespDepartment{}
		headers, err = manager.client.get(nextLink, &nextResp)
		if err != nil {
			return DepartmentSlice{}, err
		}
		output = append(output, nextResp.Departments...)
	}
	return output, nil
}

func (manager departmentManager) Create(department CreateDepartment) (Department, error) {
	output := RespDepartment{}
	jsonb, err := json.Marshal(department)
	if err != nil {
		return Department{}, err
	}
	err = manager.client.postJSON(endpoints.departments.create, jsonb, &output, http.StatusCreated)
	if err != nil {
		return Department{}, err
	}
	return output.Departments[0], nil
}

func (manager departmentManager) Update(id int64, department CreateDepartment) (Department, error) {
	output := RespDepartment{}
	jsonb, err := json.Marshal(department)
	if err != nil {
		return Department{}, err
	}
	err = manager.client.put(endpoints.departments.update(id), jsonb, &output, http.StatusOK)
	if err != nil {
		return Department{}, err
	}
	return output.Departments[0], nil
}
