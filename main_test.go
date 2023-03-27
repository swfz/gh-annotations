package main

import (
	"github.com/cli/cli/v2/pkg/httpmock"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/cli/go-gh/pkg/api"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type Workflow struct {
	Id   int
	Name string
}
type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	//BaseRepo   func() (ghrepo.Interface, error)

	PlainOutput bool

	All   bool
	Limit int
}

func TestHoge(t *testing.T) {
	workflowRuns := []Run{
		{
			WorkflowId:   1,
			JobsUrl:      "https://example.com/jobs/1",
			Id:           1001,
			Event:        "push",
			DisplayTitle: "Sample Run",
			HeadBranch:   "main",
			HtmlUrl:      "https://example.com/runs/1001",
			Name:         "Run 1001",
			Path:         "/path/to/run",
			Status:       "completed",
			StartedAt:    "2023-03-20T10:00:00Z",
			CreatedAt:    "2023-03-20T09:50:00Z",
			UpdatedAt:    "2023-03-20T10:05:00Z",
		},
		//{
		//	WorkflowId:   2,
		//	JobsUrl:      "https://example.com/jobs/1",
		//	Id:           1002,
		//	Event:        "push",
		//	DisplayTitle: "Sample Run2",
		//	HeadBranch:   "main",
		//	HtmlUrl:      "https://example.com/runs/1001",
		//	Name:         "Run 1002",
		//	Path:         "/path/to/run",
		//	Status:       "completed",
		//	StartedAt:    "2023-03-20T10:00:00Z",
		//	CreatedAt:    "2023-03-20T09:50:00Z",
		//	UpdatedAt:    "2023-03-20T10:05:00Z",
		//},
	}

	workflowRunPayload := WorkflowRuns{
		TotalCount: 1,
		Runs:       workflowRuns,
	}

	workflowJobs := []Job{
		{
			Id:          2001,
			Name:        "Sample Job",
			CheckRunUrl: "https://example.com/check_runs/2001",
			Conclusion:  "success",
			StartedAt:   "2023-03-20T10:00:00Z",
			CompletedAt: "2023-03-20T10:02:00Z",
			Status:      "completed",
			HtmlUrl:     "https://example.com/jobs/2001",
		},
	}

	workflowJobsPayload := WorkflowJobs{
		TotalCount: 1,
		Jobs:       workflowJobs,
	}

	annotationsPayload := []Annotation{
		{
			Path:            "/path/to/file",
			BlobHref:        "https://example.com/blob/1",
			Title:           "Sample Annotation",
			Message:         "This is a sample annotation",
			AnnotationLevel: "warning",
			RawDetails:      "Detailed information",
			StartLine:       10,
			StartColumn:     5,
			EndLine:         12,
			EndColumn:       8,
		},
	}

	tests := []struct {
		name       string
		fuga       string
		options    Options
		stubs      func(*httpmock.Registry)
		wantOut    string
		wantErrOut string
	}{
		{
			name:    "hoge",
			fuga:    "fuga",
			options: Options{},
			wantOut: "Repository  Workflow  Event  Job  JobStartedAt  JobCompletedAt  Conclusion  AnnotationLevel  Message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &httpmock.Registry{}
			defer reg.Verify(t)
			if tt.stubs == nil {
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs"),
					httpmock.JSONResponse(workflowRunPayload),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs/1001/jobs"),
					httpmock.JSONResponse(workflowJobsPayload),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/check-runs/2001/annotations"),
					httpmock.JSONResponse(annotationsPayload),
				)
			} else {
				tt.stubs(reg)
			}

			var httpOptions = api.ClientOptions{
				Transport: reg,
			}

			tt.options.HttpOptions = httpOptions

			ios, _, stdout, stderr := iostreams.Test()
			ios.SetStdoutTTY(true)
			tt.options.IO = ios

			run(tt.options)

			assert.Equal(t, tt.wantOut, stdout.String())
			assert.Equal(t, tt.wantErrOut, stderr.String())
		})
	}
}
