# nodeprobe — build, test and install the agent plus its systemd unit.
#
# `make install` writes outside the repo and needs root:
#     sudo make install

BINARY  := bin/nodeprobe
BINDIR  := /usr/local/bin
UNITDIR := /etc/systemd/system
UNIT    := deploy/systemd/nodeprobe.service

.PHONY: build test install clean

build:
	go build -o $(BINARY) ./cmd/nodeprobe

test:
	go test ./... -race

install: build
	install -m 0755 $(BINARY) $(BINDIR)/nodeprobe
	install -m 0644 $(UNIT) $(UNITDIR)/nodeprobe.service
	systemctl daemon-reload

clean:
	rm -rf bin
