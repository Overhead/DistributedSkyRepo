0;136;0c#FROM centos:centos6
FROM ubuntu:12.04

#WORKDIR /var/lib/docker/work

ADD hostname /usr/local/bin/
ADD main     /usr/local/bin/
#ADD sudo     /usr/bin
#ADD ufw      /usr/sbin

#RUN ["sudo", "ufw", "disable"]
#RUN ["ufw", "disable"]
RUN sudo ls
RUN sudo ufw disable

EXPOSE 1075/udp
EXPOSE 8080/tcp
EXPOSE 8001/udp
EXPOSE 8002/udp
EXPOSE 8003/udp
EXPOSE 8004/udp
EXPOSE 8005/udp
EXPOSE 8006/udp
EXPOSE 8007/udp
EXPOSE 8008/udp
EXPOSE 8009/udp
EXPOSE 8010/udp
EXPOSE 8011/udp
EXPOSE 8012/udp
EXPOSE 8013/udp
EXPOSE 8014/udp
EXPOSE 8015/udp
EXPOSE 8016/udp
EXPOSE 8017/udp
EXPOSE 8018/udp
EXPOSE 8019/udp
EXPOSE 8020/udp
EXPOSE 8021/udp
EXPOSE 8022/udp
EXPOSE 8023/udp
EXPOSE 8024/udp
EXPOSE 12500/udp

CMD ["main"]
