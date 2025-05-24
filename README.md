# Fileshare is a lightweight, grpc based centralized file server
Fileshare is designed for lightweight file server. Grpc is used for fast transfer.

# How to use?
Each fileshare needs a `settings.yml` file, which should contains below parts

``` yaml
address: 127.0.0.1:60011
database: bin/client/client.db
```

- for address, make sure that client and server has same ip address that can be accessed
- for database, just make sure the parent directory of xxx.db exists
    - for example, `bin/client/client.db` just need to make sure `bin/client` exists

## Example Structures
below is a example structure of client and server structure
```
.
├── client
│   ├── client.db
│   ├── fileshare
│   ├── llvm-2.2.tar.gz
│   └── settings.yml
└── server
    ├── fileshare
    ├── server.db
    └── settings.yml

3 directories, 8 files

```

``` yaml
# config for client/settings.yml
address: 127.0.0.1:60011
database: client.db
```
``` yaml
# config for server/settings.yml
address: 127.0.0.1:60011
database: server.db
```

## Example cmd usages:
### Client
``` sh
cd client
./fileshare upload llvm-2.2.tar.gz
```

### Server
``` sh
cd server
./fileshare server
```

# Pictures
## Upload
![](docs/pictures/upload.png)
## Download
![](docs/pictures/download.png)

## Final Structure be like:
![](docs/pictures/final-structure.png)