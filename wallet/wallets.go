package wallet

import (
	"os"
	"strings"
)

type Wallets struct {
	Wallets map[string]*Wallet
}

// func (ws *Wallets) SaveWalletToFile() {

// 	// var content bytes.Buffer

// 	// gob.Register(elliptic.P256())

// 	// encoder := gob.NewEncoder(&content)
// 	// err := encoder.Encode(ws)
// 	// if err != nil {
// 	// 	log.Panic(err)
// 	// }

// 	// err = os.WriteFile(walletFile, content.Bytes(), 0644)
// 	// if err != nil {
// 	// 	log.Panic(err)
// 	// }
// }

func (ws *Wallets) loadWalletsFromFile() error {

	entries, err := os.ReadDir(walletsFilePath)
	HandleErr(err)

	// iterate through all files
	for _, e := range entries {
		if !e.IsDir() {
			address := strings.Replace(e.Name(), ".wlt", "", 1)
			ws.Wallets[address] = LoadWalletFromFile(address)
		}
	}
	return nil
}

func LoadWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.loadWalletsFromFile()
	return &wallets, err
}

func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws *Wallets) AddWallet() string {
	newWallet := MakeWallet()
	// address := fmt.Sprint(newWallet.Address())
	address := string(newWallet.Address())

	ws.Wallets[address] = newWallet
	newWallet.SaveWalletToFile()

	return address
}
