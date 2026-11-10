# nodeprobe — build, test and install the agent plus its systemd unit.
#
# `make install` writes outside the repo and needs root. Build first as your
# own user, so root needs no Go toolchain:
#     make build
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

# Deliberately does not depend on `build`: run `make build` as your user,
# then `sudo make install`, so root never needs Go on its PATH.
install:
	@test -x $(BINARY) || { echo "$(BINARY) not found - run 'make build' first (as your own user)"; exit 1; }
	install -m 0755 $(BINARY) $(BINDIR)/nodeprobe
	install -m 0644 $(UNIT) $(UNITDIR)/nodeprobe.service
	systemctl daemon-reload

clean:
	rm -rf bin
