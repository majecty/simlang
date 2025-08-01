package evaluator

import (
	"fmt"

	"simlang/types"
)

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

func Eval(ast *types.AST) (any, error) {
	defaultEnv := &Env{EnvMap: make(map[string]any)}
	defaultEnv.EnvMap["+"] = func(args []any) (any, error) {
		resultSum := 0.0
		for _, arg := range args {
			if num, ok := arg.(float64); ok {
				resultSum += num
			}
		}
		return resultSum, nil
	}

	if result, err := evalSingle(ast.Root, defaultEnv); err != nil {
		return nil, fmt.Errorf("failed to eval: %w, original input is %+v", err, ast)
	} else {
		return result, nil
	}
}

func evalSingle(item types.ASTNode, env *Env) (any, error) {
	switch v := item.(type) {
	case *types.NumberNode:
		return v.Value, nil
	case *types.SymbolNode:
		return env.Get(v.Name), nil
	case *types.CallNode:
		switch function := v.Function.(type) {
		case *types.SymbolNode:
			evalutedArgs := make([]any, 0)
			for _, arg := range v.Args {
				evaluatedArg, err := evalSingle(arg, env)
				if err != nil {
					return nil, fmt.Errorf("failed to eval call: %w", err)
				}
				evalutedArgs = append(evalutedArgs, evaluatedArg)
			}

			f := env.Get(function.Name)
			if f == nil {
				return nil, fmt.Errorf("failed to eval call, function %s not found", function.Name)
			}

			return f.(func([]any) (any, error))(evalutedArgs)
		default:
			return item, nil
		}
	case *types.LetNode:
		executedEnv := make(map[string]any)
		for key, value := range v.LetEnv {
			evalResult, err := evalSingle(value, &Env{EnvMap: executedEnv, parent: env})
			if err != nil {
				return nil, fmt.Errorf("failed to eval let value: %w", err)
			}
			executedEnv[key] = evalResult
		}
		return evalSingle(v.Body, &Env{EnvMap: executedEnv, parent: env})
	case *types.LambdaNode:
		return func(args []any) (any, error) {
			if len(args) != len(v.Args) {
				return nil, fmt.Errorf("expected %d arguments, got %d, args: %+v", len(v.Args), len(args), args)
			}
			executedEnv := make(map[string]any)
			for i, arg := range args {
				executedEnv[v.Args[i].Name] = arg
			}
			return evalSingle(v.Body, &Env{EnvMap: executedEnv, parent: env})
		}, nil
	default:
		return v, nil
	}
}
