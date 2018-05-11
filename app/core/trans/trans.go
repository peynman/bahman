package trans

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"github.com/hjson/hjson-go"
	"avalanche/app/core/config"
	"peyman/config"
	"avalanche/app/core/logger"
	"github.com/sirupsen/logrus"
	"avalanche/app/core/app"
)

var bundle *i18n.Bundle
var localize *i18n.Localizer

func Initialize() {
	bundle = &i18n.Bundle {
		DefaultLanguage:language.English,
	}
	bundle.RegisterUnmarshalFunc("hjson", hjson.Unmarshal)

	locale := config.GetString("core.locale", "en")
	langFilesPath := app.ResourcesPath("lang/" + locale)
	transMessages := nemo.New(langFilesPath, "", nil)
	iterateMapForMessages(getLanguageTagFromString(locale), bundle, transMessages.ConfigsMap, "")

	localize = i18n.NewLocalizer(bundle, locale)
}

func L(key string) string {
	return LP(key, nil)
}

func LP(key string, params map[string]string) string {
	localized, err := localize.Localize(&i18n.LocalizeConfig{
		MessageID: key,
		TemplateData: params,
	})

	if err != nil {
		logger.ErrorFields("Localize failed", logrus.Fields {
			"error": err,
			"key": key,
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
			iterateMapForMessages(tag, bundle, inner, prefix + "." + key)
		} else {
			strVal, ok := value.(string)
			if ok {
				bundle.AddMessages(tag, &i18n.Message {
					ID: prefix + "." + key,
					Other: strVal,
				})
			}
		}
	}
}