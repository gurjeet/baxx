#!/bin/sh

tar -cf - /etc/letsencrypt | \
    encrypt -k /root/.pw | \
    curl --data-binary @- https://baxx.dev/io/$BAXX_TOKEN/letsencrypt.tar
