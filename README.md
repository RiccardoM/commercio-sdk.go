# `commercio-sdk.go`

Commercio.network SDK written in Go.

# Warning! 

This software is work-in-progress, no guarantees on interface stability and feature completeness.

# Install

```bash
go get github.com/commercionetwork/commercio-sdk.go
```

# Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/commercionetwork/commercio-sdk.go"
)

func mightFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	m := "cover safe brass same salad raccoon expect rigid service brush ski amateur sample emerge actress oblige camp business three awkward absent peasant kitchen pool"

	c := commercio.DefaultConfig

	// only used for demonstration purposes
	c.Mode = commercio.TxModeBlock

	s, err := commercio.NewSDK(m, commercio.DefaultConfig)
	mightFatal(err)

	jack, err := commercio.Address("did:com:1l9rr5ck7ed30ny3ex4uj75ezrt03gfp96z7nen")
	mightFatal(err)

	janet, err := commercio.Address("did:com:1zla8arsc5rju9wekz00yz54zguj20a96jn9cy6")
	mightFatal(err)

	send := commercio.MsgSend{
		FromAddress: jack,
		ToAddress:   janet,
		Amount:      commercio.Amount(1000),
	}

	hash, err := s.SendTransaction(send)
	mightFatal(err)

	fmt.Println("transaction hash:", hash)
}

```
