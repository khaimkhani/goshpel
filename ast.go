package main

// Stage 2: Need a better system to track and manipulate final
// repr of source code to be executed.
// Advantages to using an AST:
// 1. Consistent injection of vars as they are needed
// 2. Same with imports
// 3. Can recursively modify sub func defs to be valid

const (
	FUNCDEC = "FUNCDEC"
	// Short hand var dec
	SHVARDEC = "SHVARDEC"
	EXPR     = "EXPR"
	IMPORT   = "IMPORT"
	GOROUT   = "GOROUT"
	ROOT     = "ROOT"
)

const (
	// Node types
	BLOCK = "BLOCK"
	// This can include vars or literals
	VAL = "VAL"
	// a := 1
	VARDEC = "VARDEC"
	// a = 1
	VARASSIGN = "VARASSIGN"
	// math
	OP = "OP"
	// fmt.Println()
	FUNCCALL = "FUNCCALL"
)

type ASTNode struct {
	// this should be repr i think
	children []*ASTNode
	// node type
	ntype string
	val   string
}

func NewRootAst() *ASTNode {
	return NewAst(ROOT, ROOT)
}

func NewAst(ntype string, name string) *ASTNode {
	return &ASTNode{make([]*ASTNode, 0), ntype, name}
}

func (n *ASTNode) AddChildren(children []*ASTNode) *ASTNode {
	n.children = append(n.children, children...)
	return n
}

func (n *ASTNode) AddFuncDefNode(funcDef []string) *ASTNode {
	// New func def
	// Inject staged imports
	// break down expr
	return n
}

func (n *ASTNode) AddImportNode(imports []string) *ASTNode {
	// Inject import
	// Only if there is a var decl for that node
	return n
}

func (n *ASTNode) AddExprNode(expr []string) *ASTNode {
	// this operates a bit different to the two above.
	// called from it's immediate parent and returns itself.

	// op should be the vals in the expr
	node := NewAst("", EXPR, EXPR)

	return node
}

func (n *ASTNode) AddShvarDecNode(shvarDec []string) *ASTNode {

	node := NewAst("", SHVARDEC, SHVARDEC)
	return node
}

func GetStatementTypes(text string) ([]string, error) {

	return make([]string, 0), nil
}
