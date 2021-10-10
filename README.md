
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
 * let add = func(x, y){return x + y};
 * (fn(x) { return x }(5) + 10 ) * 10

if式もある。resultはtrueになる。
 * let result = if(10>5){ true} else {false};


### 式構文解析の用語
 * 前置演算子(prefix operator) --5
 * 後置演算子(postfix operator) foobar++
 * 中置演算子(infix operator) 5 + 5 演算子が二つのオペランドを持つ。とも言える。

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

## 中置演算子(infix expression)

```5 + 6;```オペランドが2つ。+がOperator, 5はleft, 6がrightという構造にする。

まず5をIDENTとして解析。  
次のTokenを見て、中置演算子かどうかを見極める。見極めには優先順位(precedance)を使う。  
```5```は優先順位1、```+``` は優先順位4になるため、中置演算子を含む文とみなしている。  
```＋```はpeekPrecedance()にて次のTokenの優先順位を覗き見している。
優先順位で決める部分を本書では「くっつきやすさ」と表現していた。

```go
func (p *Parser) parseExpression(precedance int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParserFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedance < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedences := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedences)
	return expression
}
```
優先順位は自分より高いものが次のトークンとして出現する限り読み出しが続くことになる。
1 + 2 + 3;

0. 
parseSatement > parseExpressionStatement > parseExpressionと進む
parseExpressionは最初優先順位1(最低)で入る。

1. 
"1"と"+"はcurTokenとpeekToken。
"1"はprefixをつけられるIDENTなので、parseIntegerLiteralを通してast.IntegerLiteralとなる。一旦leftExpへ入れる。
中置演算子の分を解析するループ入り口。curとpeekの優先順は 1("1") < 4("+")なので、InfixExpressionが期待されループ内へ。(leftは1。operatorは+)。トークンを先に進めて、rightにくるトークンをleftExp = infix(leftExp)で読む。infixはparseInfixExpressionが呼び出される。
leftExp = infix(leftExp)は　既存のleftExp("1")をleftにして、新たにRightを生やしたInfixExpressionを返すのがポイント。

2. 
parseInfixExpression内部で、新たにInfixExpresionが生成され,leftは"1" operatorは"+"となる。
leftとoperatorを引数として受け取りつつ、新たにInfixExpresionを生成して、rightとする方式になっている。
operator"+"の優先順位4を保存しつつ、Tokenを読む。Rightは2が入ることになる。
Rightを決めるべく呼び出されたp.parseExpression内では curToken"2"なので4("2") < 4("+")　でループは回らない。
結果2がRightとなって、1.のInfixExpressionが完成(Leftは1、Operatorは+、Rightは2)。1のparseExpresionはreturn。
図2-6
leftExp = infix(leftExp)でleftExpに完全なInfixExpresionが入って、ループ最終行

```go
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// この関数で新たにInfixExpressionを生成し、引数で受け取ったExpressionは左に格納する。
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	// 1 + 2 + 3が(1 + 2) + 3 と左からくっつくように左の優先順位("-")を右のトークンを読み進める時に使えるよう保存しておく。
	// この時点でcurTokenは1つめの"+"
	precedences := p.curPrecedence()

	p.nextToken()
	// この時点でcurTokenは2、だがprecedensesは4としてRightを読む
	expression.Right = p.parseExpression(precedences)
	return expression
}
```

3. 
ループ先頭に戻る。
curTokenは"2"のまま1の続き。peekTokenは+なので、1("2") < 4("+") 
leftExp にはInfixExpression(1 + 2)が入っている。このままinfix(leftExp)でRightをparseInfixExpressionで読む。
parseInfixExpressionにて、LeftにInfixExpression(1 + 2),Operatorを+,rightに3を入れる。
parseInfixExpressionはreturn。leftExpに(1 + 2) + 3が出来上がる。

4. 
ループは
セミコロンなので終わり


## 関数リテラル

```fn (<parameters>)<block statement>```
```<parameters>```は```(<params>, <params>, ....)```で空になることもある。

fn() {return foobar + barfoo;}
let myFunction = fn(x, y) { return x + y; }

returnでさらに関数を呼び出すこともできる。

```
fn() {
	return fn(x, y) {
		return x > y;
	};
}
```
パラメーターの1つとして、関数リテラルも使える。
```myFunc(x, y, fn(x, y) { return x > y; });```


### 関数呼び出し
add(1, 2)
add(1, 1 + 2 + 3)
addは関数なので、以下も有効となる
fn(x, y){ x * y;}(1, 2)

関数は値を返すので、これも有効
callsFunction(2, 3, fn(x, y){ x * y;})

構造は
<expression>(<comma separated expression>)

# 評価
インタプリタとコンパイラ
コンパイラは実行可能な成果物を残すという考え方もあるが、実用的なプログラミング言語ではそうでもない。
ASTの扱い方で典型的なのは、そのまま解釈すること「tree walking 型。実装によって再起や繰り返しを実行するのにより適した中間表現(IR; intermidiate representation)に変換したりする。
前もって AST を辿りバイトコードに変換するインタプリタもある。
バイトコードは ネイティブの機械語ではないし、アセンブリ言語でもない。
OSやインタプリタが動作している CPU上では実行できない。そうではなく、インタプリタの一部である仮想マシンで解釈される。
ソースコードを構文解析し、ASTを構築しバイトコードに変換する。バイトコードで規定された命令を実行の直前に仮想マシンがジャストインタイムでネイティブの機械語に変換するJIT(Just In Time)インタプリタ/コンパイラというものもある。

Monkeyはtree walkingインタプリタとする。
tree-walking評価器とホスト言語GoでMonkeyを評価(eval)する方法、2つが必要になる。

```js

function eval(astNode) {
	if (astNode is integerliteral) {
		return astNode.integerValue
	} else if (astNode is booleanLiteral) {
		return astNode.booleanValue
	} else if (astNode is infixExpression) {
		leftEvaluated = eval(astNode.Left)
		rightEvaluated = eval(astNode.Right)
		if astNode.Operator == "+" {
			return leftEvaluated + rightEvaluated
		} else if ast.Operator == "-" {
			return leftEvaluated - rightEvaluated
		}
	}
}
```

## evalが何を返すか？

ASTが表現する値や、ASTを評価した際メモリ上に生成する値を表現できるシステムが必要となる。
例えば、```let a = 5;```のあと、```a + a```を実行するとき、aにある値にアクセスする必要がある。

一般にインタプリタ言語において、値の内部表現を構築するには、様々な選択肢がある。
 * ホスト言語のネイティブ型(整数、真偽値など)をそのまま使う
 * 値やオブジェクトはポインタとして表現する。
 * ネイティブ型とポインタを混在して用いる

3.5.4
!trueはfalseを返すが、!5が何を返すかは言語自体のデザインになる。今回は5はtruthyな振る舞いに。

3.5.6
中置式は真偽値を生成する場合と、そうでない場合がある。 ```5-1```, ```5 > 1```