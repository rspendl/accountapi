package data

import (
	"time"

	"github.com/google/uuid"
)

// HealthResponse for status of health endpoint.
type HealthResponse struct {
	Status string `json:"status"`
}

// AccountData is an extended data struct for communicating with server.
// The pointer types allow omitting optional data.
type AccountData struct {
	ID             uuid.UUID   `json:"id"`
	OrganisationID uuid.UUID   `json:"organisation_id"`
	Type           *RecordType `json:"type,omitempty"`
	Version        *int        `json:"version,omitempty"`
	CreatedOn      *time.Time  `json:"created_on,omitempty"`
	ModifiedOn     *time.Time  `json:"modified_on,omitempty"`
	Attributes     Attributes  `json:"attributes,omitempty"`
}

// ResponseData contains account data and links, a response from account service.
type ResponseData struct {
	Data  AccountData   `json:"data,omitempty"`
	Links ResponseLinks `json:"links,omitempty"`
}

// ResponseDataList contains data with an array of AccountData, a response from account service List request.
type ResponseDataList struct {
	Data  []AccountData `json:"data,omitempty"`
	Links ResponseLinks `json:"links,omitempty"`
}

// ResponseLinks ...
type ResponseLinks struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
}

// RequestCreate is passed to server to create an account.
type RequestCreate struct {
	Data AccountData `json:"data"`
}

// RequestDelete is passed to server to delete an account.
type RequestDelete struct {
	Data RequestDeleteData `json:"data"`
}

// RequestDeleteData always contains id and version.
type RequestDeleteData struct {
	ID      uuid.UUID `json:"id"`
	Version int       `json:"version"`
}

// ErrorMessage is returned by server with error status codes.
type ErrorMessage struct {
	Message string `json:"error_message"`
}
