package account

import (
	"accountapi/data"
	"accountapi/lib"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Config contains configuration parameters for the client.
type Config struct {
	Server             string
	MaxConnections     int
	MaxIdleConnections int
	Timeout            int
}

// Client enables access to web service.
type Client struct {
	server     string
	httpClient *http.Client
}

// New creates a new client, used to connect to web account service. Communication parameters can be set to optimize
// number of open connections.
func New(cfg Config) (*Client, error) {
	transport := &http.Transport{}
	if cfg.MaxIdleConnections > 0 {
		transport.MaxConnsPerHost = cfg.MaxConnections
		transport.MaxIdleConns = cfg.MaxIdleConnections
		transport.MaxIdleConnsPerHost = cfg.MaxIdleConnections
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(cfg.Timeout) * time.Second,
	}
	return &Client{
		server:     cfg.Server,
		httpClient: client,
	}, nil
}

func (c *Client) get(endpoint string, jsonResponse interface{}) error {
	return c.doRequest("GET", endpoint, nil, jsonResponse)
}

func (c *Client) post(endpoint string, jsonRequest io.Reader, jsonResponse interface{}) error {
	return c.doRequest("POST", endpoint, jsonRequest, jsonResponse)
}

func (c *Client) delete(endpoint string, jsonRequest io.Reader, jsonResponse interface{}) error {
	return c.doRequest("DELETE", endpoint, jsonRequest, jsonResponse)
}

// doRequest sends HTTP request to the server and either inserts the response into jsonResponse
// or provides error message as ErrorAPI as a return value. Other http errors are returned
// when there's a communication error.
func (c *Client) doRequest(method string, endpoint string, jsonRequest io.Reader, jsonResponse interface{}) error {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.server, endpoint), jsonRequest)
	if err != nil {
		return err
	}
	if method != "GET" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if isErrorStatus(resp.StatusCode) {
		errMessage := data.ErrorMessage{}
		err = json.NewDecoder(resp.Body).Decode(&errMessage)
		if err != nil {
			return err
		}
		return lib.NewErrorAPI(errMessage.Message)
	}
	if method == "DELETE" { // DELETE does not return anything if it succeeds.
		return nil
	}
	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	return err
}

// Health checks server connectivity.
func (c *Client) Health() bool {
	jResult := data.HealthResponse{}
	err := c.get("/v1/health", &jResult)
	if err != nil {
		return false
	}
	return jResult.Status == "up"
}

// Create creates an account in the accont service. On success, it returns the
// account data, returned by the server, otherwise it returns and ErrorAPI error that
// includes the error message, returned by the server.
func (c *Client) Create(account *data.Account) (*data.Account, error) {
	requestType := data.Accounts // Force type "accounts" in every create request.
	jResult := data.ResponseData{}
	jRequest := data.RequestCreate{
		Data: data.AccountData{
			ID:             account.ID,
			OrganisationID: account.OrganisationID,
			Type:           &requestType,
			Attributes:     account.Attributes,
		},
	}
	bin, err := json.Marshal(&jRequest)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(bin)
	err = c.post("/v1/organisation/accounts", reader, &jResult)
	if err != nil {
		return nil, err
	}
	return accountFromResponse(&jResult), nil
}

// Fetch fetches an account by id. On success, it returns the account data,
// returned by the server, otherwise it returns and ErrorAPI error that
// includes the error message, returned by the server.
func (c *Client) Fetch(id uuid.UUID) (*data.Account, error) {
	jResult := data.ResponseData{}
	err := c.get(fmt.Sprintf("%s%s", "/v1/organisation/accounts/", id.String()), &jResult)
	if err != nil {
		return nil, err
	}
	return accountFromResponse(&jResult), nil
}

// List retrieves an array of accounts on pageNumber where the size of a page is defined by pageSize.
// If pageNumber is omitted (defaults to 0), the function has to pass data.PNNone, and pageSize
// can be omitted with data.PSNone. pageNumber can have two special values data.First or data.Last,
// representing the first or the last page.
// On success, List returns an array of accounts, returned by the server (can be empty array, but not nil),
// otherwise it returns and ErrorAPI error that includes the error message, returned by the server.
func (c *Client) List(pageNumber lib.PageNumber, pageSize lib.PageSize) (*[]data.Account, error) {
	jResult := data.ResponseDataList{}
	query := ""
	switch pageNumber {
	case lib.PNNone:
		break
	case lib.First:
		query = "page[number]=first"
	case lib.Last:
		query = "page[number]=last"
	default:
		if int(pageNumber) > 0 {
			query = fmt.Sprintf("page[number]=%d", pageNumber)
		} else {
			return nil, lib.NewErrorInvalidArgument(fmt.Sprintf("pageNumber=%d", pageNumber))
		}
	}
	if pageSize != lib.PSNone {
		query = query + "&"
	}
	switch pageSize {
	case lib.PSNone:
		break
	default:
		if int(pageSize) > 0 {
			query = query + fmt.Sprintf("page[size]=%d", pageSize)
		} else {
			return nil, lib.NewErrorInvalidArgument(fmt.Sprintf("pageSize=%d", pageNumber))
		}
	}

	err := c.get(fmt.Sprintf("%s%s", "/v1/organisation/accounts/?", query), &jResult)
	if err != nil {
		return nil, err
	}
	return accountListFromResponse(&jResult), nil
}

// Delete deletes the account identified by id and version.
// No data but potential ErrorAPI error with the error message is returned.
func (c *Client) Delete(id uuid.UUID, version int) error {
	// As id is always a valid UUID, there's no need to perform additional validation before contacting the server.
	jResult := data.ResponseData{}
	jRequest := data.RequestDelete{
		Data: data.RequestDeleteData{
			ID:      id,
			Version: version,
		},
	}
	bin, err := json.Marshal(&jRequest)
	reader := bytes.NewReader(bin)
	err = c.delete(fmt.Sprintf("/v1/organisation/accounts/%s?version=%d", id.String(), version), reader, &jResult)
	if err != nil {
		return err
	}
	return nil
}

// accountFromResponse copies the values from response into Account structure.
func accountFromResponse(r *data.ResponseData) *data.Account {
	if r == nil {
		return nil
	}
	if *r.Data.Type != data.Accounts {
		panic("response data can't be handled for type " + r.Data.Type.String())
	}
	return &data.Account{
		ID:             r.Data.ID,
		OrganisationID: r.Data.OrganisationID,
		Version:        *r.Data.Version,
		Attributes:     r.Data.Attributes,
	}
}

// accountListFromResponse copies the values from List response into an array of Accounts.
func accountListFromResponse(r *data.ResponseDataList) *[]data.Account {
	if r == nil {
		return nil
	}
	accList := []data.Account{}
	for _, d := range r.Data {
		accList = append(accList, data.Account{
			ID:             d.ID,
			OrganisationID: d.OrganisationID,
			Version:        *d.Version,
			Attributes:     d.Attributes,
		})
	}
	return &accList
}

// isErrorStatus returns true for response status codes that are not 2xx.
func isErrorStatus(statusCode int) bool {
	return statusCode < 200 || statusCode > 299
}
