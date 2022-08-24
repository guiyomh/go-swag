# go-swag

> Work in progress


Go-swag is a [OpenApi](https://www.openapis.org/) file generator based on [kin-openapi](https://github.com/getkin/kin-openapi)
and inspired by [swag](https://github.com/swaggo/swag)

## Summary
- [Getting started](#getting-started)
- [Comments format](#comments-format)
  - [API Info](#api-info)
  - [API Operation](#api-operation)
  - [Mime Types](#mime-types)

## Getting started

1. Add comments to your API
2. Download go-swag by using:

``` bash
$ go install github.com/guiyomh/go-swag/cmd/go-swag@latest
```

## Comments format

### API Info

### API Operation

| annotation  | description                                                         |
| :---------: | :------------------------------------------------------------------ |
| description | A verbose explanationn of the operation behaviour.                  |
|     id      | A unique string used to identify the operation.                     |
|    tags     | A list of tags separated by commas.                                 |
|   summary   | A short summary of what the operation does.                         |
|   accept    | A listof MIME types the APIS can consume.                           |
|   produce   | A listof MIME types the APIS can produce.                           |
|  response   | Response that separated by space<br />`statusCode dataType comment` |
|   router    | Path definition that separated by space `path httpMethod`           |
| deprecated  | Mark the endpoint as deprecated                                     |

### Mime Types

`go-swag` accepts all MIME Types which are in the correct format, that is, match */*. Besides that, `go-swag` also accepts aliases for some MIME Types as follows:

|         Alias         | MIME Type                         |
| :-------------------: | :-------------------------------- |
|         json          | application/json                  |
|          xml          | text/xml                          |
|         plain         | text/plain                        |
|         html          | text/html                         |
|         mpfd          | multipart/form-data               |
| x-www-form-urlencoded | application/x-www-form-urlencoded |
|       json-api        | application/vnd.api+json          |
|      json-stream      | application/x-json-stream         |
|     octet-stream      | application/octet-stream          |
|          png          | image/png                         |
|         jpeg          | image/jpeg                        |
|          gif          | image/gif                         |