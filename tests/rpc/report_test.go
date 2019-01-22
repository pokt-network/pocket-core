package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/rpc/relay"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"github.com/pokt-network/pocket-core/service"

	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestReport(t *testing.T) {
	report := service.Report{GID: "test", Message: "foo"}
	// Start server instance
	go http.ListenAndServe(":"+config.GetInstance().RRPCPort, shared.NewRouter(relay.Routes()))
	// @ Url
	u := "http://localhost:" + config.GetInstance().RRPCPort + "/v1/report"
	j, err := json.Marshal(report)
	if err != nil {
		t.Fatalf(err.Error())
	}
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(j))
	// hanlde error
	if err != nil {
		t.Fatalf("Cannot create post request " + err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	// setup http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Unable to do post request " + err.Error())
	}
	// get body of response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unable to unmarshal response: " + err.Error())
	}
	expectedBody := "\"Okay! The node has been successfully reported to our servers and will be reviewed! Thank you!\""
	fmt.Println("Expected Body:", expectedBody)
	fmt.Println("Received Body", string(body))
	if expectedBody != string(body) {
		log.Fatalf("Body is not as expected")
	}
	t.Log(string(body))
	b, err := ioutil.ReadFile(_const.REPORTFILENAME)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(string(b))
}
