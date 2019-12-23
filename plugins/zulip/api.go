package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type client struct {
	address string
	email   string
	token   string
}

type Client interface {
	SendMessage(stream string, topic string, message string) (int, error)
	UpdateMessage(id int, message string) error
}

func NewClient(address string, botEmail string, botToken string) Client {
	return &client{
		address,
		botEmail,
		botToken,
	}
}

type MessageId struct {
	Id int `json:"id"`
}

func (client *client) SendMessage(stream string, topic string, message string) (int, error) {
	hc := http.Client{}

	form := url.Values{}
	form.Add("type", "stream")
	form.Add("to", stream)
	form.Add("subject", topic)
	form.Add("content", message)
	f := form.Encode()

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/messages?%s", client.address, f), nil)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(client.email, client.token)

	res, err := hc.Do(req)
	if err != nil {
		return -1, err
	}
	if res.StatusCode != 200 {
		return -1, errors.New(fmt.Sprintf("Err %d: %s", res.StatusCode, res.Body))
	}

	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return -1, err
	}

	var re MessageId
	err = json.Unmarshal(d, &re)
	if err != nil {
		return -1, err
	}

	return re.Id, nil
}

func (client *client) UpdateMessage(id int, message string) error {
	hc := http.Client{}

	form := url.Values{}
	form.Add("content", message)
	f := form.Encode()

	req, err := http.NewRequest("PATCH",
		fmt.Sprintf("%s/api/v1/messages/%d?%s", client.address, id, f), nil)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(client.email, client.token)

	res, err := hc.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Err %d: %s", res.StatusCode, res.Body))
	}

	return nil
}
