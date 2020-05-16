package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"doubtnut.com/checkTheseOut/common"
	"doubtnut.com/checkTheseOut/config"
	"github.com/jung-kurt/gofpdf"
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

func toPdf(userID common.UserID) {
	var filename string
	filename = "By_" + string(userID) + "_at_" + time.Now().Format("2006-01-02 15:04:05") + ".pdf"

	var pdf = gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	for i := 0; i < len(userData[userID]); i++ {
		pdf.CellFormat(190, 7, "Question "+": "+string(userData[userID][i]), "0", i+1, "AL", false, 0, "")
	}
	log.Info(pdf.OutputFileAndClose("pdfLogs/" + filename))
}

func inactivityStopWatch(userID common.UserID) {
	var timeout = time.Duration(conf.Functional.InactivityTimeInSec) * time.Minute
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
			go toPdf(userID)

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
