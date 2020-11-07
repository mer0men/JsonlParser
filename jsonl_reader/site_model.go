package jsonl_reader

type Site struct {
	Url             string   `json:"url"`
	State           string   `json:"state"`
	Categories      []string `json:"categories"`
	CategoryAnother string   `json:"category_another"`
	ForMainPage     bool     `json:"for_main_page"`
	Ctime           int64    `json:"ctime"`
}
