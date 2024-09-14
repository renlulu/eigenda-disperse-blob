package eigenda_disperse_blob

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/ethereum/go-ethereum/crypto"
	"time"
)

var (
	ErrEigenDADisperseFailed  = errors.New("disperse blob failed")
	ErrEigenDADisperseTimeout = errors.New("disperse blob timeout")
)

func Disperse(ctx context.Context, client disperser.DisperserClient, privateKey string) (*disperser.RetrieveBlobRequest, error) {
	disperseBlobReply, err := auth(ctx, client, privateKey)
	if err != nil {
		return nil, err
	}
	ticket := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ctx.Done():
			return nil, ErrEigenDADisperseTimeout
		case <-ticket.C:
			statusReply, err := client.GetBlobStatus(ctx, &disperser.BlobStatusRequest{
				RequestId: disperseBlobReply.RequestId,
			})
			if err != nil {
				return nil, err
			}
			switch statusReply.GetStatus() {
			case disperser.BlobStatus_CONFIRMED, disperser.BlobStatus_FINALIZED:
				return &disperser.RetrieveBlobRequest{
					BatchHeaderHash: statusReply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash(),
					BlobIndex:       statusReply.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
				}, nil
			case disperser.BlobStatus_FAILED:
				return nil, ErrEigenDADisperseFailed
			default:
				// waiting for confirmation
				continue
			}
		}
	}

}

func auth(ctx context.Context, client disperser.DisperserClient, privateKey string) (*disperser.DisperseBlobReply, error) {
	var r *disperser.DisperseBlobRequest
	priv, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	pulblicKey, _, err := GetPubkeyFromPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	r.AccountId = pulblicKey
	authClient, err := client.DisperseBlobAuthenticated(ctx)
	if err != nil {
		return nil, err
	}
	err = authClient.Send(&disperser.AuthenticatedRequest{
		Payload: &disperser.AuthenticatedRequest_DisperseRequest{
			DisperseRequest: r,
		},
	})
	if err != nil {
		return nil, err
	}
	ticket := time.NewTicker(time.Second * 10)
	defer ticket.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticket.C:
			// per grpc, this function is block utils it receives a message, or the steam is done, so we can
			// return the error here
			authReply, err := authClient.Recv()
			if err != nil {
				return nil, err
			}
			// the authReply is either AuthenticatedReply_BlobAuthHeader or AuthenticatedReply_DisperseReply
			authHeader, ok := authReply.Payload.(*disperser.AuthenticatedReply_BlobAuthHeader)
			if ok {
				buf := make([]byte, 4)
				binary.BigEndian.PutUint32(buf, authHeader.BlobAuthHeader.ChallengeParameter)
				hash := crypto.Keccak256(buf)
				signed, err := crypto.Sign(hash, priv)
				if err != nil {
					return nil, err
				}
				err = authClient.Send(&disperser.AuthenticatedRequest{
					Payload: &disperser.AuthenticatedRequest_AuthenticationData{
						AuthenticationData: &disperser.AuthenticationData{
							AuthenticationData: signed,
						},
					},
				})

			} else {
				disperseReply, ok := authReply.Payload.(*disperser.AuthenticatedReply_DisperseReply)
				if !ok {
					return nil, fmt.Errorf("invalid reply type")
				}
				return disperseReply.DisperseReply, nil

			}

		}

	}
}
