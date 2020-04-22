FROM golang:latest

ENV executable="executable"

RUN mkdir /service
WORKDIR /service
COPY $executable .
COPY microtransaction_offers.json .

ENTRYPOINT ./$executable