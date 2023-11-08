FROM alpine:latest

RUN echo "Asia/shanghai" > /etc/timezone && \
    ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ADD bin/* /opt/zzauth/
ADD docs/* /opt/zzauth/docs/
ADD static/* /opt/zzauth/static/
WORKDIR /opt/zzauth
EXPOSE 9900
CMD [ "/opt/zzauth/authgate" ]
