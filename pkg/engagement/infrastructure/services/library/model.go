package library

import "time"

// GhostCMSPost is the body of the post sourced from Ghost CMS
type GhostCMSPost struct {
	ID            string           `json:"id"`
	UUID          string           `json:"uuid"`
	Slug          string           `json:"slug"`
	Title         string           `json:"title"`
	HTML          string           `json:"html"`
	Excerpt       string           `json:"excerpt"`
	URL           string           `json:"url"`
	FeatureImage  string           `json:"feature_image"`
	Featured      bool             `json:"featured"`
	Visibility    string           `json:"visibility"`
	ReadingTime   int              `json:"reading_time"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	PublishedAt   time.Time        `json:"published_at"`
	CommentID     string           `json:"comment_id"`
	Tags          []GhostCMSTag    `json:"tags"`
	Authors       []GhostCMSAuthor `json:"authors"`
	PrimaryAuthor GhostCMSAuthor   `json:"primary_author"`
	PrimaryTag    GhostCMSTag      `json:"primary_tag"`
}

// IsEntity marks a Ghost CMS post as an entity
func (g GhostCMSPost) IsEntity() {}

// GhostCMSTag represents the structure of a Ghost CMS tag
type GhostCMSTag struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	Visibility  string  `json:"visibility"`
	URL         string  `json:"url"`
}

// IsEntity marks a Ghost CMS tag as an entity
func (g GhostCMSTag) IsEntity() {}

// GhostCMSAuthor is used to serialize authors of Ghost CMS posts
type GhostCMSAuthor struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ProfileImage string `json:"profile_image"`
	Website      string `json:"website"`
	Location     string `json:"location"`
	Facebook     string `json:"facebook"`
	Twitter      string `json:"twitter"`
	URL          string `json:"url"`
}

// GhostCMSServerResponse assembles the posts fetched from a Ghost server, for
// serialization
type GhostCMSServerResponse struct {
	Posts []*GhostCMSPost `json:"posts"`
}
