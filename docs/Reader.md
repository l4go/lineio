# type Reader
selectでの非同期処理に対応した、io.Readerからの行ベースでの読み取り処理を行うライブラリです。

## import
```go
import "github.com/l4go/lineio"
```
vendoringして使うことを推奨します。

## 利用サンプル

- [NewReader() example](../examples/ex_lineio/ex_lineio.go)
- [NewReaderByDelim() example](../examples/ex_lineio_nul/ex_lineio_nul.go)

## メソッド概略

### func NewReader(rio io.Reader) \*Reader
Readerを生成します。行の区切り文字には、'\n'(0x0A)が使われます。

### func NewReaderByDelim(rio io.Reader, delim byte) \*Reader
行区切りの文字(byte)を指定して、 Readerを生成します。

### func (r \*Reader) Recv() <-chan \[\]byte

読み取ったデータを行単位で返すchannelを返します。  
読み取り完了もしくは、読み取りのエラーが発生した場合に、channelがクローズ状態になります。
エラーの発生は、クローズ状態のあと、**Err()**メソッドを利用して取得します。

### func (r \*Reader) Err() error
読み込みがcloseした場合(EOFやErrClosed)には、nilを返します。それ以外の場合は、該当のエラーの値を返します。
読み取りが正常終了を判定するために使います。

### func (r \*Reader) Close()
Readerを開放するための後処理をします。
