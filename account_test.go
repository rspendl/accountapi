// Client library tests assume an account API service is running, no mock server is implemented.
// Account service docker image: 'form3tech/interview-accountapi:v1.0.0-4-g63cf8434'
package account_test

import (
	account "accountapi"
	"accountapi/data"
	"accountapi/lib"
	"accountapi/lib/test"
	"encoding/json"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/biter777/countries"
	"github.com/google/uuid"
	"golang.org/x/text/currency"
)

const (
	SERVER = "http://127.0.0.1:8080"
	// TestMaxConnections has to be somewhat lower than the number of available PostgreSQL connections
	// at the server (100 by default) or the server will issue "server_errors".
	TestMaxConnections     = 80
	TestMaxIdleConnections = 16
	TestTimeout            = 3
	TestLongTimeout        = 10 // Timeout for complete duration of tests of parallel execution.
	TestCountry            = countries.UnitedKingdom
	TestCurrency           = "GBP"
)

func newClient() (*account.Client, error) {
	server := os.Getenv("APISERVICE")
	if server == "" {
		server = SERVER
	}
	return account.New(account.Config{
		Server:             server,
		MaxConnections:     TestMaxConnections,
		MaxIdleConnections: TestMaxIdleConnections,
		Timeout:            TestTimeout,
	})
}

// generateBasicAccount generates the smallest valid account with only id, organisation_id and country fields set.
func generateBasicAccount() *data.Account {
	id, _ := uuid.NewRandom()
	testCountry := data.NewCountryCode(TestCountry)
	return &data.Account{
		ID:             id,
		OrganisationID: id,
		Attributes: data.Attributes{
			Country: testCountry,
		},
	}

}

func TestBasicOperation(t *testing.T) {
	acc := generateBasicAccount()
	testBasicOperation(t, acc)
}

func TestFullAccountData(t *testing.T) {
	id, _ := uuid.NewRandom()
	testCountry := data.NewCountryCode(TestCountry)
	testCurrency := data.NewCurrency(currency.MustParseISO(TestCurrency))

	// TODO: Name and AlternativeNames are not used by the service API.
	// Switched and Status are returned as false / confirmed, the test sets both fields
	// to enable simple verification if the data in the response match the request.
	acc := &data.Account{
		ID:             id,
		OrganisationID: id,
		Attributes: data.Attributes{
			Country:       testCountry,
			BaseCurrency:  testCurrency,
			AccountNumber: "123",
			BankID:        "THEBANK",
			BankIDCode:    "THEBANKCODE",
			BIC:           "SOMEBIC9",
			IBAN:          "XY42SOMEIBAN123456",
			// Name:                    []string{"Account Holder", "Another Name"},
			// AlternativeNames:        []string{"AltFirst Name", "AltMiddle Name", "AltLast Name"},
			AccountClassification:   data.Business,
			JointAccount:            true,
			AccountMatchingOptOut:   true,
			SecondaryIdentification: "2ID",
			Switched:                false,
			Status:                  data.Confirmed,
		},
	}
	testBasicOperation(t, acc)
}

// testBasicOperation creates an account with minimal data (id and country), fetches the data and deletes it.
func testBasicOperation(t *testing.T, acc *data.Account) {
	client, err := newClient()
	if err != nil {
		t.Fail()
	}
	if !client.Health() {
		t.Fatal("Can't connect to server")
	}
	createdAccount, err := client.Create(acc)
	if err != nil {
		t.Errorf("Error creating account %s: %s", acc.ID, err)
		t.Fail()
	} else {
		if createdAccount.ID != acc.ID {
			t.Error("Expected account ID", acc.ID, "received", createdAccount.ID)
			t.Fail()
		}
	}
	fetchedAccount, err := client.Fetch(acc.ID)
	if err != nil {
		t.Fatalf("Error fetching account %s: %s", acc.ID, err)
	} else if !reflect.DeepEqual(fetchedAccount.Attributes, acc.Attributes) {
		t.Errorf("Attributes of fetched account don't match the input account: %+v", fetchedAccount.Attributes)
		t.Fail()
	}
	s, _ := json.Marshal(fetchedAccount)
	t.Logf(string(s))
	err = client.Delete(fetchedAccount.ID, fetchedAccount.Version)
	if err != nil {
		t.Errorf("Error deleting account %s: %s", fetchedAccount.ID, err)
		t.Fail()
	}
}

func TestList(t *testing.T) {
	const (
		NACCOUNTS                = 1100             // Number of accounts, created to test lists. Has to be at least 1000, i.e. default page[size].
		DEFAULT_SERVICE_PAGESIZE = 1000             // Default page[size] for List request. TODO: Documentation says the default is 100, but it's 1000.
		PAGESIZE                 = lib.PageSize(10) // Page size for tests.
	)
	if !test.IsLong() { // Only run when -tags=long flag was set.
		return
	}
	client, err := newClient()
	if err != nil {
		t.Fail()
	}
	if !client.Health() {
		t.Fatal("Can't connect to server")
	}
	cleanDatabase(t, client) // Make sure the database is empty before the test.
	acc := generateBasicAccount()
	accountIDs := []uuid.UUID{}
	for i := 1; i <= NACCOUNTS; i++ {
		acc.ID, _ = uuid.NewRandom()
		createdAccount, err := client.Create(acc)
		if err != nil {
			t.Fatalf("Error creating account %s", acc.ID.String())
		}
		accountIDs = append(accountIDs, createdAccount.ID)
	}

	// Check a list with the default pageNumber and pageSize.
	fetchedList, err := client.List(lib.PNNone, lib.PSNone)
	if err != nil {
		t.Errorf("Error retrieving the list with defaults: %v", err)
		t.Fail()
	} else if fetchedList != nil {
		if len(*fetchedList) != DEFAULT_SERVICE_PAGESIZE {
			t.Errorf("Expected a list of 1000 (default) accounts, retrieved %d.", len(*fetchedList))
			t.Fail()
		}
	} else {
		t.Error("Listing with defaults returned a nil list without an error.")
		t.Fail()
	}
	// Check the first, the last and a custom page.
	pageSize := PAGESIZE
	pageNumber := lib.PageNumber(2)
	fetchedList, err = client.List(lib.First, pageSize)
	if err != nil {
		t.Errorf("Error retrieving the list for the first page[number] and page[size]=%d: %v", int(pageSize), err)
		t.Fail()
	} else if fetchedList != nil {
		if len(*fetchedList) != int(pageSize) {
			t.Errorf("Expected a list of the first %d accounts, retrieved %d.", pageSize, len(*fetchedList))
			t.Fail()
		}
	} else {
		t.Error("Listing the first page returned a nil list without an error.")
		t.Fail()
	}
	fetchedList, err = client.List(lib.Last, pageSize)
	if err != nil {
		t.Errorf("Error retrieving the list for the last page[number] and page[size]=%d: %v", int(pageSize), err)
		t.Fail()
	} else if fetchedList != nil {
		if len(*fetchedList) != int(pageSize) {
			t.Errorf("Expected a list of the last %d accounts, retrieved %d.", pageSize, len(*fetchedList))
			t.Fail()
		}
	} else {
		t.Error("Listing the last page returned a nil list without an error.")
		t.Fail()
	}

	fetchedList, err = client.List(pageNumber, pageSize)
	if err != nil {
		t.Errorf("Error retrieving the list for page[number]=%d and page[size]=%d: %v", int(pageNumber), int(pageSize), err)
		t.Fail()
	} else if fetchedList != nil {
		if len(*fetchedList) != int(pageSize) {
			t.Errorf("Expected a list of %d accounts, retrieved %d.", pageSize, len(*fetchedList))
			t.Fail()
		}
	} else {
		t.Errorf("Listing page %d returned a nil list without an error.", int(pageNumber))
		t.Fail()
	}
	// TODO: verify if there is any assurance that the list of accounts is returned in the same order as they were created.
	for i, id := range accountIDs[int(pageNumber)*int(pageSize) : (int(pageNumber)+1)*int(pageSize)-1] {
		if id != (*fetchedList)[i].ID {
			t.Errorf("Account #%d in the list expected with id %s, got %s", i, id, (*fetchedList)[i].ID)
			t.Fail()
		}
	}
	// Check an out-of-range response. Retrieval of an account beyond the last should succeed but return an empty list.
	fetchedList, err = client.List(lib.PageNumber(NACCOUNTS+1), lib.PageSize(1))
	if err != nil {
		t.Errorf("Error retrieving the list for out-of-range page: %v", err)
		t.Fail()
	}
	if len(*fetchedList) > 0 {
		t.Errorf("Expected an empty list of accounts, retrieved %d accounts.", len(*fetchedList))
		t.Fail()
	}
	// Check calls with invalid arguments.
	fetchedList, err = client.List(lib.PageNumber(-99), lib.PageSize(1))
	if err == nil {
		t.Errorf("Calling List with invalid page number should return an error.")
		t.Fail()
	}
	if !lib.IsErrorInvalidArgument(err) {
		t.Errorf("Calling List with invalid page number should return an ErrorInvalidArgument.")
		t.Fail()
	}
	if fetchedList != nil {
		t.Errorf("Failed function call should return a nil list of accounts.")
		t.Fail()
	}
	// Cleanup.
	for _, id := range accountIDs {
		err := client.Delete(id, acc.Version)
		if err != nil {
			t.Errorf("Error deleting account %s", id)
			t.Fail()
		}
	}
	// The database should be empy after the cleanup.
	fetchedList, err = client.List(lib.PNNone, lib.PSNone)
	if err != nil {
		t.Errorf("Error retrieving the list of any leftover accounts: %v", err)
		t.Fail()
	}
	if fetchedList != nil && len(*fetchedList) != 0 {
		t.Errorf("The list is not empty after cleanup, %d leftover accounts.", len(*fetchedList))
		t.Fail()
		// Do the cleanup of leftover accounts.
		cleanDatabase(t, client)
	}
}

func TestParallelRequests(t *testing.T) {
	const (
		NACCOUNTS = 300
	)
	if !test.IsLong() {
		return
	}
	// Create a new client with long timeout.
	server := os.Getenv("APISERVICE")
	if server == "" {
		server = SERVER
	}
	client, err := account.New(account.Config{
		Server:             server,
		MaxConnections:     TestMaxConnections,
		MaxIdleConnections: TestMaxIdleConnections,
		Timeout:            TestLongTimeout,
	})

	if err != nil {
		t.Fail()
	}
	if !client.Health() {
		t.Fatal("Can't connect to server")
	}
	cleanDatabase(t, client) // Make sure the database is empty before the test.
	accountIDs := []uuid.UUID{}
	accIDchan := make(chan uuid.UUID, NACCOUNTS)
	for i := 0; i < NACCOUNTS; i++ {
		acc := generateBasicAccount()
		acc.ID, _ = uuid.NewRandom()
		go func(a *data.Account, accIDs chan uuid.UUID) {
			createdAccount, err := client.Create(a)
			if err != nil {
				t.Errorf("Error creating account %s: %v", acc.ID.String(), err)
				t.Fail()
				accIDs <- a.ID // Write id into channel to avoid hangup/timeout panic when collecting the data.
			} else {
				accIDs <- createdAccount.ID
			}
		}(acc, accIDchan)
	}
	for i := 0; i < NACCOUNTS; i++ {
		accountIDs = append(accountIDs, <-accIDchan)
	}

	// There should be NACCOUNTS accounts in the database now.
	fetchedList, err := client.List(lib.PNNone, lib.PageSize(NACCOUNTS))
	if err != nil {
		t.Errorf("Error retrieving the accounts, created by parallel statements: %v", err)
		t.Fail()
	} else if fetchedList != nil {
		if len(*fetchedList) != NACCOUNTS {
			t.Errorf("Expected a list of %d accounts, retrieved %d.", NACCOUNTS, len(*fetchedList))
			t.Fail()
		}
	} else {
		t.Error("Listing of accounts returned a nil list without an error.")
		t.Fail()
	}

	// Now delete all accounts in parallel and check if the database is empty.
	wg := sync.WaitGroup{}
	for i := 0; i < NACCOUNTS; i++ {
		wg.Add(1)
		go func(id uuid.UUID) {
			err := client.Delete(id, 0) // All created accounts have version=0.
			if err != nil {
				t.Errorf("Error deleting account in parallel %s: %v", id.String(), err)
				t.Fail()
			}
			wg.Done()
		}(accountIDs[i])
	}
	wg.Wait()

	// There should be no accounts in the database left.
	fetchedList, err = client.List(lib.PNNone, lib.PageSize(NACCOUNTS))
	if err != nil {
		t.Errorf("Error retrieving the leftover accounts after parallel delete: %v", err)
		t.Fail()
	} else if fetchedList != nil {
		if len(*fetchedList) != 0 {
			t.Errorf("Expected an empty list of accounts, retrieved %d.", len(*fetchedList))
			t.Fail()
		}
	} else {
		t.Error("Listing of accounts returned a nil list without an error.")
		t.Fail()
	}
	cleanDatabase(t, client)
}

// cleanDatabase deletes all accounts from the database.
func cleanDatabase(t *testing.T, client *account.Client) {
	for { // As there may be more accounts than the default pageSize, loop until all accounts are deleted.
		fetchedList, err := client.List(lib.PNNone, lib.PSNone)
		if err != nil {
			return
		}
		if fetchedList == nil { // This should never happen, if err==nil, fetchedList should not be nil.
			return
		}
		if len(*fetchedList) == 0 { // Empty database, success.
			return
		}
		t.Logf("Cleanup of %d accounts from the database", len(*fetchedList))
		for _, acc := range *fetchedList {
			err := client.Delete(acc.ID, acc.Version)
			if err != nil {
				t.Errorf("Error deleting leftover account %s", acc.ID)
				t.Fail()
			}
		}
	}
}
