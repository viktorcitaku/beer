# beer ![Travis (.org)](https://img.shields.io/travis/viktorcitaku/beer?label=tests)
Simple Golang Application using a public API to pickup beers

### System requirements

```
Docker: 19.03.1
Golang: 1.15.2
Browser: Chrome (latest), Firefox (latest)
```

### How to build and run the application?
#### The docker way

If you have docker installed the following commands have to be executed one after the other in the current root folder (Please tweak the docker commands accordingly).

1. `docker build -t beer .`
2. `docker run -it -p 8080:8080 beer`

#### The standard way

If you have golang already installed, then in the current root folder run the following command `go run main.go`

### Live Demo

[Beer](https://obscure-spire-53165.herokuapp.com/)
