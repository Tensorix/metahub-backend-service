# metahub-backend-service

## Setup

```shell
# Install dependencies
go mod download

# Generate protobuf files
buf generate

# Create secret image
dd if=/dev/urandom of=secret.img bs=512 count=1024

# Init the database
sqlite3 mbs.sqlite
sqlite> .read init.sql
```
