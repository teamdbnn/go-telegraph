package telegraph

type response struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type responseAccount struct {
	response
	Result *Account `json:"result,omitempty"`
}

type responsePage struct {
	response
	Result *Page `json:"result,omitempty"`
}
type responsePageList struct {
	response
	Result *PageList `json:"result,omitempty"`
}

type responsePageViews struct {
	response
	Result *PageViews `json:"result,omitempty"`
}

type responseUpload struct {
	Path string `json:"src"`
}

// Account represents a Telegraph account.
type Account struct {
	// Only returned by the createAccount and revokeAccessToken method. Access token of the Telegraph
	// account.
	AccessToken string `json:"access_token"`

	// URL to authorize a browser on telegra.ph and connect it to a Telegraph account. This URL is valid
	// for only one use and for 5 minutes only.
	AuthURL string `json:"auth_url,omitempty"`

	// Account name, helps users with several accounts remember which they are currently using. Displayed
	// to the user above the "Edit/Publish" button on Telegra.ph, other users don't see this name.
	ShortName string `json:"short_name"`

	// Default author name used when creating new articles.
	AuthorName string `json:"author_name"`

	// Profile link, opened when users click on the author's name below the title. Can be any link, not
	// necessarily to a Telegram profile or channel.
	AuthorURL string `json:"author_url"`

	// Number of pages belonging to the Telegraph account.
	PageCount int `json:"page_count,omitempty"`
}

// PageList represents a list of Telegraph articles belonging to an account. Most recently created articles
// first.
type PageList struct {
	// Total number of pages belonging to the target Telegraph account.
	TotalCount int `json:"total_count"`

	// Requested pages of the target Telegraph account.
	Pages []Page `json:"pages"`
}

// Page represents a page on Telegraph.
type Page struct {
	// Path to the page.
	Path string `json:"path"`

	// URL of the page.
	URL string `json:"url"`

	// Title of the page.
	Title string `json:"title"`

	// Description of the page.
	Description string `json:"description"`

	// Name of the author, displayed below the title.
	AuthorName string `json:"author_name,omitempty"`

	// Profile link, opened when users click on the author's name below the title. Can be any link, not
	// necessarily to a Telegram profile or channel.
	AuthorURL string `json:"author_url,omitempty"`

	// Image URL of the page.
	ImageURL string `json:"image_url,omitempty"`

	// Content of the page.
	Content []Node `json:"content,omitempty"`

	// Number of page views for the page.
	Views int `json:"views"`

	// Only returned if access_token passed. True, if the target Telegraph account can edit the page.
	CanEdit bool `json:"can_edit,omitempty"`
}

// PageViews represents the number of page views for a Telegraph article.
type PageViews struct {
	// Number of page views for the target page.
	Views int `json:"views"`
}

// Node is abstract object represents a DOM Node. It can be a String which represents a DOM text node or a
// NodeElement object.
type Node interface{}

var _ Node = &NodeElement{}

// NodeElement represents a DOM element node.
type NodeElement struct {
	// Name of the DOM element. Available tags: a, aside, b, blockquote, br, code, em, figcaption, figure,
	// h3, h4, hr, i, iframe, img, li, ol, p, pre, s, strong, u, ul, video.
	Tag string `json:"tag"`

	// Attributes of the DOM element. Key of object represents name of attribute, value represents value
	// of attribute. Available attributes: href, src.
	Attrs map[string]string `json:"attrs,omitempty"`

	// List of child nodes for the DOM element.
	Children []Node `json:"children,omitempty"`
}
