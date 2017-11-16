package jc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type SystemClient struct {
	Config     map[string]interface{}
	ConfigPath string

	HttpClient    http.Client
	ClientKey     ClientPrivateKey
	ClientKeyPath string
	SystemKey     string
}

// Returns a new system client with the loaded configuration file + client key, ready to sign and send API requests
func NewSystemClient(configPath, clientKeyPath string) (client SystemClient, err error) {
	if configPath == "" {
		configPath = "/opt/jc/jcagent.conf"
	}

	if clientKeyPath == "" {
		clientKeyPath = "/opt/jc/client.key"
	}

	client = SystemClient{
		ConfigPath:    configPath,
		ClientKeyPath: clientKeyPath,
	}

	err = client.LoadConfig()
	if err != nil {
		return
	}

	err = client.LoadPrivateKey()

	client.HttpClient = http.Client{}

	return
}

func (c *SystemClient) AssociateGroup(groupId string) (result string, err error) {
	endpoint := "/api/v2/systemgroups/" + groupId + "/members"
	json := "{ \"op\": \"add\", \"type\": \"system\", \"id\": \"" + c.SystemKey + "\" }"

	_, err := c.Do(endpoint, "POST", json)

	return "Successfully associated to " + groupId, err
}

// Sends an API request to the endpoint specified for the system key
func (c *SystemClient) Do(endpoint, httpMethod, body string) (resp *http.Response, err error) {
	if endpoint == "" {
		endpoint = "/api/systems/" + c.SystemKey
	}

	url := c.url(endpoint)

	req, err := http.NewRequest(httpMethod, url, bytes.NewReader([]byte(body)))

	time := getTime()

	requestSigHeader, err := c.getAuthSignature(time, httpMethod, endpoint)

	if err != nil {
		log.Println("Error retrieving signature for request", err)
		return
	}

	req.Header.Set("Authorization", requestSigHeader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Date", time)

	resp, err = c.HttpClient.Do(req)

	return
}

// Gets an Auth Signature header based off the client.key file and system key
func (c *SystemClient) getAuthSignature(time, httpMethod, endpoint string) (header string, err error) {
	header = "Signature "

	requestSig, err := c.ClientKey.SignatureForRequest(time, httpMethod, endpoint)

	if err != nil {
		return
	}

	requestSigMap := map[string]string{
		"keyId":     "system/" + c.SystemKey,
		"headers":   "request-line date",
		"algorithm": "rsa-sha256",
		"signature": requestSig,
	}

	for key, val := range requestSigMap {
		header = header + key + "=\"" + val + "\","
	}

	return
}

// Loads the jcagent.conf file from disk or wherever the systemClient is configured for
func (c *SystemClient) LoadConfig() (err error) {
	configFile, err := ioutil.ReadFile(c.ConfigPath)
	if err != nil {
		return
	}

	var config map[string]interface{}
	json.Unmarshal(configFile, &config)

	c.Config = config

	if c.Config["systemKey"] == nil {
		err = errors.New("config: missing systemKey")
	} else {
		c.SystemKey = c.Config["systemKey"].(string)
	}

	return
}

// Loads the client.key file from disk or wherever the systemClient is configured for
func (c *SystemClient) LoadPrivateKey() (err error) {
	pk, err := LoadClientPrivateKeyFromFile(c.ClientKeyPath)

	if err != nil {
		return
	}

	c.ClientKey = pk

	return
}

// Returns a JumpCloud API formatted url
func (c *SystemClient) url(uri string) (url string) {
	return "https://console.jumpcloud.com" + uri
}
