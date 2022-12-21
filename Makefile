GO_VERSION := 1.19
GO_ZIP_ARM64_LATEST := go$(GO_VERSION).linux-arm64.tar.gz

setup: install-go init-go

install-go:
	wget "https://go.dev/dl/${GO_ZIP_ARM64_LATEST}"
	sudo tar -C /usr/local -xvf ${GO_ZIP_ARM64_LATEST}
	rm ${GO_ZIP_ARM64_LATEST}

init-go:
	echo 'export PATH=$$PATH:/usr/local/go/bin' >> $${HOME}.bashrc
	echo 'export PATH=$$PATH:$${HOME}/go/bin' >> $${HOME}.bashrc

upgrade-go:
	sudo rm -rf /usr/bin/go
	wget "https://go.dev/dl/${GO_ZIP_ARM64_LATEST}"
	sudo  tar -C /usr/local -xzf ${GO_ZIP_ARM64_LATEST}
	rm ${GO_ZIP_ARM64_LATEST}

build:
	go build -o api cmd/main.go

test:
	go test ./... -coverprofile=coverage.out

coverage:
	go tool cover -func coverage.out | grep "total:" | \
	awk '{print ((int($$3) > 80 != 1) }'

report:
	go tool cover -html=coverage.out -o cover.html