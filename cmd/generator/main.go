package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
		api.AddPlugin(NewImportPlugin(engagementSources, nil, true, "pkg/engagement/presentation")),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}

// NewImportPlugin initializes a new import plugin
//
// early sources are the source files to add before loading the service schema
// late sources are the source files to add after loading the service schema
// generate is a flag to determine whether to generate schema files or not
// path represents the path to store the imported schema files the folder name is `exported`
func NewImportPlugin(earlySources, lateSources []*ast.Source, generate bool, path string) plugin.Plugin {

	p := &ImportPlugin{
		earlySources: earlySources,
		lateSources:  lateSources,
		generate:     generate,
	}

	if generate {
		p.directory = p.Directory(path)
	}

	return p
}

// ImportPlugin is a gqlgen plugin that hooks into the gqlgen code generation lifecycle
// and adds schema definations from an imported library
type ImportPlugin struct {
	// the additional sources i.e "graphql files"
	earlySources, lateSources []*ast.Source
	directory                 string
	generate                  bool
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

		if m.generate {
			m.GenerateSchemaFile(m.directory, source)
		}

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

		if m.generate {
			m.GenerateSchemaFile(m.directory, source)
		}
	}

	return &o
}

func (m *ImportPlugin) Directory(path string) string {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	dir = filepath.Join(dir, path, "imported")

	// remove the old generated files if they exist
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		err = os.RemoveAll(dir)
		if err != nil {
			log.Println(err)
		}
	}

	// create a new generated folder
	err = os.Mkdir(dir, 0755)
	if err != nil {
		log.Println(err)
	}

	return dir
}

func (m *ImportPlugin) GenerateSchemaFile(dir string, source *ast.Source) {

	fileName := filepath.Base(source.Name)
	file := filepath.Join(dir, fileName)

	f, err := os.Create(file)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	_, err = f.WriteString(source.Input)
	if err != nil {
		log.Println(err)
	}
}
