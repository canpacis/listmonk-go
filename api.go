package listmonkgo

import (
	"context"
	"encoding/json"
	"errors"
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

type UpdateListMembershipsParams struct {
	// Array of user IDs to be modified.
	IDs []int `json:"ids"`
	// Action to be applied: add, remove, or unsubscribe.
	Acion string `json:"action"`
	// Array of list IDs to be modified.
	TargetListIDs []int `json:"target_list_ids"`
	// Required for add	Subscriber status: confirmed, unconfirmed, or unsubscribed.
	Status string `json:"status"`
}

// Modify subscriber list memberships.
func (c *Client) UpdateListMemberships(ctx context.Context, params *UpdateListMembershipsParams) (bool, error) {
	path := "/api/subscribers/lists"
	resp, err := request[Response[bool]](c, ctx, "PUT", path, params)
	if err != nil {
		return false, err
	}
	return resp.Data, nil
}

// Update a specific subscriber.
// Note: All parameters must be set, if not, the subscriber will be removed from all previously assigned lists.
func (c *Client) UpdateSubscriber(ctx context.Context, id int, params *CreateSubscriberParams) (*CreateSubscriberResponse, error) {
	path := fmt.Sprintf("/api/subscribers/%d", id)
	resp, err := request[Response[*CreateSubscriberResponse]](c, ctx, "PUT", path, params)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

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
	Delimeter string `json:"delim"`
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
	Delimeter string     `json:"delim"`
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

type CampaignType string

const (
	RegularCampaign CampaignType = "regular"
	OptinCampaign   CampaignType = "optin"
)

type Campaign struct {
	ID          int          `json:"id"`
	TemplateID  int          `json:"template_id"`
	UUID        uuid.UUID    `json:"uuid"`
	Type        CampaignType `json:"type"`
	Messenger   string       `json:"messenger"`
	ContentType string       `json:"content_type"`

	Name       string   `json:"name"`
	Subject    string   `json:"subject"`
	FromEmail  string   `json:"from_email"`
	Body       string   `json:"body"`
	BodySource string   `json:"body_source"`
	AltBody    string   `json:"alt_body"`
	Status     string   `json:"status"`
	Tags       []string `json:"tags"`

	Media []struct {
		ID       int    `json:"id"`
		Filename string `json:"filename"`
	} `json:"media"`

	Lists []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"lists"`

	Headers           []map[string]any `json:"headers"`
	Archive           bool             `json:"archive"`
	ArchiveSlug       string           `json:"archive_slug"`
	ArchiveTemplateID int              `json:"archive_template_id"`
	ArchiveMeta       map[string]any   `json:"archive_meta"`

	Views   int `json:"views"`
	Clicks  int `json:"clicks"`
	Bounces int `json:"bounces"`
	Sent    int `json:"sent"`
	ToSend  int `json:"to_send"`

	SendAt    time.Time `json:"send_at"`
	StartedAt time.Time `json:"started_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetCampaignParams struct {
	// Sorting order: ASC for ascending, DESC for descending.
	Order string `url:"order"`
	// Result sorting field. Options: name, status, created_at, updated_at.
	OrderBy string `url:"order_by"`
	//SQL query expression to filter campaigns.
	Query string `url:"query"`
	// Status to filter campaigns. Repeat in the query for multiple values.
	Status []string `url:"status"`
	// Tags to filter campaigns. Repeat in the query for multiple values.
	Tags []string `url:"tags"`
	// Page number for paginated results.
	Page int `url:"page"`
	// Results per page. Set as 'all' for all results.
	PerPage int `url:"per_page"`
	// When set to true, returns response without body content.
	NoBody bool `url:"no_body"`
}

type GetCampaignResponse struct {
	Results []Campaign `json:"results"`
	Search  string     `json:"search"`
	Query   string     `json:"query"`
	Total   int        `json:"total"`
	Page    int        `json:"page"`
	PerPage int        `json:"per_page"`
}

// Retrieve all campaigns.
func (c *Client) GetCampaigns(ctx context.Context, params *GetCampaignParams) (*GetCampaignResponse, error) {
	path := "/api/campaigns"
	resp, err := request[Response[*GetCampaignResponse]](c, ctx, "GET", path, params)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Retrieve a specific campaign.
func (c *Client) GetCampaign(ctx context.Context, id int, noBody bool) (*Campaign, error) {
	path := fmt.Sprintf("/api/campaigns/%d", id)
	type params struct {
		NoBody bool `url:"no_body"`
	}
	resp, err := request[Response[*Campaign]](c, ctx, "GET", path, params{NoBody: noBody})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Retrieve preview of a campaign.
func (c *Client) GetCampaignPreview(ctx context.Context, id int) (string, error) {
	path := fmt.Sprintf("/api/campaigns/%d/preview", id)
	resp, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		decoder := json.NewDecoder(resp.Body)
		data := new(ErrorResponse)
		if err := decoder.Decode(data); err != nil {
			return "", err
		}
		return "", errors.New(data.Message)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Retrieve stats of specified campaigns.
// API Docs do not provide any meaningful documentation for the response of this action
func (c *Client) GetCampaignsStats(ctx context.Context, ids []int) ([]any, error) {
	path := "/api/campaigns/running/stats"
	type params struct {
		IDs []int `url:"id"`
	}
	resp, err := request[Response[[]any]](c, ctx, "GET", path, params{IDs: ids})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

type CampaignStatType string

const (
	ViewCampaignStat   CampaignStatType = "views"
	ClickCampaignStat  CampaignStatType = "clicks"
	LinkCampaignStat   CampaignStatType = "links"
	BounceCampaignStat CampaignStatType = "bounces"
)

type GetCampaignViewsParams struct {
	// Campaign IDs to get stats for.
	IDs []int `url:"id"`
	//	Analytics type: views, links, clicks, bounces
	Type CampaignStatType `url:"-"`
	// Start value of date range.
	From time.Time `url:"from"`
	// End value of date range.
	To time.Time `url:"to"`
}

type TimeSeriesData []map[string]any

// Retrieve view counts for a campaign.
func (c *Client) GetCampaignViews(ctx context.Context, params *GetCampaignViewsParams) (TimeSeriesData, error) {
	path := fmt.Sprintf("/api/campaigns/analytics/%s", params.Type)
	resp, err := request[Response[TimeSeriesData]](c, ctx, "GET", path, params)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

type CreateCampaignParams struct {
	// Campaign name.
	Name string `json:"name"`
	// Campaign email subject.
	Subject string `json:"subject"`
	// List IDs to send campaign to.
	Lists []int `json:"lists"`
	// 'From' email in campaign emails. Defaults to value from settings if not provided.
	FromEmail string `json:"from_email"`
	// Campaign type: 'regular' or 'optin'.
	Type CampaignType `json:"type"`
	// Content type: 'richtext', 'html', 'markdown', 'plain', 'visual'.
	ContentType string `json:"content_type"`
	// Content body of campaign.
	Body string `json:"body"`
	// If content_type is visual, the JSON block source of the body.
	BodySource string `json:"body_source"`
	// Alternate plain text body for HTML (and richtext) emails.
	Altbody string `json:"altbody"`
	// Timestamp to schedule campaign. Format: 'YYYY-MM-DDTHH:MM:SSZ'.
	SendAt time.Time `json:"send_at"`
	// 'email' or a custom messenger defined in settings. Defaults to 'email' if not provided.
	Messenger string `json:"messenger"`
	// Template ID to use. Defaults to default template if not provided.
	TemplateID int `json:"template_id"`
	// Tags to mark campaign.
	Tags []string `json:"tags"`
	// Key-value pairs to send as SMTP headers. Example: [{"x-custom-header": "value"}].
	Headers []map[string]any `json:"headers"`
}

// Create a new campaign.
func (c *Client) CreateCampaign(ctx context.Context, params *CreateCampaignParams) (*Campaign, error) {
	path := "/api/campaigns"
	resp, err := request[Response[*Campaign]](c, ctx, "POST", path, params)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Test campaign with arbitrary subscribers.
func (c *Client) TestCampaign(ctx context.Context, id int, subscribers []string) error {
	path := fmt.Sprintf("/api/campaigns/%d/test", id)
	type params struct {
		Subscribers []string `json:"subscribers"`
	}
	_, err := request[any](c, ctx, "POST", path, params{Subscribers: subscribers})
	if err != nil {
		return err
	}
	return nil
}

// Update a campaign.
func (c *Client) UpdateCampaign(ctx context.Context, id int, params *CreateCampaignParams) (*Campaign, error) {
	path := fmt.Sprintf("/api/campaigns/%d", id)
	resp, err := request[Response[*Campaign]](c, ctx, "PUT", path, params)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

type CampaignStatus string

const (
	CampaignStatusScheduled CampaignStatus = "scheduled"
	CampaignStatusRunning   CampaignStatus = "running"
	CampaignStatusPaused    CampaignStatus = "paused"
	CampaignStatusCancelled CampaignStatus = "cancelled"
)

// Change status of a campaign.
func (c *Client) ChangeCampaignStatus(ctx context.Context, id int, status CampaignStatus) (*Campaign, error) {
	path := fmt.Sprintf("/api/campaigns/%d/status", id)
	type params struct {
		Status CampaignStatus `json:"status"`
	}
	resp, err := request[Response[*Campaign]](c, ctx, "PUT", path, params{Status: status})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

type ArchiveCampaignParams struct {
	// State of the public archive.
	Archive bool `json:"archive"`
	// Archive template id. Defaults to 0.
	ArchiveTemplateID int `json:"archive_template_id"`
	// Optional Metadata to use in campaign message or template. Eg: name, email.
	ArchiveMeta map[string]any `json:"archive_meta"`
	// Name for page to be used in public archive URL
	ArchiveSlug string `json:"archive_slug"`
}

type ArchiveCampaignResponse struct {
	Archive           bool           `json:"archive"`
	ArchiveTemplateID int            `json:"archive_template_id"`
	ArchiveMeta       map[string]any `json:"archive_meta"`
	ArchiveSlug       string         `json:"archive_slug"`
}

// Publish campaign to public archive.
func (c *Client) ArchiveCampaign(ctx context.Context, id int, params *ArchiveCampaignParams) (*ArchiveCampaignParams, error) {
	path := fmt.Sprintf("/api/campaigns/%d/archive", id)
	resp, err := request[Response[*ArchiveCampaignParams]](c, ctx, "PUT", path, params)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Delete a campaign.
func (c *Client) DeleteCampaign(ctx context.Context, id int) (bool, error) {
	path := fmt.Sprintf("/api/campaigns/%d", id)
	resp, err := request[Response[bool]](c, ctx, "DELETE", path, nil)
	if err != nil {
		return false, err
	}
	return resp.Data, nil
}

type GetMediaListResponse struct {
	ID        int       `json:"id"`
	UUID      uuid.UUID `json:"uuid"`
	Filename  string    `json:"filename"`
	URI       string    `json:"uri"`
	ThumbURL  string    `json:"thumb_url"`
	CreatedAt time.Time `json:"created_at"`
}

// Get uploaded media files.
func (c *Client) GetMediaList(ctx context.Context) ([]GetMediaListResponse, error) {
	path := "/api/media"
	resp, err := request[Response[[]GetMediaListResponse]](c, ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

type Media struct {
	ID          int            `json:"id"`
	UUID        uuid.UUID      `json:"uuid"`
	Filename    string         `json:"filename"`
	ContentType string         `json:"content_type"`
	ThumbURL    string         `json:"thumb_url"`
	URL         string         `json:"url"`
	Provider    string         `json:"provider"`
	Meta        map[string]any `json:"meta"`
	CreatedAt   time.Time      `json:"created_at"`
}

// Get specific uploaded media file.
func (c *Client) GetMedia(ctx context.Context, id int) (*Media, error) {
	path := fmt.Sprintf("/api/media/%d", id)
	resp, err := request[Response[*Media]](c, ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

type UploadMediaResponse struct {
	ID        int       `json:"id"`
	UUID      uuid.UUID `json:"uuid"`
	Filename  string    `json:"filename"`
	URI       string    `json:"uri"`
	ThumbURI  string    `json:"thumb_uri"`
	CreatedAt time.Time `json:"created_at"`
}

// Upload media file.
func (c *Client) UploadMedia(ctx context.Context, file io.Reader) (*UploadMediaResponse, error) {
	path := "/api/media"
	resp, err := c.multipart(ctx, path, map[string]string{}, map[string]io.Reader{"file": file})
	if err != nil {
		return nil, err
	}
	data, err := decode[Response[*UploadMediaResponse]](resp)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

// Delete uploaded media file.
func (c *Client) DeleteMedia(ctx context.Context, id int) (bool, error) {
	path := fmt.Sprintf("/api/media/%d", id)
	resp, err := request[Response[bool]](c, ctx, "DELETE", path, nil)
	if err != nil {
		return false, err
	}
	return resp.Data, nil
}

type Template struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Body      string    `json:"body"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// BodySource string `json:"body_source"`
}

// Retrieve all templates.
func (c *Client) GetTemplates(ctx context.Context) ([]Template, error) {
	path := "/api/templates"
	resp, err := request[Response[[]Template]](c, ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Retrieve a template.
func (c *Client) GetTemplate(ctx context.Context, id int) (*Template, error) {
	path := fmt.Sprintf("/api/templates/%d", id)
	resp, err := request[Response[*Template]](c, ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Retrieve template HTML preview.
func (c *Client) GetTemplatePreview(ctx context.Context, id int) (string, error) {
	path := fmt.Sprintf("/api/templates/%d/preview", id)
	resp, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		decoder := json.NewDecoder(resp.Body)
		data := new(ErrorResponse)
		if err := decoder.Decode(data); err != nil {
			return "", err
		}
		return "", errors.New(data.Message)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

type CreateTemplateParams struct {
	// Name of the template
	Name string `json:"name"`
	// Type of the template (campaign, campaign_visual, or tx)
	Type string `json:"type"`
	// Subject line for the template (only for tx)
	Subject string `json:"subject"`
	// HTML body of the template
	Body string `json:"body"`
	// If type is campaign_visual, the JSON source for the email-builder tempalate
	BodySource string `json:"body_source"`
}

// Create a template.
func (c *Client) CreateTemplate(ctx context.Context, params *CreateTemplateParams) (*Template, error) {
	path := "/api/templates"
	resp, err := request[Response[[]Template]](c, ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Data[0], nil
}

// Update a template.
func (c *Client) UpdateTemplate(ctx context.Context, id int) (*Template, error) {
	path := fmt.Sprintf("/api/templates/%d", id)
	resp, err := request[Response[[]Template]](c, ctx, "PUT", path, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Data[0], nil
}

// Set a template as the default.
func (c *Client) SetDefaultTemplate(ctx context.Context, id int) (*Template, error) {
	path := fmt.Sprintf("/api/templates/%d/default", id)
	resp, err := request[Response[*Template]](c, ctx, "PUT", path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Delete a template.
func (c *Client) DeleteTemplate(ctx context.Context, id int) (bool, error) {
	path := fmt.Sprintf("/api/templates/%d", id)
	resp, err := request[Response[bool]](c, ctx, "DELETE", path, nil)
	if err != nil {
		return false, err
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

type GetBouncesParams struct {
	// Bounce record retrieval for particular campaign id
	CompaignID int `url:"campaign_id"`
	// Page number for pagination.
	Page int `url:"page"`
	// Results per page. Set to 'all' to return all results.
	PerPage int `url:"per_page"`
	//
	Source string `url:"source"`
	// Fields by which bounce records are ordered. Options:"email", "campaign_name", "source", "created_at".
	OrderBy string `url:"order_by"`
	// Sorts the result. Allowed values: 'asc','desc'
	Order string `url:"order"`
}

type Bounce struct {
	ID             int            `json:"id"`
	Type           string         `json:"type"`
	Source         string         `json:"source"`
	Email          string         `json:"email"`
	SubscriberID   int            `json:"subscriber_id"`
	SubscriberUUID uuid.UUID      `json:"subscriber_uuid"`
	Meta           map[string]any `json:"meta"`
	Campaign       struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"campaign"`
	CreatedAt time.Time `json:"created_at"`
}

type GetBouncesResponse struct {
	Results []Bounce `json:"results"`
	Query   string   `json:"query"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
}

// Retrieve bounce records.
func (c *Client) GetBounces(ctx context.Context, params *GetBouncesParams) (*GetBouncesResponse, error) {
	path := "/api/bounces"
	resp, err := request[Response[*GetBouncesResponse]](c, ctx, "GET", path, params)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Delete all bounce records.
func (c *Client) DeleteAllBounces(ctx context.Context) (bool, error) {
	path := "/api/bounces"
	type params struct {
		All bool `url:"all"`
	}
	resp, err := request[Response[bool]](c, ctx, "DELETE", path, params{All: true})
	if err != nil {
		return false, err
	}
	return resp.Data, nil
}

// Delete multiple bounce records.
func (c *Client) DeleteBounces(ctx context.Context, ids []int) (bool, error) {
	path := "/api/bounces"
	type params struct {
		IDs []int `url:"id"`
	}
	resp, err := request[Response[bool]](c, ctx, "DELETE", path, params{IDs: ids})
	if err != nil {
		return false, err
	}
	return resp.Data, nil
}

// Delete specific bounce record.
func (c *Client) DeleteBounce(ctx context.Context, id int) (bool, error) {
	path := fmt.Sprintf("/api/bounces/%d", id)
	resp, err := request[Response[bool]](c, ctx, "DELETE", path, nil)
	if err != nil {
		return false, err
	}
	return resp.Data, nil
}
