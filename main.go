package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"doubtnut.com/checkTheseOut/common"
	"doubtnut.com/checkTheseOut/config"
	log "github.com/sirupsen/logrus"
)

var conf *config.SystemSettings

var sessions = map[common.UserID](chan int){}
var userData = map[common.UserID]([]common.Question){}

func init() {
	conf = config.Load("config.json")

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(conf.SetLogLevel())
}

func inactivityStopWatch(userID common.UserID) {
	var timeout = time.Duration(conf.Functional.InactivityTimeInSec) * time.Second
	var timer = time.NewTimer(timeout)
	var resetTimer = sessions[userID]

	for alive := true; alive; {
		select {
		case <-resetTimer:
			log.Info("Timer of user ", userID, " has been reset")
			timer.Reset(timeout)
		case <-timer.C:
			alive = false
			delete(sessions, userID)
			log.Info(userData[userID])
		}
	}
}

func generatePdf(w http.ResponseWriter, r *http.Request) {
	var request common.RequestFormat
	requestFormatError := json.NewDecoder(r.Body).Decode(&request)

	if requestFormatError != nil {
		log.Error("Failed to marshall request ", requestFormatError)
		w.WriteHeader(http.StatusBadRequest)

	} else {
		w.WriteHeader(http.StatusOK)
		if inputChannel, ok := sessions[request.UserID]; ok {
			delete(userData, request.UserID)
			userData[request.UserID] = request.Questions
			inputChannel <- 1 // to reset the timer

		} else {
			sessions[request.UserID] = make(chan int)
			userData[request.UserID] = request.Questions
			go inactivityStopWatch(request.UserID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
}

func main() {

	listenPort := conf.Service.Port
	apiVersion := conf.Service.APIVersion

	http.HandleFunc("/"+apiVersion+"/generatePdf", generatePdf)
	log.Info("HTTP Server is listening at ", conf.Service.Port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil))
}
