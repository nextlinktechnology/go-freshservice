package freshdesk

import "fmt"

type departmentEndpoints struct {
	all    string
	create string
	update func(int64) string
}

type requesterEndpoints struct {
	all    string
	create string
	search func(string) string
	update func(int64) string
}

type ticketEndpoints struct {
	all             string
	create          string
	view            func(int64) string
	search          func(string) string
	reply           func(int64) string
	conversations   func(int64) string
	updatedSinceAll func(string) string
}

type servicerequestEndpoints struct {
	create func(int64) string
	view   func(int64) string
}

var endpoints = struct {
	departments    departmentEndpoints
	requesters     requesterEndpoints
	tickets        ticketEndpoints
	servicerequest servicerequestEndpoints
}{
	departments: departmentEndpoints{
		all:    "/api/v2/departments",
		create: "/api/v2/departments",
		update: func(id int64) string { return fmt.Sprintf("/api/v2/departments/%d", id) },
	},
	requesters: requesterEndpoints{
		all:    "/api/v2/requesters",
		create: "/api/v2/requesters",
		update: func(id int64) string { return fmt.Sprintf("/api/v2/requesters/%d", id) },
		search: func(query string) string { return fmt.Sprintf("/api/v2/requesters?%s", query) },
	},
	tickets: ticketEndpoints{
		all:           "/api/v2/tickets",
		create:        "/api/v2/tickets",
		view:          func(id int64) string { return fmt.Sprintf("/api/v2/tickets/%d", id) },
		search:        func(query string) string { return fmt.Sprintf("/api/v2/tickets?%s", query) },
		reply:         func(id int64) string { return fmt.Sprintf("/api/v2/tickets/%d/reply", id) },
		conversations: func(id int64) string { return fmt.Sprintf("/api/v2/tickets/%d/conversations", id) },
		updatedSinceAll: func(timeString string) string {
			return fmt.Sprintf("/api/v2/tickets?updated_since=%s", timeString)
		},
	},
	servicerequest: servicerequestEndpoints{
		create: func(id int64) string { return fmt.Sprintf("/api/v2/service_catalog/items/%d/place_request", id) },
		view:   func(id int64) string { return fmt.Sprintf("/api/v2/tickets/%d/requested_items", id) },
	},
}
