package telegraph

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type CreateAccountParams struct {
	// AuthorName (String, 0-128 characters) Default author name used when creating new articles.
	AuthorName string
	// AuthorURL (String, 0-512 characters) Default profile link, opened when users click on the author's name below the title.
	// Can be any link, not necessarily to a Telegram profile or channel.
	AuthorURL string
}

// CreateAccount Use this method to create a new Telegraph account.
// Most users only need one account, but this can be useful for channel administrators who would like to keep individual author names and profile links for each of their channels.
// On success, returns an Account object with the regular fields and an additional access_token field.
// https://telegra.ph/api#createAccount
func (c *Client) CreateAccount(ctx context.Context, shortName string, params *CreateAccountParams, opts ...RequestOption) (*Account, error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "createAccount",
	}
	r.setFormParam("short_name", shortName)

	if params != nil {
		if params.AuthorURL != "" {
			r.setFormParam("author_url", params.AuthorURL)
		}
		if params.AuthorName != "" {
			r.setFormParam("author_name", params.AuthorName)
		}
	}

	resp, err := c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	acc := new(responseAccount)
	if err = json.Unmarshal(resp, acc); err != nil {
		return nil, err
	}

	if !acc.OK {
		return acc.Result, errors.New(acc.Error)
	}

	return acc.Result, nil
}

type EditAccountInfoParams struct {
	ShortName string
	// AuthorName (String, 0-128 characters) Default author name used when creating new articles.
	AuthorName string
	// AuthorURL (String, 0-512 characters) Default profile link, opened when users click on the author's name below the title.
	// Can be any link, not necessarily to a Telegram profile or channel.
	AuthorURL string
}

// EditAccountInfo Use this method to update information about a Telegraph account.
// Pass only the parameters that you want to edit. On success, returns an Account object with the default fields.
// https://telegra.ph/api#editAccountInfo
func (c *Client) EditAccountInfo(ctx context.Context, params *EditAccountInfoParams, opts ...RequestOption) (*Account, error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "editAccountInfo",
		secured:  true,
	}
	if params != nil {
		if params.ShortName != "" {
			r.setFormParam("short_name", params.ShortName)
		}
		if params.AuthorName != "" {
			r.setFormParam("author_name", params.AuthorName)
		}
		if params.AuthorURL != "" {
			r.setFormParam("author_url", params.AuthorURL)
		}
	}

	resp, err := c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	acc := new(responseAccount)
	if err = json.Unmarshal(resp, acc); err != nil {
		return nil, err
	}
	if !acc.OK {
		return nil, errors.New(acc.Error)
	}
	return acc.Result, nil
}

type GetAccountInfoOption struct {
	// Fields (Array of String, default = ["short_name","author_name","author_url"])
	// List of account fields to return. Available fields: short_name, author_name, author_url, auth_url, page_count.
	Fields []string
}

// GetAccountInfo Use this method to get information about a Telegraph account. Returns an Account object on success.
// https://telegra.ph/api#getAccountInfo
func (c *Client) GetAccountInfo(ctx context.Context, option *GetAccountInfoOption, opts ...RequestOption) (*Account, error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "getAccountInfo",
		secured:  true,
	}
	if option != nil {
		if len(option.Fields) > 0 {
			fields, err := json.Marshal(option.Fields)
			if err != nil {
				return nil, err
			}
			r.setFormParam("fields", string(fields))
		}
	}
	resp, err := c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	acc := new(responseAccount)
	if err = json.Unmarshal(resp, acc); err != nil {
		return nil, err
	}
	if !acc.OK {
		return nil, errors.New(acc.Error)
	}
	return acc.Result, nil
}

// RevokeAccessToken Use this method to revoke access_token and generate a new one,
// for example, if the user would like to reset all connected sessions, or you have reasons to believe the token was compromised.
// On success, returns an Account object with new access_token and auth_url fields.
// https://telegra.ph/api#revokeAccessToken
func (c *Client) RevokeAccessToken(ctx context.Context, opts ...RequestOption) (account *Account, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "revokeAccessToken",
		secured:  true,
	}

	resp, err := c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	acc := new(responseAccount)
	if err = json.Unmarshal(resp, acc); err != nil {
		return nil, err
	}
	if !acc.OK {
		return nil, errors.New(acc.Error)
	}
	return acc.Result, nil
}
