package eigenda_disperse_blob

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"testing"
)

func Test_Disperse(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	pri := fmt.Sprintf("%x", crypto.FromECDSA(privateKey))
	ctx := context.Background()
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})
	conn, err := grpc.Dial("disperser.eigenda.xyz:443", grpc.WithTransportCredentials(creds))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	disperse := disperser.NewDisperserClient(conn)
	data := []byte("hello world")
	req, err := Disperse(ctx, disperse, pri, data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("header hash: %s, blob index: %d", hex.EncodeToString(req.BatchHeaderHash), req.BlobIndex)
}
