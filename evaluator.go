package main

type Env struct {
	EnvMap map[string]any
	parent *Env
}

func (e *Env) Get(name string) any {
  if val, ok := e.EnvMap[name]; ok {
    return val
  }
  if e.parent != nil {
    return e.parent.Get(name)
  }
  return nil
}

func eval(ast *AST) any {
	return evalSingle(ast.Root, &Env{EnvMap: make(map[string]any)})
}

func evalSingle(item ASTNode, env *Env) any {
	switch v := item.(type) {
	case *NumberNode:
		return v.Value
	case *SymbolNode:
		return env.Get(v.Name)
	case *CallNode:
		switch function := v.Function.(type) {
		case *SymbolNode:
			switch function.Name {
			case "+":
				return evalAdd(v.Args, env)
			default:
				return item
			}
		default:
			return item
		}
	case *LetNode:
		executedEnv := make(map[string]any)
		for key, value := range v.LetEnv {
      executedEnv[key] = evalSingle(value, &Env{EnvMap: executedEnv, parent: env})
		}
    return evalSingle(v.Body, &Env{EnvMap: executedEnv, parent: env})
	default:
		return v
	}
}

func evalAdd(args []ASTNode, env *Env) float64 {
	var sum float64
	sum = 0
	for _, arg := range args {
		evaluated := evalSingle(arg, env)
		if num, ok := evaluated.(float64); ok {
			sum += num
		}
	}
	return sum
}
