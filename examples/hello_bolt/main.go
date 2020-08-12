package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/cayleygraph/quad"
	gateway "github.com/epik-protocol/epik-gateway-backend"
	"github.com/epik-protocol/epik-gateway-backend/graph"
	_ "github.com/epik-protocol/epik-gateway-backend/graph/kv/bolt"
)

func main() {
	// File for your new BoltDB. Use path to regular file and not temporary in the real world
	tmpdir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(tmpdir) // clean up

	// Initialize the database
	err = graph.InitQuadStore("bolt", tmpdir, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Open and use the database
	store, err := gateway.NewGraph("bolt", tmpdir, nil)
	if err != nil {
		log.Fatalln(err)
	}

	store.AddQuad(quad.Make("phrase of the day", "is of course", "Hello BoltDB!", "demo graph"))

	// Now we create the path, to get to our data
	p := gateway.StartPath(store, quad.String("phrase of the day")).Out(quad.String("is of course"))

	// This is more advanced example of the query.
	// Simpler equivalent can be found in hello_world example.

	ctx := context.TODO()
	// Now we get an iterator for the path and optimize it.
	// The second return is if it was optimized, but we don't care for now.
	its, _ := p.BuildIterator(ctx).Optimize(ctx)
	it := its.Iterate()

	// remember to cleanup after yourself
	defer it.Close()

	// While we have items
	for it.Next(ctx) {
		token := it.Result()                // get a ref to a node (backend-specific)
		value := store.NameOf(token)        // get the value in the node (RDF)
		nativeValue := quad.NativeOf(value) // convert value to normal Go type

		fmt.Println(nativeValue) // print it!
	}
	if err := it.Err(); err != nil {
		log.Fatalln(err)
	}
}
