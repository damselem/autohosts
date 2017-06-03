.PHONY: crosscompile
crosscompile:
	GOOS=linux GOARCH=amd64 go build -o bin/autohosts-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -o bin/autohosts-darwin-amd64

.PHONY: release
release: crosscompile
	upx --ultra-brute bin/*
	git tag ${VERSION} && git push origin master --tags
	github-release release --user damselem --repo autohosts --tag ${VERSION} --name "autohosts ${VERSION}" --description "Version ${VERSION}" -s ${TOKEN}
	github-release upload --user damselem --repo autohosts --tag ${VERSION} --name "autohosts-linux-amd64" --file bin/autohosts-linux-amd64 -s ${TOKEN}
	github-release upload --user damselem --repo autohosts --tag ${VERSION} --name "autohosts-darwin-amd64" --file bin/autohosts-darwin-amd64 -s ${TOKEN}
