package routes

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

type statusBody struct {
	BackEnd   string `json:"backend"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
}

func Resources(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got the request")
}

func RequestStatusHandler(w http.ResponseWriter, req *http.Request) {

	if req.Body != nil {
		var bodyBytes []byte
		bodyBytes, _ = ioutil.ReadAll(req.Body)
		bodyForStatus := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		decoder := json.NewDecoder(bodyForStatus)
		var statusData statusBody
		err := decoder.Decode(&statusData)

		if err != nil {
			fmt.Println("Error occurred while decoding the /status body ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	return 

}
