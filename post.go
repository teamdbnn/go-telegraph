package telegraph

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

type CreatePageParams struct {
	// AuthorName Author name, displayed below the article's title.
	AuthorName string
	// AuthorURL Profile link, opened when users click on the author's name below the title. Can be any link, not necessarily to a Telegram profile or channel
	AuthorURL string
	// ReturnContent If true, a content field will be returned in the Page object.
	ReturnContent bool
}

// CreatePage Use this method to create a new Telegraph page. On success, returns a Page object.
// https://telegra.ph/api#createPage
func (c *Client) CreatePage(ctx context.Context, title string, content []Node, params *CreatePageParams, opts ...RequestOption) (*Page, error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "createPage",
		secured:  true,
	}

	r.setFormParam("title", title)
	contentData, err := json.Marshal(&content)
	if err != nil {
		return nil, err
	}
	r.setFormParam("content", string(contentData))
	if params != nil {
		if params.AuthorName != "" {
			r.setFormParam("author_name", params.AuthorName)
		}
		if params.AuthorURL != "" {
			r.setFormParam("author_url", params.AuthorURL)
		}
		if params.ReturnContent {
			r.setFormParam("return_content", params.ReturnContent)
		}
	}

	resp, err := c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	page := new(responsePage)
	if err = json.Unmarshal(resp, page); err != nil {
		return nil, err
	}
	if !page.OK {
		return nil, errors.New(page.Error)
	}
	return page.Result, nil
}

type EditPageParams struct {
	// AuthorName Author name, displayed below the article's title.
	AuthorName string
	// AuthorURL Profile link, opened when users click on the author's name below the title. Can be any link, not necessarily to a Telegram profile or channel
	AuthorURL string
	// ReturnContent If true, a content field will be returned in the Page object.
	ReturnContent bool
}

// EditPage Use this method to edit an existing Telegraph page. On success, returns a Page object.
// https://telegra.ph/api#editPage
func (c *Client) EditPage(ctx context.Context, path, title string, content []Node, params *EditPageParams, opts ...RequestOption) (*Page, error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: fmt.Sprintf("%v/%v", "editPage", path),
	}
	r.setFormParam("title", title)
	contentData, err := json.Marshal(&content)
	if err != nil {
		return nil, err
	}
	r.setFormParam("content", string(contentData))
	if params != nil {
		if params.AuthorName != "" {
			r.setFormParam("author_name", params.AuthorName)
		}
		if params.AuthorURL != "" {
			r.setFormParam("author_url", params.AuthorURL)
		}
		if params.ReturnContent {
			r.setFormParam("return_content", params.ReturnContent)
		}
	}
	httpResponse, err := c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	page := new(responsePage)
	if err = json.Unmarshal(httpResponse, page); err != nil {
		return nil, err
	}
	if !page.OK {
		return nil, errors.New(page.Error)
	}
	return page.Result, nil
}

type GetPageParams struct {
	// ReturnContent If true, content field will be returned in Page object.
	ReturnContent bool
}

// GetPage Use this method to get a Telegraph page. Returns a Page object on success.
// https://telegra.ph/api#getPage
func (c *Client) GetPage(ctx context.Context, path string, option *GetPageParams, opts ...RequestOption) (*Page, error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: fmt.Sprintf("%v/%v", "getPage", path),
	}

	if option != nil {
		if option.ReturnContent {
			r.setFormParam("return_content", option.ReturnContent)
		}
	}
	resp, err := c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	page := new(responsePage)
	if err = json.Unmarshal(resp, page); err != nil {
		return nil, err
	}
	if !page.OK {
		return nil, errors.New(page.Error)
	}
	return page.Result, nil
}

type GetPageListParams struct {
	// Offset Sequential number of the first page to be returned.
	Offset int
	// Limit Limits the number of pages to be retrieved.
	Limit int
}

// GetPageList Use this method to get a list of pages belonging to a Telegraph account. Returns a PageList object, sorted by most recently created pages first.
// https://telegra.ph/api#getPageList
func (c *Client) GetPageList(ctx context.Context, params *GetPageListParams, opts ...RequestOption) (*PageList, error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "getPageList",
		secured:  true,
	}

	if params != nil {
		if params.Offset > 1 {
			r.setFormParam("offset", params.Offset)
		}
		if params.Limit > 1 {
			r.setFormParam("limit", params.Limit)
		}
	}

	resp, err := c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	pageList := new(responsePageList)
	if err = json.Unmarshal(resp, pageList); err != nil {
		return nil, err
	}
	if !pageList.OK {
		return nil, errors.New(pageList.Error)
	}
	return pageList.Result, nil
}

type GetViewsParams struct {
	// Year (Integer, 2000-2100)
	// Required if month is passed. If passed, the number of page views for the requested year will be returned.
	Year int64
	// Month (Integer, 1-12)
	// Required if day is passed. If passed, the number of page views for the requested month will be returned.
	Month int64
	// Day (Integer, 1-31)
	// Required if hour is passed. If passed, the number of page views for the requested day will be returned.
	Day int64
	// Hour (0-24)
	// If passed, the number of page views for the requested hour will be returned.
	Hour int64
}

// GetViews Use this method to get the number of views for a Telegraph article.
// Returns a PageViews object on success. By default, the total number of page views will be returned.
// https://telegra.ph/api#getViews
func (c *Client) GetViews(ctx context.Context, path string, option *GetViewsParams, opts ...RequestOption) (*PageViews, error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: fmt.Sprintf("%v/%v", "getViews", path),
	}

	if option != nil {
		if option.Year > 0 {
			r.setFormParam("year", option.Year)
		}
		if option.Month > 0 {
			r.setFormParam("month", option.Month)
		}
		if option.Day > 0 {
			r.setFormParam("day", option.Day)
		}
		if option.Hour > 0 {
			r.setFormParam("hour", option.Hour)
		}
	}
	resp, err := c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	pageView := new(responsePageViews)
	if err = json.Unmarshal(resp, pageView); err != nil {
		return nil, err
	}
	if !pageView.OK {
		return nil, errors.New(pageView.Error)
	}
	return pageView.Result, nil
}

func (c *Client) Upload(ctx context.Context, filenames []string, opts ...RequestOption) ([]string, error) {
	files := make([]*os.File, 0, len(filenames))

	// Close the file handle finished processing.
	defer func() {
		for _, file := range files {
			file.Close()
		}
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		files = append(files, file)

		part, err := writer.CreateFormFile(fmt.Sprintf("%x", sha256.Sum256([]byte(filename))), filename)
		if err != nil {
			return nil, err
		}

		if _, err = io.Copy(part, file); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	c.BaseURL = baseURL
	r := &request{
		method:   http.MethodPost,
		endpoint: "upload",
		body:     body,
	}

	opts = append(opts, WithHeader("Content-Type", writer.FormDataContentType(), false))

	resp, err := c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	upload := make([]responseUpload, 0)
	if err = json.Unmarshal(resp, &upload); err != nil {
		m := map[string]string{}
		if err = json.Unmarshal(resp, &m); err != nil {
			return nil, err
		}

		return nil, errors.New(m["error"])
	}

	paths := make([]string, 0, len(upload))
	for _, u := range upload {
		paths = append(paths, u.Path)
	}

	return paths, nil
}
