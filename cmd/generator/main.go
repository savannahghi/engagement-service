package main

import (
	"fmt"
	"os"

	"github.com/99designs/gqlgen/plugin"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/vektah/gqlparser/v2/ast"

	openSource "github.com/savannahghi/engagement/pkg/engagement/presentation/graph/generated"
)

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	}

	err = api.Generate(cfg,
		api.AddPlugin(NewCustomPlugin()),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}

func NewCustomPlugin() plugin.Plugin {
	return &CustomPlugin{}
}

type CustomPlugin struct{}

func (m *CustomPlugin) Name() string {
	return "Custom Plugin"
}

func (m *CustomPlugin) MutateConfig(cfg *config.Config) error {
	return nil
}

func (m *CustomPlugin) InjectSourceEarly() *ast.Source {
	o := ast.Source{
		Name:    "opensource.graphql",
		Input:   "",
		BuiltIn: false,
	}

	for _, source := range openSource.Sources() {
		if source.Name == "federation/directives.graphql" || source.Name == "federation/entity.graphql" {
			continue
		}
		o.Input += source.Input
	}

	return &o
}

func (m *CustomPlugin) InjectSourceLate(schema *ast.Schema) *ast.Source {
	return nil
}
