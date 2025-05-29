# Fileshare is a lightweight, grpc based centralized file server
Fileshare is designed for lightweight file server. Grpc is used for fast transfer.

Fileshare auto check the validity of the file transferred.

On both server and client side, both download and upload, **fileshare** will check the `sha256sum` value automatically

Fileshare also provides web api for monitoring sqlite data, see examples below

Fileshare monitors upload, linkgen, download actions at server side,

# How to use?
Each fileshare needs a `settings.yml` file in the `same folder with fileshare`, which should contains below parts

``` yaml
grpc_address: 0.0.0.0:60011
web_address: 0.0.0.0:8080
database: server.db
share_code_length: 8
cache_directory: .cache
download_directory: .download
certs_path: certs
valid_days: 30
blocked_ips:
  - 127.0.0.1
```

## Configuration files explained

- for `grpc address` and `web address`, make sure that client and server has same ip address that can be accessed
- for database, just make sure the parent directory of xxx.db exists
    - for example, `client/client.db` just need to make sure `client` exists
- for share_code_length, make sure this is `not set` to the default length of sha256 (which is 64 by default)
- for cache_directory, cache file chunks. if not set, then use `$HOME/.fileshare`
- for download_directory, where download file is stored. if not set, then use `$HOME/Downloads`
- for valid_days: set the default valid days for a share link, if not set, then default is `7`, lives for a week
- for blocked_ips, all requests from this ip addr will be blocked

### Examples
#### Server
Notice that share_code_length `cannot` be set to the length of the size of sha256, which is `64`
``` yaml
# config for server/settings.yml
grpc_address: 0.0.0.0:60011
web_address: 0.0.0.0:8080
database: server.db

# this config determines the length of sharelink code
# any length is ok, except the fixed size of sha256, which is 64
share_code_length: 8
cache_directory: .cache
download_directory: .download

# below configurations are server side only
certs_path: certs
valid_days: 30
blocked_ips:
  - 127.0.0.1
```

#### Client
``` yaml
# config for client/settings.yml
grpc_address: 0.0.0.0:60011
web_address: 0.0.0.0:8080
database: client.db
share_code_length: 8
cache_directory: .cache
download_directory: .download
```

## LinkCode generating:
### Exciting ability introduced! If u wanna share a file with your friends, you can generate linkcode by doing this:

#### Below command will generate a linkcode like

``` sh
fileshare linkgen llvm-2.2.tar.gz 788d871aec139e0c61d49533d0252b21c4cd030e91405491ee8cb9b2d0311072
```

``` sh
INFO[0000] Generated Code is: [fzHghSyr]
```

## Example Structures
below is a example structure of client and server structure
```
.
├── client
│   ├── client.db
│   ├── fileshare
│   ├── kafka_2.13-4.0.0.tgz
│   ├── llvm-2.2.tar.gz
│   └── settings.yml
└── server
    ├── fileshare
    ├── server.db
    └── settings.yml

3 directories, 8 files
```

## Example Usages
### Upload
![](docs/pictures/upload.png)

### Download
![](docs/pictures/download.png)

### LinkGen
![](docs/pictures/linkgen.png)

### Final Structure be like:
![](docs/pictures/final-structure.png)

### Cache Clean
![](docs/pictures/cache-clean.png)

### Cmd usages:

#### Server
``` sh
fileshare server
```

#### Client Upload
``` sh
fileshare upload llvm-2.2.tar.gz
```

#### Client Download
- Use the linkcode shared by your friends, and download with this code is ok!
``` sh
fileshare download fzHghSyr
```

- Optional Usages: Notice that `following hash` is the `checksum` of the file using **sha256sum**
``` sh
fileshare download 788d871aec139e0c61d49533d0252b21c4cd030e91405491ee8cb9b2d0311072
```

#### Client Gen Code
Notice that the parameters are `filename` `checksum256` `valid days`

Below cmd will make a sharelink which will live for 300 days
``` sh
fileshare linkgen llvm-2.2.tar.gz 788d871aec139e0c61d49533d0252b21c4cd030e91405491ee8cb9b2d0311072 300
```

This cmd do not specify valid days, then server will generate `according to settings.yml at server side`
``` sh
fileshare linkgen llvm-2.2.tar.gz 788d871aec139e0c61d49533d0252b21c4cd030e91405491ee8cb9b2d0311072
```

#### Client / Server clean cache
Clean cache command can be used both at server and client
``` sh
fileshare cache clean
```

### Web Apis:
``` sh
curl 0.0.0.0:8080/fileinfo
```

``` sh
curl 0.0.0.0:8080/sharelink
```

``` sh
curl 0.0.0.0:8080/record
```
responses are structure below:
``` json
{
  "data": [
    ...
  ]
}
```

#### Example Output
``` json
{
    "data": [
        {
            "Filename": "fileshare",
            "Sha256": "e21645144861413cfd06a268fb3ff6d6a65da0f002034c1667d4607f664faee3",
            "ChunkSize": 1048576,
            "ChunkNumber": 19,
            "FileSize": 19709952,
            "UploadedChunks": "[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18]",
            "Link": [
                {
                    "Sha256": "e21645144861413cfd06a268fb3ff6d6a65da0f002034c1667d4607f664faee3",
                    "LinkCode": "DoBsLlwu",
                    "CreatedBy": "172.16.33.118:9837",
                    "CreatedAt": "2025-05-29T13:25:58.7261564+08:00",
                    "OutdatedAt": "2026-03-25T13:25:58.7261564+08:00"
                }
            ],
            "Record": [
                {
                    "Sha256": "e21645144861413cfd06a268fb3ff6d6a65da0f002034c1667d4607f664faee3",
                    "InteractAction": "upload",
                    "ClientIp": "172.16.33.118:9836",
                    "Os": "linux,amd64,ethan-archlinux",
                    "Time": "2025-05-29T13:25:55.6178337+08:00"
                },
                {
                    "Sha256": "e21645144861413cfd06a268fb3ff6d6a65da0f002034c1667d4607f664faee3",
                    "InteractAction": "linkgen",
                    "ClientIp": "172.16.33.118:9837",
                    "Os": "linux,amd64,ethan-archlinux",
                    "Time": "2025-05-29T13:25:58.7405099+08:00"
                },
                {
                    "Sha256": "e21645144861413cfd06a268fb3ff6d6a65da0f002034c1667d4607f664faee3",
                    "InteractAction": "download",
                    "ClientIp": "172.16.33.118:9838",
                    "Os": "linux,amd64,ethan-archlinux",
                    "Time": "2025-05-29T13:26:06.8239799+08:00"
                }
            ]
        }
    ]
}
```

## Using Docker?
First download `fileshare.docker.zip` from releases and import this zip file to your docker

And download binary from `fileshare.tar.gz`, extract to fileshare

Then run following commands:
``` sh
docker run -d --name fileshare -p 60011:60011 -p 8080:8080 fileshare:0.1.4
```
