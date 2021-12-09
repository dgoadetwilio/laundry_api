FROM centos:7.6.1810
RUN mkdir -p /opt/tweek
COPY bin/* /opt/tweek

WORKDIR /opt/tweek
ENTRYPOINT ["./laundry_api"]

EXPOSE 8090
