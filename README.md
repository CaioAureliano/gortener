# gortener
a URL Shortener api with Go

## About
Rest api allow the user to shorten a url and view the stats.

#### Techs
 - Golang `1.17.8`
 - MongoDB  `5.0.8`
 - Redis `6.2.7`
 
## Setup
Project need a docker(and docker-compose) to initialize

``` bash
docker -v
docker-compose -v
```

but, before run docker, the steps are necessary:

### Variables

Rename/Copy `.env.example` to `.env` 

| Variable | Description |
| -------- | ----------- |
| `PORT` | tcp port to expose from API |
| `MONGO_DB_NAME` | database name from application |

### Run
_and build_
``` bash
docker-compose up -d --build
```

> to see and follow logs
``` bash
docker-compose logs -f api
```
