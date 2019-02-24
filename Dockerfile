FROM postgres:latest

RUN \
	echo "[http]\n\tsslverify = false" >> /root/.gitconfig

RUN \
	apt-get update && \
	apt-get install -y --no-install-recommends \
	curl \
	make \
	git \
	gcc \
	ca-certificates \
	build-essential \
	dnsutils \
	less \
	&& rm -rf /var/lib/apt/lists/*

RUN \
	curl -k -o /tmp/go.tgz https://dl.google.com/go/go1.11.3.linux-amd64.tar.gz && \
	tar -C /usr/local -xzf /tmp/go.tgz && \
	rm /tmp/go.tgz

#USER postgres

#RUN /etc/init.d/postgresql start

RUN \
	mkdir -p /usr/lib/postgresql/data/ && \
	#cp /usr/share/postgresql/postgresql.conf.sample /usr/lib/postgresql/data/postgresql.conf && \
	chown postgres: /usr/lib/postgresql/data && \
	chmod 755 /usr/lib/postgresql/data/

#RUN mkdir /home/postgres
#RUN chown postgres: /home/postgres
RUN usermod -m -d /home/postgres postgres

USER postgres

ENV GOPATH=/home/postgres/go GOROOT=/usr/local/go PATH=/usr/local/go/bin:/home/postgres/go/bin:/usr/lib/postgresql/11/bin:$PATH

RUN echo $GOPATH

RUN \
	#mkdir -p $GOPATH && \
	go get golang.org/x/lint/golint && \
	go get github.com/mjibson/esc && \
	chmod -R o+rwx $GOPATH

RUN ls -ltr /usr/lib/postgresql
