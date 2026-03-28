// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"strings"
)

// StatementRef stores statement pointer with class scope path and source span.
type StatementRef struct {
	// ClassPath is parent class path where statement is declared.
	// Empty path means top-level scope.
	ClassPath []string `json:"class_path,omitempty" yaml:"class_path,omitempty"`

	// Statement is AST statement node reference.
	Statement *Statement `json:"statement,omitempty" yaml:"statement,omitempty"`

	// Start is source start position.
	Start Position `json:"start,omitzero" yaml:"start,omitempty"`

	// End is source end position.
	End Position `json:"end,omitzero" yaml:"end,omitempty"`
}

// WalkClasses iterates class declarations in source order.
// Callback receives full class path (top -> nested).
// Return false from callback to stop traversal.
func (file File) WalkClasses(visit func(path []string, classDecl *ClassDecl) bool) {
	if visit == nil {
		return
	}

	walkClassStatements(file.Statements, nil, visit)
}

// WalkStatements iterates all statements in source order.
// Callback receives statement reference with scope path and source span.
// Return false from callback to stop traversal.
func (file File) WalkStatements(visit func(ref StatementRef) bool) {
	if visit == nil {
		return
	}

	walkStatements(file.Statements, nil, visit)
}

// FindClass resolves class by full path, for example: CfgVehicles -> Car.
func (file File) FindClass(path ...string) (*ClassDecl, bool) {
	if len(path) == 0 {
		return nil, false
	}

	current := file.Statements
	var classDecl *ClassDecl

	for _, className := range path {
		next, ok := findNestedClass(current, className)
		if !ok {
			return nil, false
		}

		classDecl = next
		current = next.Body
	}

	return classDecl, true
}

// FindClass resolves direct nested class by name.
func (classDecl *ClassDecl) FindClass(name string) (*ClassDecl, bool) {
	if classDecl == nil {
		return nil, false
	}

	return findNestedClass(classDecl.Body, name)
}

// FindProperty resolves direct class property assignment by name.
func (classDecl *ClassDecl) FindProperty(name string) (*PropertyAssign, bool) {
	if classDecl == nil {
		return nil, false
	}

	for idx := range classDecl.Body {
		stmt := &classDecl.Body[idx]
		if stmt.Kind != NodeProperty || stmt.Property == nil || stmt.Property.Name != name {
			continue
		}

		return stmt.Property, true
	}

	return nil, false
}

// FindArrayAssign resolves direct class array assignment by name.
func (classDecl *ClassDecl) FindArrayAssign(name string) (*ArrayAssign, bool) {
	if classDecl == nil {
		return nil, false
	}

	for idx := range classDecl.Body {
		stmt := &classDecl.Body[idx]
		if stmt.Kind != NodeArrayAssign || stmt.ArrayAssign == nil || stmt.ArrayAssign.Name != name {
			continue
		}

		return stmt.ArrayAssign, true
	}

	return nil, false
}

// PathString joins class path using "/" separator.
func (ref StatementRef) PathString() string {
	if len(ref.ClassPath) == 0 {
		return ""
	}

	joined := ref.ClassPath[0]
	var joinedSb115 strings.Builder
	for idx := 1; idx < len(ref.ClassPath); idx++ {
		joinedSb115.WriteString("/" + ref.ClassPath[idx])
	}
	joined += joinedSb115.String()

	return joined
}

// walkClassStatements walks class declarations recursively.
func walkClassStatements(statements []Statement, parentPath []string, visit func(path []string, classDecl *ClassDecl) bool) bool {
	for idx := range statements {
		stmt := &statements[idx]
		if stmt.Kind != NodeClass || stmt.Class == nil {
			continue
		}

		nextPath := append([]string(nil), parentPath...)
		nextPath = append(nextPath, stmt.Class.Name)
		path := append([]string(nil), nextPath...)
		if !visit(path, stmt.Class) {
			return false
		}

		if !walkClassStatements(stmt.Class.Body, nextPath, visit) {
			return false
		}
	}

	return true
}

// walkStatements walks all statements recursively in source order.
func walkStatements(statements []Statement, classPath []string, visit func(ref StatementRef) bool) bool {
	for idx := range statements {
		stmt := &statements[idx]
		ref := StatementRef{
			ClassPath: append([]string(nil), classPath...),
			Statement: stmt,
			Start:     stmt.Start,
			End:       stmt.End,
		}
		if !visit(ref) {
			return false
		}

		if stmt.Kind != NodeClass || stmt.Class == nil {
			continue
		}

		nextPath := append([]string(nil), classPath...)
		nextPath = append(nextPath, stmt.Class.Name)
		if !walkStatements(stmt.Class.Body, nextPath, visit) {
			return false
		}
	}

	return true
}

// findNestedClass resolves class declaration by name in direct statement scope.
func findNestedClass(statements []Statement, name string) (*ClassDecl, bool) {
	for idx := range statements {
		stmt := &statements[idx]
		if stmt.Kind != NodeClass || stmt.Class == nil || stmt.Class.Name != name {
			continue
		}

		return stmt.Class, true
	}

	return nil, false
}
