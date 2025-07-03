package main

func eval(ast *AST) any {
	return evalSingle(ast.Root)
}

func evalSingle(item ASTNode) any {
	switch v := item.(type) {
	case *NumberNode:
		return v.Value
	case *SymbolNode:
		return v.Name
	case *CallNode:
		switch function := v.Function.(type) {
		case *SymbolNode:
			switch function.Name {
			case "+":
				return evalAdd(v.Args)
			default:
				return item
			}
		default:
			return item
		}
	default:
		return v
	}
}

func evalAdd(args []ASTNode) float64 {
	var sum float64
	sum = 0
	for _, arg := range args {
		evaluated := evalSingle(arg)
		if num, ok := evaluated.(float64); ok {
			sum += num
		}
	}
	return sum
}
