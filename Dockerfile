FROM leochan007/ubuntu1604_base

LABEL MAINTAINER leo chan <leochan007@163.com>

ENV DEBIAN_FRONTEND noninteractive

COPY updater /root

RUN chmod a+x /root/updater

RUN  echo 20190311 > /root/version

WORKDIR /root

CMD /root/updater
