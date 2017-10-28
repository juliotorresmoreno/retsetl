#/bin/sh

CC = go
CFLAGS = build
FROM = .
TO = ./bin
BIN = retsetl
INSTALL = /opt/retsetl
SERVICE_FROM = ./etc/retsetl.service
SERVICE_TO = /lib/systemd/system
USER=retsetl

all: $(OBJ)
	go get
	$(CC) $(CFLAGS) $(FROM)
	mv $(BIN) $(TO)

install: $(OBJ)
	@if [ `id -u` = 0 ]; then\
		`useradd retsetl --no-create-home --system`;\
		`mkdir $(INSTALL)`;\
		`mkdir $(INSTALL)/config`;\
		`cp $(TO)/$(BIN) $(INSTALL) -f`;\
		`cp ./config/config.json $(INSTALL)/config -f`;\
		`cp ./data $(INSTALL) -f -r`;\
		`cp ./graphiql $(INSTALL) -f -r`;\
		`cp ./public $(INSTALL) -f -r`;\
		`cp $(SERVICE_FROM) $(SERVICE_TO) -f`;\
		systemctl daemon-reload;\
		systemctl reenable retsetl.service;\
		chown $(USER):$($USER) $(INSTALL)/config/config.json;\
		echo "Installed application";\
		systemctl restart retsetl.service;\
		systemctl status retsetl.service;\
	else\
		echo "You are not root";\
    fi
