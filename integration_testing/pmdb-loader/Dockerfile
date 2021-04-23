FROM arti.dev.cray.com/third-party-docker-stable-local/postgres:11-alpine
COPY pmdb_csvs.tar.gz /
COPY load.sh /
RUN mkdir /csv 
RUN tar -xvzf pmdb_csvs.tar.gz -C /csv
RUN rm /pmdb_csvs.tar.gz
ENV PMDB_HOSTNAME postgresql
ENV PMDB_PORT 5432
ENTRYPOINT [ "/load.sh" ] 
