package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/clfs/perftfuzz/engine"
	_ "github.com/clfs/perftfuzz/engine/stockfish"
)

func main() {
	var (
		pathFlag = flag.String("path", "", "absolute path to the chess engine")
		nameFlag = flag.String("name", "", "name of the chess engine")
	)
	flag.Parse()

	if *pathFlag == "" || *nameFlag == "" {
		flag.Usage()
		return
	}

	ctx := context.Background()

	e, err := engine.New(ctx, *nameFlag, *pathFlag)
	if err != nil {
		log.Fatalf("new engine: %v", err)
	}
	defer e.Close()

	http.ListenAndServe(":8080", e)
}
