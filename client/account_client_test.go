package client

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client = InitializeAccountsClient("http://localhost:8080", "/v1/organisation/accounts")

// verify account is created
func TestCreate(t *testing.T) {
	tests := map[string]struct {
		File string
	}{
		"without payee confirmation": {
			"./test_input/account_without_payee_confirmation.json",
		},
		"with payee confirmation": {
			"./test_input/account_with_payee_confirmation.json",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			account := testInput(test.File)
			createdAccount, err := client.Create(account)
			assert.NotEqual(t, AccountData{}, createdAccount)
			assert.Nil(t, err)
		})
	}
}

// verify duplicate accounts not created
func TestCreate_Conflict(t *testing.T) {
	account := testInput("./test_input/account.json")
	_, err := client.Create(account)
	assert.Nil(t, err)

	createdAccount, err := client.Create(account)
	assert.Equal(t, AccountData{}, createdAccount)
	assert.NotNil(t, err)
	assert.Equal(t, err, fmt.Errorf("account %s already exists", account.ID))
}

// verify account is deleted
func TestDelete(t *testing.T) {
	account := testInput("./test_input/account.json")
	account.ID = "ad27e265-9605-4b4b-a0e5-3003ea9cc4de"
	_, err := client.Create(account)
	assert.Nil(t, err)

	err = client.Delete(account.ID, "0")
	assert.Nil(t, err)
}

// verify account is not found
func TestDelete_Account_Not_Found(t *testing.T) {
	accountID := "ad27e265-9605-4b4b-a0e5-3003ea9cc4df"
	err := client.Delete(accountID, "0")
	assert.NotNil(t, err)
	assert.Equal(t, err, fmt.Errorf("account %s not found", accountID))
}

func testInput(file string) AccountData {
	bytes, _ := os.ReadFile(file)
	var account AccountData
	json.Unmarshal(bytes, &account)
	return account
}
