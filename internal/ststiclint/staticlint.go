package ststiclint

import (
	"strings"

	"github.com/gnieto/mulint/mulint"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/staticcheck"
)

func Analyzer() []*analysis.Analyzer {

	// Стандартные статические анализаторы пакета golang.org/x/tools/go/analysis/passes
	checks := standardStatic()

	// Добавляем все анализаторов класса SA пакета staticCheck
	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix("SA", v.Analyzer.Name) {
			checks = append(checks, v.Analyzer)
		}
	}

	// Добавляем проверки Code simplifications
	// Подробнее: https://staticcheck.io/docs/checks#S
	sChecks := map[string]bool{
		"S1003": true, // Replace call to strings.Index with strings.Contains
		"S1004": true, // Replace call to bytes.Compare with bytes.Equal
	}
	for _, v := range staticcheck.Analyzers {
		if sChecks[v.Analyzer.Name] {
			checks = append(checks, v.Analyzer)
		}
	}

	// Добавляем анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакета main
	checks = append(checks, osExitAnalyzer)

	// Добавление публичных анализаторов
	checks = append(checks, errcheck.Analyzer)
	checks = append(checks, mulint.Mulint)

	return checks
}
