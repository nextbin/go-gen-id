package main

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
}

func TestHttpGin(t *testing.T) {
	url := "http://localhost:11001/genId"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{"body": string(body)}).Info("request http-gin success")
}
