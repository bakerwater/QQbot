FROM alpine
WORKDIR /apps
COPY QQbot .
EXPOSE 5071
ENTRYPOINT ["./QQbot"]