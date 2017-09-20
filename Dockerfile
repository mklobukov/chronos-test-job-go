FROM alpine

COPY src/testjob/cmd/testjob/testjob /usr/local/bin/

CMD ["usr/local/bin/testjob", "config.json"]
