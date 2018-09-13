package wallet

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrjson"
	"github.com/decred/dcrd/rpcclient"
)

type Wallet struct {
	client *rpcclient.Client
}

/*type TicketCacheRef struct {
	TxID string
	DbID uint64
}*/

func NewWalletClient(host, user, pass, cert string, disableTLS bool) *Wallet {
	return &Wallet{
		client: connectWalletRPC(host, user, pass, cert, disableTLS),
	}
}

func connectWalletRPC(host, user, pass, cert string, disableTLS bool) *rpcclient.Client {
	ntfnHandlers := rpcclient.NotificationHandlers{}

	certs, err := ioutil.ReadFile(cert)
	if err != nil {
		log.Fatal(err)
	}

	//Connect to local dcrwallet RPC server using websockets
	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		Endpoint:     "ws",
		User:         user,
		Pass:         pass,
		Certificates: certs,
		DisableTLS:   disableTLS,
	}
	client, err := rpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (scope *Wallet) GetUnspent(w http.ResponseWriter, r *http.Request) {
	unspent, err := scope.client.ListUnspent()
	if err != nil {
		log.Fatal(err)
	}
	scope.json(w, unspent)
}

func (scope *Wallet) GetAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := scope.client.ListAccounts()
	if err != nil {
		log.Fatal(err)
	}
	scope.json(w, accounts)
}

func (scope *Wallet) GetBalance(w http.ResponseWriter, r *http.Request) {
	balance, err := scope.client.GetBalance("*")
	if err != nil {
		log.Fatal(err)
	}
	scope.json(w, balance)
}

func (scope *Wallet) GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := scope.client.ListTransactionsCount("*", 999999)
	if err != nil {
		log.Fatal(err)
	}
	scope.json(w, transactions)
}

func (scope *Wallet) GetTransaction(txid string) *dcrjson.GetTransactionResult {
	txHash, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		log.Fatal(err)
	}
	transaction, err := scope.client.GetTransaction(txHash)
	if err != nil {
		log.Fatal(err)
	}
	return transaction
}

func (scope *Wallet) ListTransactions() []dcrjson.ListTransactionsResult {
	transactions, err := scope.client.ListTransactionsCount("*", 999999)
	if err != nil {
		log.Fatal(err)
	}
	return transactions
}

func (scope *Wallet) GetTickets(includeImmature bool) []*chainhash.Hash {
	tickets, err := scope.client.GetTickets(includeImmature)
	if err != nil {
		log.Fatal(err)
	}
	return tickets
}

/*
func GetTicketHashes(db *sql.DB, wallets []*rpcclient.Client) []*TicketCacheRef {
	tickets := make([]string, 0)
	ticketCache := make([]*TicketCacheRef, 0)
	for i, wallet := range wallets {
		theseTickets, err := wallet.GetTickets(true)
		if err != nil {
			log.Fatal(err)
		}
		for i, ticket := range theseTickets {
			tickets = append(tickets, ticket.String())
		}
	}
	dbIDs, err := *dcrpg.RetrieveTicketIDsByHashes(db, tickets)
	if err != nil {
		log.Fatal(err)
	}
	for i, dbID := range dbIDs {
		var thisPair = *TicketCacheRef{
			TxID: tickets[i],
			DbID: dbID,
		}
		ticketCache = append(ticketCache, thisPair)
	}
	return ticketCache
}
*/
func (scope *Wallet) json(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	d, err := json.Marshal(data)

	if err == nil {
		w.Write(d)
	}
}
