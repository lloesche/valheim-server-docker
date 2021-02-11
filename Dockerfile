FROM debian:stable
COPY valheim-server /usr/local/bin/
COPY valheim-updater /usr/local/bin/
COPY valheim-backup /usr/local/bin/
ADD https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz /tmp
RUN dpkg --add-architecture i386 \
    && apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get -y install apt-utils \
    && DEBIAN_FRONTEND=noninteractive apt-get -y dist-upgrade \
    && DEBIAN_FRONTEND=noninteractive apt-get -y install \
        lib32gcc1 \
        libsdl2-2.0-0 \
        libsdl2-2.0-0:i386 \
        ca-certificates \
        supervisor \
        procps \
        locales \
        unzip \
        zip \
        rsync \
    && echo 'LANG="en_US.UTF-8"' > /etc/default/locale \
    && echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
    && locale-gen \
    && apt-get clean \
    && mkdir -p /var/log/supervisor /opt/valheim /opt/valheim_dl /opt/steamcmd /root/.config/unity3d/IronGate /config \
    && ln -s /config /root/.config/unity3d/IronGate/Valheim \
    && tar xzvf /tmp/steamcmd_linux.tar.gz -C /opt/steamcmd/ \
    && chown -R root:root /opt/steamcmd \
    && chmod 755 /opt/steamcmd/steamcmd.sh /opt/steamcmd/linux32/steamcmd /opt/steamcmd/linux32/steamerrorreporter \
    && chmod +x /usr/local/bin/valheim-* \
    && cd "/opt/steamcmd" \
    && ./steamcmd.sh +login anonymous +quit \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
COPY supervisord.conf /etc/supervisor/supervisord.conf

VOLUME ["/config", "/opt/valheim_dl"]
EXPOSE 2456-2458/udp
WORKDIR /
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]
