package main

import (
	"context"
	"fmt"
	"go/ast"

	"github.com/liut0/gomultilinter/api"
)

type deadcodeLinter struct {
}

var LinterFactory api.LinterFactory = &deadcodeLinter{}

func (l *deadcodeLinter) NewLinterConfig() api.LinterConfig {
	return &deadcodeLinter{}
}

func (l *deadcodeLinter) NewLinter() (api.Linter, error) {
	return l, nil
}

func (*deadcodeLinter) Name() string {
	return "deadcode"
}

func (l *deadcodeLinter) LintPackage(ctx context.Context, pkg *api.Package, reporter api.IssueReporter) error {
	files := make(map[string]*ast.File, len(pkg.PkgInfo.Files))
	for _, f := range pkg.PkgInfo.Files {
		pos := pkg.FSet.Position(f.Pos())
		files[pos.Filename] = f
	}

	astPkg := &ast.Package{
		Files: files,
		Name:  pkg.PkgInfo.Pkg.Name(),
	}

	for _, report := range doPackage(pkg.FSet, astPkg) {
		reporter.Report(&api.Issue{
			Position: pkg.FSet.Position(report.pos),
			Category: "deadcode",
			Message:  fmt.Sprintf("%s is unused", report.name),
			Severity: api.SeverityWarning,
		})
	}

	return nil
}
