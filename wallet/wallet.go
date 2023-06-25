package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ripemd160"
)

const (
	checkSumLength  = 4
	version         = byte(0x00)
	walletsFilePath = "./tmp/wallets/"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func (w Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version}, pubHash...)
	checksum := CheckSum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)

	// fmt.Printf("pub key: %x\n", w.PublicKey)
	// fmt.Printf("pub hash: %x\n", pubHash)
	// fmt.Printf("addresss: %x\n", address)

	return address
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pub
}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()
	return &Wallet{private, public}
}

func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}

	publicRipMD := hasher.Sum(nil)

	return publicRipMD
}

func CheckSum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checkSumLength]
}

func (w *Wallet) SaveWalletToFile() {
	filename := walletsFilePath + string(w.Address()) + ".wlt"

	privateKeyBytes, err := x509.MarshalECPrivateKey(&w.PrivateKey)
	if err != nil {
		log.Panic(err)
	}
	privateKeyFile, err := os.Create(filename)
	if err != nil {
		log.Panic(err)
	}
	defer privateKeyFile.Close()
	err = pem.Encode(privateKeyFile, &pem.Block{
		Bytes: privateKeyBytes,
	})
	if err != nil {
		log.Panic(err)
	}
}

func LoadWalletFromFile(address string) *Wallet {
	filename := filepath.Join(walletsFilePath + address + ".wlt")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Panic(errors.New("no wallet with such address"))
	}
	privKeyFile, err := os.ReadFile(filename)
	HandleErr(err)
	pemBlock, _ := pem.Decode(privKeyFile)
	HandleErr(err)
	privKey, err := x509.ParseECPrivateKey(pemBlock.Bytes)
	HandleErr(err)
	publicKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	return &Wallet{
		PrivateKey: *privKey,
		PublicKey:  publicKey,
	}
}

func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-checkSumLength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checkSumLength]
	targetChecksum := CheckSum(append([]byte{version}, pubKeyHash...))

	return true
	return bytes.Compare(actualChecksum, targetChecksum) == 0
}
