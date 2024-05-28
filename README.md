![giphy](https://github.com/dorodango-maker/tic-tac-toe/assets/86306494/8e849b57-27d0-494c-b0e1-0c2c37d99227)


# OXゲーム（tic-tac-toe）
Ebitengineで作成した最初のゲームです。
現在は2人対戦のみ対応。気が向いたらCPU対戦モード作るかも。

# 起動方法
GoとEbitengineがインストールされていない場合は、下記記事参考にインストールしてください。
https://zenn.dev/dorodango_maker/articles/12e67170e4dc72

```
$ cd 任意のディレクトリ
$ git clone git@github.com:dorodango-maker/tic-tac-toe.git
$ go run .
```

# 操作方法
- 左クリック：OXの配置
- 右クリック：（ゲーム終了時）ゲームリスタート

# ルール
それぞれの手番で4個目を設置する時に、一番古いOXが消えます。
（自分の手番になると次に消えるものが半透明になります。）
それ以外は普通のOXゲームと同じです。
