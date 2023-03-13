package utils

import (
	"context"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translation "github.com/go-playground/validator/v10/translations/en"
	zh_translation "github.com/go-playground/validator/v10/translations/zh"
	"github.com/go-playground/validator/v10/translations/zh_tw"
	"github.com/olongfen/toolkit/multi/xerror"
	"github.com/olongfen/toolkit/scontext"
	"strings"
)

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()

}

// ValidateForm validate
func ValidateForm(ctx context.Context, form interface{}) error {
	var (
		errs = xerror.ValidateError{}
	)
	language := scontext.GetLanguage(ctx)
	err := validate.Struct(form)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			//_err := &scontext.Error{}
			//_err.Detail = e.Translate(translate(language))
			//_err.Title = e.Title()
			//_err.Failed = e.Field()
			//_err.Value = e.Value()
			errs[strings.ToLower(e.Field()[:1])+e.Field()[1:]] = e.Translate(translate(language))
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func translate(language string) ut.Translator {
	var trans ut.Translator
	switch language {
	case "en":
		defaultEn := en.New()
		uni := ut.New(defaultEn, defaultEn)
		trans, _ = uni.GetTranslator(language)
		_ = en_translation.RegisterDefaultTranslations(validate, trans)
	case "zh-tw":
		defaultZhTw := zh_Hant_TW.New()
		uni := ut.New(defaultZhTw, defaultZhTw)
		trans, _ = uni.GetTranslator(language)
		_ = zh_tw.RegisterDefaultTranslations(validate, trans)
	default:
		defaultZh := zh.New()
		uni := ut.New(defaultZh, defaultZh)
		trans, _ = uni.GetTranslator(language)
		_ = zh_translation.RegisterDefaultTranslations(validate, trans)
	}
	return trans
}
