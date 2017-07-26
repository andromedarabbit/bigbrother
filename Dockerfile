FROM alpine

RUN apk add --update ca-certificates && \
    rm -rf /var/cache/apk/* /tmp/*

ENTRYPOINT ["/opt/app/bigbrother"]
RUN mkdir -p /opt/app
WORKDIR /opt/app

COPY bin/bigbrother ./

ENV LOGGING_LEVEL WARNING

COPY docker-entrypoint.sh ./
ENTRYPOINT ["./docker-entrypoint.sh"]

CMD ["./bigbrother"]
