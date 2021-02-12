# lloesche/valheim-server Docker image
![Valheim](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/Logo_valheim.png "Valheim")

Valheim Server in a Docker Container  

# Basic Docker Usage

The name of the Docker image is `lloesche/valheim-server`.

Volume mount the server config directory to `/config` within the Docker container.

If you have an existing world on a Windows system you can copy it from e.g.  
  `C:\Users\Lukas\AppData\LocalLow\IronGate\Valheim\worlds`  
to e.g.  
  `$HOME/valheim-server-config/worlds`  
and run the image with `$HOME/valheim-server-config` volume mounted to `/config` inside the container.

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

Warning: `SERVER_PASS` must be at least 5 characters long. Otherwise `valheim_server.x86_64` will refuse to start!

A fresh start will take several minutes depending on your Internet connection speed as the container will download the Valheim dedicated server from Steam (~1 GB).

Do not forget to modify `WORLD_NAME` to reflect the name of your world! For existing worlds that is the filename in the `worlds/` folder without the `.db/.fwl` extension.

If you want to play with friends over the Internet and are behind NAT make sure that UDP ports 2456-2458 are forwarded to the container host.
Also ensure they are publicly accessible in any firewall.

If your server name does not show up in the server list  a couple of minutes after startup you likely have a firewall issue.

There is more info in section [Finding Your Server](#finding-your-server).

For LAN-only play see section [Steam Server Favorites & LAN Play](#steam-server-favorites--lan-play)


# Environment Variables
| Name | Default | Purpose |
|----------|----------|-------|
|`SERVER_NAME` | `My Server` | Name that will be shown in the server browser |
|`SERVER_PORT` | `2456` | UDP start port that the server will listen on |
|`WORLD_NAME` | `Dedicated` | Name of the world without `.db/.fwl` file extension |
|`SERVER_PASS` | `secret` | Password for logging into the server - min. 5 characters! |
|`SERVER_PUBLIC` | `1` | Whether the server should be listed in the server browser (`1`) or not (`0`) |
|`UPDATE_INTERVAL` | `900` | How often we check Steam for an updated server version in seconds |
|`BACKUPS_INTERVAL` | `3600` | Interval in seconds between backup runs |
|`BACKUPS_DIRECTORY` | `/config/backups` | Path to the backups directory |
|`BACKUPS_MAX_AGE` | `3` | Age in days after which old backups are flushed |
|`BACKUPS_DIRECTORY_PERMISSIONS` | `755` | Unix permissions for the backup directory |
|`BACKUPS_FILE_PERMISSIONS` | `644` | Unix permissions for the backup zip files |
|`CONFIG_DIRECTORY_PERMISSIONS` | `755` | Unix permissions for the /config directory |
|`WORLDS_DIRECTORY_PERMISSIONS` | `755` | Unix permissions for the /config/worlds directory |
|`WORLDS_FILE_PERMISSIONS` | `644` | Unix permissions for the files in /config/worlds |


# Docker+systemd Example
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


# Finding Your Server
Once the server is up and running and the log says something like
```
02/09/2021 10:42:24: Game server connected
```
it can still be challenging to actually find the server.

There are three ways of getting to your server. Either using the Steam server browser, adding the IP manually or using the in-game `Community` server list.

## In-game
When in-game, click on `Join Game` and select `Community`. Wait for the game to load the list of all 4000+ servers.
Only 200 servers will be shown at a time so we will have to enter part of our server name to filter the view.
![in-game server browser](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/find1.png "in-game server browser")

## Steam Server Browser
When using the Steam server browser, in Steam go to `View -> Servers`. Click on `CHANGE FILTERS` and select Game `Valheim`.
Wait for Steam to load all 4000+ Servers then sort the `SERVERS` column by clicking on its title. Scroll down until you find your server.
![Steam server browser](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/find2.png "Steam server browser")
From there you can right-click it and add as a favourite.

Note that in my tests when connecting to the server via the Steam server browser I had to enter the server password twice. Once in Steam and once in-game.

## Steam Server Favorites & LAN Play
A third option within Steam is to add the server manually by IP. This also allows for LAN play without the need to open or forward any firewall ports.

Steps:
1) Within Steam click on `View -> Servers`
2) `FAVORITES`
3) `ADD SERVER`
4) Enter Server IP and port+1. So if the server is running on UDP port `2456` enter `ip:2457`
5) `FIND GAMES AT THIS ADDRESS...`
6) `ADD SELECTED GAME SERVER TO FAV...`

![Add server manually](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/find3.png "Add server manually")

Do not use the `ADD THIS ADDRESS TO FAVORITES` button at this point.

NOTE: Sometimes I will get the following error when trying to connect to a LAN server:
![Steam Server Browser Error](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/find4.png "Steam Server Browser Error")

In those cases it sometimes helped to add the server again, but this time using port `2456` and now pressing the `ADD THIS ADDRESS TO FAVORITES` button.
It will not generate a new entry in the favourites list but seemingly just update the existing one that was originally discovered on port `2457`.

Sometimes it also helps to press the `REFRESH` button and then immediately double click on the Server.

Overall LAN play via the Steam Server Browser has been a bit hit and miss for me while online play using the in-game search has resulted in the most consistent success.


# Admin Commands
Upon startup the server will create a file `/config/adminlist.txt`. In it you can list the IDs of all administrator users.

The ID of a user can be gotten either in-game by pressing ***F2***
![User ID in-game](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/admin2.png "User ID in-game")

or in the server logs when a user connects.
![User ID in logs](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/admin1.png "User ID in logs")

Administrators can press ***F5*** to open the in-game console and use commands like `ban` and `kick`.
![Kick a user](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/admin3.png "Kick a user")


# Synology Help
This is not an extensive tutorial, but I hope these screenshots can be helpful.
Beware that the server can use multiple GB of RAM and produces a lot of CPU load.

![Step 1](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/step1.png "Step 1")
![Step 2](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/step2.png "Step 2")
![Step 3](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/step3.png "Step 3")
![Step 4a](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/step4.png "Step 4a")
![Step 4b](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/step4b.png "Step 4b")
![Step 5](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/step5.png "Step 5")
![Step 6](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/step6.png "Step 6")
![Step 7](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/step7.png "Step 7")
