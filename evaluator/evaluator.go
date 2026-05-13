// Created and Developed by Farhan Kertadiwangsa
package evaluator

import (
	"fmt"

	"ky/ast"
	"ky/object"
)

// Singleton untuk nilai yang sering dipakai (zero-allocation)
var (
	BENAR  = &object.Boolean{Value: true}
	SALAH  = &object.Boolean{Value: false}
	KOSONG = &object.Kosong{}
)

// Eval mengevaluasi node AST dan mengembalikan nilai Object
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Program
	case *ast.Program:
		return evalProgram(node, env)

	// Statements
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.BuatStatement:
		val := Eval(node.Value, env)
		if object.IsError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return KOSONG

	case *ast.AssignStatement:
		val := Eval(node.Value, env)
		if object.IsError(val) {
			return val
		}
		if !env.Update(node.Name.Value, val) {
			return object.NewError("baris %d: variabel '%s' belum dideklarasikan (gunakan 'buat')",
				node.Token.Line, node.Name.Value)
		}
		return KOSONG

	case *ast.BalikStatement:
		if node.ReturnValue == nil {
			return &object.ReturnValue{Value: KOSONG}
		}
		val := Eval(node.ReturnValue, env)
		if object.IsError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.HentiStatement:
		return &object.BreakSignal{}

	case *ast.LanjutStatement:
		return &object.ContinueSignal{}

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.BooleanLiteral:
		return nativeBool(node.Value)

	case *ast.KosongLiteral:
		return KOSONG

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if object.IsError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if object.IsError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if object.IsError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.JikaExpression:
		return evalJikaExpression(node, env)

	case *ast.SelamaExpression:
		return evalSelamaExpression(node, env)

	case *ast.UlangExpression:
		return evalUlangExpression(node, env)

	case *ast.FungsiLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
			Name:       node.Name,
		}

	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if object.IsError(fn) {
			return fn
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && object.IsError(args[0]) {
			return args[0]
		}
		return applyFunction(fn, args)

	case *ast.ArrayLiteral:
		elems := evalExpressions(node.Elements, env)
		if len(elems) == 1 && object.IsError(elems[0]) {
			return elems[0]
		}
		return &object.Array{Elements: elems}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if object.IsError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if object.IsError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.MapLiteral:
		return evalMapLiteral(node, env)
	}

	return object.NewError("node tidak dikenal: %T", node)
}

// Program & Block

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range program.Statements {
		result = Eval(stmt, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range block.Statements {
		result = Eval(stmt, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_OBJ || rt == object.ERROR_OBJ ||
				rt == object.BREAK_OBJ || rt == object.CONTINUE_OBJ {
				return result
			}
		}
	}
	return result
}

// Identifier

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := Builtins[node.Value]; ok {
		return builtin
	}
	return object.NewError("baris %d: '%s' tidak ditemukan", node.Token.Line, node.Value)
}

// Prefix

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "-":
		return evalMinusPrefix(right)
	case "!", "tidak", "bukan":
		return nativeBool(!isTruthy(right))
	default:
		return object.NewError("operator prefix tidak dikenal: %s%s", op, right.Type())
	}
}

func evalMinusPrefix(right object.Object) object.Object {
	switch v := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -v.Value}
	case *object.Float:
		return &object.Float{Value: -v.Value}
	default:
		return object.NewError("operator '-' tidak bisa diterapkan ke %s", right.Type())
	}
}

// Infix

func evalInfixExpression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfix(op, left.(*object.Integer), right.(*object.Integer))
	case left.Type() == object.FLOAT_OBJ || right.Type() == object.FLOAT_OBJ:
		return evalFloatInfix(op, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfix(op, left.(*object.String), right.(*object.String))
	case op == "==":
		return nativeBool(left == right)
	case op == "!=":
		return nativeBool(left != right)
	case op == "dan":
		return nativeBool(isTruthy(left) && isTruthy(right))
	case op == "atau":
		return nativeBool(isTruthy(left) || isTruthy(right))
	case left.Type() != right.Type():
		return object.NewError("tipe tidak cocok: %s %s %s", left.Type(), op, right.Type())
	default:
		return object.NewError("operator tidak dikenal: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalIntegerInfix(op string, left, right *object.Integer) object.Object {
	l, r := left.Value, right.Value
	switch op {
	case "+":
		return &object.Integer{Value: l + r}
	case "-":
		return &object.Integer{Value: l - r}
	case "*":
		return &object.Integer{Value: l * r}
	case "/":
		if r == 0 {
			return object.NewError("pembagian dengan nol!")
		}
		return &object.Integer{Value: l / r}
	case "%":
		if r == 0 {
			return object.NewError("modulo dengan nol!")
		}
		return &object.Integer{Value: l % r}
	case "==":
		return nativeBool(l == r)
	case "!=":
		return nativeBool(l != r)
	case "<":
		return nativeBool(l < r)
	case ">":
		return nativeBool(l > r)
	case "<=":
		return nativeBool(l <= r)
	case ">=":
		return nativeBool(l >= r)
	default:
		return object.NewError("operator '%s' tidak dikenal untuk bilangan", op)
	}
}

func evalFloatInfix(op string, left, right object.Object) object.Object {
	l := toFloatVal(left)
	r := toFloatVal(right)
	switch op {
	case "+":
		return &object.Float{Value: l + r}
	case "-":
		return &object.Float{Value: l - r}
	case "*":
		return &object.Float{Value: l * r}
	case "/":
		if r == 0 {
			return object.NewError("pembagian dengan nol!")
		}
		return &object.Float{Value: l / r}
	case "==":
		return nativeBool(l == r)
	case "!=":
		return nativeBool(l != r)
	case "<":
		return nativeBool(l < r)
	case ">":
		return nativeBool(l > r)
	case "<=":
		return nativeBool(l <= r)
	case ">=":
		return nativeBool(l >= r)
	default:
		return object.NewError("operator '%s' tidak dikenal untuk desimal", op)
	}
}

func evalStringInfix(op string, left, right *object.String) object.Object {
	switch op {
	case "+":
		return &object.String{Value: left.Value + right.Value}
	case "==":
		return nativeBool(left.Value == right.Value)
	case "!=":
		return nativeBool(left.Value != right.Value)
	case "<":
		return nativeBool(left.Value < right.Value)
	case ">":
		return nativeBool(left.Value > right.Value)
	default:
		return object.NewError("operator '%s' tidak berlaku untuk teks", op)
	}
}

//  Kondisi & Perulangan

func evalJikaExpression(node *ast.JikaExpression, env *object.Environment) object.Object {
	cond := Eval(node.Condition, env)
	if object.IsError(cond) {
		return cond
	}
	if isTruthy(cond) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	}
	return KOSONG
}

func evalSelamaExpression(node *ast.SelamaExpression, env *object.Environment) object.Object {
	var result object.Object = KOSONG
	for {
		cond := Eval(node.Condition, env)
		if object.IsError(cond) {
			return cond
		}
		if !isTruthy(cond) {
			break
		}
		result = Eval(node.Body, env)
		if result != nil {
			switch result.Type() {
			case object.RETURN_OBJ, object.ERROR_OBJ:
				return result
			case object.BREAK_OBJ:
				return KOSONG
			case object.CONTINUE_OBJ:
				continue
			}
		}
	}
	return result
}

func evalUlangExpression(node *ast.UlangExpression, env *object.Environment) object.Object {
	loopEnv := object.NewEnclosedEnvironment(env)
	if node.Init != nil {
		val := Eval(node.Init, loopEnv)
		if object.IsError(val) {
			return val
		}
	}
	var result object.Object = KOSONG
	for {
		if node.Condition != nil {
			cond := Eval(node.Condition, loopEnv)
			if object.IsError(cond) {
				return cond
			}
			if !isTruthy(cond) {
				break
			}
		}
		result = Eval(node.Body, loopEnv)
		if result != nil {
			switch result.Type() {
			case object.RETURN_OBJ, object.ERROR_OBJ:
				return result
			case object.BREAK_OBJ:
				return KOSONG
			case object.CONTINUE_OBJ:
				// jalankan post dulu, lalu lanjut
			}
		}
		if node.Post != nil {
			val := Eval(node.Post, loopEnv)
			if object.IsError(val) {
				return val
			}
		}
	}
	return result
}

// Fungsi & Pemanggilan

func evalExpressions(exprs []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exprs {
		val := Eval(e, env)
		if object.IsError(val) {
			return []object.Object{val}
		}
		result = append(result, val)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return object.NewError("fungsi '%s' butuh %d argumen, dapat %d",
				fn.Name, len(fn.Parameters), len(args))
		}
		extEnv := object.NewEnclosedEnvironment(fn.Env)
		for i, param := range fn.Parameters {
			extEnv.Set(param.Value, args[i])
		}
		result := Eval(fn.Body, extEnv)
		if rv, ok := result.(*object.ReturnValue); ok {
			return rv.Value
		}
		return result
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return object.NewError("bukan sebuah fungsi: %s", fn.Type())
	}
}

// Array & Map

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		arr := left.(*object.Array)
		i := index.(*object.Integer).Value
		if i < 0 {
			i = int64(len(arr.Elements)) + i // indeks negatif dari belakang
		}
		if i < 0 || int(i) >= len(arr.Elements) {
			return object.NewError("indeks %d di luar batas (panjang: %d)", index.(*object.Integer).Value, len(arr.Elements))
		}
		return arr.Elements[i]
	case left.Type() == object.MAP_OBJ:
		m := left.(*object.Map)
		key, ok := index.(object.Hashable)
		if !ok {
			return object.NewError("tipe %s tidak bisa dijadikan kunci map", index.Type())
		}
		pair, ok := m.Pairs[key.HashKey()]
		if !ok {
			return KOSONG
		}
		return pair.Value
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		s := []rune(left.(*object.String).Value)
		i := index.(*object.Integer).Value
		if i < 0 {
			i = int64(len(s)) + i
		}
		if i < 0 || int(i) >= len(s) {
			return object.NewError("indeks %d di luar batas teks", index.(*object.Integer).Value)
		}
		return &object.String{Value: string(s[i])}
	default:
		return object.NewError("operasi indeks tidak didukung: %s[%s]", left.Type(), index.Type())
	}
}

func evalMapLiteral(node *ast.MapLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.MapPair)
	for keyNode, valNode := range node.Pairs {
		key := Eval(keyNode, env)
		if object.IsError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return object.NewError("tipe %s tidak bisa dijadikan kunci map", key.Type())
		}
		val := Eval(valNode, env)
		if object.IsError(val) {
			return val
		}
		pairs[hashKey.HashKey()] = object.MapPair{Key: key, Value: val}
	}
	return &object.Map{Pairs: pairs}
}

// Helper

func isTruthy(obj object.Object) bool {
	switch v := obj.(type) {
	case *object.Boolean:
		return v.Value
	case *object.Kosong:
		return false
	case *object.Integer:
		return v.Value != 0
	case *object.Float:
		return v.Value != 0
	case *object.String:
		return v.Value != ""
	case *object.Array:
		return len(v.Elements) > 0
	default:
		return true
	}
}

func nativeBool(b bool) *object.Boolean {
	if b {
		return BENAR
	}
	return SALAH
}

func toFloatVal(obj object.Object) float64 {
	switch v := obj.(type) {
	case *object.Integer:
		return float64(v.Value)
	case *object.Float:
		return v.Value
	default:
		return 0
	}
}

// FormatError memformat pesan error
func FormatError(err *object.Error, _ string) string {
	return fmt.Sprintf("\n[.ky ERROR] %s\n", err.Message)
}
