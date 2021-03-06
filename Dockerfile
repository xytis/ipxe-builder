FROM ubuntu:latest

RUN apt update && apt install --yes build-essential liblzma-dev mtools mkisofs syslinux

COPY ipxe/src/ /ipxe/
WORKDIR /ipxe
RUN make

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
COPY ipxe-builder /usr/local/bin/ipxe-builder
ENTRYPOINT ["entrypoint.sh"]
