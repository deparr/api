package model

type Language struct {
	Name    string `json:"name,omitempty"`
	Color   string `json:"color,omitempty"`
	Percent int    `json:"percent,omitempty"`
}

type RepoLangs []Language

type Repository struct {
	Name     string    `json:"name,omitempty"`
	Owner    string    `json:"owner,omitempty"`
	Desc     string    `json:"desc,omitempty"`
	Url      string    `json:"url,omitempty"`
	PushedAt string    `json:"pushed_at,omitempty"`
	Language RepoLangs `json:"language,omitempty"`
	Stars    int       `json:"stars,omitempty"`
}
