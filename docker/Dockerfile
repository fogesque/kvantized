# Use Ubuntu 24.04 as the base for native GCC-14 support
FROM ubuntu:24.04

# Prevent interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install system dependencies and build tools
RUN apt-get update && apt-get install -y \
    build-essential \
    gcc-14 \
    g++-14 \
    cmake \
    ninja-build \
    git \
    curl \
    zip \
    unzip \
    tar \
    pkg-config \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Set GCC-14 as the default compiler
RUN update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-14 100 \
    --slave /usr/bin/g++ g++ /usr/bin/g++-14 \
    --slave /usr/bin/gcov gcov /usr/bin/gcov-14

# Setup vcpkg
ENV VCPKG_ROOT=/opt/vcpkg
ENV PATH="${VCPKG_ROOT}:${PATH}"

RUN git clone https://github.com/microsoft/vcpkg.git ${VCPKG_ROOT} \
    && ${VCPKG_ROOT}/bootstrap-vcpkg.sh -disableMetrics

# Set working directory
WORKDIR /workspace

# Default to bash
CMD ["/bin/bash"]