

# github.com/prometheus/client_golang でメトリクスを取る

この記事はGoアドベントカレンダー向けの記事です。

この記事では `github.com/prometheus/client_golang` を使って
基本的なメトリクスを取ると同時に、HTTPハンドラのメトリクスを追加してみます。
また、HTTPハンドラのメトリクスで、99.9 percentileのメトリクスを取れるようにしてみます。

## 基本的なメトリクスの取得とエンドポイントの公開

[`github.com/prometheus/client_golang`](https://github.com/prometheus/client_golang) を使えばメトリクスを取ることが出来ます。
prometheus向けのメトリクスのエンドポイントも生やせます。

下のような形で簡単にメトリクスのエンドポイントを生やす事ができます。

```go
package main

import (
    "log"
    "net/http"

    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    http.Handle("/metrics", promhttp.Handler())
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

デフォルトでは、以下の2つのメトリクスのコレクターが登録されています。

* Process Collector
  * CPU, メモリ, ファイルディスクリプタなどのメトリクスが取れます。
* Go Collector
  * GCやゴルーチン、OSのスレッド数などのメトリクスが取れます。

以下の記事が詳しいです。

* Exploring Prometheus Go client metrics - https://povilasv.me/prometheus-go-metrics/

## HTTPハンドラのメトリクスを取る

HTTPハンドラのメトリクスを取ってみました。
以下のようなメトリクスが取れます。

デフォルトでは、50 percentile, 90 percentile, 99 percentileのメトリクスが取れます。

```javascript
http_request_duration_microseconds{handler="hello_world",quantile="0.5"} 82897.181
http_request_duration_microseconds{handler="hello_world",quantile="0.9"} 96527.057
http_request_duration_microseconds{handler="hello_world",quantile="0.99"} 98644.142
http_request_duration_microseconds_sum{handler="hello_world"} 1.0489467330000002e+06
http_request_duration_microseconds_count{handler="hello_world"}
```

ソースは下のようになりました。

```go
package main

import (
	"time"
	"log"
	"net/http"
	"fmt"
	"math/rand"
	
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", prometheus.InstrumentHandlerFunc("hello_world", hello))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
	
func hello(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello World")
}
```

簡単ですね。

### 99.9 percentileのメトリクスを取ってみる

99.9 percentileのメトリクスを取りたい時はどうすれば良いでしょうか。
ライブラリの中のコードを参考にすると、下のような形で取ることが出来ます。

```go
	http.Handle("/", prometheus.InstrumentHandlerFuncWithOpts(
		prometheus.SummaryOpts{
			Subsystem:   "http",
			ConstLabels: prometheus.Labels{"handler": "hello_world"},
			Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001, 0.999: 0.0001},
		},
		hello))
```

メトリクスは以下のような形で取れます。
増えましたね。

```diff
http_request_duration_microseconds{handler="hello_world",quantile="0.5"} 50336.671
http_request_duration_microseconds{handler="hello_world",quantile="0.9"} 91294.148
http_request_duration_microseconds{handler="hello_world",quantile="0.99"} 100699.414
+ http_request_duration_microseconds{handler="hello_world",quantile="0.999"} 104040.94
```

## まとめ

`github.com/prometheus/client_golang` を使うことで簡単なメトリクスがデフォルト設定されていることを確認しました。
また、InstrumentHandlerFuncやInstrumentHandlerFuncWithOptsなどを使うことで
HTTPエンドポイントのメトリクスが自動で取得されることが分かります。

エンドポイント毎のレスポンスタイムの劣化などは非常に便利な情報です。
HTTPサーバを書いているなら、メトリクスを取って監視してみてはいかがでしょうか。

この記事では紹介しませんが
他にも自分でコードを書くことでメトリクスを取ることが可能です。

終わり。