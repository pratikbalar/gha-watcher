package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type RunningWorkflow struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	StartedAt string `json:"run_started_at"`
	CancelUrl string `json:"cancel_url"`
	ReRunUrl  string `json:"rerun_url"`
}

type RunningActionResponse struct {
	TotalCount   int               `json:"total_count"`
	WorkflowRuns []RunningWorkflow `json:"workflow_runs"`
}

func apiCall(method string, url string, body io.Reader) []byte {
	username := os.Getenv("GITHUB_USERNAME")
	token := os.Getenv("GITHUB_TOKEN")
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(username, token)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return data
}

func restartAction(flow *RunningWorkflow) {

	// Cancel workflow
	apiCall("POST", flow.CancelUrl, nil)

	// ReRun workflow
	apiCall("POST", flow.ReRunUrl, nil)

}

func ghActionCheck() {

	body := apiCall("GET", "https://api.github.com/repos/Improwised/apricot-2/actions/runs?status=in_progress", nil)

	fres := RunningActionResponse{}
	if err := json.Unmarshal(body, &fres); err != nil {
		log.Fatal(err)
	}

	for i, item := range fres.WorkflowRuns {
		t, e := time.Parse(time.RFC3339, item.StartedAt)
		if e != nil {
			panic(e)
		}
		fmt.Println(i, item.Name)
		// calculate diff
		diff := time.Now().Sub(t).Minutes()

		fmt.Println(t, diff)
		if diff > 12 {
			restartAction(&item)
		}

	}

}

func main() {
	godotenv.Load()
	ghActionCheck()
	// for {
	// 	time.Sleep(1 * time.Minute)
	// }
}
