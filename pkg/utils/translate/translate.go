package translate

import (
	custom_log "admin-phone-shop-api/pkg/custom_log"
	errors "admin-phone-shop-api/pkg/utils/error"
	"fmt"
	"log"
	"path/filepath"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

var bundle *i18n.Bundle

func Init() *errors.ErrorResponse {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	localeFiles := []string{
		"pkg/translates/localize/i18n/en.yaml",
		"pkg/translates/localize/i18n/km.yaml",
	}

	for _, file := range localeFiles {
		_, err := bundle.LoadMessageFile(filepath.Join(file))
		if err != nil {
			log.Printf("error loading local file %s: %v", file, err)
			custom_log.NewCustomLog("translate_error", err.Error(), "error")
			return &errors.ErrorResponse{
				MessageID: "ErrorLoadMessage",
				Err:       err,
			}
		}
	}
	return nil
}

func TranslateWithError(c *fiber.Ctx, key string, templateData ...map[string]interface{}) (string, *errors.ErrorResponse) {
	if bundle == nil {
		custom_log.NewCustomLog("I18nNotInit", Init().ErrorString(), "error")
		return "", &errors.ErrorResponse{
			MessageID: key,
			Err:       fmt.Errorf("translation service is unavailable for MessageID: %s", key),
		}
	}

	lang := c.Get("Accept-Language", "en")
	localizer := i18n.NewLocalizer(bundle, lang)

	data := map[string]interface{}{}
	if len(templateData) > 0 && templateData[0] != nil {
		data = templateData[0]
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	})
	if err != nil {
		log.Printf("error localizing message ID %s: %v", key, err)
		custom_log.NewCustomLog("TranslationNotFound", err.Error(), "error")
		return "", &errors.ErrorResponse{
			MessageID: key,
			Err:       fmt.Errorf("translation not found for MessageID: %s", key),
		}
	}
	return msg, nil
}


func Translate(MessageID string, param *string, c *fiber.Ctx) string {
	var translate string
	if param != nil {
		translate = fiberi18n.MustLocalize(c, &i18n.LocalizeConfig{
			MessageID: MessageID,
			TemplateData: map[string]interface{}{
				"name": param,
			},
		})
	} else {
		translate = fiberi18n.MustLocalize(c, &i18n.LocalizeConfig{
			MessageID: MessageID,
		})
	}
	return translate

}
