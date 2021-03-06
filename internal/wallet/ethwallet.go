package wallet

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

// etherum wallet type
type EthWallet struct {
	path     string
	password string
	key      *keystore.Key
}

func NewEthWallet(path string, password string) Wallet {
	return &EthWallet{path: path, password: password, key: nil}
}

// return wallet address in hex mode
func (w *EthWallet) Address() string {
	if w.key == nil {
		return common.BigToAddress(common.Big0).Hex()
	}

	return w.key.Address.Hex()
}

func (w *EthWallet) PrivateKey() string {
	if w.key == nil {
		return ""
	}

	return "0x" + common.Bytes2Hex(crypto.FromECDSA(w.key.PrivateKey))
}

func (w *EthWallet) PublicKey() string {
	if w.key == nil {
		return ""
	}

	return "0x" + common.Bytes2Hex(crypto.FromECDSAPub(&w.key.PrivateKey.PublicKey))
}

func (w *EthWallet) FilePath() string {
	return w.path
}

func (w *EthWallet) Password() string {
	return w.password
}

// generate key for wallet
func (w *EthWallet) GenerateKey() error {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return err
	}

	UUID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	w.key = &keystore.Key{
		Id:         UUID,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}

	return nil
}

func (w *EthWallet) SetPrivateKey(key string) error {
	privateKey, err := crypto.ToECDSA(common.FromHex(key))
	if err != nil {
		return err
	}

	UUID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	w.key = &keystore.Key{
		Id:         UUID,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}

	return nil
}

func (w *EthWallet) SaveKey() error {
	if w.path == "" {
		return errors.New("Not set the key file path")
	}

	if w.key == nil {
		return errors.New("Not set the private key for wallet")
	}

	// encrypt private key to keystore
	keyjson, err := keystore.EncryptKey(w.key, w.password, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(w.path, keyjson, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (w *EthWallet) LoadKey() error {
	if w.path == "" {
		return errors.New("Not set the key file path")
	}

	keyjson, err := ioutil.ReadFile(w.path)
	if err != nil {
		return err
	}

	key, err := keystore.DecryptKey(keyjson, w.password)
	if err != nil {
		return err
	}

	w.key = key
	return nil
}

func (w *EthWallet) IsKeyFile(fi os.FileInfo) bool {
	// Skip editor backups and UNIX-style hidden files.
	if strings.HasSuffix(fi.Name(), "~") || strings.HasPrefix(fi.Name(), ".") {
		return false
	}

	// Skip misc special files, directories (yes, symlinks too).
	if fi.IsDir() || fi.Mode()&os.ModeType != 0 {
		return false
	}

	return true
}
