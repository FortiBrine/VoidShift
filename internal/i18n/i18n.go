package i18n

import (
	"embed"
	"encoding/json"
	"fmt"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFS embed.FS

type Localizer = goi18n.Localizer

var bundle = mustLoadBundle()

func mustLoadBundle() *goi18n.Bundle {
	b := goi18n.NewBundle(language.Ukrainian)
	b.RegisterUnmarshalFunc("json", json.Unmarshal)

	for _, path := range []string{"locales/uk.json", "locales/en.json"} {
		if _, err := b.LoadMessageFileFS(localeFS, path); err != nil {
			panic(fmt.Errorf("i18n: loading %s: %w", path, err))
		}
	}

	return b
}

func FromAcceptLanguage(header string) *Localizer {
	return goi18n.NewLocalizer(bundle, header)
}

func T(l *Localizer, messageID string, data ...map[string]any) string {
	cfg := &goi18n.LocalizeConfig{MessageID: messageID}
	if len(data) > 0 {
		cfg.TemplateData = data[0]
	}

	msg, err := l.Localize(cfg)
	if err != nil {
		return messageID
	}
	return msg
}
