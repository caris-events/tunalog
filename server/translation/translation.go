package translation

import (
	"fmt"
	"html/template"
	"log"

	"github.com/teacat/i18n"
)

var Instance *Translation

type Translation struct {
	l *i18n.Locale
}

func init() {
	i := i18n.New("zh-tw")
	if err := i.LoadGlob("translation/*.json"); err != nil {
		log.Fatalln("load glob:", err)
	}
	Instance = &Translation{
		l: i.NewLocale("zh-tw"),
	}
}

func Str(v string) template.HTML {
	return template.HTML(Instance.l.String(v))
}

func Strf(v string, args ...interface{}) template.HTML {
	return template.HTML(fmt.Sprintf(Instance.l.String(v), args...))
}
