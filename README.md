# GORM 編（Go言語）サンプル

## 使い方

### Docker環境の起動方法
- 筆者は、MacとWSLで確認しています。  
- Dockerなどは使えるようにしておいてください。  
- リポジトリからダウンロードしたら以下を実行してください。
    ```sh
    $ cd gorm-introduction  # ダウンロードして展開したディレクトリへ
    $ docker compose up  # コンテナ起動

    # 以下は、別タブか別ウインドウで
    $ docker ps  # 3つのコンテナが Up 状態になっていればOKです。
    $ docker exec -it gorm-go go mod tidy  # 必要なパッケージを取得します
    ```

### MySQLとphpMyAdmin
- DBはMySQLです。  
- DBの中身を確認したいときは、phpMyAdminを利用できるようにしています。(もちろん他のツールを使っても大丈夫です)  
    - Docker環境を起動して、Webブラウザで http://localhost:60080/ にアクセスすると利用できます。  
- usersテーブルの定義は、[init.sql](./initdb.d/init.sql)に記載しました。  

### サンプルコード[main.go](main.go)の実行方法
```sh
$ docker exec -it gorm-go go run main.go
    # または、
$ docker exec -it gorm-go ash  # コンテナのシェルを起動
$ go run main.go  # コンテナのシェルにて入力
```

