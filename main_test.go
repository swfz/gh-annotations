package main

import (
	"flag"
	"github.com/cli/cli/v2/pkg/httpmock"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
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
			options: Options{repo: "swfz/gh-annotations"},
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
			wantOut: strings.Join([]string{
				"Repository\tWorkflow\tEvent\tJob\tJobStartedAt\tJobCompletedAt\tConclusion\tAnnotationLevel\tMessage",
				"swfz/gh-annotations\tSample Workflow Run\tpush\tSample Job\t2023-03-20T10:00:00Z\t2023-03-20T10:02:00Z\tsuccess\twarning\tThis is a sample annotation",
				"",
			}, "\n"),
		},
		{
			name:    "1workflow, 1run, 1job, 0annotation",
			options: Options{repo: "swfz/gh-annotations"},
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
			wantOut: "Repository\tWorkflow\tEvent\tJob\tJobStartedAt\tJobCompletedAt\tConclusion\tAnnotationLevel\tMessage\n",
		},
		{
			name:    "1workflow no runs",
			options: Options{repo: "swfz/gh-annotations"},
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs"),
					httpmock.FileResponse("./fixtures/workflow_run_norun.json"),
				)
			},
			wantOut: "Repository\tWorkflow\tEvent\tJob\tJobStartedAt\tJobCompletedAt\tConclusion\tAnnotationLevel\tMessage\n",
		},
		{
			name:    "1workflow, 2run, 2job. last run has no annotation",
			options: Options{repo: "swfz/gh-annotations"},
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
			wantOut: "Repository\tWorkflow\tEvent\tJob\tJobStartedAt\tJobCompletedAt\tConclusion\tAnnotationLevel\tMessage\n",
		},
		{
			name:    "1workflow, 1run, 1job, 2annotation",
			options: Options{repo: "swfz/gh-annotations"},
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
			wantOut: strings.Join([]string{
				"Repository\tWorkflow\tEvent\tJob\tJobStartedAt\tJobCompletedAt\tConclusion\tAnnotationLevel\tMessage",
				"swfz/gh-annotations\tSample Workflow Run\tpush\tSample Job\t2023-03-20T10:00:00Z\t2023-03-20T10:02:00Z\tsuccess\twarning\tThis is a sample annotation",
				"swfz/gh-annotations\tSample Workflow Run\tpush\tSample Job\t2023-03-20T10:00:00Z\t2023-03-20T10:02:00Z\tsuccess\twarning\tannotation in line",
				"",
			}, "\n"),
		},
		{
			name:    "2workflow x 2run, 3job, 4annotation",
			options: Options{repo: "swfz/gh-annotations"},
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
			wantOut: strings.Join([]string{
				"Repository\tWorkflow\tEvent\tJob\tJobStartedAt\tJobCompletedAt\tConclusion\tAnnotationLevel\tMessage",
				"swfz/gh-annotations\tSample Workflow Run\tpush\tSample Job\t2023-03-20T10:00:00Z\t2023-03-20T10:02:00Z\tsuccess\twarning\tThis is a sample annotation",
				"swfz/gh-annotations\tAwesome Workflow\tpush\tAwesome First Job\t2023-02-20T10:00:00Z\t2023-02-20T10:02:00Z\tsuccess\twarning\tThis Method is deplicated",
				"swfz/gh-annotations\tAwesome Workflow\tpush\tAwesome First Job\t2023-02-20T10:00:00Z\t2023-02-20T10:02:00Z\tsuccess\twarning\tdeplicated",
				"swfz/gh-annotations\tAwesome Workflow\tpush\tAwesome Second Job\t2023-02-20T10:00:30Z\t2023-02-20T10:02:40Z\tfailure\tfailure\tProcess completed with exit code 1.",
				"",
			}, "\n"),
		},
		{
			name: "json output",
			options: Options{
				repo: "swfz/gh-annotations",
				json: true,
			},
			wantOut: strings.Join([]string{
				"[",
				"  {",
				`    "repository": "swfz/gh-annotations",`,
				`    "workflow_name": "Sample Workflow Run",`,
				`    "workflow_event": "push",`,
				`    "workflow_path": "/path/to/run",`,
				`    "workflow_url": "https://example.com/actions/runs/1001",`,
				`    "workflow_run_started_at": "",`,
				`    "workflow_created_at": "2023-03-20T09:50:00Z",`,
				`    "workflow_updated_at": "2023-03-20T10:05:00Z",`,
				`    "job_name": "Sample Job",`,
				`    "job_conclusion": "success",`,
				`    "job_started_at": "2023-03-20T10:00:00Z",`,
				`    "job_completed_at": "2023-03-20T10:02:00Z",`,
				`    "annotation_level": "warning",`,
				`    "message": "This is a sample annotation"`,
				"  }",
				"]",
				"",
			}, "\n"),
		},
		{
			name: "json output, empty value",
			options: Options{
				repo: "swfz/gh-annotations",
				json: true,
			},
			stubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("GET", "repos/swfz/gh-annotations/actions/runs"),
					httpmock.FileResponse("./fixtures/workflow_run_norun.json"),
				)
			},
			wantOut: "[]\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skipf("%s. case skipped.", tt.name)
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
				AuthToken: "dummy-token",
				Transport: reg,
			}

			tt.options.HttpOptions = httpOptions

			ios, _, stdout, stderr := iostreams.Test()
			ios.SetStdoutTTY(false)
			tt.options.IO = ios

			run(tt.options)

			assert.Equal(t, tt.wantOut, stdout.String())
			assert.Equal(t, tt.wantErrOut, stderr.String())
		})
	}
}

func Test_checkArgs(t *testing.T) {
	t.Run("位置引数が無ければ続行する", func(t *testing.T) {
		msg, ok := checkArgs("")
		assert.True(t, ok)
		assert.Equal(t, "", msg)
	})

	t.Run("位置引数があれば -repo での指定を促して止める", func(t *testing.T) {
		msg, ok := checkArgs("swfz/gh-annotations")
		assert.False(t, ok)
		assert.Equal(t, "Did you mean this?\n\tgh annotations -repo swfz/gh-annotations\n", msg)
	})
}

type exitCalled struct{ code int }

func Test_main_positionalArg(t *testing.T) {
	origArgs, origExit, origFlags := os.Args, osExit, flag.CommandLine
	t.Cleanup(func() { os.Args, osExit, flag.CommandLine = origArgs, origExit, origFlags })
	flag.CommandLine = flag.NewFlagSet("gh-annotations", flag.ContinueOnError)
	os.Args = []string{"gh-annotations", "swfz/gh-annotations"}
	osExit = func(code int) { panic(exitCalled{code: code}) }

	defer func() {
		r := recover()
		assert.Equal(t, exitCalled{code: 1}, r)
	}()
	main()
	t.Fatal("main should exit")
}

func newMockClient(t *testing.T, path string, responder httpmock.Responder) *api.RESTClient {
	t.Helper()
	reg := &httpmock.Registry{}
	t.Cleanup(func() { reg.Verify(t) })
	reg.Register(httpmock.REST("GET", path), responder)
	client, err := api.NewRESTClient(api.ClientOptions{AuthToken: "dummy-token", Transport: reg})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func Test_fetchErrors(t *testing.T) {
	const repo = "swfz/gh-annotations"
	runsPath := "repos/" + repo + "/actions/runs"
	jobsPath := "repos/" + repo + "/actions/runs/1001/jobs"
	annotationsPath := "repos/" + repo + "/check-runs/10001/annotations"

	tests := []struct {
		name      string
		path      string
		responder httpmock.Responder
		fetch     func(client *api.RESTClient) error
	}{
		{
			name:      "fetchWorkflowRuns: HTTP 500",
			path:      runsPath,
			responder: httpmock.StatusStringResponse(500, `{"message":"error"}`),
			fetch:     func(c *api.RESTClient) error { _, err := fetchWorkflowRuns(c, repo); return err },
		},
		{
			name:      "fetchWorkflowRuns: 型が合わない応答",
			path:      runsPath,
			responder: httpmock.StringResponse(`{"total_count":"many"}`),
			fetch:     func(c *api.RESTClient) error { _, err := fetchWorkflowRuns(c, repo); return err },
		},
		{
			name:      "fetchJobs: HTTP 500",
			path:      jobsPath,
			responder: httpmock.StatusStringResponse(500, `{"message":"error"}`),
			fetch:     func(c *api.RESTClient) error { _, err := fetchJobs(c, repo, Run{Id: 1001}); return err },
		},
		{
			name:      "fetchJobs: 型が合わない応答",
			path:      jobsPath,
			responder: httpmock.StringResponse(`{"total_count":"many"}`),
			fetch:     func(c *api.RESTClient) error { _, err := fetchJobs(c, repo, Run{Id: 1001}); return err },
		},
		{
			name:      "fetchAnnotations: HTTP 500",
			path:      annotationsPath,
			responder: httpmock.StatusStringResponse(500, `{"message":"error"}`),
			fetch:     func(c *api.RESTClient) error { _, err := fetchAnnotations(c, repo, Job{Id: 10001}); return err },
		},
		{
			name:      "fetchAnnotations: 型が合わない応答",
			path:      annotationsPath,
			responder: httpmock.StringResponse(`[{"start_line":"first"}]`),
			fetch:     func(c *api.RESTClient) error { _, err := fetchAnnotations(c, repo, Job{Id: 10001}); return err },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newMockClient(t, tt.path, tt.responder)
			assert.Error(t, tt.fetch(client))
		})
	}
}
