# Form3 Take Home Exercise: Client library for fake account service
Robert Špendl

robert.spendl@chronos.si

https://github.com/rspendl

## Description
The client library for fake account service implements the `Create`, `Fetch`, `List` and `Delete` operations on the `accounts` resource.

### Example

```go
import (
	account "accountapi"
	"accountapi/data"
	"accountapi/lib"
	"fmt"

	"github.com/biter777/countries"
	"github.com/google/uuid"
	"golang.org/x/text/currency"
)

func main() {
	client, err := account.New(account.Config{
		Server:             "http://127.0.0.1:8080",
		MaxConnections:     64,
		MaxIdleConnections: 16,
		Timeout:            3,
	})
	id, _ := uuid.NewRandom()
	org, _ := uuid.NewRandom()
	acc := data.Account{
		ID:             id,
		OrganisationID: org,
		Attributes: data.Attributes{
			Country:      data.NewCountryCode(countries.UnitedKingdom),
			BaseCurrency: data.NewCurrency(currency.GBP),
		},
	}
	createdAccount, err := client.Create(&acc)
	if err != nil {
		fmt.Printf("Error creating account: %v", err)
	}
	fmt.Printf("Account %s created", createdAccount.ID.String())
	fetchedAccount, err := client.Fetch(acc.ID)
	if err != nil {
		fmt.Printf("Error fetching account: %v", err)
	}
	fmt.Printf("Account %s fetched", fetchedAccount.ID.String())
	listedAccounts, err := client.List(lib.PageNumber(0), lib.PageSize(10))
	if err != nil {
		fmt.Printf("Error listing accounts: %v", err)
	}
	fmt.Printf("%d accounts listed", len(*listedAccounts))
	err = client.Delete(acc.ID, 0) // Delete version 0.
	if err != nil {
		fmt.Printf("Error deleting account: %v", err)
	}
}
```

## How to run tests
```docker-compose up```

or run tests directly

```go test ./...```

this includes list tests that take >10s

```go test ./... -tags=long```

## Notes

- The tests expect an empty account database when they start: running the tests will **erase** all accounts, including the ones that were not created by the tests.
- The tests read APISERVICE environment variable for account service URL, the default is http://127.0.0.1:8080.
- Long tests (lists, parallel requests) can be run with *go test -tags=long* .
- For currency codes, the client library uses golang.org/x/text/currency and for country codes github.com/biter777/countries . The libraries are retrieved using `go get`, for production environment libraries should be managed using `dep`. 
- Client library does not implement cancellable requests as the responses are quick so setting a reasonable timeout (~3s in tests) is sufficient.
- The account server, provided for the exercise, does not limit the number of connections to the SQL database, multiple parallel requests (~100) exhaust the connection pool and cause server errors. The tests limit the number of connections to the server to 80 to avoid these errors.
- The following differences between the [documentation](http://api-docs.form3.tech/api.html#organisation-accounts) and running service were found:
  - default page[size] parameter is documented to be 100, the service implements 1000;
  - "name" and "alternative_names" are not implemented in the service (documentation only states that private_identification and relationships are missing); the client library still sends these fields, but they are omitted in the tests as the values can't be fetched.
