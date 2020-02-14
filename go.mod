module github.com/assetsadapterstore/velas-adapter

go 1.12

require (
	github.com/DataDog/zstd v1.4.4 // indirect
	github.com/Sereal/Sereal v0.0.0-20200210135736-180ff2394e8a // indirect
	github.com/asdine/storm v2.1.2+incompatible
	github.com/astaxie/beego v1.12.0
	github.com/blocktree/go-owcdrivers v1.2.0
	github.com/blocktree/go-owcrypt v1.1.2
	github.com/blocktree/openwallet v1.7.0
	github.com/btcsuite/btcutil v0.0.0-20191219182022-e17c9730c422
	github.com/ethereum/go-ethereum v1.9.9
	github.com/go-errors/errors v1.0.1
	github.com/shopspring/decimal v0.0.0-20200105231215-408a2507e114
	gopkg.in/resty.v1 v1.12.0
)

replace (
	golang.org/x/crypto v0.0.0-20190404164418-38d8ce5564a5 => github.com/golang/crypto v0.0.0-20190404164418-38d8ce5564a5
	golang.org/x/net v0.0.0-20181220203305-927f97764cc3 => github.com/golang/net v0.0.0-20181220203305-927f97764cc3
)
