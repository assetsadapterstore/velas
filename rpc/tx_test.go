package rpc

import (
	"reflect"
	"testing"
)

func TestTx_GetByHashList(t *testing.T) {
	type fields struct {
		bk *BaseClient
	}
	type args struct {
		hashes []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []TxResponse
		wantErr bool
	}{
		{
			name: "Normal test",
			args: args{
				hashes: []string{"c63d9c73a711f0f51b80d3e71b98d75170f1fa259cda5b5bb588700ac4a0f4f0",
					"dbec977c3c03530655c63f0debb92533c9d6fe9dfa9cf7ccf4b91752d8c3f920"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Tx{
				bk: tt.fields.bk,
			}
			bk := newBaseClient(Url)
			tx.bk = bk
			got, err := tx.GetByHashList(tt.args.hashes...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tx.GetByHashList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tx.GetByHashList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTx_GetHashListByHeight(t *testing.T) {
	type fields struct {
		bk *BaseClient
	}
	type args struct {
		height int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "get txs by height",
			args: args{
				height: 905281,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Tx{
				bk: tt.fields.bk,
			}
			bk := newBaseClient(Url)
			tx.bk = bk
			got, err := tx.GetHashListByHeight(tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tx.GetHashListByHeight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tx.GetHashListByHeight() = %v, want %v", got, tt.want)
			}
		})
	}
}
