# lloesche/valheim-server Docker image
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


# Backups
The container will on startup and periodically create a backup of the `worlds/` directory.

The default is once per hour but can be changed using the `BACKUPS_INTERVAL` environment variable.
The number is in seconds. Meaning for hourly backups set `BACKUPS_INTERVAL=3600`.

Default backup directory is `/config/backups/` within the container. A different directory can be set using the `BACKUPS_DIRECTORY` environment variable.
It makes sense to have this directory be a volume mount from the host.
Warning: do not make the backup directory a subfolder of `/config/worlds/`. Otherwise each backup will backup all previous backups.

By default 3 days worth of backups will be kept. A different number can be configured using `BACKUPS_MAX_AGE`. The value is in days.

Beware that backups are performed while the server is running. As such files might be in an open state when the backup runs.
However the `worlds/` directory also contains a `.db.old` file for each world which should always be closed and in a consistent state.


# Environment Variables
| Name | Default | Purpose |
|----------|----------|-------|
|`SERVER_NAME` | `My Server` | Name that will be shown in the server browser |
|`SERVER_PORT` | `2456` | UDP start port that the server will listen on |
|`WORLD_NAME` | `Dedicated` | Name of the world without `.db/.fwl` file extension |
|`SERVER_PASS` | `secret` | Password for logging into the server |
|`SERVER_PUBLIC` | `1` | Whether the server should be listed in the server browser (`1`) or not (`0`) |
|`UPDATE_INTERVAL` | `900` | How often we check Steam for an updated server version in seconds |
|`BACKUPS_INTERVAL` | `3600` | Interval in seconds between backup runs |
|`BACKUPS_DIRECTORY` | `/config/backups` | Path to the backups directory |
|`BACKUPS_MAX_AGE` | `7` | Age in days after which old backups are flushed |
|`BACKUPS_DIRECTORY_PERMISSIONS` | `755` | Unix permissions for the backup directory |
|`BACKUPS_FILE_PERMISSIONS` | `644` | Unix permissions for the backup zip files |
|`CONFIG_DIRECTORY_PERMISSIONS` | `755` | Unix permissions for the /config directory |
|`WORLDS_DIRECTORY_PERMISSIONS` | `755` | Unix permissions for the /config/worlds directory |
|`WORLDS_FILE_PERMISSIONS` | `644` | Unix permissions for the files in /config/worlds |


# Synology Help
This is not an extensive tutorial, but I hope these screenshots help somebody getting Valheim server running on their Synology.
Beware that the server can use multiple GB of RAM and produces a lot of CPU load.
![Step 1](https://raw.githubusercontent.com/lloesche/valheim-server-docker/master/misc/step1.png "Step 1")
![Step 2](https://raw.githubusercontent.com/lloesche/valheim-server-docker/master/misc/step2.png "Step 2")
![Step 3](https://raw.githubusercontent.com/lloesche/valheim-server-docker/master/misc/step3.png "Step 3")
![Step 4](https://raw.githubusercontent.com/lloesche/valheim-server-docker/master/misc/step4.png "Step 4")
![Step 5](https://raw.githubusercontent.com/lloesche/valheim-server-docker/master/misc/step5.png "Step 5")
![Step 6](https://raw.githubusercontent.com/lloesche/valheim-server-docker/master/misc/step6.png "Step 6")
![Step 7](https://raw.githubusercontent.com/lloesche/valheim-server-docker/master/misc/step7.png "Step 7")
