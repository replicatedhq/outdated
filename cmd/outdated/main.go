package main

import (
	"github.com/replicatedhq/outdated/cmd/outdated/cli"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	cli.InitAndExecute()
}
