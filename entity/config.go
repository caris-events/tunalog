package entity

type ColorScheme string

const (
	ColorSchemeLight ColorScheme = "light"
	ColorSchemeDark  ColorScheme = "dark"
	ColorSchemeAuto  ColorScheme = ""
)

type FontFamily string

const (
	FontFamilyNotoSans  FontFamily = "sans"
	FontFamilyNotoSerif FontFamily = "serif"
)

type AuthorBlock string

const (
	AuthorBlockNone  AuthorBlock = "none"
	AuthorBlockStart AuthorBlock = "start"
	AuthorBlockEnd   AuthorBlock = "end"
)

type Config struct {
	Name              string      `json:"name"`
	Description       string      `json:"description"`
	IsPublic          bool        `json:"is_public"`
	DateFormat        string      `json:"date_format"`
	TimeFormat        string      `json:"time_format"`
	Timezone          int         `json:"timezone"`
	InjectedHead      string      `json:"injected_head"`
	InjectedFoot      string      `json:"injected_foot"`
	InjectedPostStart string      `json:"injected_post_start"`
	InjectedPostEnd   string      `json:"injected_post_end"`
	FooterText        string      `json:"footer_text"`
	ColorScheme       ColorScheme `json:"color_scheme"`
	ContainerWidth    string      `json:"container_width"`
	FontFamily        FontFamily  `json:"font_family"`
	FontSize          string      `json:"font_size"`
	HighlightJS       bool        `json:"highlight_js"`
	AuthorBlock       AuthorBlock `json:"author_block"`
	PostsPerPage      int         `json:"posts_per_page"`
	Theme             string      `json:"theme"`
	Locale            string      `json:"locale"`
}

func (c *Config) IsCustomTimeFormat() bool {
	return c.TimeFormat != "PM 03:04" && c.TimeFormat != "15:04" && c.TimeFormat != "03:04 PM"
}

func (c *Config) IsCustomDateFormat() bool {
	return c.DateFormat != "2006-01-02" && c.DateFormat != "01/02/2006" && c.DateFormat != "02/01/2006"
}
