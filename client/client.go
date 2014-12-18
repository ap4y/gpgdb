package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/openpgp"

	"github.com/ap4y/gpgdb/lib"
)

type Client struct {
	Host   string
	entity *lib.Entity
	http.Client
}

func NewClient(host string, entity *lib.Entity) *Client {
	return &Client{host, entity, http.Client{}}
}

func (c *Client) Put(key, value string) error {
	body, err := c.Encrypt(value)
	if err != nil {
		return err
	}

	request, err := lib.NewRequest("PUT", c.Host+"/v1/keys/"+key, body)
	if err != nil {
		return err
	}

	_, err = c.do(request)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Keys() ([]string, error) {
	request, err := lib.NewRequest("GET", c.Host+"/v1/keys", nil)
	if err != nil {
		return nil, err
	}

	response, err := c.do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Unable to get keys")
	}

	decoder := json.NewDecoder(response.Body)
	result := map[string][]string{}
	if err = decoder.Decode(&result); err != nil {
		return nil, err
	}

	return result["keys"], err
}

func (c *Client) Get(key string) ([]byte, error) {
	request, err := lib.NewRequest("GET", c.Host+"/v1/keys/"+key, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.do(request)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(response.Body)
	result := map[string]string{}
	if err = decoder.Decode(&result); err != nil {
		return nil, err
	}

	value, err := c.Decrypt([]byte(result["value"]))
	if err != nil {
		return nil, err
	}

	return value, err
}

func (c *Client) Delete(key string) error {
	request, err := lib.NewRequest("DELETE", c.Host+"/v1/keys/"+key, nil)
	if err != nil {
		return err
	}

	_, err = c.do(request)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Encrypt(value string) ([]byte, error) {
	buffer := new(bytes.Buffer)
	writer, err := openpgp.Encrypt(buffer, c.entity.EntityList, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	if _, err = writer.Write([]byte(value)); err != nil {
		return nil, err
	}

	if err = writer.Close(); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(buffer)
	if err != nil {
		return nil, err
	}

	base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(base64Text, data)

	return base64Text, nil
}

func (c *Client) Decrypt(value []byte) ([]byte, error) {
	reader := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer(value))
	md, err := openpgp.ReadMessage(reader, c.entity.EntityList, nil, nil)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) do(request *lib.Request) (*http.Response, error) {
	if err := request.Sign(c.entity); err != nil {
		return nil, err
	}

	response, err := c.Do(request.Request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		error, _ := ioutil.ReadAll(response.Body)
		message := fmt.Sprintf("Unable to write value, Code: %s, message: %s",
			response.Status, error)
		return nil, errors.New(message)
	}

	return response, err
}
