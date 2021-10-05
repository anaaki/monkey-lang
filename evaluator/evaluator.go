package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	//特定の値しか取らないオブジェクトはインスタンスを1つつくって参照で使う。
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// 文
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	// 式
	case *ast.Boolean:
		return nativeBooltoBooleanObject(node.Value)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}
	return nil
}

func evalStatements(stmt []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range stmt {
		result = Eval(statement)
	}
	return result
}

func nativeBooltoBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
