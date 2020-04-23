# wolhttp
wolhttp は Wake-on-LAN を行う HTTP サーバーです。

## 設定ファイル
wolhttp を動作させるには設定ファイルが必要です。設定ファイルは JSON となっており、`-c` コマンドライン引数を使うことで読み込ませることができます。以下に設定ファイルの例を示します。

```json
{
    "http_addr": ":8080",
    "target": {
        "name_1": {
            "hardware_addr": "00:11:22:33:44:55"
        },
        "name_2": {
            "network": "udp",
            "local_addr": "192.168.0.1:12345",
            "remote_addr": "255.255.255.255:9",
            "hardware_addr": "00:11:22:33:44:55"
        }
    }
}
```

設定項目は以下の通りです。
- `"http_addr"`: HTTP 待ち受けアドレス。省略時は `":http"`。
- `"target"`: オブジェクト。省略不可。
    - キーは、任意の名前。ただし、空文字は不可。
    - 値は、オブジェクト。
        - `"network"`: `"udp"` または `"udp4"`（IPv4 のみ）または `"udp6"`（IPv6 のみ）。省略時は `"udp"`。
        - `"local_addr"`: ローカルアドレス。省略時は自動的に決まります。
        - `"remote_addr"`: 送信先アドレス。省略時は `"255.255.255.255:9"`。
        - `"hardware_addr"`: MAC アドレス。省略不可。

## HTTP API
クエリパラメータで `"target"` の名前を指定すると、対応するパラメータで Wake-on-LAN が実行されます。以下に `curl` コマンドでの例を示します。

```sh
curl "https://127.0.0.1:8080/?target=name_1"
```
