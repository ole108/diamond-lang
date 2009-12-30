package common


// Define call types as 'enumeration':
type CallTypeEnum int
const (
  FREE_CALL = iota;
  BIND_CALL;
  CLOSURE_CALL;
)

/// AstNode - Basic interface of all AST nodes.
type AstNode interface {
  Error(msg string);
  String() string;
}

/// ExprAST - Basic interface of all expression AST nodes.
type ExprAst interface {
  AstNode;
  IsExpr() bool;
}

/// NumberExprAst - Interface of expressions for numeric literals like 123.
type NumberExprAst interface {
  ExprAst;
  Value() int64;
}

/// BoolExprAst - Interface of expressions for boolean literals like TRUE.
type BoolExprAst interface {
  ExprAst;
  Value() bool;
}

/// CharExprAst - Interface of expressions for character literals like 'a'.
type CharExprAst interface {
  ExprAst;
  Value() byte;
}

/// StringExprAst - Interface of expressions for string literals like `bla`.
type StringExprAst interface {
  ExprAst;
  Value() string;
}

/// ValueExprAst - Interface of expressions for referencing a value, like 'a'.
type ValueExprAst interface {
  ExprAst;
  Module() string;
  Value() string;
  Subs() []string;
}

/// ConstantExprAst - Interface of expressions for referencing a constant, like 'PI'.
type ConstantExprAst interface {
  ExprAst;
  Module() string;
  Constant() string;
  Subs() []string;
}

/// CallExprAst - Interface of expressions for function calls.
type CallExprAst interface {
  ExprAst;
  Module() string;
  Func() string;
  Type() CallTypeEnum;
  Args() []ExprAst;
}

/// AssignmentAst - Interface of assignment statements line: value = expr
type AssignmentAst interface {
  AstNode;
  Value() ValueExprAst;
  Expr() ExprAst;
}

/// BlockExprAst - Interface of expressions for function calls.
type BlockExprAst interface {
  ExprAst;
  Assignments() []AssignmentAst;
  Expr() ExprAst;
}

/// PrototypeAst - Interface of a 'prototype' for a function,
/// which captures its name, and its argument names (thus implicitly the number
/// of arguments the function takes).
type PrototypeAst interface {
  AstNode;
  Func() string;
  Args() []string;
}

/// FunctionAst - Interface of a function definition itself.
type FunctionAst interface {
  AstNode;
  Prototype() PrototypeAst;
  Body() ExprAst;
}

