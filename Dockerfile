FROM golang:1.19.1-buster as builder

RUN DEBIAN_FRONTEND=noninteractive \
	apt-get update && apt-get install -y build-essential 

ADD . /yq
WORKDIR /yq
#RUN go mod tidy
RUN go build
ADD corpus/*.json /corpus/

FROM golang:1.19.1-buster
COPY --from=builder /yq/yq /
COPY --from=builder /corpus/*.json /testsuite/

ENTRYPOINT ["/yq", "@@"]
