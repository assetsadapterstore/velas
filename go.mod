module github.com/assetsadapterstore/velas-adapter

go 1.12

require (
	github.com/asdine/storm v2.1.2+incompatible
	github.com/astaxie/beego v1.11.1
	github.com/blocktree/go-owcdrivers v1.0.16
	github.com/blocktree/go-owcrypt v1.0.1
	github.com/blocktree/openwallet v1.4.3
	github.com/btcsuite/btcutil v0.0.0-20190316010144-3ac1210f4b38
	github.com/ethereum/go-ethereum v1.8.25
	github.com/go-errors/errors v1.0.1
	github.com/imroc/req v0.2.3
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24
	github.com/tidwall/gjson v1.2.1
	golang.org/x/crypto v0.0.0-20190404164418-38d8ce5564a5
	gopkg.in/resty.v1 v1.12.0
)

replace (
	golang.org/x/crypto v0.0.0-20190404164418-38d8ce5564a5 => github.com/golang/crypto v0.0.0-20190404164418-38d8ce5564a5
	golang.org/x/net v0.0.0-20181220203305-927f97764cc3 => github.com/golang/net v0.0.0-20181220203305-927f97764cc3
)
