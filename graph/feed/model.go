package feed

import "time"

// GhostCMSPost is the body of the post sourced from Ghost CMs. While the output from the API call is
// quite detaied, we will only cherr-pick what we need.
type GhostCMSPost struct {
	ID           *string       `json:"id,omitempty"`
	CreatedAt    time.Time     `json:"created_at,omitempty"`
	Excerpt      string        `json:"excerpt,omitempty"`
	FeatureImage string        `json:"feature_image,omitempty"`
	HTML         string        `json:"html,omitempty"`
	PublishedAt  time.Time     `json:"published_at,omitempty"`
	Slug         string        `json:"slug,omitempty"`
	Title        string        `json:"title,omitempty"`
	URL          string        `json:"url,omitempty"`
	ReadingTime  string        `json:"reading_time,omitempty"`
	Tags         []GhostCMSTag `json:"tags,omitempty"`
}

// IsEntity ...
func (g GhostCMSPost) IsEntity() {}

// GhostCMSTag represemt the structure of a tag. We cherry-pick only what we need
type GhostCMSTag struct {
	ID   *string `json:"id,omitempty"`
	Name string  `json:"name,omitempty"`
	Slug string  `json:"slug,omitempty"`
	URL  string  `json:"url,omitempty"`
}

// IsEntity ...
func (g GhostCMSTag) IsEntity() {}
