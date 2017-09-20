package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type CallbackPostRequestBody struct {
	JobID             string `json:"job_id" bson:"job_id"`
	JobName           string `json:"job_name" bson:"job_name"`
	JobContainerID    string `json:"job_container_id" bson:"job_container_id"`
	JobInstanceID     string `json:"job_instance_id" bson:"job_instance_id"`
	State             int    `json:"state" bson:"state"`
	Status            int    `json:"status" bson:"status"`
	StatusDescription string `json:"status_description" bson:"status_description"`
}

func jobCallbackMethod(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	paramsID := params["id"]

	var jobInfo CallbackPostRequestBody
	if err := json.NewDecoder(req.Body).Decode(&jobInfo); err != nil {
		fmt.Println("Error decoding JSON: ", err)
		return
	}

	//DELETE THIS later
	fmt.Println("job info: ", jobInfo)

	fmt.Printf("Received job completion at: %v\n", time.Now())
	fmt.Printf("Params ID: %s\n", paramsID)
	fmt.Printf("Job ID: %s\nJob Name: %s\nJob Container ID: %s\nJob Instance ID: %s\n",
						jobInfo.JobID, jobInfo.JobName, jobInfo.JobContainerID, jobInfo.JobInstanceID)
	fmt.Printf("State: %d\nStatus: %d\nStatus Description: %s\n",
						jobInfo.State, jobInfo.Status, jobInfo.StatusDescription)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/jobcallback/{id}", jobCallbackMethod).Methods("POST")
	log.Fatal(http.ListenAndServe(":3007", router))
}
