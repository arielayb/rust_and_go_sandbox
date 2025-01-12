package application

import (
	"encoding/json"
	"net/http"
)

type Application struct {
	dataPacket SafeStore
}

func (app *Application) postAlert(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "application/json")
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//lastID++
	//tasks[lastID] = task
	//w.WriteHeader(http.StatusCreated)
	//json.NewEncoder(w).Encode(task)

}
