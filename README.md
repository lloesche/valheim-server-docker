# lloesche/valheim-server Docker image
![Valheim](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/Logo_valheim.png "Valheim")

Valheim Server in a Docker Container (with [ValheimPlus](#valheimplus) support)  

[![Docker Badge](https://img.shields.io/docker/pulls/lloesche/valheim-server.svg)](https://hub.docker.com/r/lloesche/valheim-server)

# Basic Docker Usage

The name of the Docker image is `lloesche/valheim-server`.

Volume mount the server config directory to `/config` within the Docker container.

If you have an existing world on a Windows system you can copy it from e.g.  
  `C:\Users\Lukas\AppData\LocalLow\IronGate\Valheim\worlds`
to e.g.  
  `$HOME/valheim-server/config/worlds`
and run the image with `$HOME/valheim-server/config` volume mounted to `/config` inside the container.
The container directory `/opt/valheim` contains the downloaded server. It can optionally be volume mounted to avoid having to download the server on each fresh start.

```
$ mkdir -p $HOME/valheim-server/config/worlds $HOME/valheim-server/data
# copy existing world
$ docker run -d \
    --name valheim-server \
    -p 2456-2458:2456-2458/udp \
    -v $HOME/valheim-server/config:/config \
    -v $HOME/valheim-server/data:/opt/valheim \
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

There is more info in section [Finding Your Server](#finding-your-server).

For LAN-only play see section [Steam Server Favorites & LAN Play](#steam-server-favorites--lan-play)

For more deployment options see the [Deployment section](#deployment). 


# Environment Variables
| Name | Default | Purpose |
|----------|----------|-------|
| `SERVER_NAME` | `My Server` | Name that will be shown in the server browser |
| `SERVER_PORT` | `2456` | UDP start port that the server will listen on |
| `WORLD_NAME` | `Dedicated` | Name of the world without `.db/.fwl` file extension |
| `SERVER_PASS` | `secret` | Password for logging into the server - min. 5 characters! |
| `SERVER_PUBLIC` | `true` | Whether the server should be listed in the server browser (`true`) or not (`false`) |
| `SERVER_ARGS` |  | Additional Valheim server CLI arguments |
| `UPDATE_CRON` | `*/15 * * * *` | [Cron schedule](https://en.wikipedia.org/wiki/Cron#Overview) for update checks (disabled if set to an empty string or if the legacy `UPDATE_INTERVAL` is set) |
| `UPDATE_IF_IDLE` | `true` | Only run update check if no players are connected to the server |
| `RESTART_CRON` | `0 5 * * *` | [Cron schedule](https://en.wikipedia.org/wiki/Cron#Overview) for server restarts (disabled if set to an empty string) |
| `TZ` | `Etc/UTC` | Container [time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) |
| `BACKUPS` | `true` | Whether the server should create periodic backups (`true` or `false`) |
| `BACKUPS_CRON` | `0 * * * *` | [Cron schedule](https://en.wikipedia.org/wiki/Cron#Overview) for world backups (disabled if set to an empty string or if the legacy `BACKUPS_INTERVAL` is set) |
| `BACKUPS_DIRECTORY` | `/config/backups` | Path to the backups directory |
| `BACKUPS_MAX_AGE` | `3` | Age in days after which old backups are flushed |
| `PERMISSIONS_UMASK` | `022` | [Umask](https://en.wikipedia.org/wiki/Umask) to use for backups, config files and directories |
| `STEAMCMD_ARGS` | `validate` | Additional steamcmd CLI arguments |
| `VALHEIM_PLUS` | `false` | Whether [ValheimPlus](https://github.com/valheimPlus/ValheimPlus) mod should be loaded (config in `/config/valheimplus`) |

There are a few undocumented environment variables that could break things if configured wrong. They can be found in [`defaults`](defaults).


## Event hooks
The following environment variables can be populated to run commands whenever specific events happen.

| Name | Default | Purpose |
|----------|----------|-------|
| `PRE_BOOTSTRAP_HOOK` |  | Command to be executed before bootstrapping is done. Startup is blocked until this command returns. |
| `POST_BOOTSTRAP_HOOK` |  | Command to be executed after bootstrapping is done and before the server or any services are started. Can be used to install additional packages or perform additional system setup. Startup is blocked until this command returns. |
| `PRE_BACKUP_HOOK` |  | Command to be executed before a backup is created. The string `@BACKUP_FILE@` will be replaced by the full path of the future backup zip file. Backups are blocked until this command returns. See [Post backup hook](#post-backup-hook) for details. |
| `POST_BACKUP_HOOK` |  | Command to be executed after a backup is created. The string `@BACKUP_FILE@` will be replaced by the full path of the backup zip file. Backups are blocked until this command returns. See [Post backup hook](#post-backup-hook) for details. |
| `PRE_UPDATE_CHECK_HOOK` |  | Command to be executed before an update check is performed. Current update is blocked until this command returns. |
| `POST_UPDATE_CHECK_HOOK` |  | Command to be executed after an update check was performed. Future updates are blocked until this command returns. |
| `PRE_START_HOOK` |  | Command to be executed before the first server start is performed. Current start is blocked until this command returns. |
| `POST_START_HOOK` |  | Command to be executed after the first server start was performed. Future restarts and update checks are blocked until this command returns. |
| `PRE_RESTART_HOOK` |  | Command to be executed before a server restart is performed. Current restart is blocked until this command returns. |
| `POST_RESTART_HOOK` |  | Command to be executed after a server restart was performed. Future restarts and update checks are blocked until this command returns. |
| `PRE_SERVER_RUN_HOOK` |  | Command to be executed before the server is started. Server startup is blocked until this command returns. |
| `POST_SERVER_RUN_HOOK` |  | Command to be executed after the server has finished running. Server shutdown is blocked until this command returns or a shutdown timeout is triggered after 29 seconds. |


## ValheimPlus config from Environment Variables
ValheimPlus config can be specified in environment variables using the syntax `VPCFG_<section>_<variable>=<value>`.

Example:
```
-e VPCFG_Server_enabled=true -e VPCFG_Server_enforceMod=false -e VPCFG_Server_dataRate=500
```

turns into
```
[Server]
enabled=true
enforceMod=false
dataRate=500
```
All existing configuration in /config/valheimplus/valheim_plus.cfg is retained and a backup of the old config is created as /config/valheimplus/valheim_plus.old before writing the new config file.


# System requirements
On our system while idle with no players connected Valheim server consumes around 2.8 GB RSS and 10 GB VSZ. All the while using around 30% of one CPU Core on a 2.40GHz Intel Xeon E5-2620 v3. Valheim server is making use of many threads with two of them seemingly doing the bulk of the work each responsible for around 8-10% of the 30% of idle load.

The picture changes when players connect. The first player increased overall load to 42%, the second player to 53%. In the thread view we see that a thread that was previously consuming 10% is now hovering around 38%. Meaning while Valheim server creates 50 threads on our system it looks like there is a single thread doing the bulk of all work (~70%) with no way for the Kernel to distribute the load to many cores.

Therefor our minimum requirements would be a dual core system with 4 GB of RAM and 8 GB of Swap. And our recommended system would be a high clocked 4 core server with 16 GB of RAM. A few very high clocked cores will be more beneficial than having many cores. I.e. two 5 GHz cores will yield better performance than six 2 GHz cores.
This holds especially true the more players are connected to the system.


# Deployment

## Deploying with Docker and systemd
Create an optional config file `/etc/sysconfig/valheim-server`
```
SERVER_NAME="My Server"
SERVER_PORT=2456
WORLD_NAME=Dedicated
SERVER_PASS=secret
SERVER_PUBLIC=true
```

Then enable the Docker container on system boot
```
$ sudo mkdir -p /etc/valheim /opt/valheim
$ sudo curl -o /etc/systemd/system/valheim.service https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/valheim.service
$ sudo systemctl daemon-reload
$ sudo systemctl enable valheim.service
$ sudo systemctl start valheim.service
```

## Deploying to Kubernetes
Kubernetes manifests using this container image, along with a helm chart, are available from the following repository:
[https://github.com/Addyvan/valheim-k8s](https://github.com/Addyvan/valheim-k8s)

The chart is also available directly using:
```bash
helm repo add valheim-k8s https://addyvan.github.io/valheim-k8s/
helm repo update
helm install valheim-server valheim-k8s/valheim-k8s # see repo for full config
```

## Deploying to AWS ECS
CDK Project for spinning up a Valheim game server on AWS Using ECS Fargate and Amazon EFS is available here:
[https://github.com/rileydakota/valheim-ecs-fargate-cdk](https://github.com/rileydakota/valheim-ecs-fargate-cdk)


# Updates
By default the container will check for Valheim server updates every 15 minutes.
If an update is found it is downloaded and the server restarted.
This update schedule can be changed using the `UPDATE_CRON` environment variable.


# Backups
The container will on startup and periodically create a backup of the `worlds/` directory.

The default is once per hour but can be changed using the `BACKUPS_CRON` environment variable.

Default backup directory is `/config/backups/` within the container. A different directory can be set using the `BACKUPS_DIRECTORY` environment variable.
It makes sense to have this directory be a volume mount from the host.
Warning: do not make the backup directory a subfolder of `/config/worlds/`. Otherwise each backup will backup all previous backups.

By default 3 days worth of backups will be kept. A different number can be configured using `BACKUPS_MAX_AGE`. The value is in days.

Beware that backups are performed while the server is running. As such files might be in an open state when the backup runs.
However the `worlds/` directory also contains a `.db.old` file for each world which should always be closed and in a consistent state.

## Post backup hook
After a backup ZIP has been created the command specified by `$POST_BACKUP_HOOK` will be executed if set to a non-zero string.
Within that command the string `@BACKUP_FILE@` will be replaced by the full path to the just created ZIP file.

### Example
```
-v $HOME/.ssh/id_rsa:/root/.ssh/id_rsa \
-v $HOME/.ssh/known_hosts:/root/.ssh/known_hosts \
-e POST_BACKUP_HOOK='timeout 300 scp @BACKUP_FILE@ myself@example.com:~/backups/$(basename @BACKUP_FILE@)'
```

If the post backup hook requires additional packages like e.g. `awscli` the `POST_BOOTSTRAP_HOOK` environment variable could be used to install those.


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


# ValheimPlus
[ValheimPlus](https://github.com/valheimPlus/ValheimPlus) is a popular Valheim mod.
It has been incorporated into this container. To enable V+ provide the env variable `VALHEIM_PLUS=true`.
Upon first start V+ will create a new directory `/config/valheimplus` where its config files are located.
As a user you are mainly concerned with the values in `/config/valheimplus/valheim_plus.cfg`.
For most modifications the mod has to be installed both, on the server as well as all the clients that connect to the server.
A few modifications, like for example changing the `dataRate` can be done server only.

## Updates
ValheimPlus is automatically being updated using the same `UPDATE_CRON` schedule the Valheim server uses to check for updates. If an update of either
Valheim server or ValheimPlus is found it is being downloaded, configured and the server automatically restarted.
This also means your clients always need to run the latest ValheimPlus version or won't be able to connect. If this is undesired the schedule could be changed to only check for updates once per day. Example  `UPDATE_CRON='0 6 * * *'` would only check at 6 AM.

## Configuration
See [ValheimPlus config from Environment Variables](#valheimplus-config-from-environment-variables)

### Server data rate
A popular change is to increase the server send rate.

To do so enable ValheimPlus (`VALHEIM_PLUS=true`) and configure the following section in `/config/valheimplus/valheim_plus.cfg`
```
[Server]
enabled=true
enforceMod=false
dataRate=600
```
(Or whatever `dataRate` value you require. The value is in kb/s with a default of 60.)

Alternatively start with `-e VPCFG_Server_enabled=true -e VPCFG_Server_enforceMod=false -e VPCFG_Server_dataRate=600`.

### Disable server password
Another popular mod for LAN play that does not require the clients to run ValheimPlus is to turn off password authentication.

To do so enable ValheimPlus (`VALHEIM_PLUS=true`), set an empty password (`SERVER_PASS=""`), make the server non-public (`SERVER_PUBLIC=false`) and configure the following section in `/config/valheimplus/valheim_plus.cfg`
```
[Server]
enabled=true
enforceMod=false
disableServerPassword=true
```
Alternatively start with `-e VPCFG_Server_enabled=true -e VPCFG_Server_enforceMod=false -e VPCFG_Server_disableServerPassword=true`.

Ensure that the server can not be accessed from the public Internet. If you like to have the LAN experience but over the Internet I can highly recommend [ZeroTier](https://www.zerotier.com/). It is an open source VPN service where you can create a virtual network switch that you and your friends can join. It is like Hamachi but free and open source. They do have a paid product for Businesses with more than 50 users. So for more than 50 users you could either get their Business product or alternatively would have to host the VPN controller yourself.


# Synology Help
## First install
This is not an extensive tutorial, but I hope these screenshots can be helpful.
Beware that the server can use multiple GB of RAM and produces a lot of CPU load.

![Step 1](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/step1.png "Step 1")
![Step 2](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/step2.png "Step 2")
![Step 3](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/step3.png "Step 3")
![Step 4](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/step4.png "Step 4")
![Step 5](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/step5.png "Step 5")
![Step 6](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/step6.png "Step 6")
![Step 7](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/step7.png "Step 7")
![Step 8](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/step8.png "Step 8")

## Updating the container image to the latest version
The process of updating the image clears all data stored inside the container. So before doing a container image upgrade, make absolutely sure that `/config`, which contains your world, is an external volume stored on your NAS (Step 4 of the [First install](#first-install) process). It is also a good idea to copy the latest version of the world backup to another location, like your PC.
![Update Step 1](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/update1.png "Update Step 1")
![Update Step 2](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/update2.png "Update Step 2")
![Update Step 3](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/update3.png "Update Step 3")
![Update Step 4](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/update4.png "Update Step 4")
![Update Step 5](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/update5.png "Update Step 5")
![Update Step 6](https://raw.githubusercontent.com/lloesche/valheim-server-docker/main/misc/update6.png "Update Step 6")
