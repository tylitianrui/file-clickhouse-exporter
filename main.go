package main

import (
	"log"

	"github.com/tylitianrui/file-clickhouse-exporter/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		// handle  err
		log.Fatal(err)
	}
}
