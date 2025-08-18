# Listmonk Go SDK

**Unofficial** Go SDK for the [Listmonk](https://listmonk.app/) newsletter and mailing list manager

## Install

```
go get github.com/canpacis/listmonk-go
```

## Usage

```go
// Import the package
import listmonkgo "github.com/canpacis/listmonk-go"

// Create a client
client := listmonkgo.New(
  listmonkgo.WithBaseURL(/* Base URL */),
  listmonkgo.WithAPIUser(/* API User */),
  listmonkgo.WithToken(/* API Token */),
  // You can optionally pass a custom http client
  listmonkgo.WithHTTPClient(/* Your client */),
)
```

### Example

This example creates a client and retrieves the mailing lists and prints them.

```go
package main

import (
  "context"
  "log"

  listmonkgo "github.com/canpacis/listmonk-go"
)

func main() {
  // Load variables from environment
  client := listmonkgo.New(
		listmonkgo.WithBaseURL(os.Getenv("LISTMONK_URL")),
		listmonkgo.WithAPIUser(os.Getenv("API_USER")),
		listmonkgo.WithToken(os.Getenv("API_TOKEN")),
	)

  // Get lists, optionally pass your parameters
  lists, err := client.GetLists(context.Background(), &listmonkgo.GetListsParams{})
  if err != nil {
    log.Fatal(err)
  }

  for _, list := range lists.Results {
    fmt.Println(list)
  }
}
```

> You can find the API reference in [Listmonk](https://listmonk.app/docs/) website.

## Roadmap

I'm currently developing the SDK and trying to react some parity with the [Listmonk API](https://listmonk.app/docs/). Tests are not enough at the moment as well.

- [x] Subscribers
- [x] Lists
- [x] Import
- [ ] Campaigns
- [ ] Media
- [ ] Templates
- [x] Transactional
- [ ] Bounces

Any help is welcome.