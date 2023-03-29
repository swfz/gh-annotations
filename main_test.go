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
	}{
		{
			name:    "1workflow, 1job, 1annotation",
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
Repository           Workflow  Event  Job         JobStartedAt          JobCompletedAt        Conclusion  AnnotationLevel  Message
swfz/gh-annotations  Run 1001  push   Sample Job  2023-03-20T10:00:00Z  2023-03-20T10:02:00Z  success     warning          This is a sample annotation
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
