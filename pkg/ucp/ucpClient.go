package ucp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client - Is the basic Client struct
type Client struct {
	Username   string
	Password   string
	UCPURL     string
	IgnoreCert bool
}

var token string

// NewBasicAuthClient - Creates a basic client to connecto the UCP
func NewBasicAuthClient(username, password, url string, ignoreCert bool) *Client {
	return &Client{
		Username:   username,
		Password:   password,
		UCPURL:     url,
		IgnoreCert: ignoreCert,
	}
}

// Connect - Will attempt to connect to UCP
func (c *Client) Connect() error {
	// Add the /auth/log to the URL
	url := fmt.Sprintf("%s/auth/login", c.UCPURL)

	data := map[string]string{
		"username": c.Username,
		"password": c.Password,
	}
	b, err := json.Marshal(data)

	if err != nil {
		return err
	}

	response, err := c.postRequest(url, b)
	if err != nil {
		return err
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(response, &responseData)
	if err != nil {
		return err
	}
	if responseData["auth_token"] != nil {
		token = responseData["auth_token"].(string)
	} else {
		return fmt.Errorf("No Authorisation token returned")
	}
	return nil
}

// Disconnect -
func (c *Client) Disconnect() error {

	return nil
}

// POST data to the server and return the response as bytes
func (c *Client) postRequest(url string, d []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(d))
	if err != nil {
		return nil, err
	}
	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// GET data from the server and return the response as bytes
func (c *Client) getRequest(url string, d []byte) ([]byte, error) {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(d))
	if len(token) != 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	if err != nil {
		return nil, err
	}
	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.IgnoreCert},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}
