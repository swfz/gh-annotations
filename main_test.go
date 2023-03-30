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
			name: "1workflow, 1run, 2job, 2annotation",
			skip: true,
		},
		{
			name: "1workflow, 2run, 2job, 1annotation",
			skip: true,
		},
		{
			name: "2workflow, 2run, 2job, 2annotation",
			skip: true,
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
			name: "json output",
			skip: true,
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
