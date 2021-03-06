FROM debian:stable as build-env
ARG TESTS
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get -y install apt-utils
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y build-essential curl python3 python3-pip shellcheck
WORKDIR /build/busybox
RUN curl -L -o /tmp/busybox.tar.bz2 https://busybox.net/downloads/busybox-1.32.1.tar.bz2
RUN tar xjvf /tmp/busybox.tar.bz2 --strip-components=1 -C /build/busybox
RUN make defconfig
RUN make install
COPY ./vpenvconf/ /build/vpenvconf/
WORKDIR /build/vpenvconf
RUN if [ "${TESTS:-true}" = true ]; then pip3 install tox && tox; fi
RUN python3 setup.py bdist --format=gztar
COPY bootstrap /usr/local/sbin/
COPY valheim-* /usr/local/bin/
COPY defaults /usr/local/etc/valheim/
COPY common /usr/local/etc/valheim/
COPY contrib/* /usr/local/share/valheim/contrib/
RUN if [ "${TESTS:-true}" = true ]; then shellcheck -a -x -s bash -e SC2034 /usr/local/sbin/bootstrap /usr/local/bin/valheim-* /usr/local/share/valheim/contrib/*.sh; fi

FROM debian:stable
COPY --from=build-env /build/busybox/_install/bin/busybox /usr/local/bin/busybox
COPY --from=build-env /build/vpenvconf/dist/vpenvconf-*.linux-x86_64.tar.gz /tmp/vpenvconf.tar.gz
COPY bootstrap /usr/local/sbin/
COPY valheim-* /usr/local/bin/
COPY defaults /usr/local/etc/valheim/
COPY common /usr/local/etc/valheim/
COPY contrib/* /usr/local/share/valheim/contrib/
ADD https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz /tmp/
RUN dpkg --add-architecture i386 \
    && apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get -y install apt-utils \
    && DEBIAN_FRONTEND=noninteractive apt-get -y dist-upgrade \
    && DEBIAN_FRONTEND=noninteractive apt-get -y install \
        libc6-dev \
        lib32gcc1 \
        libsdl2-2.0-0 \
        libsdl2-2.0-0:i386 \
        curl \
        libcurl4 \
        libcurl4:i386 \
        tcpdump \
        ca-certificates \
        supervisor \
        procps \
        locales \
        unzip \
        zip \
        rsync \
        openssh-client \
        jq \
        python3-minimal \
        python3-pkg-resources \
    && echo 'LANG="en_US.UTF-8"' > /etc/default/locale \
    && echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
    && ln -s /usr/local/bin/busybox /usr/local/sbin/syslogd \
    && ln -s /usr/local/bin/busybox /usr/local/sbin/crond \
    && ln -s /usr/local/bin/busybox /usr/local/bin/crontab \
    && ln -s /usr/local/bin/busybox /usr/local/bin/vi \
    && ln -s /usr/local/bin/busybox /usr/local/bin/wget \
    && ln -s /usr/local/bin/busybox /usr/local/bin/less \
    && rm -f /bin/sh \
    && ln -s /bin/bash /bin/sh \
    && cd / \
    && tar xzvf /tmp/vpenvconf.tar.gz \
    && locale-gen \
    && apt-get clean \
    && mkdir -p /var/spool/cron/crontabs /var/log/supervisor /opt/valheim /opt/steamcmd /root/.config/unity3d/IronGate /config \
    && ln -s /config /root/.config/unity3d/IronGate/Valheim \
    && tar xzvf /tmp/steamcmd_linux.tar.gz -C /opt/steamcmd/ \
    && chown -R root:root /opt/steamcmd \
    && chmod 755 /opt/steamcmd/steamcmd.sh \
        /opt/steamcmd/linux32/steamcmd \
        /opt/steamcmd/linux32/steamerrorreporter \
        /usr/local/sbin/bootstrap \
        /usr/local/bin/valheim-* \
    && cd "/opt/steamcmd" \
    && ./steamcmd.sh +login anonymous +quit \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
COPY supervisord.conf /etc/supervisor/supervisord.conf

ENV TZ=Etc/UTC
VOLUME ["/config", "/opt/valheim"]
EXPOSE 2456-2458/udp
EXPOSE 9001/tcp
WORKDIR /
CMD ["/usr/local/sbin/bootstrap"]
