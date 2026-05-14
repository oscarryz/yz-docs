package main_check

import (
	"fmt"
	"yz/internal/ir"
	"yz/internal/lexer"
	"yz/internal/parser"
	"yz/internal/sema"
)

func main() {
	src := `
Account: { balance Int }
loader #(acc Account) { acc }
user #(acc Account) {
    loaded: loader(acc)
    print(loaded.balance)
}
main: {
    a: Account(42)
    user(a)
}
`
	_ = lexer.Tokenize([]byte(src))
	p := parser.New([]byte(src))
	sf, err := p.ParseFile()
	if err != nil {
		panic(err)
	}
	a := sema.NewAnalyzer()
	if err := a.AnalyzeFile(sf); err != nil {
		panic(err)
	}
	f := ir.Lower(sf, a, "main")
	for _, decl := range f.Decls {
		sd, ok := decl.(*ir.SingletonDecl)
		if !ok {
			continue
		}
		for _, md := range sd.Methods {
			if md.Name != "Call" {
				continue
			}
			fmt.Printf("Method %s on %s:\n", md.Name, md.RecvType)
			for _, stmt := range md.Body {
				if es, ok := stmt.(*ir.ExprStmt); ok {
					if th, ok := es.Expr.(*ir.ThunkExpr); ok {
						fmt.Printf("  ThunkExpr ResultType=%q RecvCown=%q ExtraCowns=%v\n", th.ResultType, th.RecvCown, th.ExtraCowns)
						for i, s := range th.Body {
							if ds, ok := s.(*ir.DeclStmt); ok {
								fmt.Printf("    [%d] DeclStmt Name=%q Type=%q IsThunk=%v\n", i, ds.Name, ds.Type, ds.IsThunk)
							} else {
								fmt.Printf("    [%d] %T\n", i, s)
							}
						}
					}
				}
			}
		}
	}
}
