package main

import (
	"fmt"
	"os"

	"github.com/99designs/gqlgen/plugin"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	engagementLib "github.com/savannahghi/engagement/pkg/engagement/presentation/graph/generated"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	}

	engagementSources := engagementLib.Sources()

	err = api.Generate(cfg,
		api.AddPlugin(NewImportPlugin(engagementSources, nil)),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}

// NewImportPlugin initializes a new import plugin
func NewImportPlugin(earlySources, lateSources []*ast.Source) plugin.Plugin {
	return &ImportPlugin{
		earlySources: earlySources,
		lateSources:  lateSources,
	}
}

// ImportPlugin is a gqlgen plugin that hooks into the gqlgen code generation lifecycle
// and adds schema definations from an imported library
type ImportPlugin struct {
	// the additional sources i.e "graphql files"
	earlySources, lateSources []*ast.Source
}

// Name is the name of the plugin
func (m *ImportPlugin) Name() string {
	return "import plugin"
}

// MutateConfig implements the ConfigMutator interface
func (m *ImportPlugin) MutateConfig(cfg *config.Config) error {
	return nil
}

// InjectSourceEarly is used to inject the library schema before loading the service schema.
func (m *ImportPlugin) InjectSourceEarly() *ast.Source {
	// check if there are sources
	if m.earlySources == nil {
		return nil
	}

	// initialize a graphql file that holds the imported schema as it's own source file
	o := ast.Source{
		Name:    "imported.graphql",
		Input:   "",
		BuiltIn: false,
	}

	for _, source := range m.earlySources {
		// federation directives and entities are already provided using the federation plugin
		// They should be skipped to avoid conflict with/from the federation plugin
		if source.Name == "federation/directives.graphql" || source.Name == "federation/entity.graphql" {
			continue
		}
		// Contents of the source file
		o.Input += source.Input
	}

	return &o
}

// InjectSourceLate is used to inject more sources after loading the service souces
func (m *ImportPlugin) InjectSourceLate(schema *ast.Schema) *ast.Source {
	// check if there are late sources
	if m.lateSources == nil {
		return nil
	}

	// initialize a graphql file that holds the imported schema as it's own source file
	o := ast.Source{
		Name:    "imported.graphql",
		Input:   "",
		BuiltIn: false,
	}

	for _, source := range m.earlySources {
		// federation directives and entities are already provided using the federation plugin
		// They should be skipped to avoid conflict with the federation one
		if source.Name == "federation/directives.graphql" || source.Name == "federation/entity.graphql" {
			continue
		}
		// Contents of the source file
		o.Input += source.Input
	}

	return &o
}
