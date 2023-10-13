package main

import (
	"context"
	"log"

	"github.com/teamdbnn/go-telegraph"
)

func main() {
	// sandbox token d3b25feccb89e508a9114afb82aa421fe2a9712b963b387cc5ad71e58722
	token := "<TOKEN HERE>"

	client := telegraph.NewClient(token)

	log.Printf("> Loaded client: %#+v", client)

	ctx := context.TODO()

	// GetAccountInfo
	if account, err := client.GetAccountInfo(ctx, nil); err == nil {
		log.Printf("> GetAccountInfo result: %#+v", account)
	} else {
		log.Printf("* GetAccountInfo error: %s", err)
	}

	// EditAccountInfo
	if account, err := client.EditAccountInfo(ctx, &telegraph.EditAccountInfoParams{
		ShortName:  "Sandbox",
		AuthorName: "Anonymous",
	}); err == nil {
		log.Printf("> EditAccountInfo result: %#+v", account)
	} else {
		log.Printf("* EditAccountInfo error: %s", err)
	}

	// CreatePage
	content, err := telegraph.ContentFormat("<p>Hello, World!</p>")
	if err != nil {
		log.Fatalf("* ContentFormat error: %v", err)
	}
	page, err := client.CreatePage(ctx, "Test page", content, &telegraph.PageParams{
		AuthorName:    "",
		AuthorURL:     "",
		ReturnContent: true,
	})
	if err != nil {
		log.Printf("> CreatePage result: %#+v", page)
		log.Printf("> Created page url: %s", page.URL)
	}

	// GetPage
	page, err = client.GetPage(ctx, page.Path, &telegraph.GetPageParams{ReturnContent: true})
	if err == nil {
		log.Printf("> GetPage result: %#+v", page)
	} else {
		log.Printf("* GetPage error: %s", err)
	}

	// EditPage
	content, err = telegraph.ContentFormat("<p>Hello, New World!</p>")
	if err != nil {
		log.Fatalf("* ContentFormat error: %v", err)
	}

	page, err = client.EditPage(ctx, page.Path, "Test page (edited)", content, &telegraph.PageParams{
		ReturnContent: true,
	})
	if err == nil {
		log.Printf("> EditPage result: %#+v", page)
		log.Printf("> Edited page url: %s", page.URL)
	} else {
		log.Printf("* EditPage error: %s", err)
	}

	// GetPageList
	pages, err := client.GetPageList(ctx, &telegraph.GetPageListParams{
		Offset: 0,
		Limit:  50,
	})
	if err != nil {
		log.Fatalf("* GetPageList error: %s", err)
	}

	log.Printf("> GetPageList result: %#+v", pages)

	for _, p := range pages.Pages {
		// GetViews
		views, err0 := client.GetViews(ctx, p.Path, &telegraph.GetViewsParams{
			Year:  2016,
			Month: 1,
			Day:   1,
			Hour:  3,
		})
		if err0 != nil {
			log.Fatalf("* GetViews error: %s", err0)
		}
		log.Printf("> GetViews result for %s: %#+v", p.Path, views)
	}
	// RevokeAccessToken
	account, err := client.RevokeAccessToken(ctx)
	if err != nil {
		log.Fatalf("* RevokeAccessToken error: %s", err)
	}
	log.Printf("> RevokeAccessToken result: %#+v", account)
}
