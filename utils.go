package eigenda_disperse_blob

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
)

func GetPubkeyFromPrivateKey(privateKey string) (string, string, error) {
	priv, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", "", err
	}
	publicKey := priv.PublicKey
	publicKeyBytes := crypto.FromECDSAPub(&publicKey)
	address := crypto.PubkeyToAddress(publicKey).Hex()

	return "0x" + hex.EncodeToString(publicKeyBytes), address, nil

}
