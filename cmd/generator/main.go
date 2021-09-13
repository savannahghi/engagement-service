package main

import (
	"fmt"
	"os"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	engagementLib "github.com/savannahghi/engagementcore/pkg/engagement/presentation/graph/generated"
	"github.com/savannahghi/serverutils"
)

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	}

	engagementSources := engagementLib.Sources()

	err = api.Generate(cfg,
		api.AddPlugin(serverutils.NewImportPlugin(engagementSources, nil, true, "pkg/engagement/presentation")),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}
