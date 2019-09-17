FROM alpine:3.10.2

ARG ARCH="amd64"
ARG OS="linux"
ADD .build .build
COPY .build/${OS}-${ARCH}/logstash-exporter        /bin/logstash-exporter

USER nobody
EXPOSE 9649
ENTRYPOINT ["/bin/logstash-exporter"]