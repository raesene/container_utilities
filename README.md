# Container Utilities

Some little programs that I've used an LLM to write to help with container tasks. Just putting them here in case they're useful.

To get them working, have a valid go environment and run `go build <filename>.go` in the directory. or `go run <filename>.go` to run them directly.

## Image Timestamps

Connects to Docker hub and lists the tags on a repository along with their last built date.

## Resourcer

Just a re-implementation of [this script](https://stackoverflow.com/a/51289417) in go. Connects to a Kubernetes cluster via curl and lists the resource, sub-resources and verbs supported by the cluster. Easiest way to get it working is to run `kubectl proxy` and then run the program with the IP:port of the proxy. It will connect to the proxy and list the resources.