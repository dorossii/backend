
dorossii

認証方法のテンプレートをまとめたリポジトリ

## 使用技術

- 言語: Golang
- 環境: Docker
- DB: MySQL
- ダッシュボード: React
- リバースプロキシ: nginx

## セットアップ方法

1. [taskfile](https://taskfile.dev/installation/) をインストールする
2. ``task --version`` を実行して taskfile を確認する
3. configs 配下にある *env_temolate をコピーして *.env を作成する
   各種 Secret などは [こちら](https://www.graviness.com/app/pwg/?l=64&n=1&m=1&r=3&s=1&c=0-9A-Za-z!%22%23%24%25%26'()*%2B%2C%5C-.%2F%3A%3B%3C%3D%3E%3F%40%5B%5C%5D%5E_%60%7B%7C%7D~) などで作成しておく
4. ``task setup`` を実行する
   注意: すでに設置アップ済みの場合データベースの中身が削除されます
5. mysqlコンテナの起動が完了したら データベースコンテナを再起動する
   ``task restart``
6. https://localhost:8947/auth/_/ にアクセスして管理ユーザーを作成する
7. ダッシュボードで各種プロバイダの設定をする
8. https://localhost:8947/statics/ で確かめてみる
9. 終わり

## セットアップの動作

- alpine と openssl コンテナが起動します
  - nginx 用の自己証明書が発行されます
  - jwt 用の秘密鍵が発行されます
  - jwt 用の公開鍵が発行されます
- mysql コンテナが起動します
- auth コンテナが起動します
- app コンテナが起動します
- nginx コンテナが起動します

## ディレクトリ構成

- configs : 設定ファイル .env が格納されています
- database : 設定ファイル my.cnf が格納されています
- openssl : nginx 周りの jwt 秘密鍵や公開鍵が格納されています
- nginx : nginx 周りの設定ファイルが格納されています

## 各種コマンド

- ``task setup`` : セットアップ
- ``task clean`` : コンテナ落として全て削除
- ``task down`` : コンテナ落とす
