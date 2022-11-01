package main

import (
	"encoding/json"
	"fmt"
	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/repository"
	"strconv"
)

type Run struct {
	WorkflowId int    `json:"workflow_id"`
	JobsUrl    string `json:"jobs_url"`
	Id         int    `json:"id"`
}

type WorkflowRuns struct {
	TotalCount int   `json:"total_count"`
	Runs       []Run `json:"workflow_runs"`
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
	for _, run := range workflowRuns.Runs {
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

func getRuns(client api.RESTClient, repository repository.Repository) WorkflowRuns {
	var res map[string]interface{}
	path := "repos/" + repository.Owner() + "/" + repository.Name() + "/actions/runs"
	client.Get(path, &res)
	//fmt.Printf("%+v", res)

	jsonStr, _ := json.Marshal(res)
	var workflowRunsRes WorkflowRuns
	if err := json.Unmarshal([]byte(jsonStr), &workflowRunsRes); err != nil {
		panic(err)
	}

	return workflowRunsRes
}

func getJobs(client api.RESTClient, repository repository.Repository, run Run) WorkflowJobs {
	var res map[string]interface{}
	path := "repos/" + repository.Owner() + "/" + repository.Name() + "/actions/runs/" + strconv.Itoa(run.Id) + "/jobs"
	client.Get(path, &res)

	jsonStr, _ := json.Marshal(res)
	var jobs WorkflowJobs
	if err := json.Unmarshal([]byte(jsonStr), &jobs); err != nil {
		panic(err)
	}

	return jobs
}

func getAnnotations(client api.RESTClient, repository repository.Repository, job Job) Annotations {
	var res []interface{}
	path := "repos/" + repository.Owner() + "/" + repository.Name() + "/check-runs/" + strconv.Itoa(job.Id) + "/annotations"
	client.Get(path, &res)

	jsonStr, _ := json.Marshal(res)
	var annotations Annotations
	if err := json.Unmarshal([]byte(jsonStr), &annotations); err != nil {
		panic(err)
	}

	return annotations
}

func main() {
	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	currentRepository, _ := gh.CurrentRepository()

	//fmt.Printf("%+v\n", workflowRunsRes)

	workflowRuns := getRuns(client, currentRepository)
	latestRuns := latest(workflowRuns)
	fmt.Printf("%+v\n", latestRuns)

	for _, run := range latestRuns {
		jobs := getJobs(client, currentRepository, run)

		for _, job := range jobs.Jobs {
			annotations := getAnnotations(client, currentRepository, job)
			fmt.Print("\n===========================================\n")
			fmt.Printf("%+v\n", annotations)
		}
	}
}
