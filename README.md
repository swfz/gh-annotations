# gh-annotations

Extension to output the list of annotations from the recently executed Workflow in the current repository.

## Usage

```shell
gh annotations
```

### Optional Args

```shell
$ gh annotations --help
  -repo string
        Optional Repository Name eg) owner/repo
  -json bool 
        Output JSON Format
```

### Example

```shell
$ gh annotations -repo swfz/ngx-sample
Repository       Workflow  Event  Job       JobStartedAt          JobCompletedAt        Conclusion  AnnotationLevel  Message
swfz/ngx-sample  ci        push   prettier  2023-02-12T01:35:36Z  2023-02-12T01:38:13Z  success     warning          Node.js 12 actions are deprecated. Please update the following actions t...
swfz/ngx-sample  ci        push   prettier  2023-02-12T01:35:36Z  2023-02-12T01:38:13Z  success     warning          The `save-state` command is deprecated and will be disabled soon. Please...
swfz/ngx-sample  ci        push   lint      2023-02-12T01:35:35Z  2023-02-12T01:38:04Z  success     warning          Node.js 12 actions are deprecated. Please update the following actions t...
swfz/ngx-sample  ci        push   lint      2023-02-12T01:35:35Z  2023-02-12T01:38:04Z  success     warning          The `save-state` command is deprecated and will be disabled soon. Please...
swfz/ngx-sample  ci        push   test      2023-02-12T01:35:35Z  2023-02-12T01:37:32Z  success     warning          Node.js 12 actions are deprecated. Please update the following actions t...
swfz/ngx-sample  ci        push   test      2023-02-12T01:35:35Z  2023-02-12T01:37:32Z  success     warning          The `save-state` command is deprecated and will be disabled soon. Please...
```

```shell
$ gh annotations -json | jq
[
  {
    "repository": "swfz/gh-annotations",
    "workflow_name": "release",
    "workflow_event": "push",
    "workflow_path": ".github/workflows/release.yml",
    "workflow_url": "https://github.com/swfz/gh-annotations/actions/runs/3371495473",
    "workflow_run_started_at": "2022-11-02T02:31:36Z",
    "workflow_created_at": "2022-11-02T02:31:36Z",
    "workflow_updated_at": "2022-11-02T02:32:13Z",
    "job_name": "release",
    "job_conclusion": "success",
    "job_started_at": "2022-11-02T02:31:36Z",
    "job_completed_at": "2022-11-02T02:32:13Z",
    "annotation_level": "warning",
    "message": "Node.js 12 actions are deprecated. For more information see: https://github.blog/changelog/2022-09-22-github-actions-all-actions-will-begin-running-on-node16-instead-of-node12/. Please update the following actions to use Node.js 16: actions/checkout, actions/checkout"
  }
]
```

## Install

```shell
gh extension install swfz/gh-annotations
```
