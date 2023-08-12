.PHONY : deps
deps:
	go mod tidy
	go mod vendor


.PHONY : build
build:
	go  build  -o file-clickhouse-exporter  main.go 