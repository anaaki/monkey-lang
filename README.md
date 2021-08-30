
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
 * バッファリング、保存はしない。```NextToken()```を呼ぶことで、ソースコードを読み進めて、トークンを返す。
 * ソースコードはstringとしてあつかう。本当はファイル名や行番号があったほうがデバッグしやすいけど、シンプルにしたいので。

 