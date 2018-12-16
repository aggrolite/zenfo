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
	&& rm -rf /var/lib/apt/lists/*

ENV GOPATH=/go \
	GOROOT=/usr/local/go \
	PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

RUN \
	curl -k -o /tmp/go.tgz https://dl.google.com/go/go1.11.3.linux-amd64.tar.gz && \
	tar -C /usr/local -xzf /tmp/go.tgz && \
	rm /tmp/go.tgz

RUN \
	mkdir -p $GOPATH && \
	go get -insecure golang.org/x/lint/golint && \
	chmod -R o+rwx $GOPATH
