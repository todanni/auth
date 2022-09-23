package keys

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	vault "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

const (
	signingKeysVaultPath = "signing-keys"
)

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
	privkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
	return string(privkeyPem)
}

func WriteKeyToVault(ctx context.Context, client *vault.Client, privateKey, keyId string) error {
	privateKeyMap := map[string]interface{}{
		keyId: privateKey,
	}

	_, err := client.KVv2(signingKeysVaultPath).Put(ctx, keyId, privateKeyMap)
	if err != nil {
		log.Fatalf("Unable to write secret: %v to the vault", err)
		return err
	}
	return nil
}

func GetSigningKey(ctx context.Context, client *vault.Client, keyID string) (*rsa.PrivateKey, error) {
	privateKeyData, err := client.KVv2(signingKeysVaultPath).Get(ctx, keyID)
	if err != nil {
		log.Fatalf(
			"Unable to read the secret from  vault: %v",
			err,
		)
		return nil, err
	}

	privateKeyString := privateKeyData.Data[keyID].(string)
	block, _ := pem.Decode([]byte(privateKeyString))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
