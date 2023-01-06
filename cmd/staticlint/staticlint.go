package main

import (
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
	"honnef.co/go/tools/staticcheck"
)

type Analyzer struct {
	//Doc      *Documentation
	Analyzer *analysis.Analyzer
}

func main() {
	// Определяем map подключаемых правил всех анализаторов класса SA пакета staticcheck.io не менее одного анализатора
	// остальных классов пакета staticcheck.io.
	checks := map[string]bool{
		"SA1000": true,
		"SA1001": true,
		"SA1002": true,
		"SA1003": true,
		"SA1004": true,
		"SA1005": true,
		"SA1006": true,
		"SA1007": true,
		"SA1008": true,
		"SA1010": true,
		"SA1011": true,
		"SA1012": true,
		"SA1013": true,
		"SA1014": true,
		"SA1015": true,
		"SA1016": true,
		"SA1017": true,
		"SA1018": true,
		"SA1019": true,
		"SA1020": true,
		"SA1021": true,
		"SA1023": true,
		"SA1024": true,
		"SA1025": true,
		"SA1026": true,
		"SA1027": true,
		"SA1028": true,
		"SA1029": true,
		"SA1030": true,
		"SA2000": true,
		"SA2001": true,
		"SA2002": true,
		"SA2003": true,
		"SA3000": true,
		"SA3001": true,
		"SA4000": true,
		"SA4001": true,
		"SA4003": true,
		"SA4004": true,
		"SA4005": true,
		"SA4006": true,
		"SA4008": true,
		"SA4009": true,
		"SA4010": true,
		"SA4011": true,
		"SA4012": true,
		"SA4013": true,
		"SA4014": true,
		"SA4015": true,
		"SA4016": true,
		"SA4017": true,
		"SA4018": true,
		"SA4019": true,
		"SA4020": true,
		"SA4021": true,
		"SA4022": true,
		"SA4023": true,
		"SA4024": true,
		"SA4025": true,
		"SA4026": true,
		"SA4027": true,
		"SA4028": true,
		"SA4029": true,
		"SA4030": true,
		"SA4031": true,
		"SA5000": true,
		"SA5001": true,
		"SA5002": true,
		"SA5003": true,
		"SA5004": true,
		"SA5005": true,
		"SA5007": true,
		"SA5008": true,
		"SA5009": true,
		"SA5010": true,
		"SA5011": true,
		"SA5012": true,
		"SA6000": true,
		"SA6001": true,
		"SA6002": true,
		"SA6003": true,
		"SA6005": true,
		"SA9001": true,
		"SA9002": true,
		"SA9003": true,
		"SA9004": true,
		"SA9005": true,
		"SA9006": true,
		"SA9007": true,
		"SA9008": true,
		"S1000":  true,
		"S1001":  true,
		"S1002":  true,
		"S1003":  true,
		"S1004":  true,
		"S1005":  true,
		"S1006":  true,
		"S1007":  true,
		"S1008":  true,
		"S1009":  true,
		"S1010":  true,
		"S1011":  true,
		"S1012":  true,
		"S1016":  true,
		"S1017":  true,
		"S1018":  true,
		"S1019":  true,
		"S1020":  true,
		"S1021":  true,
		"S1023":  true,
		"S1024":  true,
		"S1025":  true,
		"S1028":  true,
		"S1029":  true,
		"S1030":  true,
		"S1031":  true,
		"S1032":  true,
		"S1033":  true,
		"S1034":  true,
		"S1035":  true,
		"S1036":  true,
		"S1037":  true,
		"S1038":  true,
		"S1039":  true,
		"S1040":  true,
		"ST1000": true,
		"ST1001": true,
		"ST1003": true,
		"ST1005": true,
		"ST1006": true,
		"ST1008": true,
		"ST1011": true,
		"ST1012": true,
		"ST1013": true,
		"ST1015": true,
		"ST1016": true,
		"ST1017": true,
		"ST1018": true,
		"ST1019": true,
		"ST1020": true,
		"ST1021": true,
		"ST1022": true,
		"ST1023": true,
	}
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
	}

	var mychecks []*analysis.Analyzer

	for _, v := range staticcheck.Analyzers {
		// добавляем в массив нужные проверки
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	for k := range multichecks {
		// добавляем в массив нужные проверки
		mychecks = append(mychecks, k)
	}

	multichecker.Main(
		mychecks...,
	)
}
