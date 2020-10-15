package model

// FaqGraphqlResponse represents an Faq Item Query graphql response
type FaqGraphqlResponse struct {
	Data struct {
		Faq []Faq `json:"faqs"`
	} `json:"data"`
}

// LibraryItemGraphqlResponse represents a Library Item Query graphql response
type LibraryItemGraphqlResponse struct {
	Data struct {
		LibraryItem []LibraryItem `json:"libraryItems"`
	} `json:"data"`
}

// FeedItemGraphqlResponse represents a feed Item's graphql response
type FeedItemGraphqlResponse struct {
	Data struct {
		FeedItem []FeedItem `json:"feedItems"`
	} `json:"data"`
}
