// Package main предназначен для анализа кода.
// Применяются следующие анализаторы:
//  1. Стандартные статических анализаторы пакета golang.org/x/tools/go/analysis/passes
//  2. Все анализаторы класса SA пакета staticcheck.io;
//  3. Анализаторы пакета staticcheck.io:
//     3.1 S1003 Replace call to strings.Index with strings.Contains
//     3.2 "S1004  Replace call to bytes.Compare with bytes.Equal
//  4. Проверка необработанных ошибок github.com/kisielk/errcheck
//  5. Проверка потенциальных дедлоков github.com/gnieto/mulint/mulint
//  6. Проверка на использование os.Exit в функции mai.main
//
// Пример: staticlint ./...
package main

import (
	"goAdvancedTpl/internal/ststiclint"

	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	// определяем map подключаемых правил
	myChecks := ststiclint.Analyzer()

	multichecker.Main(
		myChecks...,
	)
}
