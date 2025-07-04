package system

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"slices"
	"time"

	"github.com/caris-events/tunalog/entity"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/teacat/i18n"
)

const Version = "1.0.2"

var (
	Config *entity.Config

	localeBase *i18n.I18n
	Locale     *i18n.Locale

	themeLocaleBase *i18n.I18n
	themeLocale     *i18n.Locale

	IndexTmpl    *template.Template
	SingularTmpl *template.Template
	NotFoundTmpl *template.Template

	//go:embed locales
	LocalesFS embed.FS
	//go:embed themes/default
	ThemesFS embed.FS

	funcs = template.FuncMap{
		"add": func(x, y int) int {
			return x + y
		},
		"sub": func(x, y int) int {
			return x - y
		},
		"html": func(v string) template.HTML {
			return template.HTML(v)
		},
		"unix2date": func(v int64) string {
			return time.Unix(v, 0).Format(Config.DateFormat)
		},
		"markdown": func(v string) template.HTML {
			p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Footnotes | parser.SuperSubscript | parser.MathJax)
			doc := p.Parse([]byte(v))
			renderer := html.NewRenderer(html.RendererOptions{
				Flags: html.HrefTargetBlank,
			})
			return template.HTML(markdown.Render(doc, renderer))
		},
		"__": func(v string) template.HTML {
			return template.HTML(themeLocale.String(v))
		},
		"_f": func(v string, data ...any) string {
			return fmt.Sprintf(themeLocale.String(v), data...)
		},
	}
)

func init() {
	if err := os.MkdirAll("data/uploads/images", 0755); err != nil {
		log.Fatalln(err)
	}
	if err := os.MkdirAll("data/uploads/covers", 0755); err != nil {
		log.Fatalln(err)
	}
	if err := os.MkdirAll("data/themes", 0755); err != nil {
		log.Fatalln(err)
	}
	if _, err := os.Stat("data/themes/default"); os.IsNotExist(err) {
		f, err := fs.Sub(ThemesFS, "themes")
		if err != nil {
			log.Fatalln(err)
		}
		if err := os.CopyFS("data/themes", f); err != nil {
			log.Fatalln(err)
		}
	}

	// init locale
	localeBase = i18n.New("en_US")
	localeBase.LoadFS(LocalesFS, "locales/*.json")

	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		return // non-exist triggers wizard
	}
	b, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	if err := json.Unmarshal(b, &Config); err != nil {
		log.Fatalln(err)
	}

	if Config != nil {
		if err := SaveConfig(); err != nil {
			log.Fatalln(err)
		}
	}
}

func SaveConfig() error {
	b, err := json.MarshalIndent(Config, "", "    ")
	if err != nil {
		return err
	}
	if err := os.WriteFile("config.json", b, 0644); err != nil {
		return err
	}

	IndexTmpl, err = template.
		New("index.html").
		Funcs(funcs).ParseGlob(fmt.Sprintf("data/themes/%s/*.html", Config.Theme))
	if err != nil {
		return err
	}
	SingularTmpl, err = template.
		New("singular.html").
		Funcs(funcs).ParseGlob(fmt.Sprintf("data/themes/%s/*.html", Config.Theme))
	if err != nil {
		return err
	}
	// load 404.html optionally
	NotFoundTmpl = nil
	if _, err := os.Stat(fmt.Sprintf("data/themes/%s/404.html", Config.Theme)); !os.IsNotExist(err) {
		NotFoundTmpl, err = template.
			New("404.html").
			Funcs(funcs).ParseGlob(fmt.Sprintf("data/themes/%s/*.html", Config.Theme))
		if err != nil {
			return err
		}
	}
	// load theme locales, or skip if not exists
	themeLocaleBase = i18n.New("default")
	if _, err := os.Stat(fmt.Sprintf("data/themes/%s/locales", Config.Theme)); !os.IsNotExist(err) {
		themeLocaleBase.LoadGlob(fmt.Sprintf("data/themes/%s/locales/*.json", Config.Theme))
		themeLocale = themeLocaleBase.NewLocale(Config.Locale)
	}

	ReloadLocale(Config.Locale)
	return nil
}

// ===============================
// Locale
// ===============================

func ReloadLocale(v ...string) {
	Locale = localeBase.NewLocale(v...)
}

// ===============================
// Themes
// ===============================

func Themes() (themes []string) {
	files, err := os.ReadDir("data/themes")
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			themes = append(themes, file.Name())
		}
	}
	return
}

func ThemeExists(v string) bool {
	return slices.Index(Themes(), v) != -1
}
