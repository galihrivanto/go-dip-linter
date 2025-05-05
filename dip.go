package linters

import (
	"go/ast"
	"go/types"
	"path/filepath"
	"strings"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("dip", New)
}

type Settings struct {
	Include      []string `json:"include" yaml:"include"`             // Filepath include patterns
	Exclude      []string `json:"exclude" yaml:"exclude"`             // Filepath exclude patterns
	NamePatterns []string `json:"name_patterns" yaml:"name_patterns"` // Constructor name patterns
}

type Plugin struct {
	settings Settings
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[Settings](settings)
	if err != nil {
		return nil, err
	}

	return &Plugin{settings: s}, nil
}

func (f *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "dip",
			Doc:  "detects violations of the Dependency Inversion Principle",
			Run:  f.run,
		},
	}, nil
}

func (f *Plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

func (f *Plugin) run(pass *analysis.Pass) (interface{}, error) {
	includePaths := make(map[string]struct{})
	for _, path := range f.settings.Include {
		includePaths[path] = struct{}{}
	}

	excludePaths := make(map[string]struct{})
	for _, path := range f.settings.Exclude {
		excludePaths[path] = struct{}{}
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// Check if the function name suggests it's a constructor
			if !isConstructor(funcDecl.Name.Name) {
				return true
			}

			// Skip functions without a return type
			if funcDecl.Type.Results == nil || len(funcDecl.Type.Results.List) == 0 {
				return true
			}

			// Get the return type of the function
			for _, result := range funcDecl.Type.Results.List {
				resultType := pass.TypesInfo.TypeOf(result.Type)
				if resultType == nil {
					continue
				}

				// Check if the return type is a concrete type (not an interface)
				if _, isInterface := resultType.Underlying().(*types.Interface); !isInterface {
					// Check include/exclude paths
					filePath := pass.Fset.Position(funcDecl.Pos()).Filename

					// Convert absolute filePath to a relative path
					relPath, err := filepath.Rel(pass.Pkg.Path(), filePath)
					if err != nil {
						// If conversion fails, fallback to absolute path
						relPath = filePath
					}

					// Check include paths
					for includePath := range includePaths {
						if strings.HasPrefix(relPath, includePath) {
							goto CheckExclude
						}
					}
					if len(includePaths) > 0 {
						continue
					}

					// Check exclude paths
				CheckExclude:
					for excludePath := range excludePaths {
						if strings.HasPrefix(relPath, excludePath) {
							continue
						}
					}

					// Check if the function name matches any constructor name pattern
					if !f.matchesNamePattern(funcDecl.Name.Name) {
						continue
					}

					// Report the violation
					pass.Reportf(funcDecl.Pos(), "constructor '%s' returns a concrete type '%s', use an interface instead", funcDecl.Name.Name, resultType.String())
				}
			}
			return true
		})
	}
	return nil, nil
}

// Helper function to check if a function name suggests it's a constructor
func isConstructor(name string) bool {
	return len(name) > 3 && strings.HasPrefix(name, "New")
}

// Helper function to check if a function name matches any of the specified patterns
func (f *Plugin) matchesNamePattern(name string) bool {
	// Skip the check if no patterns are defined
	if len(f.settings.NamePatterns) == 0 {
		return true
	}

	// match the name with the patterns
	for _, pattern := range f.settings.NamePatterns {
		if strings.Contains(strings.ToLower(name), strings.ToLower(pattern)) {
			return true
		}
	}

	return false
}
