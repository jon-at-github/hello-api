GO_VERSION := 1.20
ARCH := amd64
GO_BINARY_LATEST := go$(GO_VERSION).linux-$(ARCH).tar.gz
TAG := $$(git describe --abbrev=0 --tags --always)
HASH := $$(git rev-parse HEAD)
DATE := $$(date +%Y-%m-%d.%H:%M:%S)
LDFLAGS := -w -X github.com/jon-at-github/hello-api/handlers.hash=$(HASH) -X github.com/jon-at-github/hello-api/handlers.tag=$(TAG) -X github.com/jon-at-github/hello-api/handlers.date=$(DATE)

setup: install-go init-go install-lint install-godog

install-go:
	sudo rm -rf /usr/local/go
	sudo tar -C /usr/local -xvf ${GO_BINARY_LATEST}
	rm ${GO_BINARY_LATEST}

init-go:
	echo 'export PATH=$$PATH:/usr/local/go/bin' >> $${HOME}.bashrc
	echo 'export PATH=$$PATH:$${HOME}/go/bin' >> $${HOME}.bashrc

upgrade-go:
	sudo rm -rf /usr/local/go
	wget "https://go.dev/dl/${GO_BINARY_LATEST}"
	sudo  tar -C /usr/local -xzf ${GO_BINARY_LATEST}
	rm ${GO_BINARY_LATEST}

build:
	go build -ldflags "$(LDFLAGS)" -o api cmd/main.go

test:
	go test ./... -coverprofile=coverage.out

coverage:
	go tool cover -func coverage.out | grep "total:" | \
	awk '{print ((int($$3) > 80) != 1) }'

report:
	go tool cover -html=coverage.out -o cover.html

check-format:
	test -z $$(go fmt ./...)

install-lint:
	sudo curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.41.1

copy-hooks:
	chmod +x scripts/hooks/
	cp -r scripts/hooks .git/.
install-godog:
	go install github.com/cucumber/godog/cmd/godog@latest
install-k8s-redis:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm install redis-cluster bitnami/redis --set password=$$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')
deploy:
	kubectl apply -f k8s
