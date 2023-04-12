package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/fatih/color"
	"net/http"
	"os"
	"sync"
	"time"
)

type Leader struct {
	Client
	Password string
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
func (l *Leader) ListenChangeView(changedDataChan chan CampAPI, logChan chan string) {
	defer func() {
		if r := recover(); r != nil {
			color.Red("Dispatcher server is stopped")
			os.Exit(0)
		}
	}()
	prevCamp := l.GetCamp()
	changedDataChan <- prevCamp
	for {
		camp := l.GetCamp()
		if yes, message := camp.Equals(prevCamp); !yes {
			select {
			case <-changedDataChan:
			default:
			}

			changedDataChan <- camp
			logChan <- message
			prevCamp = camp
		}
		time.Sleep(2 * time.Second)
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
		l.ListenChangeView(changedDataChan, logChan)
		wg.Done()
	}()
	go func() {
		l.StartLeaderView(changedDataChan, logChan)
		wg.Done()
	}()
	wg.Wait()
}
