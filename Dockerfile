FROM alpine:latest

COPY ./etc/timezone /etc/timezone

COPY ./etc/localtime /etc/localtime

COPY ./main /bin/kk-ping

RUN chmod +x /bin/kk-ping

COPY ./config /config

COPY ./app.ini /app.ini

ENV KK_ENV_CONFIG /config/env.ini

VOLUME /config

CMD kk-ping $KK_ENV_CONFIG

