online-election-manager: ./cmd/online-election-manager/main.go ./cmd/initialize-database/main.go
	go install ./cmd/initialize-database
	go install ./cmd/online-election-manager

.PHONY: clean

clean:
	rm -f $(GOPATH)/bin/online-election-manager
	rm -f $(GOPATH)/bin/initialize-database