# go-pagerduty [![Build Status](https://travis-ci.org/hudl/go-pagerduty.svg?branch=master)](https://travis-ci.org/hudl/go-pagerduty)

`go-pagerduty` is a Go client library for interfacing with the
[PagerDuty API](http://developer.pagerduty.com/), modeled after google's
awesome [google/go-github](http://github.com/google/go-github) library.

## Usage

Import the `pagerduty` package to get started.

```go
import "github.com/hudl/go-pagerduty/pagerduty"
```

Then construct a new PagerDuty client and set the PagerDuty subdomain and API
key.

```go
client := pagerduty.NewClient(nil, "subdomain", "super-secret-api-key")
```

You can use the various services registered with the client to access differnt
parts of the PagerDuty API. For exmaple, you can use the `Incidents` service to
interact with the [Incidents API](https://developer.pagerduty.com/documentation/rest/incidents):

```go
incident, resp, err := client.Incidents.List(nil)
```

Check out more detailed examples in the [`examples`](./examples) directory.

### Helpers

The `Bool()`, `Int()` and `String()` helper functions in
[`pagerduty/pagerduty.go`](./pagerduty/pagerduty.go) are used to wrap values in
pointer variants for easier translation to JSON.

For example, to make an `Team` struct:
```go
Team := &pagerduty.Email{
    ID:          String("id"),
    Name:        String("name"),
    Description: String("description"),
}
```

Pointers to values are used in many of the public types to show intent.
Meaning that if you pass a `nil` value to a struct field, it will be omitted.
Without knowing the intent of the creator, it would be impossible to
differentiate the zero-values for some of the primitive types, such as `int`
and `bool` from an intended zero-value. They would end up always be encoded to
JSON and sent to the PagerDuty API, possibly triggering API errors.

## Roadmap

This library is currently under development and has a limited subset of the
PagerDuty API implemented, specifically just the Email API. We plan to
eventually implement the entire PagerDuty API. Pull requests are welcome!

## License

This library is distributed under the MIT license found in the
[LICENSE](./LICENSE) file.
