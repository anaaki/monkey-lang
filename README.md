
#️ 字句解析でやりたいこと
インプットは文字
列
```
"let x = 5 + 5;"
```

アウトプットは「トークン」。
記号、識別子、キーワード、数字などをトークンに変換する。
```
[
    LET,
    IDENTIFIER("x"),
    EQUAL_SIGN,
    INTEGER(5),
    PLUS_SIGN,
    INTEGER(5),
    SEMICOLON
]
```

## ポイント

ホワイトスペースはトークンとして出てこない

トークンの種類は
 * タイプ（整数）
 * 識別子（変数）
 * キーワード（let、fnなど）
 * キーワード

特殊なトークンを作る
 * ILLEGALはトークンや文字が未知であることを示す
 * EOFはファイル終端をしめし、構文解析機にここで終了して良いと伝える

## 仕様

 * ソースコードを文字列として受け取り、出力としてトークンを返す
 * ソースコードはASCII文字だけに対応
 * バッファリング、保存はしない。```NextToken()```を呼ぶことで、ソースコードを読み進めて、トークンを返す。
 * ソースコードはstringとしてあつかう。本当はファイル名や行番号があったほうがデバッグしやすいけど、シンプルにしたいので。

 # Lexer

1文字つづ文字列（ソースコード）をうけとり、トークン化していく。
positionは現在の文字列、readPositionはこれから読み込む位置(positionの次に読む場所)

トークンは次の文字列が何かによって、種類が変わるのでreadPositionを設けている

2文字のトークン対応は、現在の文字と先読みした文字を個別に判定して実現する。
=の場合、次に=があれば、==と判定する。


# 構文解析

> JSON パーサーはテキストを入力として受け取り、その入力に対応するデータ構 造を生成する。
> これはプログラミング言語のパーサーがしていることと全く同じなんだ。違いはJSON パーサーの場合は入力を見れば
> データ構造がすぐにわかるというだけだ。

データ構造をプログラマが意識することはほとんどない。ほとんどのインタプリタ、コンパイラにおいて、
ソースコードの内部表現は構文木、もしくは抽象構文木。
セミコロン、改行文字、ホワイトスペース、波括弧、角括弧、丸括弧などは、言語や構文解析器によってはASTを構築する際に
構文解析器を導く役割を持つだけで、AST中には出現しない。


構文解析器が実施することは
 * ソースコードを入力として(テキストまた はトークン列として)受け取り、ソースコードを表現するようなあるデータ構造を生成する。
 * その構文解析の間、入力が期待された構造に従っているかをチェックする。

## 構文解析の種類

ボトムアップ
トップダウン(再帰下降構文解析、アーリー法、予測的構文解析)

Monkeyは再帰下降構文解析（トップダウン演算子優先順位）解析器。Pratt構文解析器

## AST疑似コード

```js
function parseProgram() {
	program = newProgramASTNode()
	advanceTokens()
	for (currentToken() != EOF_TOKEN) {
		statement = null
		if (currentToken() == LET_TOKEN) {
			statement = parseLetStatement()
		} else if (currentToken() == RETURN_TOKEN) {
			statement = parseReturnStatement()
		} else if (currentToken() == IF_TOKEN) {
			statement = parseIfStatement()
		}
		if (statement != null) {
			program.Statements.push(statement)
		}
		advanceTokens()
	}
	return program
}

function parseLetStatement() {
	advanceTokens()
	identifier = parseIdentifier()
	advanceTokens()
	if currentToken() != EQUAL_TOKEN {
		parseError("no equal sign!")
		return null
	}
	advanceTokens()
	value = parseExpression()
	variableStatement = newVariableStatementASTNode()
    variableStatement.identifier = identifier
    variableStatement.value = value
	return variableStatement
}

function parseIdentifier() {
	identifier = newIdentifierASTNode()
    identifier.token = currentToken()
    return identifier
}

function parseExpression() {
	if (currentToken() == INTEGER_TOKEN) {
		if (nextToken() == PLUS_TOKEN) {
			return parseOperatorExpression()
		} else if (nextToken() == SEMICOLON_TOKEN) {
			return parseIntegerLiteral()
		}
	} else if (currentToken() == LEFT_PAREN) {
		return parseGroupedExpression()
	}
	// [...]
}


function parseOperatorExpression() {
    operatorExpression = newOperatorExpression()
    operatorExpression.left = parseIntegerLiteral()
    advanceTokens()
    operatorExpression.operator = currentToken()
    advanceTokens()
    operatorExpression.right = parseExpression()
    return operatorExpression()
}
```


## 構文解析 let
式<expression>は値を生成し、文<identifire>はしない。
let <identifire> = <expression>;
式は例えば”5”。5を生成する。
let x = 5 は文。何も生成しないから。

## 構文解析 return
return <expression>;
return の後に式が来る。

## 式
Monkeyの式は複雑なので、expectPeekなど、現在のトークンに基づいて何かを決めるやり方はできない。
Vaughan Prattを使うことになる。

>構文解析関数(ここでは parseLetStatementメソッドを思い浮かべてほしい)を (BNFやEBNFで定義される)文法ルールに関連付けるのではなく、Prattはこれらの関数(彼は 「semmantic code」と呼んでいる)を単一のトークンタイプに関連付ける点だ。この方式の肝となるの

は、それぞれのトークンタイプに対して 2 つの構文解析関数を関連付けるところだ。これはトークンの 配置、つまり中置か前置かによる。
return let以外は全部式として扱うことになる。つまり値を生成する。
 * 前置演算子 -5, !true
 * 二項演算子 5 + 5
 * foo == bar

識別子も式になる
 * add(foo, bar)
 * foo * bar / foobar

関数リテラルも四季になる
let add = func(x, y){return x + y};
(fn(x) { return x }(5) + 10 ) * 10

if式もある。resultはtrueになる。
let result = if(10>5){ true} else {false};


### 式構文解析の用語
前置演算子(prefix operator) --5
後置演算子(postfix operator) foobar++
中置演算子(infix operator) 5 + 5 演算子が二つのオペランドを持つ。とも言える。

演算子の優先順位(operator precedance, order of operations)
5 + 5 * 10　みたいに掛け算を先に計算する的なやつ。演算子の次にくるオペランドが演算子にどの程度くっつくかを

## Prattのポイント

トークンタイプごとに構文解析関数を関連付ける。
トークン: 解析関数のマップを作っておく。

構文解析において、statement(文)はletとreturnしかない。他のものが来たらexpression(式)とする
foobar;

整数リテラルは式。そのものが値を生成している。
let x = 5;

## 前置演算子

```-5```は<prefix operator><expression>;　いかなる式でも前置演算子の後に来て、オペランドになれる。
```5 + -add(5, 5)```や```!isGreaterThanZero(3);```も有効
