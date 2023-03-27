package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/tableprinter"
	"github.com/cli/go-gh/pkg/term"
	"io"
	"log"
	"net/http"
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
	StartedAt    string `json:"run_started_at"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
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
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`
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
	WorkflowStarted string `json:"workflow_run_started_at"`
	WorkflowCreated string `json:"workflow_created_at"`
	WorkflowUpdated string `json:"workflow_updated_at"`
	JobName         string `json:"job_name"`
	JobConclusion   string `json:"job_conclusion"`
	JobStarted      string `json:"job_started_at"`
	JobCompleted    string `json:"job_completed_at"`
	AnnotationLevel string `json:"annotation_level"`
	Message         string `json:"message"`
}

// Retrieve the latest runs from the given WorkflowRuns
func latestRuns(workflowRuns WorkflowRuns) []Run {
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

// Fetch workflow runs from the given repository
func fetchWorkflowRuns(client api.RESTClient, repository string) (WorkflowRuns, error) {
	var res map[string]interface{}
	path := "repos/" + repository + "/actions/runs"
	err := client.Get(path, &res)

	if err != nil {
		return WorkflowRuns{}, err
	}

	jsonStr, _ := json.Marshal(res)
	var workflowRunsRes WorkflowRuns
	err = json.Unmarshal([]byte(jsonStr), &workflowRunsRes)
	if err != nil {
		return WorkflowRuns{}, err
	}

	//fmt.Printf("%#v\n", workflowRunsRes)

	return workflowRunsRes, nil
}

// Fetch jobs for the given repository and run
func fetchJobs(client api.RESTClient, repository string, run Run) (WorkflowJobs, error) {
	var res map[string]interface{}
	path := "repos/" + repository + "/actions/runs/" + strconv.Itoa(run.Id) + "/jobs"
	err := client.Get(path, &res)

	if err != nil {
		return WorkflowJobs{}, err
	}

	jsonStr, _ := json.Marshal(res)
	var jobs WorkflowJobs
	err = json.Unmarshal([]byte(jsonStr), &jobs)
	if err != nil {
		return WorkflowJobs{}, err
	}

	fmt.Printf("%#v\n", jobs)

	return jobs, nil
}

// Fetch annotations for the given repository and job
func fetchAnnotations(client api.RESTClient, repository string, job Job) ([]Annotation, error) {
	var res []interface{}
	path := "repos/" + repository + "/check-runs/" + strconv.Itoa(job.Id) + "/annotations"
	err := client.Get(path, &res)

	if err != nil {
		return []Annotation{}, err
	}

	jsonStr, _ := json.Marshal(res)
	var annotations []Annotation
	err = json.Unmarshal([]byte(jsonStr), &annotations)
	if err != nil {
		return []Annotation{}, err
	}

	fmt.Printf("%#v\n", annotations)

	return annotations, nil
}

// Convert the given run, job, and annotation to a Record
func toRecord(repository string, run Run, job Job, annotation Annotation) Record {
	r := Record{
		Repository:      repository,
		WorkflowName:    run.Name,
		WorkflowEvent:   run.Event,
		WorkflowPath:    run.Path,
		WorkflowUrl:     run.HtmlUrl,
		WorkflowStarted: run.StartedAt,
		WorkflowCreated: run.CreatedAt,
		WorkflowUpdated: run.UpdatedAt,
		JobName:         job.Name,
		JobStarted:      job.StartedAt,
		JobCompleted:    job.CompletedAt,
		JobConclusion:   job.Conclusion,
		AnnotationLevel: annotation.AnnotationLevel,
		Message:         annotation.Message,
	}

	return r
}

type Options struct {
	IO          *iostreams.IOStreams
	HttpClient  func() (*http.Client, error)
	HttpOptions api.ClientOptions
	repo        string
	json        bool
}

func run(options Options) {
	client, err := gh.RESTClient(&options.HttpOptions)
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

	workflowRuns, err := fetchWorkflowRuns(client, repositoryPath)
	if err != nil {
		log.Fatal(err)
	}
	latestRuns := latestRuns(workflowRuns)

	var summary []Record

	for _, run := range latestRuns {
		// Fetch jobs for the given run
		jobs, err := fetchJobs(client, repositoryPath, run)
		if err != nil {
			log.Fatal(err)
		}

		for _, job := range jobs.Jobs {
			annotations, err := fetchAnnotations(client, repositoryPath, job)
			if err != nil {
				log.Fatal(err)
			}

			for _, annotation := range annotations {
				record := toRecord(repositoryPath, run, job, annotation)
				summary = append(summary, record)
			}
		}
	}

	if options.json {
		summaryJson, _ := json.MarshalIndent(summary, "", "  ")
		fmt.Println(string(summaryJson))
	} else {

		terminal := term.FromEnv()
		termWidth, _, _ := terminal.Size()
		var out io.Writer

		if options.IO != nil {
			out = options.IO.Out
		} else {
			out = terminal.Out()
		}

		tp := tableprinter.New(out, terminal.IsTerminalOutput(), termWidth)

		tp.AddField("Repository")
		tp.AddField("Workflow")
		tp.AddField("Event")
		tp.AddField("Job")
		tp.AddField("JobStartedAt")
		tp.AddField("JobCompletedAt")
		tp.AddField("Conclusion")
		tp.AddField("AnnotationLevel")
		tp.AddField("Message")
		tp.EndRow()

		for _, row := range summary {
			tp.AddField(row.Repository)
			tp.AddField(row.WorkflowName)
			tp.AddField(row.WorkflowEvent)
			tp.AddField(row.JobName)
			tp.AddField(row.JobStarted)
			tp.AddField(row.JobCompleted)
			tp.AddField(row.JobConclusion)
			tp.AddField(row.AnnotationLevel)
			tp.AddField(row.Message)
			tp.EndRow()
		}

		tp.Render()
	}
}

func main() {
	var options Options

	flag.StringVar(&options.repo, "repo", "", "Repository Name eg) owner/repo")
	flag.BoolVar(&options.json, "json", false, "Output JSON")
	flag.Parse()

	run(options)
}
