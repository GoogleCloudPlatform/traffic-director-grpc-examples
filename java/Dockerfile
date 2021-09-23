FROM openjdk:11

RUN mkdir -p /build/install/wallet

ADD build/install/wallet /build/install/wallet

RUN ln -s /build/install/wallet/bin/account-server /account-server

RUN ln -s /build/install/wallet/bin/stats-server /stats-server

RUN ln -s /build/install/wallet/bin/wallet-server /wallet-server

RUN ln -s /build/install/wallet/bin/client /wallet-client

CMD ["/bin/sleep","inf"]

EXPOSE 8000
