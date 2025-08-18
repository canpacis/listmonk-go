package listmonkgo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type GetSubscribersParams struct {
	// Subscriber search by SQL expression.
	Query string `url:"query"`
	// ID of lists to filter by. Repeat in the query for multiple values.
	ListID []int `url:"list_id"`
	// Subscription status to filter by if there are one or more list_ids
	SubscriptionStatus string `url:"subscription_status"`
	// Result sorting field. Options: name, status, created_at, updated_at.
	OrderBy string `url:"order_by"`
	// Sorting order: ASC for ascending, DESC for descending.
	Order string `url:"order"`
	// Page number for paginated results.
	Page int `url:"page"`
	// Results per page. Set as 'all' for all results.
	PerPage int `url:"per_page"`
}

type Subscription struct {
	ID          int       `json:"id"`
	UUID        uuid.UUID `json:"uuid"`
	Type        ListType  `json:"type"`
	Optin       ListOptin `json:"optin"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	SubscriptionStatus    string         `json:"subscription_status"`
	SubscriptionCreatedAt time.Time      `json:"subscription_created_at"`
	SubscriptionUpdatedAt time.Time      `json:"subscription_updated_at"`
	SubscriptionMeta      map[string]any `json:"subscription_meta"`
}

type SubscriberStatus string

const (
	EnabledSubscriberStatus     SubscriberStatus = "enabled"
	BlocklistedSubscriberStatus SubscriberStatus = "blocklisted"
)

type Subscriber struct {
	ID         int              `json:"id"`
	UUID       uuid.UUID        `json:"uuid"`
	Email      string           `json:"email"`
	Name       string           `json:"name"`
	Status     SubscriberStatus `json:"status"`
	Attributes map[string]any   `json:"attribs"`
	Lists      []Subscription   `json:"lists"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

type GetSubscribersResponse struct {
	Results []Subscriber `json:"results"`
	Search  string       `json:"search"`
	Query   string       `json:"query"`
	Total   int          `json:"total"`
	PerPage int          `json:"per_page"`
	Page    int          `json:"page"`
}

// Query and retrieve subscribers.
func (c *Client) GetSubscribers(ctx context.Context, params *GetSubscribersParams) (*GetSubscribersResponse, error) {
	path := "/api/subscribers"
	resp, err := request[Response[*GetSubscribersResponse]](c, ctx, "GET", path, params)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Retrieve a specific subscriber.
func (c *Client) GetSubscriber(ctx context.Context, id int) (*Subscriber, error) {
	path := fmt.Sprintf("/api/subscribers/%d", id)
	resp, err := request[Response[*Subscriber]](c, ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

type ExportProfile struct {
	ID         int              `json:"id"`
	UUID       uuid.UUID        `json:"uuid"`
	Email      string           `json:"email"`
	Name       string           `json:"name"`
	Status     SubscriberStatus `json:"status"`
	Attributes map[string]any   `json:"attribs"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

type ExportSubscription struct {
	Name               string    `json:"name"`
	Type               ListType  `json:"type"`
	SubscriptionStatus string    `json:"subscription_status"`
	CreatedAt          time.Time `json:"created_at"`
}

type ExportSubscriberResponse struct {
	Profile       []ExportProfile      `json:"profile"`
	Subscriptions []ExportSubscription `json:"subscriptions"`
	CampaignViews []any                `json:"campaign_views"`
	LinkClicks    []any                `json:"link_clicks"`
}

// Export a specific subscriber data that gives profile, list subscriptions, campaign views
// and link clicks information. Names of private lists are replaced with "Private list".
func (c *Client) ExportSubscriber(ctx context.Context, id int) (*ExportSubscriberResponse, error) {
	path := fmt.Sprintf("/api/subscribers/%d/export", id)
	resp, err := request[*ExportSubscriberResponse](c, ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Retrieve a subscriber bounce records.
func (c *Client) GetSubscriberBounces(ctx context.Context, id int) ([]any, error) {
	path := fmt.Sprintf("/api/subscribers/%d/bounces", id)
	resp, err := request[Response[[]any]](c, ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

type CreateSubscriberParams struct {
	// Subscriber's email address.
	Email string `json:"email"`
	// Subscriber's name.
	Name string
	// Subscriber's status: enabled, blocklisted.
	Status SubscriberStatus
	// List of list IDs to subscribe to.
	Lists []int
	// Attributes of the new subscriber.
	Attributes map[string]any
	// If true, subscriptions are marked as confirmed and no-optin emails are sent for double opt-in lists.
	PreconfirmSubscriptions bool `json:"preconfirm_subscriptions"`
}

type CreateSubscriberResponse struct {
	ID         int              `json:"id"`
	UUID       uuid.UUID        `json:"uuid"`
	Email      string           `json:"email"`
	Name       string           `json:"name"`
	Status     SubscriberStatus `json:"status"`
	Attributes map[string]any   `json:"attribs"`
	Lists      []Subscription   `json:"lists"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

// Create a new subscriber.
func (c *Client) CreateSubscriber(ctx context.Context, params *CreateSubscriberParams) (*CreateSubscriberResponse, error) {
	path := "/api/subscribers"
	resp, err := request[Response[*CreateSubscriberResponse]](c, ctx, "POST", path, params)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Sends optin confirmation email to subscribers.
func (c *Client) SendOptinConfirmation(ctx context.Context, id int) (bool, error) {
	path := fmt.Sprintf("/api/subscribers/%d/optin", id)
	resp, err := request[Response[bool]](c, ctx, "POST", path, nil)
	if err != nil {
		return false, err
	}
	return resp.Data, nil
}

type CreateSubscriptionParams struct {
	// Subscriber's email address.
	Email string `json:"email"`
	// Subscriber's name.
	Name string `json:"name"`
	// List of list UUIDs.
	ListUUIDs []uuid.UUID `json:"list_uuids"`
}

// Create a public subscription.
func (c *Client) CreateSubscription(ctx context.Context, params *CreateSubscriptionParams) (bool, error) {
	path := "/api/public/subscription"
	_, err := request[Response[any]](c, ctx, "POST", path, params)
	if err != nil {
		return false, err
	}
	return true, nil
}

/*

// // Modify subscriber list memberships.
// func (c *Client) CreateSubscription(ctx context.Context) {
	// 	path := "/api/subscribers/lists"
	// }

	// Update a specific subscriber.
func (c *Client) UpdateSubscriber(ctx context.Context, id int) {
	path := fmt.Sprintf("/api/subscribers/%d", id)
}

*/

// Blocklist a specific subscriber.
func (c *Client) BlocklistSubscriber(ctx context.Context, id int) (bool, error) {
	path := fmt.Sprintf("/api/subscribers/%d/blocklist", id)
	resp, err := request[Response[bool]](c, ctx, "PUT", path, nil)
	if err != nil {
		return false, err
	}
	return resp.Data, nil
}

// Blocklist one or many subscribers.
func (c *Client) BlocklistSubscribers(ctx context.Context, ids []int) (bool, error) {
	path := "/api/subscribers/blocklist"
	type params struct {
		IDs []int `json:"ids"`
	}
	resp, err := request[Response[bool]](c, ctx, "PUT", path, params{IDs: ids})
	if err != nil {
		return false, err
	}
	return resp.Data, nil
}

type BlocklistSubscribersWithQueryParams struct {
	// SQL expression to filter subscribers with.
	Query string `json:"query"`
	// Optional list IDs to limit the filtering to.
	ListIDs []int `json:"list_ids"`
}

// Blocklist subscribers based on SQL expression.
func (c *Client) BlocklistSubscribersWithQuery(ctx context.Context, params *BlocklistSubscribersWithQueryParams) (bool, error) {
	path := "/api/subscribers/query/blocklist"
	resp, err := request[Response[bool]](c, ctx, "PUT", path, params)
	if err != nil {
		return false, err
	}
	return resp.Data, nil
}

// Delete a specific subscriber.
func (c *Client) DeleteSubscriber(ctx context.Context, id int) (bool, error) {
	path := fmt.Sprintf("/api/subscribers/%d", id)
	resp, err := request[Response[bool]](c, ctx, "DELETE", path, nil)
	if err != nil {
		return false, err
	}
	return resp.Data, nil
}

// Delete a specific subscriber's bounce records.
func (c *Client) DeleteSubscriberBounces(ctx context.Context, id int) (bool, error) {
	path := fmt.Sprintf("/api/subscribers/%d/bounces", id)
	resp, err := request[Response[bool]](c, ctx, "DELETE", path, nil)
	if err != nil {
		return false, err
	}
	return resp.Data, nil
}

// Delete one or more subscribers.
func (c *Client) DeleteSubscribers(ctx context.Context, ids []int) (bool, error) {
	path := "/api/subscribers"
	type params struct {
		IDs []int `url:"id"`
	}
	resp, err := request[Response[bool]](c, ctx, "DELETE", path, params{IDs: ids})
	if err != nil {
		return false, err
	}
	return resp.Data, nil
}

type DeleteSubscribersWithQueryParams struct {
	// SQL expression to filter subscribers with.
	Query string `json:"query"`
	// Optional list IDs to limit the filtering to.
	ListIDs []int `json:"list_ids"`
	// When set to true, ignores any query and deletes all subscribers.
	All bool `json:"all"`
}

// Delete subscribers based on SQL expression.
func (c *Client) DeleteSubscribersWithQuery(ctx context.Context, params *DeleteSubscribersWithQueryParams) (bool, error) {
	path := "/api/subscribers/query/delete"
	resp, err := request[Response[bool]](c, ctx, "POST", path, params)
	if err != nil {
		return false, err
	}
	return resp.Data, nil
}

type GetListsParams struct {
	// String for list name search.
	Query string `url:"query"`
	// Status to filter lists. Repeat in the query for multiple values.
	Status []string `url:"status"`
	// Tags to filter lists. Repeat in the query for multiple values.
	Tag []string `url:"tag"`
	// Sort field. Options: name, status, created_at, updated_at.
	OrderBy string `url:"order_by"`
	// Sorting order. Options: ASC, DESC.
	Order string `url:"order"`
	// Page number for pagination.
	Page int `url:"page"`
	// Results per page. Set to 'all' to return all results.
	PerPage int `url:"per_page"`
}

type ListType string

const (
	PublicTypeListEntry  ListType = "public"
	PrivateTypeListEntry ListType = "private"
)

type ListOptin string

const (
	SingleOptinListEntry ListOptin = "single"
	DoubleOptinListEntry ListOptin = "double"
)

type List struct {
	ID              int       `json:"id"`
	UUID            uuid.UUID `json:"uuid"`
	Type            ListType  `json:"type"`
	Optin           ListOptin `json:"optin"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Tags            []string  `json:"tags"`
	SubscriberCount int       `json:"subscriber_count"`
	// TODO: add SubscriberStatuses field
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetListsResponse struct {
	Results []List `json:"results"`
	Total   int    `json:"total"`
	PerPage int    `json:"per_page"`
	Page    int    `json:"page"`
}

// Retrieve lists.
func (c *Client) GetLists(ctx context.Context, params *GetListsParams) (*GetListsResponse, error) {
	path := "/api/lists"
	resp, err := request[Response[*GetListsResponse]](c, ctx, "GET", path, params)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

type PublicListEntry struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

// Retrieve public lists with name and uuid to submit a subscription. This is an unauthenticated call to enable scripting to subscription form.
func (c *Client) GetPublicLists(ctx context.Context) ([]PublicListEntry, error) {
	path := "/api/public/lists"
	resp, err := request[[]PublicListEntry](c, ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Retrieve a specific list.
func (c *Client) GetList(ctx context.Context, id int) (*List, error) {
	path := fmt.Sprintf("/api/lists/%d", id)
	resp, err := request[Response[*List]](c, ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

type CreateListParams struct {
	//	Name of the new list.
	Name string `json:"name"`
	//	Type of list. Options: private, public.
	Type ListType `json:"type"`
	//	Opt-in type. Options: single, double.
	Optin ListOptin `json:"optin"`
	//	Associated tags for a list.
	Tags []string `json:"tags"`
	//	Description of the new list.
	Description string `json:"description"`
}

// Create a new list.
func (c *Client) CreateList(ctx context.Context, params *CreateListParams) (*List, error) {
	path := "/api/lists"
	resp, err := request[Response[*List]](c, ctx, "POST", path, params)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

type UpdateListParams struct {
	//	New name for the list.
	Name string `json:"name"`
	//	Type of list. Options: private, public.
	Type ListType `json:"type"`
	//	Opt-in type. Options: single, double.
	Optin ListOptin `json:"optin"`
	//	Associated tags for a list.
	Tags []string `json:"tags"`
	//	Description of the new list.
	Description string `json:"description"`
}

func (c *Client) UpdateList(ctx context.Context, id int, params *UpdateListParams) (*List, error) {
	path := fmt.Sprintf("/api/lists/%d", id)
	resp, err := request[Response[*List]](c, ctx, "PUT", path, params)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Delete a specific list.
func (c *Client) DeleteList(ctx context.Context, id int) (bool, error) {
	path := fmt.Sprintf("/api/lists/%d", id)
	resp, err := request[Response[bool]](c, ctx, "DELETE", path, nil)
	if err != nil {
		return false, err
	}
	return resp.Data, nil
}

type TemplateContentType string

const (
	HTMLTemplate     TemplateContentType = "html"
	MarkdownTemplate TemplateContentType = "markdown"
	PlainTemplate    TemplateContentType = "plain"
)

type SendTemplateParams struct {
	// Email of the subscriber. Can substitute with subscriber_id.
	SubscriberEmail string `json:"subscriber_email"`
	// Subscriber's ID can substitute with subscriber_email.
	SubscriberID int `json:"subscriber_id"`
	// Multiple subscriber emails as alternative to subscriber_email.
	SubscriberEmails []string `json:"subscriber_emails"`
	//		Multiple subscriber IDs as an alternative to subscriber_id.
	SubscriberIDs []int `json:"subscriber_ids"`
	//	ID of the transactional template to be used for the message.
	TemplateID int `json:"template_id"`
	//	Optional sender email.
	FromEmail string `json:"from_email"`
	// Optional nested JSON map. Available in the template as {{ .Tx.Data.* }}.
	Data any `json:"data"`
	//		Optional array of email headers.
	Headers http.Header `json:"headers"`
	//	Messenger to send the message. Default is email.
	Messenger string `json:"messenger"`
	// Email format options include html, markdown, and plain.
	ContentType TemplateContentType `json:"content_type"`
}

func (p *SendTemplateParams) MarshalJSON() ([]byte, error) {
	headers := []map[string]any{}

	for key := range p.Headers {
		mp := map[string]any{}
		mp[key] = p.Headers.Get(key)
		headers = append(headers, mp)
	}

	return json.Marshal(map[string]any{
		"subscriber_email":  p.SubscriberEmail,
		"subscriber_id":     p.SubscriberID,
		"subscriber_emails": p.SubscriberEmails,
		"subscriber_ids":    p.SubscriberIDs,
		"template_id":       p.TemplateID,
		"from_email":        p.FromEmail,
		"data":              p.Data,
		"headers":           headers,
		"messenger":         p.Messenger,
		"content-type":      p.ContentType,
	})
}

type GetImportStatisticsResponse struct {
	Name     string `json:"name"`
	Total    int    `json:"total"`
	Imported int    `json:"imported"`
	Status   string `json:"status"`
}

// Retrieve import statistics.
func (c *Client) GetImportStatistics(ctx context.Context) (*GetImportStatisticsResponse, error) {
	path := "/api/import/subscribers"
	resp, err := request[Response[*GetImportStatisticsResponse]](c, ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Retrieve import logs.
func (c *Client) GetImportLogs(ctx context.Context) (string, error) {
	path := "/api/import/subscribers/logs"
	resp, err := request[Response[string]](c, ctx, "GET", path, nil)
	if err != nil {
		return "", err
	}
	return resp.Data, nil
}

type ImportMode string

const (
	ImportModeSubscribe ImportMode = "subscribe"
	ImportModeBlocklist ImportMode = "blocklist"
)

type ImportSubscribersConfig struct {
	// subscribe or blocklist
	Mode ImportMode `json:"mode"`
	// Single character indicating delimiter used in the CSV file, eg: ,
	Delimeter rune `json:"delim"`
	// List IDs to add subscribers to.
	Lists []int `json:"lists"`
	//	Whether to overwrite the subscriber parameters including subscriptions or ignore records that are already present in the database.
	Overwrite bool `json:"overwrite"`
}

type ImportSubscribersParams struct {
	Config ImportSubscribersConfig
	File   io.Reader
}

type ImportSubscribersResponse struct {
	Mode      ImportMode `json:"mode"`
	Delimeter rune       `json:"delim"`
	Lists     []int      `json:"lists"`
	Overwrite bool       `json:"overwrite"`
}

// Send a CSV (optionally ZIP compressed) file to import subscribers.
func (c *Client) ImportSubscribers(ctx context.Context, params *ImportSubscribersParams) (*ImportSubscribersResponse, error) {
	path := "/api/import/subscribers"
	config, err := json.Marshal(params.Config)
	if err != nil {
		return nil, err
	}
	resp, err := c.multipart(ctx, path, map[string]string{"params": string(config)}, map[string]io.Reader{"file": params.File})
	if err != nil {
		return nil, err
	}
	data, err := decode[Response[*ImportSubscribersResponse]](resp)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

type AbortImportResponse struct {
	Name     string `json:"name"`
	Total    int    `json:"total"`
	Imported int    `json:"imported"`
	Status   string `json:"status"`
}

// Stop and remove an import.
func (c *Client) AbortImport(ctx context.Context) (*AbortImportResponse, error) {
	path := "/api/import/subscribers"
	resp, err := request[Response[*AbortImportResponse]](c, ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Allows sending transactional messages to one or more subscribers via a preconfigured transactional template.
func (c *Client) SendTemplate(ctx context.Context, params *SendTemplateParams) (bool, error) {
	path := "/api/tx"
	resp, err := request[Response[bool]](c, ctx, "POST", path, params)
	if err != nil {
		return false, err
	}
	return resp.Data, nil
}
