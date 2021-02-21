FROM debian:stable as build-env
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get -y install apt-utils
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y build-essential curl
WORKDIR /build/busybox
RUN curl -L -o /tmp/busybox.tar.bz2 https://busybox.net/downloads/busybox-1.32.1.tar.bz2
RUN tar xjvf /tmp/busybox.tar.bz2 --strip-components=1 -C /build/busybox
RUN make defconfig
RUN make install

FROM debian:stable
COPY --from=build-env /build/busybox/_install/bin/busybox /bin/busybox
COPY valheim-* /usr/local/bin/
ADD https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz /tmp/
RUN dpkg --add-architecture i386 \
    && apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get -y install apt-utils \
    && DEBIAN_FRONTEND=noninteractive apt-get -y dist-upgrade \
    && DEBIAN_FRONTEND=noninteractive apt-get -y install \
        lib32gcc1 \
        libsdl2-2.0-0 \
        libsdl2-2.0-0:i386 \
        libcurl4 \
        libcurl4:i386 \
        libc6-dev \
        ca-certificates \
        supervisor \
        procps \
        locales \
        unzip \
        zip \
        rsync \
        wget \
        curl \
    && echo 'LANG="en_US.UTF-8"' > /etc/default/locale \
    && echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
    && ln -s /bin/busybox /sbin/syslogd \
    && ln -s /bin/busybox /usr/sbin/crond \
    && ln -s /bin/busybox /usr/bin/crontab \
    && ln -s /bin/busybox /usr/bin/vi \
    && ln -s /bin/busybox /usr/bin/wget \
    && locale-gen \
    && apt-get clean \
    && mkdir -p /var/spool/cron/crontabs /var/log/supervisor /opt/valheim /opt/valheimplus /opt/steamcmd /root/.config/unity3d/IronGate /config \
    && ln -s /config /root/.config/unity3d/IronGate/Valheim \
    && tar xzvf /tmp/steamcmd_linux.tar.gz -C /opt/steamcmd/ \
    && chown -R root:root /opt/steamcmd \
    && chmod 755 /opt/steamcmd/steamcmd.sh /opt/steamcmd/linux32/steamcmd /opt/steamcmd/linux32/steamerrorreporter /usr/local/bin/valheim-* \
    && cd "/opt/steamcmd" \
    && ./steamcmd.sh +login anonymous +quit \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
COPY supervisord.conf /etc/supervisor/supervisord.conf

ENV TZ=Etc/UTC
VOLUME ["/config", "/opt/valheim", "/opt/valheimplus"]
EXPOSE 2456-2458/udp
WORKDIR /
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]
