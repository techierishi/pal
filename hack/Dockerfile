FROM ghcr.io/goreleaser/goreleaser-cross:v1.21.6

RUN \
    set -x; \
    echo "Adding mssing dependencies" \
 && apt-get update \
 && apt-get install --no-install-recommends -y -q \
        xvfb libx11-dev libegl1-mesa-dev libgles2-mesa-dev \
 && apt -y autoremove \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* \
    /tmp/* \
    /var/tmp/* \
 && rm -rf /usr/share/man/* \
    /usr/share/doc \