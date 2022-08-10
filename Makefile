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