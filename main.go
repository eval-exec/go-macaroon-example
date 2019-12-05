package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/macaroon.v2"
	"log"
)

var err error

func MustNew(rootKey, id []byte, loc string) *macaroon.Macaroon {
	m, err := macaroon.New(rootKey, id, loc, macaroon.LatestVersion)
	if err != nil {
		panic(err)
	}
	return m
}
func main() {
	storageRootKey := []byte("storagekey1")
	storageId := []byte("storage machine 1")
	storageLoc := "Storage Beijing"
	StorageNode := MustNew(storageRootKey, storageId, storageLoc)

	StorageNodeOrder1 := StorageNode.Clone()
	inspectMacaroon(StorageNodeOrder1)
	caveats := map[string]bool{
		"/storage/order1": true,
	}

	check := func(cav string) error {
		log.Println(cav)
		if cav == "/storage/order1" && caveats[cav] == true {
			return nil
		}
		if cav == "minera" && caveats[cav] == true {
			return nil
		}
		log.Println(cav,"is cav")
		log.Println(caveats[cav])
		return fmt.Errorf("%s condition not met", cav)
	}
	for cav := range caveats {
		if err := StorageNodeOrder1.AddFirstPartyCaveat([]byte(cav)); err != nil {
			log.Fatal(err)
		}
	}
	inspectMacaroon(StorageNodeOrder1)

	minerRootKey := []byte("minerserverkey1")
	minerId := []byte("miner server1")
	minerLoc := "Miner China"

	if err = StorageNodeOrder1.AddThirdPartyCaveat(minerRootKey, minerId, minerLoc); err != nil {
		log.Fatalln(err)
	}
	log.Println("added third caveat")

	minerServer := MustNew(minerRootKey, minerId, minerLoc)

	inspectMacaroon(minerServer)


	minerServer.Bind(StorageNodeOrder1.Signature())

	log.Println("storagenode1:")
	inspectMacaroon(StorageNodeOrder1)

	if err = StorageNodeOrder1.Verify(storageRootKey, check, []*macaroon.Macaroon{minerServer}); err != nil {
		log.Fatalln(err)
	}

}

func inspectMacaroon(m *macaroon.Macaroon) {
	mj, err := m.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	mj, err = prettyprint(mj)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(mj))
}

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
func printSliccStr(ss []string) {
	for s := range ss {
		fmt.Println(s)
	}
}
