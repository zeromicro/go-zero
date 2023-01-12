package httpx

import (
	"errors"
	"net/http"
	"strings"

	"golang.org/x/text/language"

	"github.com/zeromicro/go-zero/core/logx"

	enLang "github.com/go-playground/locales/en"
	zhLang "github.com/go-playground/locales/zh_Hans"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

const xForwardedFor = "X-Forwarded-For"

var supportLang map[string]string

// GetFormValues returns the form values.
func GetFormValues(r *http.Request) (map[string]any, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		if err != http.ErrNotMultipart {
			return nil, err
		}
	}

	params := make(map[string]any, len(r.Form))
	for name := range r.Form {
		formValue := r.Form.Get(name)
		if len(formValue) > 0 {
			params[name] = formValue
		}
	}

	return params, nil
}

// GetRemoteAddr returns the peer address, supports X-Forward-For.
func GetRemoteAddr(r *http.Request) string {
	v := r.Header.Get(xForwardedFor)
	if len(v) > 0 {
		return v
	}

	return r.RemoteAddr
}

type Validator struct {
	Validator *validator.Validate
	Uni       *ut.UniversalTranslator
	Trans     map[string]ut.Translator
}

func NewValidator() *Validator {
	v := Validator{}
	en := enLang.New()
	zh := zhLang.New()
	v.Uni = ut.New(zh, en, zh)
	v.Validator = validator.New()
	enTrans, _ := v.Uni.GetTranslator("en")
	zhTrans, _ := v.Uni.GetTranslator("zh")
	v.Trans = make(map[string]ut.Translator)
	v.Trans["en"] = enTrans
	v.Trans["zh"] = zhTrans
	// add support languages
	initSupportLanguages()

	err := enTranslations.RegisterDefaultTranslations(v.Validator, enTrans)
	if err != nil {
		logx.Errorw("register English translation failed", logx.Field("detail", err.Error()))
		return nil
	}
	err = zhTranslations.RegisterDefaultTranslations(v.Validator, zhTrans)
	if err != nil {
		logx.Errorw("register Chinese translation failed", logx.Field("detail", err.Error()))

		return nil
	}

	return &v
}

func (v *Validator) Validate(data any, lang string) string {
	err := v.Validator.Struct(data)
	if err == nil {
		return ""
	}

	targetLang, parseErr := ParseAcceptLanguage(lang)
	if parseErr != nil {
		return parseErr.Error()
	}

	errs, ok := err.(validator.ValidationErrors)

	if ok {
		transData := errs.Translate(v.Trans[targetLang])
		s := strings.Builder{}
		for _, v := range transData {
			s.WriteString(v)
			s.WriteString(" ")
		}
		return s.String()
	}

	invalid, ok := err.(*validator.InvalidValidationError)
	if ok {
		return invalid.Error()
	}

	return ""
}

func ParseAcceptLanguage(lang string) (string, error) {
	tags, _, err := language.ParseAcceptLanguage(lang)
	if err != nil {
		return "", errors.New("fail to parse accept language")
	}

	for _, v := range tags {
		if val, ok := supportLang[v.String()]; ok {
			return val, nil
		}
	}

	return "zh", nil
}

func initSupportLanguages() {
	supportLang = make(map[string]string)
	supportLang["zh"] = "zh"
	supportLang["zh-CN"] = "zh"
	supportLang["en"] = "en"
	supportLang["en-US"] = "en"
}
