package trans

import (
	"github.com/hjson/hjson-go"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/peyman-abdi/conf"
	"golang.org/x/text/language"
)

type localImpl struct {
	bundle   *i18n.Bundle
	localize *i18n.Localizer
	logger   services.Logger
}
func (t *localImpl) L(key string) string {
	return t.LP(key, nil)
}

func (t *localImpl) LP(key string, params map[string]string) string {
	localized, err := t.localize.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: params,
	})

	if err != nil {
		t.logger.ErrorFields("Localize failed", map[string]interface{}{
			"error":  err,
			"key":    key,
			"params": params,
		})
		return ""
	}
	return localized
}

func getLanguageTagFromString(tag string) language.Tag {
	switch tag {
	case "fa":
	case "persian":
		return language.Persian
	}
	return language.English
}
func iterateMapForMessages(tag language.Tag, bundle *i18n.Bundle, m map[string]interface{}, prefix string) {
	for key, value := range m {
		inner, ok := value.(map[string]interface{})
		if ok {
			iterateMapForMessages(tag, bundle, inner, prefix+"."+key)
		} else {
			strVal, ok := value.(string)
			if ok {
				bundle.AddMessages(tag, &i18n.Message{
					ID:    prefix + "." + key,
					Other: strVal,
				})
			}
		}
	}
}


func New(config services.Config, app services.Application, logger services.Logger) services.Localization {
	trans := new(localImpl)
	trans.bundle = &i18n.Bundle{
		DefaultLanguage: language.English,
	}
	trans.bundle.RegisterUnmarshalFunc("hjson", hjson.Unmarshal)
	trans.logger = logger

	locale := config.GetString("services.locale", "en")
	langFilesPath := app.ResourcesPath("lang/" + locale)
	transMessages, _ := conf.New(langFilesPath, "", nil)
	iterateMapForMessages(getLanguageTagFromString(locale), trans.bundle, transMessages.ConfigsMap, "")

	trans.localize = i18n.NewLocalizer(trans.bundle, locale)
	return trans
}
