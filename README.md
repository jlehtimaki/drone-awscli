# drone-awscli
Drone plugin to execute AWS Cli commands.

## Build
Build the binary with the following commands:

```export GO111MODULE=on
go mod download
go test
go build
```

## Docker

Build the docker image with:
```
docker build --rm=true -t lehtux/drone-awscli .
```

## Usage
```
docker run --rm lehtux/drone-awscli -e AWS_ACCESS_KEY=.... -e -e AWS_SECRET_KEY=.... -e PLUGIN_COMMAND="sts get-caller-identity"
```

## Parameters
| Paramenter            | Description                   |Required|
| -------------         |:-------------:                |:-------------:|
| AWS_ACCESS_KEY        | AWS Access key                | YES
| AWS_SECRET_KEY        | AWS Access key secret         | YES
| PLUGIN_ASSUME_ROLE    | AWS Assume role               | NO
| PLUGIN_COMMAND        | AWS Client command to be run  | YES