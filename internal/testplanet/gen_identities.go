// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information

// +build ignore

package main

import (
	"bytes"
	"context"
	"encoding/pem"
	"flag"
	"fmt"
	"go/format"
	"os"

	"storj.io/storj/pkg/peertls"
	"storj.io/storj/pkg/provider"
)

func main() {
	count := flag.Int("count", 5, "number of identities to create")
	out := flag.String("out", "identities_table.go", "generated file")
	flag.Parse()

	var buf bytes.Buffer
	buf.WriteString(`
		// Copyright (C) 2018 Storj Labs, Inc.
		// See LICENSE for copying information
		
		package testplanet
		
		var pregeneratedIdentities = NewIdentities(
	`)

	for k := 0; k < *count; k++ {
		fmt.Println("Creating", k)
		ca, err := provider.NewCA(context.Background(), 14, 4)
		if err != nil {
			panic(err)
		}

		identity, err := ca.NewIdentity()
		if err != nil {
			panic(err)
		}

		var chain bytes.Buffer
		err = peertls.WriteChain(&chain, identity.Leaf, ca.Cert)
		if err != nil {
			panic(err)
		}

		var keys bytes.Buffer
		err = peertls.WriteKey(&keys, identity.Key)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(&buf, "mustParsePEM(%q, %q),\n", chain.Bytes(), keys.Bytes())
	}

	buf.WriteString(`)`)

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	file, err := os.Create(*out)
	if err != nil {
		panic(err)
	}

	if _, err := file.Write(formatted); err != nil {
		panic(err)
	}

	if err := file.Close(); err != nil {
		panic(err)
	}
}

func encodeBlocks(blocks ...*pem.Block) ([]byte, error) {
	var buf bytes.Buffer
	for _, block := range blocks {
		if err := pem.Encode(&buf, block); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}