package rvcfg

// NodeKind identifies AST statement node type.
type NodeKind string

const (
	// NodeClass is class declaration node.
	NodeClass NodeKind = "class"

	// NodeDelete is delete statement node.
	NodeDelete NodeKind = "delete"

	// NodeProperty is scalar/property assignment node.
	NodeProperty NodeKind = "property"

	// NodeArrayAssign is array assignment/append node.
	NodeArrayAssign NodeKind = "array_assign"

	// NodeExtern is extern declaration node.
	NodeExtern NodeKind = "extern"

	// NodeEnum is enum declaration node.
	NodeEnum NodeKind = "enum"
)

// ValueKind identifies AST value shape.
type ValueKind string

const (
	// ValueScalar is scalar token sequence value.
	ValueScalar ValueKind = "scalar"

	// ValueArray is array literal value.
	ValueArray ValueKind = "array"
)

// File is parsed config AST root.
type File struct {
	// Source is logical source name, usually file path.
	Source string `json:"source,omitempty" yaml:"source,omitempty"`

	// Statements contains top-level declarations in source order.
	Statements []Statement `json:"statements,omitempty" yaml:"statements,omitempty"`

	// Start is source start position.
	Start Position `json:"start,omitempty" yaml:"start,omitempty"`

	// End is source end position.
	End Position `json:"end,omitempty" yaml:"end,omitempty"`
}

// Statement stores one top-level or class-body declaration.
type Statement struct {
	// Class is class declaration payload.
	Class *ClassDecl `json:"class,omitempty" yaml:"class,omitempty"`

	// Delete is delete declaration payload.
	Delete *DeleteStmt `json:"delete,omitempty" yaml:"delete,omitempty"`

	// Property is scalar assignment payload.
	Property *PropertyAssign `json:"property,omitempty" yaml:"property,omitempty"`

	// ArrayAssign is array assignment payload.
	ArrayAssign *ArrayAssign `json:"array_assign,omitempty" yaml:"array_assign,omitempty"`

	// Extern is extern declaration payload.
	Extern *ExternDecl `json:"extern,omitempty" yaml:"extern,omitempty"`

	// Enum is enum declaration payload.
	Enum *EnumDecl `json:"enum,omitempty" yaml:"enum,omitempty"`

	// TrailingComment is optional comment attached after statement on same line.
	TrailingComment *Comment `json:"trailing_comment,omitempty" yaml:"trailing_comment,omitempty"`

	// Kind is active statement kind.
	Kind NodeKind `json:"kind,omitempty" yaml:"kind,omitempty"`

	// LeadingComments are comments directly attached before statement.
	LeadingComments []Comment `json:"leading_comments,omitempty" yaml:"leading_comments,omitempty"`

	// TrailingComments are additional comments attached after statement (usually before EOF).
	TrailingComments []Comment `json:"trailing_comments,omitempty" yaml:"trailing_comments,omitempty"`

	// Start is source start position.
	Start Position `json:"start,omitempty" yaml:"start,omitempty"`

	// End is source end position.
	End Position `json:"end,omitempty" yaml:"end,omitempty"`
}

// Comment stores original comment token text and position.
type Comment struct {
	// Text is original comment token text, including delimiters.
	Text string `json:"text,omitempty" yaml:"text,omitempty"`

	// Start is source start position.
	Start Position `json:"start,omitempty" yaml:"start,omitempty"`

	// End is source end position.
	End Position `json:"end,omitempty" yaml:"end,omitempty"`
}

// ClassDecl describes class declaration body or forward declaration.
type ClassDecl struct {
	// Name is class identifier.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Base is optional parent class expression.
	Base string `json:"base,omitempty" yaml:"base,omitempty"`

	// Body contains class members for non-forward declaration.
	Body []Statement `json:"body,omitempty" yaml:"body,omitempty"`

	// Forward marks `class Name;` form.
	Forward bool `json:"forward,omitempty" yaml:"forward,omitempty"`
}

// DeleteStmt describes delete statement.
type DeleteStmt struct {
	// Name is deleted class/property name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

// ExternDecl describes extern declaration.
type ExternDecl struct {
	// Name is external symbol name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Class indicates `extern class Name;` form.
	Class bool `json:"class,omitempty" yaml:"class,omitempty"`
}

// EnumDecl describes enum declaration.
type EnumDecl struct {
	// Name is optional enum identifier.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Items are enum members in source order.
	Items []EnumItem `json:"items,omitempty" yaml:"items,omitempty"`
}

// EnumItem describes one enum member.
type EnumItem struct {
	// Name is enum member identifier.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// ValueRaw is optional explicit member value expression.
	ValueRaw string `json:"value_raw,omitempty" yaml:"value_raw,omitempty"`
}

// PropertyAssign describes `name = value;` statement.
type PropertyAssign struct {
	// Name is assignment target.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Value is assigned expression.
	Value Value `json:"value,omitempty" yaml:"value,omitempty"`
}

// ArrayAssign describes `name[] = value;` or `name[] += value;`.
type ArrayAssign struct {
	// Name is assignment target.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Value is assigned expression, usually array literal.
	Value Value `json:"value,omitempty" yaml:"value,omitempty"`

	// Append is true for `+=`, false for `=`.
	Append bool `json:"append,omitempty" yaml:"append,omitempty"`
}

// Value describes scalar expression or array literal.
type Value struct {
	// Kind identifies scalar vs array.
	Kind ValueKind `json:"kind,omitempty" yaml:"kind,omitempty"`

	// Raw is compact scalar token sequence for ValueScalar.
	Raw string `json:"raw,omitempty" yaml:"raw,omitempty"`

	// Elements are nested items for ValueArray.
	Elements []Value `json:"elements,omitempty" yaml:"elements,omitempty"`

	// Start is source start position.
	Start Position `json:"start,omitempty" yaml:"start,omitempty"`

	// End is source end position.
	End Position `json:"end,omitempty" yaml:"end,omitempty"`
}
