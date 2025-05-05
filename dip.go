package linters

import (
	"go/ast"
	"go/types"
	"path/filepath"

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
			assignStmt, ok := n.(*ast.AssignStmt)
			if !ok {
				return true
			}

			for i, rhs := range assignStmt.Rhs {
				callExpr, ok := rhs.(*ast.CallExpr)
				if !ok {
					continue
				}

				// Check if the function being called is a constructor
				funIdent, ok := callExpr.Fun.(*ast.Ident)
				if !ok {
					continue
				}

				obj := pass.TypesInfo.ObjectOf(funIdent)
				if obj == nil {
					continue
				}

				// Skip external packages
				if pkg := obj.Pkg(); pkg != nil && pkg.Path() != pass.Pkg.Path() {
					continue
				}

				sig, ok := obj.Type().(*types.Signature)
				if !ok || sig.Results().Len() == 0 {
					continue
				}

				// Check if the return type is a concrete type (not an interface)
				named, ok := sig.Results().At(0).Type().(*types.Named)
				if !ok {
					continue
				}

				iface := named.Underlying()
				if _, isInterface := iface.(*types.Interface); isInterface {
					continue // Skip if the return type is an interface
				}

				// Check if the function name matches the constructor name pattern
				if !f.matchesNamePattern(funIdent.Name) {
					continue
				}

				// Check include/exclude paths
				filePath := pass.Fset.Position(assignStmt.Pos()).Filename
				if len(includePaths) > 0 {
					if _, included := includePaths[filePath]; !included {
						continue
					}
				}
				if _, excluded := excludePaths[filePath]; excluded {
					continue
				}

				// Report if the assigned variable is not an interface
				if len(assignStmt.Lhs) > i {
					lhs := assignStmt.Lhs[i]
					ident, ok := lhs.(*ast.Ident)
					if !ok {
						continue
					}

					lhsType := pass.TypesInfo.TypeOf(ident)
					if lhsType == nil {
						continue
					}

					if _, ok := lhsType.Underlying().(*types.Interface); !ok {
						pass.Reportf(assignStmt.Pos(), "direct initialization of concrete service '%s', use interface instead", named.Obj().Name())
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

// Helper function to check if a function name matches any of the specified patterns
func (f *Plugin) matchesNamePattern(name string) bool {
	for _, pattern := range f.settings.NamePatterns {
		if matched, _ := filepath.Match(pattern, name); matched {
			return true
		}
	}
	return false
}
