package rpc

import (
	"testing"
)

func TestBlock_GetByHash(t *testing.T) {
	type fields struct {
		baseAddress string
	}
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *BlockResponse
		wantErr bool
	}{
		{
			name:   "Normal test",
			fields: fields{baseAddress: Url},
			args: args{
				hash: "9e546bc9b274f08d053be29643cc76e7cc69ebb9fe8ba6465014ebe4c339b79a",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.baseAddress)
			got, err := client.Block.GetByHash(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("nil response")
			}
			/*if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByHash() got = %v, want %v", got, tt.want)
			}
			*/
			if got.Header.TxnCount != uint32(len(got.Transactions)) {
				t.Error("incorrect transactions count")
			}
			t.Logf("%+v", *got)
			t.Logf("%+v", got.Header)
			for _, tx := range got.Transactions {
				t.Logf("%+v", tx)
			}
		})
	}
}
