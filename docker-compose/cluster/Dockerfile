FROM redis/redis-stack-server:7.2.0-M01

RUN mkdir -p /redis

WORKDIR /redis

COPY redis.conf .
COPY startup.sh /usr/local/bin/
COPY startup_cluster.sh /usr/local/bin/

RUN useradd redis
RUN chown redis:redis /redis/* 
RUN chmod +x /usr/local/bin/startup.sh
RUN chmod +x /usr/local/bin/startup_cluster.sh

EXPOSE 6379

ENTRYPOINT [ "/usr/local/bin/startup.sh" ]