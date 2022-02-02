# cosmologger


## Development

To have a development environment just run the following commands:
```bash
git clone git@github.com:archway-network/cosmologger.git
cd cosmologger
docker-compose up -d
```

Then enter the shell of the running dev-container:

```bash
docker exec -it cosmologger sh
```

To build and try:

```bash
clear && go build . && ./cosmologger
```