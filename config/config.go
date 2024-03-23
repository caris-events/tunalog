package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

const Version = "1.0.0"

var Instance *Config

type ColorScheme string

const (
	ColorSchemeLight ColorScheme = "light"
	ColorSchemeDark  ColorScheme = "dark"
	ColorSchemeAuto  ColorScheme = "auto"
)

type FontFamily string

const (
	FontFamilyNotoSans  FontFamily = "noto_sans"
	FontFamilyNotoSerif FontFamily = "noto_serif"
)

type AuthorBlock string

const (
	AuthorBlockNone  AuthorBlock = "none"
	AuthorBlockStart AuthorBlock = "start"
	AuthorBlockEnd   AuthorBlock = "end"
)

type Config struct {
	Name              string      `json:"name"`
	URL               string      `json:"url"`
	Description       string      `json:"description"`
	IsPublic          bool        `json:"is_public"`
	DateFormat        string      `json:"date_format"`
	TimeFormat        string      `json:"time_format"`
	Timezone          int         `json:"timezone"`
	InjectedHead      string      `json:"injected_head"`
	InjectedFoot      string      `json:"injected_foot"`
	InjectedPostStart string      `json:"injected_post_start"`
	InjectedPostEnd   string      `json:"injected_post_end"`
	IsPoweredByShown  bool        `json:"is_powered_by_shown"`
	IsCopyrightShown  bool        `json:"is_copyright_shown"`
	ColorScheme       ColorScheme `json:"color_scheme"`
	ContainerWidth    string      `json:"container_width"`
	FontFamily        FontFamily  `json:"font_family"`
	FontSize          string      `json:"font_size"`
	HighlightJS       bool        `json:"highlight_js"`
	AuthorBlock       AuthorBlock `json:"author_block"`
	PostsPerPage      int         `json:"posts_per_page"`
	AuthSecret        string      `json:"auth_secret"`
	Theme             string      `json:"theme"`
}

func init() {
	if _, err := os.Stat("config.json"); errors.Is(err, os.ErrNotExist) {
		return // non-exist should trigger wizard, not an error
	}
	b, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalln("read file:", err)
	}
	if err := json.Unmarshal(b, &Instance); err != nil {
		log.Fatalln("unmarshal:", err)
	}
}

func (c *Config) Save() error {
	b, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	if err := os.WriteFile("config.json", b, 0644); err != nil {
		return err
	}
	Instance = c
	return nil
}
