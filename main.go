package main

import (
	"encoding/json"
	"fmt"
	"github.com/cli/go-gh"
	"strconv"
)

type Run struct {
	WorkflowId int    `json:"workflow_id"`
	JobsUrl    string `json:"jobs_url"`
	Id         int    `json:"id"`
}

type WorkflowRuns struct {
	TotalCount   int   `json:"total_count"`
	WorkflowRuns []Run `json:"workflow_runs"`
}

type Job struct {
	Id          int    `json:"id"`
	CheckRunUrl string `json:"check_run_url"`
	Conclusion  string `json:"conclusion"`
	Status      string `json:"status"`
	HtmlUrl     string `json:"html_url"`
}

type WorkflowJobs struct {
	TotalCount int   `json:"total_count"`
	Jobs       []Job `json:"jobs"`
}

type Annotations []struct {
	Path            string `json:"path"`
	BlobHref        string `json:"blob_href"`
	Title           string `json:"title"`
	Message         string `json:"message"`
	AnnotationLevel string `json:"annotation_level"`
	RawDetails      string `json:"raw_details"`
	StartLine       int    `json:"start_line"`
	StartColumn     int    `json:"start_column"`
	EndLine         int    `json:"end_line"`
	EndColumn       int    `json:"end_column"`
}

func latest(workflowRuns WorkflowRuns) []Run {
	var latestRuns []Run
	for _, run := range workflowRuns.WorkflowRuns {

		var existRun bool = false

		for _, latest := range latestRuns {
			if latest.WorkflowId == run.WorkflowId {
				existRun = true
				break
			}
		}
		if !existRun {
			latestRuns = append(latestRuns, run)
		}
	}

	return latestRuns
}

func main() {
	fmt.Println("hi world, this is the gh-annotations extension!")

	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	response := struct{ Login string }{}
	err = client.Get("user", &response)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("running as %s\n", response.Login)

	currentRepository, _ := gh.CurrentRepository()
	fmt.Printf("%+v", currentRepository)

	// run api request
	var res map[string]interface{}
	apiPath := "repos/" + currentRepository.Owner() + "/" + currentRepository.Name() + "/actions/runs"
	client.Get(apiPath, &res)
	//fmt.Printf("%+v", res)

	jsonStr, _ := json.Marshal(res)
	var workflowRunsRes WorkflowRuns
	if err := json.Unmarshal([]byte(jsonStr), &workflowRunsRes); err != nil {
		panic(err)
	}

	//fmt.Printf("%+v\n", workflowRunsRes)

	latestRuns := latest(workflowRunsRes)
	fmt.Printf("%+v\n", latestRuns)

	for _, run := range latestRuns {
		var jobsRes map[string]interface{}
		jobsPath := "repos/" + currentRepository.Owner() + "/" + currentRepository.Name() + "/actions/runs/" + strconv.Itoa(run.Id) + "/jobs"
		client.Get(jobsPath, &jobsRes)

		jobsJsonStr, _ := json.Marshal(jobsRes)
		var jobs WorkflowJobs
		if err := json.Unmarshal([]byte(jobsJsonStr), &jobs); err != nil {
			panic(err)
		}

		for _, job := range jobs.Jobs {
			var annotationRes []interface{}
			annotationPath := "repos/" + currentRepository.Owner() + "/" + currentRepository.Name() + "/check-runs/" + strconv.Itoa(job.Id) + "/annotations"
			fmt.Printf("%+v\n", annotationPath)
			fmt.Printf("%+v\n", job.CheckRunUrl)
			client.Get(annotationPath, &annotationRes)

			annotationJsonStr, _ := json.Marshal(annotationRes)
			var annotations Annotations
			if err := json.Unmarshal([]byte(annotationJsonStr), &annotations); err != nil {
				panic(err)
			}
			fmt.Printf("%+v\n", annotations)
		}
	}
}
