FROM debian:stable
COPY supervisord.conf /etc/supervisor/conf.d/valheim.conf
COPY valheim-server /usr/local/bin/
COPY valheim-updater /usr/local/bin/
ADD https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz /tmp
RUN apt-get update \
    && apt-get -y dist-upgrade \
    && DEBIAN_FRONTEND=noninteractive apt-get -y install \
        lib32gcc1 \
        libsdl2-2.0-0 \
        ca-certificates \
        supervisor \
        procps \
        locales \
    && echo 'LANG="en_US.UTF-8"' > /etc/default/locale \
    && echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
    && locale-gen \
    && apt-get clean \
    && adduser \
        --home /home/steam \
        --disabled-password \
        --shell /bin/bash \
        --gecos "Steam User" \
        --quiet \
        steam \
    && mkdir -p /var/log/supervisor /opt/valheim /opt/steamcmd /home/steam/.config/unity3d/IronGate /config \
    && ln -s /config /home/steam/.config/unity3d/IronGate/Valheim \
    && chown -R steam:steam /opt/valheim /home/steam /config \
    && cd /home/steam \
    && tar xzvf /tmp/steamcmd_linux.tar.gz -C /opt/steamcmd/ \
    && chown -R root:root /opt/steamcmd \
    && chmod 755 /opt/steamcmd/steamcmd.sh /opt/steamcmd/linux32/steamcmd /opt/steamcmd/linux32/steamerrorreporter \
    && chmod +x /usr/local/bin/valheim-* \
    && cd "/opt/steamcmd" \
    && ./steamcmd.sh +login anonymous +quit \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

VOLUME /config
EXPOSE 2456-2458/udp
WORKDIR /home/steam
CMD ["/usr/bin/supervisord"]
