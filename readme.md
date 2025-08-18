# Listmonk Go SDK

**Unofficial** Go SDK for the [Listmonk](https://listmonk.app/) newsletter and mailing list manager.

## Install

```
go get github.com/canpacis/listmonk-go
```

## Usage

```go
// Import the package
import listmonkgo "github.com/canpacis/listmonk-go"

// ...

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

> You can find the API reference in [Listmonk](https://listmonk.app/docs/apis/apis/) website.

## Roadmap

I'm currently developing the SDK and trying to react some parity with the [Listmonk API](https://listmonk.app/docs/). Tests are not enough at the moment as well.

- [x] Subscribers
  - [x] GET /api/subscribers
  - [x] GET /api/subscribers/{subscriber_id}
  - [x] GET /api/subscribers/{subscriber_id}/export
  - [x] GET /api/subscribers/{subscriber_id}/bounces
  - [x] POST /api/subscribers
  - [x] POST /api/subscribers/{subscriber_id}/optin
  - [x] POST /api/public/subscription
  - [ ] PUT /api/subscribers/lists
  - [ ] PUT /api/subscribers/{subscriber_id}
  - [x] PUT /api/subscribers/{subscriber_id}/blocklist
  - [x] PUT /api/subscribers/blocklist
  - [x] PUT /api/subscribers/query/blocklist
  - [x] DELETE /api/subscribers/{subscriber_id}
  - [x] DELETE /api/subscribers/{subscriber_id}/bounces
  - [x] DELETE /api/subscribers
  - [x] POST /api/subscribers/query/delete
- [x] Lists
  - [x] GET	/api/lists
  - [x] GET	/api/public/lists
  - [x] GET	/api/lists/{list_id}
  - [x] POST /api/lists
  - [x] PUT	/api/lists/{list_id}
  - [x] DELETE /api/lists/{list_id}
- [x] Import
  - [x] GET	/api/import/subscribers
  - [x] GET	/api/import/subscribers/logs
  - [x] POST /api/import/subscribers
  - [x] DELETE /api/import/subscribers
- [x] Campaigns
  - [x] GET	/api/campaigns
  - [x] GET	/api/campaigns/{campaign_id}
  - [x] GET	/api/campaigns/{campaign_id}/preview
  - [x] GET	/api/campaigns/running/stats
  - [x] GET	/api/campaigns/analytics/{type}
  - [x] POST /api/campaigns
  - [x] POST /api/campaigns/{campaign_id}/test
  - [x] PUT	/api/campaigns/{campaign_id}
  - [x] PUT	/api/campaigns/{campaign_id}/status
  - [x] PUT	/api/campaigns/{campaign_id}/archive
  - [x] DELETE /api/campaigns/{campaign_id}
- [x] Media
  - [x] GET	/api/media
  - [x] GET	/api/media/{media_id}
  - [x] POST /api/media
  - [x] DELETE /api/media/{media_id}
- [ ] Templates
  - [ ] GET	/api/templates
  - [ ] GET	/api/templates/{template_id}
  - [ ] GET	/api/templates/{template_id}/preview
  - [ ] POST /api/templates
  - [ ] POST /api/templates/preview
  - [ ] PUT	/api/templates/{template_id}
  - [ ] PUT	/api/templates/{template_id}/default
  - [ ] DELETE /api/templates/{template_id}
- [x] Transactional
  - [x] POST /api/tx
- [x] Bounces
  - [x] GET	/api/bounces
  - [x] DELETE /api/bounces
  - [x] DELETE /api/bounces/{bounce_id}	

Any help is welcome.
