[![CircleCI](https://circleci.com/gh/friendsofgo/gopherapi/tree/master.svg?style=svg)](https://circleci.com/gh/friendsofgo/gopherapi/tree/master)

# Gopher API
The Gopher API, is a evolutive simple CRUD API for formative purpose, we're building it while writing the posts of the [blog](https://blog.friendsofgo.tech).

In this API we've learnt differents, features and patterns in Go:

* Using Gorilla Mux to create an simple API
* Using a SOLID, Hexagonal Architecture
* Testing HTTP handlers
* Integration with CircleCI
* Using Wire to build dependencies [only in v0.3.1](https://github.com/friendsofgo/gopherapi/releases/tag/v0.3.1)
* Using pattern contextkey
* Using instrumenting with Zipkin

## How can I use it?

**Install**

```sh
$ go get -u github.com/friendsofgo/gopherapi/cmd/gopherapi
```

**Usage**
Launch server with predefined data

```sh
$ gopherapi --withData
The gopher server is on tap now: http://localhost:8080
```

If you want to start the server using zipkin you will need use the next option
```sh
$ gopherapi --withTrace
```

If you want start the server using cockroachdb you will need use the next option

```sh
$gopherapi --cockroach
```

## Endpoints

Fetch all gophers

```
GET /gophers
```

Fetch a gopher by ID

```
GET /gophers/{gopher_id}
```

Add a gopher

```
POST /gophers
```

Modify a gopher
```
PUT /gophers/{gopher_id}
```

Remove a gopher
```
DELETE /gophers/{gopher_id}
```

You can import the Postman collection into `api/GopherApi.postman_collection`

## Launch Zipkin

```
docker run -d -p 9411:9411 openzipkin/zipkin
```

## Contributing
If you think that you can improve with new endpoints, and functionallities the API feel free to contribute with this project with fork this repo and send your Pull Request.

## License
MIT License, see [LICENSE](https://github.com/friendsofgo/gopherapi/blob/master/LICENSE)
