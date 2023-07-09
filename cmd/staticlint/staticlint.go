package main

import (
	"github.com/dimsonson/go-yp-shortener-url/cmd/staticlint/osexit"
	"github.com/kisielk/errcheck/errcheck"
	mnd "github.com/tommy-muehle/go-mnd/v2"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {

	// Определяем map подключаемых правил стандартных статических анализаторов пакета golang.org/x/tools/go/analysis/passes
	// и двух или более любых публичных анализаторов на ваш выбор.
	multichecks := map[*analysis.Analyzer]bool{
		asmdecl.Analyzer:             true,
		assign.Analyzer:              true,
		atomic.Analyzer:              true,
		atomicalign.Analyzer:         true,
		bools.Analyzer:               true,
		buildssa.Analyzer:            true,
		buildtag.Analyzer:            true,
		cgocall.Analyzer:             true,
		composite.Analyzer:           true,
		copylock.Analyzer:            true,
		ctrlflow.Analyzer:            true,
		deepequalerrors.Analyzer:     true,
		errorsas.Analyzer:            true,
		fieldalignment.Analyzer:      true,
		findcall.Analyzer:            true,
		framepointer.Analyzer:        true,
		httpresponse.Analyzer:        true,
		ifaceassert.Analyzer:         true,
		inspect.Analyzer:             true,
		loopclosure.Analyzer:         true,
		lostcancel.Analyzer:          true,
		nilfunc.Analyzer:             true,
		nilness.Analyzer:             true,
		pkgfact.Analyzer:             true,
		printf.Analyzer:              true,
		reflectvaluecompare.Analyzer: true,
		shadow.Analyzer:              true,
		shift.Analyzer:               true,
		sigchanyzer.Analyzer:         true,
		sortslice.Analyzer:           true,
		stdmethods.Analyzer:          true,
		stringintconv.Analyzer:       true,
		structtag.Analyzer:           true,
		testinggoroutine.Analyzer:    true,
		tests.Analyzer:               true,
		timeformat.Analyzer:          true,
		unmarshal.Analyzer:           true,
		unreachable.Analyzer:         true,
		unsafeptr.Analyzer:           true,
		unusedresult.Analyzer:        true,
		unusedwrite.Analyzer:         true,
		usesgenerics.Analyzer:        true,
		errcheck.Analyzer:            true,
		mnd.Analyzer:                 true,
		osexit.OSexitCheckAnalyzer:   true,
	}
	// Массив анализаторов.
	var mychecks []*analysis.Analyzer

	// Добавление в массив подключаемых правил всех анализаторов класса SA пакета staticcheck.io.
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}
	// Добавление в массив подключаемых правил не менее одного анализатора остальных классов пакета staticcheck.io.
	for _, v := range stylecheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}
	for _, v := range simple.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}
	// Добавление в массив подключаемых правил стандартных статических анализаторов пакета golang.org/x/tools/go/analysis/passes
	// и двух или более любых публичных анализаторов на ваш выбор
	for k, ok := range multichecks {
		if ok {
			mychecks = append(mychecks, k)
		}
	}
	// Подключение правил.
	multichecker.Main(
		mychecks...,
	)
}
