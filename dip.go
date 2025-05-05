package linters

import (
	"go/ast"
	"go/types"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("dip", New)
}

type Settings struct{}

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

				funIdent, ok := callExpr.Fun.(*ast.Ident)
				if !ok {
					continue
				}

				obj := pass.TypesInfo.ObjectOf(funIdent)
				if obj == nil {
					continue
				}

				sig, ok := obj.Type().(*types.Signature)
				if !ok || sig.Results().Len() == 0 {
					continue
				}

				named, ok := sig.Results().At(0).Type().(*types.Named)
				if !ok {
					continue
				}

				iface := named.Underlying()
				if _, isInterface := iface.(*types.Interface); isInterface {
					return true // returning an interface is okay
				}

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

					// If assigned var is not an interface
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
