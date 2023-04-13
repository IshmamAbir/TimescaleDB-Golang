### ステップ 1：Docker と TimescaleDB のインストール

まず、ローカルマシンに Docker と TimescaleDB をインストールする必要があります。以下がその手順です：

マシンに Docker をインストールします。公式サイトから Docker をダウンロードできます：https://www.docker.com/get-started

Docker がインストールされたら、次のコマンドを実行して TimescaleDB Docker イメージをダウンロードします：

```
docker pull timescale/timescaledb:latest-pg15
```

### ステップ 2：TimescaleDB コンテナの起動

TimescaleDB Docker イメージをダウンロードした後、次のコマンドでコンテナを起動できます：

```
docker run --name timescaledb -e POSTGRES_PASSWORD=password -d -p 5432:5432 --restart always timescale/timescaledb:latest-pg15
```

このコマンドは、パスワードが password の timescaledb という名前の TimescaleDB コンテナを開始し、ポート 5432 をホストマシンのポート 5432 にマップします。

### ステップ 3：TimescaleDB コンテナへの接続

コンテナを起動した後、次のコマンドでそれに接続できます：

```
docker exec -it timescaledb psql -U postgres
```

このコマンドは、コンテナ内部で PostgreSQL シェルを開き、テーブルを作成し SQL コマンドを実行できます。

### ステップ 4：PostgreSQL データベースとテーブルの作成

1. First run the timescale database inside your docker container

   ```
   docker exec -it timescaledb psql -U postgres
   ```

2. PostgreSQL データベースを作成するには、次の SQL コマンドを実行できます：

   ```
   CREATE DATABASE timescale_test;
   ```

3. PostgreSQL テーブルを作成するには、最初に次のデータベースに移動する必要があります。
   ```
   \c timescale_test
   ```
4. これで、次の SQL コマンドを実行できます：

   ```
   CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       name TEXT NOT NULL,
       email TEXT NOT NULL UNIQUE,
       password TEXT NOT NULL
   );
   ```

このコマンドは、id、name、email、および password 用の列を持つ users テーブルを作成します。id 列は自動増分の主キーであり、email 列には各メールアドレスが 1 度しか使用されないことを保証する一意の制約があります。

### ステップ 5：TimescaleDB ハイパーテーブルの作成

TimescaleDB ハイパーテーブルを作成するには、次の SQL コマンドを実行できます：

```
CREATE TABLE sensor_data (
    time TIMESTAMPTZ NOT NULL,
    sensor_id TEXT NOT NULL,
    temperature DOUBLE PRECISION,
    humidity DOUBLE PRECISION,
    PRIMARY KEY (time, sensor_id)
);
```

```
SELECT create_hypertable('sensor_data', 'time');
```

This command creates a sensor_data table with columns for time, sensor_id, temperature, and humidity. The time and sensor_id columns are used as the primary key, and the create_hypertable function is used to create a TimescaleDB hypertable for the sensor_data table.

## プロジェクトへの接続

データベースが作成されました。これを使用する準備ができました。プロジェクト内で、以下のフォーマットで文字列を記述して接続文字列変数を宣言してください。

### Timescale クラウドの場合：

Web ダッシュボードから'service URL'をコピーしてください

### ローカルホストの場合（Docker から使用する場合）

以下のフォーマットで文字列を記述してください。

```
postgres://username:password@host:port/dbname?sslmode=disable
```

<br>
これで、プロジェクトを接続して実行することができます。
