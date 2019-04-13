[![CircleCI](https://circleci.com/gh/friendsofgo/gopherapi/tree/master.svg?style=svg)](https://circleci.com/gh/friendsofgo/gopherapi/tree/master)

# Gopher API
The Gopher API, is a simple CRUD API for formative purpose, we're building it while writing the posts of the [blog](https://blog.friendsofgo.tech).

## How can I use it?

**Install**

```sh
$ GO111MODULE=off go get -u github.com/friendsofgo/gopherapi/cmd/gopherapi
```

**Usage**
Launch server with predefined data

```sh
$ gopherapi --withData
The gopher server is on tap now: http://localhost:8080
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

## Contributing
If you think that you can improve with new endpoints, and functionallities the API feel free to contribute with this project with fork this repo and send your Pull Request.

## License
MIT License, see [LICENSE](https://github.com/friendsofgo/gopherapi/blob/master/LICENSE)
