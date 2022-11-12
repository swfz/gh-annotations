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
        Repository Name eg) owner/repo
```

### Example

```shell
gh annotations | jq
[
  {
    "repository": "swfz/gh-annotations",
    "workflow_name": "release",
    "workflow_event": "push",
    "workflow_path": ".github/workflows/release.yml",
    "workflow_url": "https://github.com/swfz/gh-annotations/actions/runs/3371495473",
    "job_name": "release",
    "job_conclusion": "success",
    "annotation_level": "warning",
    "message": "Node.js 12 actions are deprecated. For more information see: https://github.blog/changelog/2022-09-22-github-actions-all-actions-will-begin-running-on-node16-instead-of-node12/. Please update the following actions to use Node.js 16: actions/checkout, actions/checkout"
  }
]
```

## Install

```shell
gh extension install swfz/gh-annotations
```
