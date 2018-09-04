FROM debian
RUN dpkg --add-architecture mipsel
RUN apt-get update
RUN apt-get install -y g++-mipsel-linux-gnu
RUN apt-get install -y libboost-system-dev:mipsel
RUN apt-get install -y libssl-dev:mipsel
RUN rm -rf /var/lib/apt/lists/*
