# jx build controller

[![Documentation](https://godoc.org/github.com/jenkins-x-plugins/jx-build-controller?status.svg)](https://pkg.go.dev/mod/github.com/jenkins-x-plugins/jx-build-controller)
[![Go Report Card](https://goreportcard.com/badge/github.com/jenkins-x-plugins/jx-build-controller)](https://goreportcard.com/report/github.com/jenkins-x-plugins/jx-build-controller)
[![Releases](https://img.shields.io/github/release-pre/jenkins-x/helmboot.svg)](https://github.com/jenkins-x-plugins/jx-build-controller/releases)
[![LICENSE](https://img.shields.io/github/license/jenkins-x/helmboot.svg)](https://github.com/jenkins-x-plugins/jx-build-controller/blob/master/LICENSE)
[![Slack Status](https://img.shields.io/badge/slack-join_chat-white.svg?logo=slack&style=social)](https://slack.k8s.io/)

`jx-build-controller` watches for `PipelineRun` resources and updates the `PipelineActivity` custom resources to reflect their status.

Also `jx-build-controller` optionally stores resources and logs to long term storage (e.g. buckets).

## Getting Started

The build controller is included OOTB inside [Jenkins X](https://jenkins-x.io/docs/v3/)

Download the [jx-build-controller binary](https://github.com/jenkins-x-plugins/jx-build-controller/releases) for your operating system and add it to your `$PATH`.

## v0.5.24 nil-Stage backport

This branch is based on upstream v0.5.24, commit
`5a3d0179ee043f4d91805b25a8859e3074fd7854`. Dependency versions in `go.mod`
and `go.sum` are unchanged, including `jx-pipeline` v0.7.21 and Tekton
v0.41.0 with its `v1beta1` API. It does not adopt the newer Tekton `v1` API.

Dependencies and their license files are committed under `vendor/`.
The only local dependency code change is a nil-Stage guard in
`vendor/github.com/jenkins-x-plugins/jx-pipeline/pkg/pipelines/pipelines.go`,
in `addTaskRunsMessage`: Preview and Promote activity entries have no Stage,
so TaskRun message updates must skip them.

With Go 1.23.3, run:

```sh
CGO_ENABLED=0 go test -mod=vendor ./pkg/cmd/controller/tekton -count=1
make linux
```

The regression `TestTaskRunMessagesWithActivityStepVariants` exercises Stage,
Preview, Promote, and an empty activity through the public conversion API.
Preview and Promote panic without the guard; Stage message conversion and
non-Stage payloads are preserved with it.

Make and GoReleaser explicitly use `-mod=vendor`. The Dockerfile copies
`build/linux/jx-build-controller`, produced by `make linux`; it does not
compile or fetch dependencies. Plain Go builds also select this vendor tree
by default, but do not override it with `-mod=mod`.

**Re-vendoring overwrites the local fix.** After `go mod vendor`, reapply
`&& pa.Spec.Steps[k1].Stage != nil` to the existing bounds check in
`addTaskRunsMessage`, then rerun the regression and build before committing.
Do not upgrade dependencies or replace checksum entries to refresh this backport.

During preparation, the current upstream `jx-api` v4.8.0 tag resolved to a
different commit from the original module archive. The vendored copy matches
the existing `go.sum` and public Go checksum database, not that changed tag.
It was recovered from the original proxy-recorded commit
`aeefa907fc9497649587855686c45352bc8f7d4b`, with the module ZIP version prefix
restored and the result verified by Go against the unchanged v4.8.0 checksum
`h1:JSQSRdVZNeWc5RKlVDU+tde+G8np3mP8cmmZq4eFQeU=`.
Building from the committed vendor tree requires no module downloads.

## Commands

See the [jx-build-controller command reference](https://github.com/jenkins-x-plugins/jx-build-controller/blob/master/docs/cmd/jx-build-controller.md#jx-build-controller)
