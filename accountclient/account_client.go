package accountclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type AccountsClient interface {
	Create(AccountData) (AccountData, error)
	Delete(accountID, version string) error
	Fetch(string) (AccountData, error)
}

type accountsClient struct {
	URL    string
	client *http.Client
}

type accountBody struct {
	Data AccountData `json:"data"`
}

func InitializeAccountsClient(baseURL, resource string) AccountsClient {
	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		err = fmt.Errorf("error intializing accounts client: %w", err)
		log.Fatal(err)
	}
	u.Path = resource

	return accountsClient{
		fmt.Sprintf("%v", u),
		&http.Client{Timeout: 10 * time.Second},
	}
}

func (ac accountsClient) Create(account AccountData) (AccountData, error) {
	accountJson, err := json.Marshal(accountBody{account})
	if err != nil {
		return AccountData{}, err
	}

	req, err := http.NewRequest(http.MethodPost, ac.URL, bytes.NewBuffer(accountJson))
	req.Header.Add("Accept", "application/json")
	if err != nil {
		return AccountData{}, err
	}

	resp, err := ac.client.Do(req)
	if err != nil {
		return AccountData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		err = fmt.Errorf("account %s already exists", account.ID)
		log.Print(err)
		return AccountData{}, err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return AccountData{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("account %s not created. response: %s", account.ID, string(respBytes))
		log.Print(err)
		return AccountData{}, err
	}

	var accountBody accountBody
	err = json.Unmarshal(respBytes, &accountBody)
	if err != nil {
		return AccountData{}, err
	}

	log.Printf("account %s created", account.ID)
	return accountBody.Data, nil
}

func (ac accountsClient) Delete(accountID, version string) error {
	deleteURL, err := url.JoinPath(ac.URL, accountID)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, deleteURL, nil)
	if err != nil {
		return err
	}

	query := req.URL.Query()
	query.Add("version", version)
	req.URL.RawQuery = query.Encode()

	resp, err := ac.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("account %s not found", accountID)
		log.Print(err)
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		respBytes, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("account %s not deleted. response: %s", accountID, string(respBytes))
		log.Print(err)
		return err
	}

	log.Printf("account %s deleted", accountID)
	return nil
}

func (ac accountsClient) Fetch(accountID string) (AccountData, error) {
	fetchURL, err := url.JoinPath(ac.URL, accountID)
	if err != nil {
		return AccountData{}, err
	}

	req, err := http.NewRequest(http.MethodGet, fetchURL, nil)
	if err != nil {
		return AccountData{}, err
	}

	resp, err := ac.client.Do(req)
	if err != nil {
		return AccountData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("account %s not fetched", accountID)
		log.Print(err)
		return AccountData{}, err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return AccountData{}, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("account %s not fetched. response: %s", accountID, string(respBytes))
		log.Print(err)
		return AccountData{}, err
	}

	var accountBody accountBody
	err = json.Unmarshal(respBytes, &accountBody)
	if err != nil {
		return AccountData{}, err
	}

	log.Printf("account %s fetched", accountID)
	return accountBody.Data, nil
}
