FROM ubuntu:18.04


RUN echo "deb http://us.archive.ubuntu.com/ubuntu/ xenial-security main" | tee -a /etc/apt/sources.list
RUN apt-get update && \
  apt-get install -y build-essential gcc make cmake libtool zlib1g-dev libpng-dev libaec-dev libjpeg-dev libeccodes-dev curl git libopenjp2-7-dev libjasper-dev

ENV GOLANG_VERSION 1.11.2
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 1dfe664fa3d8ad714bbd15a36627992effd150ddabd7523931f077b3926d736d
RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
    && echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz
ENV PATH /root/go/bin:/usr/local/go/bin:$PATH
WORKDIR /root

