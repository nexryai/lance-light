all: clean llfctl

all-arch:
	GOOS=linux GOARCH=amd64 go build -o llfctl-amd64 main.go
	GOOS=linux GOARCH=arm go build -o llfctl-arm main.go
	GOOS=linux GOARCH=arm64 go build -o llfctl-arm64 main.go
	GOOS=linux GOARCH=mips go build -o llfctl-mips main.go
	GOOS=linux GOARCH=mips64 go build -o llfctl-mips64 main.go
	GOOS=linux GOARCH=riscv64 go build -o llfctl-riscv64 main.go
	GOOS=linux GOARCH=ppc64 go build -o llfctl-ppc64 main.go
	GOOS=linux GOARCH=ppc64le go build -o llfctl-ppc64le main.go
	GOOS=linux GOARCH=s390x go build -o llfctl-s390x main.go

install:
	mv llfctl /usr/bin
	cp systemd/lance.service /etc/systemd/system/
	cp config.default.yml /etc/lance.yml
	systemctl daemon-reload

clean:
	rm -f llfctl-amd64 llfctl-arm llfctl-arm64 llfctl-mips llfctl-mips64 llfctl-riscv64 llfctl-ppc64 llfctl-ppc64le llfctl-s390x
