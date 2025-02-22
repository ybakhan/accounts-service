package accountclient

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var client = InitializeAccountClient("http://localhost:8080", "v1", 10*time.Second)

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
			createdAccount, err := client.Create(context.TODO(), account)
			assert.Equal(t, account.ID, createdAccount.ID)
			assert.Nil(t, err)
		})
	}
}

// verify duplicate accounts not created
func TestCreate_Conflict(t *testing.T) {
	account := testInput("./test_input/account.json")
	_, err := client.Create(context.TODO(), account)
	assert.Nil(t, err)

	createdAccount, err := client.Create(context.TODO(), account)
	assert.Equal(t, AccountData{}, createdAccount)
	assert.Equal(t, err, fmt.Errorf("account %s already exists", account.ID))
}

// verify account is deleted
func TestDelete(t *testing.T) {
	account := testInput("./test_input/account.json")
	account.ID = "ad27e265-9605-4b4b-a0e5-3003ea9cc4de"
	_, err := client.Create(context.TODO(), account)
	assert.Nil(t, err)

	err = client.Delete(context.TODO(), account.ID, "0")
	assert.Nil(t, err)
}

// verify account is not found
func TestDelete_Account_Not_Found(t *testing.T) {
	accountID := "ad27e265-9605-4b4b-a0e5-3003ea9cc4df"
	err := client.Delete(context.TODO(), accountID, "0")
	assert.Equal(t, err, fmt.Errorf("account %s not found", accountID))
}

// verify fetch account
func TestFetch(t *testing.T) {
	account := testInput("./test_input/account.json")
	account.ID = "ad27e265-9605-4b4b-a0e5-3003ea9cc4da"
	_, err := client.Create(context.TODO(), account)
	assert.Nil(t, err)

	fetchedAccount, err := client.Fetch(context.TODO(), account.ID)
	assert.Equal(t, account.ID, fetchedAccount.ID)
	assert.Nil(t, err)
}

// verify fetch account not found
func TestFetch_Account_Not_Found(t *testing.T) {
	accountID := "ad27e265-9605-4b4b-a0e5-3003ea9cc4df"
	account, err := client.Fetch(context.TODO(), accountID)
	assert.Equal(t, AccountData{}, account)
	assert.Equal(t, err, fmt.Errorf("account %s not found", accountID))

}

func testInput(file string) AccountData {
	bytes, _ := os.ReadFile(file)
	var account AccountData
	json.Unmarshal(bytes, &account)
	return account
}
