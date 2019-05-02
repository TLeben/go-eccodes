FROM ubuntu:18.04


RUN echo "deb http://us.archive.ubuntu.com/ubuntu/ xenial-security main" | tee -a /etc/apt/sources.list
RUN apt-get update && \
  apt-get install -y curl build-essential gcc make cmake libtool zlib1g-dev libpng-dev libaec-dev libjpeg-dev curl git libopenjp2-7-dev libjasper-dev python

ENV GOLANG_VERSION 1.12.4
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 d7d1f1f88ddfe55840712dc1747f37a790cbcaa448f6c9cf51bbe10aa65442f5
RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
    && echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz
ENV PATH /root/go/bin:/usr/local/go/bin:$PATH

WORKDIR /root
RUN curl -L "https://confluence.ecmwf.int/download/attachments/45757960/eccodes-2.10.0-Source.tar.gz?api=v2" --output eccodes.tar.gz
RUN tar -xzvf eccodes.tar.gz
RUN mkdir /root/build
WORKDIR /root/build
RUN cmake -DENABLE_GRIB_TIMER=on -DENABLE_ALIGN_MEMORY=on -DENABLE_INSTALL_ECCODES_DEFINITIONS=on -DENABLE_PNG=on -DENABLE_JPG=on -DENABLE_AEC=on  -DENABLE_ECCODES_THREADS=on -DENABLE_MEMFS=on -DBUILD_SHARED_LIBS=both -DENABLE_FORTRAN=off -DCMAKE_INSTALL_PREFIX=/usr/local ../eccodes-2.10.0-Source/
RUN make -j4
RUN ctest -j8
RUN make install
RUN ln -s /usr/local/lib/libeccodes.so /usr/lib/libeccodes.so
RUN ln -s /usr/local/lib/libeccodes_memfs.so /usr/lib/libeccodes_memfs.so
RUN ln -s /usr/local/lib/libeccodes.a /usr/lib/libeccodes.a
RUN ln -s /usr/local/lib/libeccodes_memfs.a /usr/lib/libeccodes_memfs.a
RUN rm -rf /var/lib/apt/lists/* /root/eccodes.tar.gz
WORKDIR /home/dusr/code
