# Server Sent Events

Es una API creada en golang que permite enviar eventos a los clientes conectados por navegador.

La API tiene un método POST (/sendmessage) para enviarle mensajes a todos los clientes conectados.

Un ejemplo de cliente desarrollado en Angular puede encontrarlo en https://github.com/sehogas/sse-client


## Instalación

Descargue el código fuente e instálelo usando el comando `make install`.

En forma alternativa, use Docker para ejecutar el servicio en un contenedor:

```
make build
```
```
make start
```

## Instalar en systemd como binario

sudo cp sse.service /etc/systemd/system/

# Comandos systemd

sudo systemctl start sse
sudo systemctl start sse
sudo systemctl enable sse
sudo systemctl disable sse
sudo systemctl status sse


## Autor
* [Sebastian Hogas](https://github.com/sehogas)




