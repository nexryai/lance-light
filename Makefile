all: llfctl

llfctl:
	go build -o llfctl main.go

install:
	mv llfctl /usr/bin
	cp systemd/lance.service /etc/systemd/system/
	cp config.default.yml /etc/lance.yml

clean:
	rm -f llfctl
