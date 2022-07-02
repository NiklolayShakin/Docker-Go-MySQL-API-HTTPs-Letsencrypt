# start with base image
FROM mysql

RUN apt-get update && apt-get install -y locales && rm -rf /var/lib/apt/lists/* \
	&& localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8
# import data into container
# All scripts in docker-entrypoint-initdb.d/ are automatically executed during container startup
COPY ./database/*.sql /docker-entrypoint-initdb.d/