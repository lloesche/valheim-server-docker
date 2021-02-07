# valheim-server-docker
Valheim Server in a Docker Container

# Docker example

Volume mount the server config directory to
`/config` within the Docker container.

If you have an existing world on a Windows system you
can copy it from e.g.
  `C:\Users\Lukas\AppData\LocalLow\IronGate\Valheim\worlds`
to e.g.
  `$HOME/valheim-server-config/worlds`
and run the image.

Do not forget to modify `WORLD_NAME` to reflect the name of
your world!

If behind NAT make sure that UDP ports 2456-2458 are
forwarded to the container host.
Also ensure they are publicly accessible in any firewall.

If your server name does not show up in the server list 
a couple of minutes after startup you likely have a firewall
issue.

```
$ mkdir -p $HOME/valheim-server-config/worlds
# copy existing world
$ docker run -d \
    --name valheim-server \
    -p 2456-2458:2456-2458/udp \
    -v $HOME/valheim-server-config:/config \
    -e SERVER_NAME="My Server" \
    -e WORLD_NAME="Neotopia" \
    -e SERVER_PASS="secret" \
    lloesche/valheim-server
```

A fresh start will take several minutes depending on your
Internet connection speed as the container will download
the Valheim dedicated server from Steam.

# Docker+systemd example
Create an optional config file `/etc/sysconfig/valheim-server`
```
SERVER_NAME="My Server"
SERVER_PORT=2456
WORLD_NAME=Dedicated
SERVER_PASS=secret
SERVER_PUBLIC=1
```

Then enable the Docker container on system boot
```
$ sudo mkdir /etc/valheim
$ sudo curl -o /etc/systemd/system/valheim-server.service https://raw.githubusercontent.com/lloesche/valheim-server-docker/master/valheim-server.service
$ sudo systemctl daemon-reload
$ sudo systemctl enable valheim-server.service
$ sudo systemctl start valheim-server.service
```

# Updates
The container will check for Valheim server updates every 15 minutes.
If an update is found it is downloaded and the server restarted.
