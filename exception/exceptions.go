package exception

//go:generate go install github.com/Beyond-simplechain/foundation/log
//go:generate go install github.com/Beyond-simplechain/foundation/exception
//go:generate gotemplate -outfmt "gen_%v" "github.com/Beyond-simplechain/foundation/exception/template" "StdException(Exception,StdExceptionCode,\"golang standard error\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/Beyond-simplechain/foundation/exception/template" "FcException(Exception,UnspecifiedExceptionCode,\"unspecified\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/Beyond-simplechain/foundation/exception/template" "UnHandledException(Exception,UnhandledExceptionCode,\"unhandled\")"
//go:generate gotemplate -outfmt "gen_%v" "github.com/Beyond-simplechain/foundation/exception/template" "AssertException(Exception,AssertExceptionCode,\"Assert Exception\")"

//go:generate go build .
