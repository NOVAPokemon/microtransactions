FROM golang:latest

ENV executable="executable"

RUN mkdir /service
WORKDIR /service
COPY $executable .
COPY microtransaction_offers.json .

COPY dockerize .
RUN chmod +x dockerize

CMD ["$executable"]