FROM --platform=linux/amd64 ubuntu:22.04

ARG HTTP_PROXY
ARG HTTPS_PROXY

ENV http_proxy=${HTTP_PROXY} \
   https_proxy=${HTTPS_PROXY}

RUN apt-get update && apt-get install -y --no-install-recommends \
   curl \
   wget \
   git \
   unzip \
   zip \
   cron \
   vim \
   nano \
   redis-tools \ 
   build-essential \
   make \
   locales && \
   apt-get reinstall -y ca-certificates && \
   update-ca-certificates && \
   apt-get clean && \
   curl https://dl.google.com/go/go1.22.5.linux-amd64.tar.gz -o go1.22.5.linux-amd64.tar.gz && \
   rm -rf /usr/local/go && \
   tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz && \
   rm go1.22.5.linux-amd64.tar.gz && \
   sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
   locale-gen && \
   update-locale LANG=en_US.UTF-8 LC_ALL=en_US.UTF-8 && \
   ldconfig && \
   curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

ENV LANG=en_US.UTF-8 \
   LANGUAGE=en_US:en \
   LC_ALL=en_US.UTF-8 \
   TZ="UTC" \
   PATH=$PATH:/usr/local/go/bin \
   TERM=xterm

WORKDIR /app
COPY . .

RUN make deps
RUN make build

ENV DEV_MODE=false
ENV HOST=0.0.0.0
ENV PORT=6379

EXPOSE ${PORT}

CMD if [ "$DEV_MODE" = "true" ]; then \
   make dev HOST=${HOST} PORT=${PORT}; \
   else \
   make start HOST=${HOST} PORT=${PORT}; \
   fi
