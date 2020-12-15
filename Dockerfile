FROM golang:latest 

LABEL maintainer="ted <cmd.ctrl.q@gmail.com>"

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download 

# copy app files - copy all (.) files recursively in current directory (.)
COPY . .

# copy firebase credentials file into docker container 
# COPY 

# specify env variable port 
ENV PORT=":8000" 

# add env var firebase creds into container 
ENV GOOGLE_APPLICATION_CREDENTIALS="/tmp/keys/firebase-creds.json"

# build app 
RUN go build

# specify command to run server 
CMD ["./go-mux-crash-course"]