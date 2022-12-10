package proj3

import (
	"crypto/ecdsa"
	"net/http"
	"fmt"
	"github.com/libp2p/go-libp2p/core/peer"
	"log"
)

type pubKey = ecdsa.PublicKey
type privKey = ecdsa.PrivateKey

type localNotifee string

func (l localNotifee) HandlePeerFound(p peer.AddrInfo)  {
    fmt.Printf("handling the found peer, with notifee: %v", l)
}

func (c ChatServer) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
    err := r.ParseForm()
    if err != nil {
        log.Fatalln(err)
    }
	http.ListenAndServe(":8080", c)
}