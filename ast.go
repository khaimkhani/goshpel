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

type ASTRoot struct {
	imports []*ASTNode
	main    []*ASTNode
	decls   []*ASTNode
}

type ASTNode struct {
	children []*ASTNode
	// node type
	ntype      string
	startToken string
	endToken   string
	val        string
	// When a package becomes unused/used
	inject bool
}

func NewRootAst() *ASTRoot {
	return &ASTRoot{}
}

func (r *ASTRoot) AddImports(imports string) *ASTRoot {
	r.imports = append(r.imports, CreateTree(imports)...)
	return r
}

func (r *ASTRoot) AddDecls(decls string) *ASTRoot {
	r.decls = append(r.decls, CreateTree(decls)...)
	return r
}

func (r *ASTRoot) AddMain(main string) *ASTRoot {
	r.main = append(r.main, CreateTree(main)...)
	return r
}

func (r *ASTRoot) CreateSrc() {
	// return src code from ast repr
	return
}

func CreateTree(src string) []*ASTNode {
	nodes := make([]*ASTNode, 0)
	// construct tree from source and return current root
	n := &ASTNode{}
	nodes = append(nodes, n)
	return nodes

}
