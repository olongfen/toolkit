package main

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"strings"
)

const (
	IllegalAccessToken    = "IllegalAccessToken"
	IllegalCertificate    = "IllegalCertificate"
	IllegalParameter      = "IllegalParameter"
	RecordNotFound        = "RecordNotFound"
	AlreadyExists         = "AlreadyExists"
	SortParameterMismatch = "SortParameterMismatch"
)

func SetBundle(bundle *i18n.Bundle, translationDir string) {
	if len(translationDir) == 0 {
		return
	}
	files, err := os.ReadDir(translationDir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		if ext := filepath.Ext(filename); ext != ".json" {
			continue
		}
		lang := strings.TrimSuffix(filename, ".json")
		if lang == "" {
			continue
		}
		filePath := filepath.Join(translationDir, filename)
		bundle.L(filePath)
	}
}

func main() {
	bundle := i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	SetBundle(bundle, "./")
	localizer := i18n.NewLocalizer(bundle, language.English.String())

	fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    AlreadyExists,
			Other: "Hello World!",
		},
	}))
	// Outpu
}
