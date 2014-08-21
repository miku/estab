SHELL := /bin/bash
TARGETS = estab

# http://docs.travis-ci.com/user/languages/go/#Default-Test-Script
test:
	go test -v

bench:
	go test -bench=.

imports:
	goimports -w .

fmt:
	go fmt ./...

vet:
	go vet ./...

all: fmt test
	go build

install:
	go install

clean:
	go clean
	rm -f coverage.out
	rm -f $(TARGETS)
	rm -f estab-*.x86_64.rpm
	rm -f debian/estab*.deb
	rm -rf debian/estab/usr

cover:
	go get -d && go test -v	-coverprofile=coverage.out
	go tool cover -html=coverage.out

estab:
	go build cmd/estab/estab.go

# ==== packaging

deb: $(TARGETS)
	mkdir -p debian/estab/usr/sbin
	cp $(TARGETS) debian/estab/usr/sbin
	cd debian && fakeroot dpkg-deb --build estab .

REPOPATH = /usr/share/nginx/html/repo/CentOS/6/x86_64

publish: rpm
	cp estab-*.rpm $(REPOPATH)
	createrepo $(REPOPATH)

rpm: $(TARGETS)
	mkdir -p $(HOME)/rpmbuild/{BUILD,SOURCES,SPECS,RPMS}
	cp ./packaging/estab.spec $(HOME)/rpmbuild/SPECS
	cp $(TARGETS) $(HOME)/rpmbuild/BUILD
	./packaging/buildrpm.sh estab
	cp $(HOME)/rpmbuild/RPMS/x86_64/estab*.rpm .
