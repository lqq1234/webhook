FROM alpine:3.5
USER root
ADD imgeventfilter /opt/imgeventfilter
RUN chmod +x /opt/imgeventfilter
WORKDIR /opt
ENTRYPOINT ["/opt/imgeventfilter"]
