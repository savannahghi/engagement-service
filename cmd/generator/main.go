package main

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	eng_generated "github.com/savannahghi/engagement/pkg/engagement/presentation/graph/generated"
	"github.com/vektah/gqlparser/v2/ast"
)

func generate() error {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		log.Fatalf("can't load config: %s", err)
	}

	sources := []*ast.Source{}
	sources = append(sources, eng_generated.Sources()...)

	for _, src := range sources {
		// append all other sources apart from federation directives
		if src.Name != "federation/directives.graphql" {
			cfg.Sources = append(cfg.Sources, src)
		}
	}

	for _, src := range cfg.Sources {
		log.Println(src.Name)
	}

	if err = api.Generate(cfg); err != nil {
		return err
	}
	return nil
}

func main() {
	err := generate()
	if err != nil {
		log.Printf("failed to generate: %s", err)
		os.Exit(3)
	}
}
