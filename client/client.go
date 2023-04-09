package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Leader struct {
	Client
	Password string
}

type Client struct {
	http.Client      `json:"-"`
	Name             string `name:"name"`
	DispatcherServer string `json:"-"`
}

func (c *Client) GetCamp() (CampAPI, error) {
	//get request to dispatcher server /camp

	rq, err := http.NewRequest("GET", c.DispatcherServer+"/camp", nil)
	if err != nil {
		return CampAPI{}, err
	}
	do, err := c.Do(rq)
	defer do.Body.Close()

	if err != nil {
		var msg string
		_ = json.NewDecoder(do.Body).Decode(&msg)
		return CampAPI{}, errors.New(msg)
	}
	var camp CampAPI
	err = json.NewDecoder(do.Body).Decode(&camp)
	if err != nil {
		return CampAPI{}, err
	}
	return camp, nil
}

func (c *Client) JoinCamp() error {
	//post request to dispatcher server /camp/
	//with body {name: c.Name}
	rq, err := http.NewRequest("POST", c.DispatcherServer+"/camp/", nil)
	if err != nil {
		return err
	}
	do, err := c.Do(rq)
	defer do.Body.Close()
	if err != nil {
		return err
	}
	if do.StatusCode == http.StatusOK {
		return nil
	} else {
		//return the body as error message
		var msg string
		_ = json.NewDecoder(do.Body).Decode(&msg)
		return errors.New(msg)
	}
}

func (c *Leader) RemoveFromCamp(password string) error {
	//delete request to dispatcher server /camp/
	//with body {name: c.Name}
	rq, err := http.NewRequest("DELETE", c.DispatcherServer+"/camp/", nil)
	rq.Header.Add("Authorization", password)

	if err != nil {
		return err
	}
	do, err := c.Do(rq)
	defer do.Body.Close()
	if err != nil {
		return err
	}
	if do.StatusCode == http.StatusOK {
		return nil
	} else {
		//return the body as error message
		var msg string
		_ = json.NewDecoder(do.Body).Decode(&msg)
		return errors.New(msg)
	}
}

func (c *Leader) UpdateCampSettings(settings CampSettings, password string) error {
	//put request to dispatcher server /camp/
	//convert settings to json and put it in the body
	jr, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	rq, err := http.NewRequest("PUT", c.DispatcherServer+"/camp/", bytes.NewReader(jr))
	if err != nil {
		return err
	}
	rq.Header.Add("Authorization", password)

	do, err := c.Do(rq)
	defer do.Body.Close()
	if err != nil {
		return err
	}
	if do.StatusCode == http.StatusOK {
		return nil
	} else {
		var msg string
		_ = json.NewDecoder(do.Body).Decode(&msg)
		return errors.New(msg)
	}
}
