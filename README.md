# valheim-server-docker
Valheim Server in a Docker Container

# Docker example
```
$ docker run -d --rm \
    -p 27015-27050:27015-27050/tcp \
    -p 27015-27050:27015-27050/udp \
    -p 7777-7780:7777-7780/udp \
    -p 3478:3478/udp \
    -p 4379-4380:4379-4380/udp \
    -e MaxPlayers=64
    lloesche/valheim-server
```

# Docker+systemd example
Create an optional config file `/etc/sysconfig/valheim-server`
```
SERVER_NAME="My Server"
SERVER_PORT=2456
WORLD_NAME=Dedicated
SERVER_PASS=secret
SERVER_PUBLIC=0
```

```
$ sudo curl -o /etc/systemd/system/valheim-server.service https://raw.githubusercontent.com/lloesche/valheim-server-docker/master/valheim-server.service
$ sudo systemctl daemon-reload
$ sudo systemctl enable valheim-server.service
$ sudo systemctl start valheim-server.service
```
