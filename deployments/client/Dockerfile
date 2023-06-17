FROM fyneio/fyne-cross-images:darwin-sdk-extractor as sdk

WORKDIR /mnt

RUN apt update  && apt install -y python3-pip &&  pip install gdown && gdown --id 1pMJbMq6Shcmj7eSOMRXzip-2pnWJKkoB

RUN sh /usr/local/bin/darwin-sdk-extractor.sh



FROM ubuntu:latest

WORKDIR /app

RUN apt update && apt install -y build-essential wget gcc libgl1-mesa-dev xorg-dev curl && \
#    # Docker
    install -m 0755 -d /etc/apt/keyrings && \
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg && \
    chmod a+r /etc/apt/keyrings/docker.gpg && \
    echo "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
    "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
     apt update &&  apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin && \
    # Golang 1.20.5
     wget https://dl.google.com/go/go1.20.5.linux-amd64.tar.gz && \
     tar -C /opt -xzf go1.20.5.linux-amd64.tar.gz && \
     rm go1.20.5.linux-amd64.tar.gz


ENV PATH=$PATH:/opt/go/bin

COPY --from=sdk /mnt/SDKs ./SDKs
#
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod tidy && \
    go get fyne.io/fyne/v2@latest && \
    go install fyne.io/fyne/v2/cmd/fyne@latest &&  \
    go install github.com/fyne-io/fyne-cross@latest

COPY . .
