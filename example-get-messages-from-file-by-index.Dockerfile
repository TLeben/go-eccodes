FROM zachaller/go-eccodes:latest

RUN go get github.com/zachaller/go-eccodes/cmd/example-get-messages-from-file-by-index

ADD ./test-data/ARPEGE_0.1_SP1_00H12H_201709290000.grib2 data.grib2

CMD ["example-get-messages-from-file-by-index", "-file", "./data.grib2"]
