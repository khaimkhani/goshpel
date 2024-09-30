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
)

type ASTNode struct {
	// this should be repr i think
	children []*ASTNode
	op       string
	injected bool
	// node type
	ntype string
	name  string
}
