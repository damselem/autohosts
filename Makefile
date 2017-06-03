all:
	go build .

release: all
	upx --ultra-brute autohosts
	git tag ${VERSION} && git push origin master --tags
	github-release release --user damselem --repo autohosts --tag ${VERSION} --name "autohosts ${VERSION}" --description "Version ${VERSION}" -s ${TOKEN}
	github-release upload --user damselem --repo autohosts --tag ${VERSION} --name "autohosts" --file autohosts -s ${TOKEN}
