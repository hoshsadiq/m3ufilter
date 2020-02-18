package m3u

import (
	"errors"
	"fmt"
	"github.com/hoshsadiq/m3ufilter/cache"
	"github.com/maja42/goval"
	"strings"
)

var evaluator = goval.NewEvaluator()

func evaluate(ms *Stream, expr string) (result interface{}, err error) {
	var debug bool
	if expr[0] == '?' {
		debug = true
		expr = strings.TrimSpace(expr[1:])
	}

	if expr[0] == '=' {
		return strings.TrimSpace(expr[1:]), nil
	}

	variables := map[string]interface{}{
		"ChNo":     ms.ChNo,
		"Id":       ms.Id,
		"Name":     ms.Name,
		"Uri":      ms.Uri,
		"Duration": ms.Duration,
		"Logo":     ms.Logo,
		"Group":    ms.Group,
	}

	expr = cache.Expr(expr)

	//fmt.Printf("Evaluating `%s` using vars %v\n", expr, variables)
	res, err := evaluator.Evaluate(expr, variables, getEvaluatorFunctions())
	if debug {
		log.Infof("Debugging expr %s, res = %s, vars = %v", expr, res, variables)
	}
	return res, err
}

func evaluateBool(ms *Stream, expr string) (result bool, err error) {
	res, err := evaluate(ms, expr)
	if err != nil {
		return false, err
	}

	switch v := res.(type) {
	case bool:
		//fmt.Printf("Result: %v", v)
		return v, nil
	default:
		return false, errors.New(fmt.Sprintf("unexpected type %T, expected bool for expr: %s", v, expr))
	}
}

func evaluateStr(ms *Stream, expr string) (result string, err error) {
	res, err := evaluate(ms, expr)
	if err != nil {
		return "", err
	}

	switch v := res.(type) {
	case string:
		return v, nil
	default:
		return "", errors.New(fmt.Sprintf("unexpected type %T, expected string", v))
	}
}

func getEvaluatorFunctions() map[string]goval.ExpressionFunction {
	return map[string]goval.ExpressionFunction{
		"strlen":      evaluatorStrlen,
		"match":       evaluatorMatch,
		"replace":     evaluatorReplace,
		"tvg_id":      evaluatorToTvgId,
		"title":       evaluatorTitle,
		"upper_words": evaluatorUpperWord,
	}
}

func evaluatorStrlen(args ...interface{}) (interface{}, error) {
	length := len(args[0].(string))
	return (float64)(length), nil
}

func evaluatorMatch(args ...interface{}) (interface{}, error) {
	subject := args[0].(string)
	regexString := args[1].(string)

	re := cache.Regexp(regexString)
	return re.MatchString(subject), nil
}
func evaluatorReplace(args ...interface{}) (interface{}, error) {
	subject := args[0].(string)
	refind := args[1].(string)
	replace := args[2].(string)

	re := cache.Regexp(refind)
	return re.ReplaceAllString(subject, replace), nil
}
func evaluatorToTvgId(args ...interface{}) (interface{}, error) {
	subject := args[0].(string)
	subject = regexWordCallback(subject, definitions, removeWord)

	subject = strings.Replace(subject, "&", "and", -1)
	subject = strings.TrimSpace(subject)

	re := cache.Regexp(`[^a-zA-Z0-9]`)
	tvgId := re.ReplaceAllString(subject, "")

	return tvgId, nil
}

func evaluatorTitle(args ...interface{}) (interface{}, error) {
	subject := strings.ToLower(args[0].(string))

	subject = regexWordCallback(subject, definitions, removeWord)
	subject = regexWordCallback(subject, countries, removeWord)

	subject = strings.Title(subject)
	return strings.TrimSpace(subject), nil
}

func evaluatorUpperWord(args ...interface{}) (interface{}, error) {
	subject := args[0].(string)

	sargs := make([]string, len(args))
	for i := range args {
		sargs[i] = args[i].(string)
	}

	subject = regexWordCallback(subject, strings.Join(sargs[1:], "|"), strings.ToUpper)
	return strings.TrimSpace(subject), nil
}

func regexWordCallback(subject string, word string, callback func(string) string) string {
	re := cache.Regexp(`(?i)\b(` + word + `)\b`)

	subject = re.ReplaceAllStringFunc(subject, func(s string) string {
		return callback(s)
	})

	return subject
}

func removeWord(s string) string {
	return ""
}
