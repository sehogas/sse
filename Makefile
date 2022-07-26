########################################################
override TARGET=sse
VERSION=1.0
OS=linux
ARCH=amd64
FLAGS="-s -w"
CGO=0
########################################################

run:
	@echo Ejecutando programa...
	go run main.go

bin:
	@echo Generando binario ...
	CGO_ENABLED=$(CGO) GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags=$(FLAGS) -o $(TARGET) .

install: 
	@echo Instalando binario binario ...
	CGO_ENABLED=$(CGO) GOOS=$(OS) GOARCH=$(ARCH) go install -ldflags=$(FLAGS) 

build:
	@echo Construyendo imagen docker $(TARGET):$(VERSION) ...
	docker build -t $(TARGET):$(VERSION) .
	docker tag $(TARGET):$(VERSION) $(TARGET):latest

start:
	@echo Ejecutando contenedor docker $(TARGET):$(VERSION) ...
	docker run --rm -d --name $(TARGET) -p 3003:3003 $(TARGET):latest

stop:
	docker stop $(TARGET)

clean:
	@echo Borrando binario ...
	rm -rf $(TARGET)

.PHONY: clean run install build start stop
.DEFAULT: 
	@echo 'No hay disponible ninguna regla para este destino'