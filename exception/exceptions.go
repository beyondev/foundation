package exception

//go:generate go install foundation/log/...
//go:generate go install foundation/exception/...
//go:generate gotemplate -outfmt "gen_%v" "foundation/exception/template" "StdException(Exception,StdExceptionCode,\"golang standard error\")"
//go:generate gotemplate -outfmt "gen_%v" "foundation/exception/template" "FcException(Exception,UnspecifiedExceptionCode,\"unspecified\")"
//go:generate gotemplate -outfmt "gen_%v" "foundation/exception/template" "UnHandledException(Exception,UnhandledExceptionCode,\"unhandled\")"
//go:generate gotemplate -outfmt "gen_%v" "foundation/exception/template" "AssertException(Exception,AssertExceptionCode,\"Assert Exception\")"

//go:generate go build .
