FROM nova-server-base:latest

ENV executable="executable"
COPY $executable .
COPY microtransaction_offers.json .

CMD ["sh", "-c", "./$executable"]