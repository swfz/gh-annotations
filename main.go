package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"log"
	"strconv"
)

type Run struct {
	WorkflowId   int    `json:"workflow_id"`
	JobsUrl      string `json:"jobs_url"`
	Id           int    `json:"id"`
	Event        string `json:"event"`
	DisplayTitle string `json:"display_title"`
	HeadBranch   string `json:"head_branch"`
	HtmlUrl      string `json:"html_url"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	Status       string `json:"status"`
}

type WorkflowRuns struct {
	TotalCount int   `json:"total_count"`
	Runs       []Run `json:"workflow_runs"`
}

type Job struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	CheckRunUrl string `json:"check_run_url"`
	Conclusion  string `json:"conclusion"`
	Status      string `json:"status"`
	HtmlUrl     string `json:"html_url"`
}

type WorkflowJobs struct {
	TotalCount int   `json:"total_count"`
	Jobs       []Job `json:"jobs"`
}

type Annotation struct {
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

type Record struct {
	Repository      string `json:"repository"`
	WorkflowName    string `json:"workflow_name"`
	WorkflowEvent   string `json:"workflow_event"`
	WorkflowPath    string `json:"workflow_path"`
	WorkflowUrl     string `json:"workflow_url"`
	JobName         string `json:"job_name"`
	JobConclusion   string `json:"job_conclusion"`
	AnnotationLevel string `json:"annotation_level"`
	Message         string `json:"message"`
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

func getRuns(client api.RESTClient, repository string) WorkflowRuns {
	var res map[string]interface{}
	path := "repos/" + repository + "/actions/runs"
	client.Get(path, &res)
	//fmt.Printf("%+v", res)

	jsonStr, _ := json.Marshal(res)
	var workflowRunsRes WorkflowRuns
	if err := json.Unmarshal([]byte(jsonStr), &workflowRunsRes); err != nil {
		panic(err)
	}

	return workflowRunsRes
}

func getJobs(client api.RESTClient, repository string, run Run) WorkflowJobs {
	var res map[string]interface{}
	path := "repos/" + repository + "/actions/runs/" + strconv.Itoa(run.Id) + "/jobs"
	client.Get(path, &res)

	jsonStr, _ := json.Marshal(res)
	var jobs WorkflowJobs
	if err := json.Unmarshal([]byte(jsonStr), &jobs); err != nil {
		panic(err)
	}

	return jobs
}

func getAnnotations(client api.RESTClient, repository string, job Job) []Annotation {
	var res []interface{}
	path := "repos/" + repository + "/check-runs/" + strconv.Itoa(job.Id) + "/annotations"
	client.Get(path, &res)

	jsonStr, _ := json.Marshal(res)
	var annotations []Annotation
	if err := json.Unmarshal([]byte(jsonStr), &annotations); err != nil {
		panic(err)
	}

	return annotations
}

func toRecord(repository string, run Run, job Job, annotation Annotation) Record {
	r := Record{
		Repository:      repository,
		WorkflowName:    run.Name,
		WorkflowEvent:   run.Event,
		WorkflowPath:    run.Path,
		WorkflowUrl:     run.HtmlUrl,
		JobName:         job.Name,
		JobConclusion:   job.Conclusion,
		AnnotationLevel: annotation.AnnotationLevel,
		Message:         annotation.Message,
	}

	return r
}

func main() {
	var options struct {
		repo string
	}

	flag.StringVar(&options.repo, "repo", "", "Repository Name eg) owner/repo")
	flag.Parse()

	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	var repositoryPath string

	if options.repo != "" {
		repositoryPath = options.repo
	} else {
		currentRepository, _ := gh.CurrentRepository()
		repositoryPath = currentRepository.Owner() + "/" + currentRepository.Name()
	}

	workflowRuns := getRuns(client, repositoryPath)
	latestRuns := latest(workflowRuns)

	var summary []Record

	//fmt.Printf("Repository: %s/%s\n", currentRepository.Owner(), currentRepository.Name())
	for _, run := range latestRuns {
		//fmt.Printf("Workflow(%s): %s(%s)\n", run.Event, run.Name, run.Path)

		jobs := getJobs(client, repositoryPath, run)
		for _, job := range jobs.Jobs {
			//fmt.Printf("\tJob name: %s, %s\n", job.Name, job.Conclusion)
			annotations := getAnnotations(client, repositoryPath, job)
			for _, annotation := range annotations {
				//fmt.Printf("\t\t%s: %s\n", annotation.AnnotationLevel, annotation.Message)
				r := toRecord(repositoryPath, run, job, annotation)
				summary = append(summary, r)
			}
		}
	}

	j, err := json.Marshal(summary)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", string(j))
}
