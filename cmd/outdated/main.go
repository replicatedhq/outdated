package main

import "github.com/replicatedhq/outdated/cmd/outdated/cli"

// Required to use this tool with GCP
import _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

func main() {
	cli.InitAndExecute()
}
