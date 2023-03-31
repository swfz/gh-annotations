package main

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/pkg/httpmock"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/cli/go-gh/pkg/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_run(t *testing.T) {
	tests := []struct {
		name       string
		options    Options
		stubs      func(*httpmock.Registry)
		wantOut    string
		wantErrOut string
		skip       bool
	}{
		{
			name:    "1workflow, 1run, 1job, 1annotation",
			options: Options{},
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs"),
					httpmock.FileResponse("./fixtures/workflow_run.json"),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs/1001/jobs"),
					httpmock.FileResponse("./fixtures/runs_1001_jobs.json"),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/check-runs/10001/annotations"),
					httpmock.FileResponse("./fixtures/check_runs_10001_annotations.json"),
				)
			},
			wantOut: heredoc.Doc(`
Repository           Workflow             Event  Job         JobStartedAt          JobCompletedAt        Conclusion  AnnotationLevel  Message
swfz/gh-annotations  Sample Workflow Run  push   Sample Job  2023-03-20T10:00:00Z  2023-03-20T10:02:00Z  success     warning          This is a sample annotation
`),
		},
		{
			name:    "1workflow, 1run, 1job, 0annotation",
			options: Options{},
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs"),
					httpmock.FileResponse("./fixtures/workflow_run.json"),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs/1001/jobs"),
					httpmock.FileResponse("./fixtures/runs_1001_jobs.json"),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/check-runs/10001/annotations"),
					httpmock.FileResponse("./fixtures/check_runs_10001_annotations_0.json"),
				)
			},
			wantOut: heredoc.Doc(`
Repository  Workflow  Event  Job  JobStartedAt  JobCompletedAt  Conclusion  AnnotationLevel  Message
`),
		},
		{
			name:    "1workflow no runs",
			options: Options{},
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs"),
					httpmock.FileResponse("./fixtures/workflow_run_norun.json"),
				)
			},
			wantOut: heredoc.Doc(`
Repository  Workflow  Event  Job  JobStartedAt  JobCompletedAt  Conclusion  AnnotationLevel  Message
`),
		},
		{
			name:    "1workflow, 2run, 2job. last run has no annotation",
			options: Options{},
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs"),
					httpmock.FileResponse("./fixtures/workflow_run_2run.json"),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs/1001/jobs"),
					httpmock.FileResponse("./fixtures/runs_1001_jobs.json"),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/check-runs/10001/annotations"),
					httpmock.FileResponse("./fixtures/check_runs_10001_annotations_0.json"),
				)
			},
			wantOut: heredoc.Doc(`
Repository  Workflow  Event  Job  JobStartedAt  JobCompletedAt  Conclusion  AnnotationLevel  Message
`),
		},
		{
			name:    "1workflow, 1run, 1job, 2annotation",
			options: Options{},
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs"),
					httpmock.FileResponse("./fixtures/workflow_run_2run.json"),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs/1001/jobs"),
					httpmock.FileResponse("./fixtures/runs_1001_jobs.json"),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/check-runs/10001/annotations"),
					httpmock.FileResponse("./fixtures/check_runs_10001_annotations_2annotations.json"),
				)
			},
			wantOut: heredoc.Doc(`
Repository           Workflow             Event  Job         JobStartedAt          JobCompletedAt        Conclusion  AnnotationLevel  Message
swfz/gh-annotations  Sample Workflow Run  push   Sample Job  2023-03-20T10:00:00Z  2023-03-20T10:02:00Z  success     warning          This is a sample annotation
swfz/gh-annotations  Sample Workflow Run  push   Sample Job  2023-03-20T10:00:00Z  2023-03-20T10:02:00Z  success     warning          annotation in line
`),
		},
		{
			name:    "2workflow x 2run, 3job, 4annotation",
			options: Options{},
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs"),
					httpmock.FileResponse("./fixtures/workflow_run_2x2run.json"),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs/1001/jobs"),
					httpmock.FileResponse("./fixtures/runs_1001_jobs.json"),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs/2001/jobs"),
					httpmock.FileResponse("./fixtures/runs_2001_jobs.json"),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/check-runs/10001/annotations"),
					httpmock.FileResponse("./fixtures/check_runs_10001_annotations.json"),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/check-runs/20001/annotations"),
					httpmock.FileResponse("./fixtures/check_runs_20001_annotations_2annotations.json"),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/check-runs/20002/annotations"),
					httpmock.FileResponse("./fixtures/check_runs_20002_annotations.json"),
				)
			},
			wantOut: heredoc.Doc(`
Repository           Workflow             Event  Job                 JobStartedAt          JobCompletedAt        Conclusion  AnnotationLevel  Message
swfz/gh-annotations  Sample Workflow Run  push   Sample Job          2023-03-20T10:00:00Z  2023-03-20T10:02:00Z  success     warning          This is a sample annotation
swfz/gh-annotations  Awesome Workflow     push   Awesome First Job   2023-02-20T10:00:00Z  2023-02-20T10:02:00Z  success     warning          This Method is deplicated
swfz/gh-annotations  Awesome Workflow     push   Awesome First Job   2023-02-20T10:00:00Z  2023-02-20T10:02:00Z  success     warning          deplicated
swfz/gh-annotations  Awesome Workflow     push   Awesome Second Job  2023-02-20T10:00:30Z  2023-02-20T10:02:40Z  failure     failure          Process completed with exit code 1.
`),
		},
		{
			name: "json output",
			options: Options{
				json: true,
			},
			wantOut: heredoc.Doc(`
[
  {
    "repository": "swfz/gh-annotations",
    "workflow_name": "Sample Workflow Run",
    "workflow_event": "push",
    "workflow_path": "/path/to/run",
    "workflow_url": "https://example.com/actions/runs/1001",
    "workflow_run_started_at": "",
    "workflow_created_at": "2023-03-20T09:50:00Z",
    "workflow_updated_at": "2023-03-20T10:05:00Z",
    "job_name": "Sample Job",
    "job_conclusion": "success",
    "job_started_at": "2023-03-20T10:00:00Z",
    "job_completed_at": "2023-03-20T10:02:00Z",
    "annotation_level": "warning",
    "message": "This is a sample annotation"
  }
]
`),
		},
		{
			name: "json output, empty value",
			options: Options{
				json: true,
			},
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs"),
					httpmock.FileResponse("./fixtures/workflow_run_norun.json"),
				)
			},
			wantOut: heredoc.Doc(`[]
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skipf(tt.name + ". case skipped.")
			}

			reg := &httpmock.Registry{}
			defer reg.Verify(t)
			if tt.stubs == nil {
				// TODO: すべてからレスポンスのほうがいいかも
				// TODO: 正規表現行けるなら指定してすべてのリクエストを差し替える
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs"),
					httpmock.FileResponse("./fixtures/workflow_run.json"),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs/1001/jobs"),
					httpmock.FileResponse("./fixtures/runs_1001_jobs.json"),
				)
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/check-runs/10001/annotations"),
					httpmock.FileResponse("./fixtures/check_runs_10001_annotations.json"),
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
