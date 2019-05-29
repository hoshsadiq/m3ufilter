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
	if expr[0] == '=' {
		return strings.TrimSpace(expr[1:]), nil
	}

	debug := false
	if expr[0] == '?' {
		debug = true
		expr = strings.TrimSpace(expr[1:])
	}

	variables := map[string]interface{}{
		"Name":     ms.Name,
		"Uri":      ms.Uri,
		"Duration": ms.Duration,
		"Id":       ms.Id,
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
		"strlen":  evaluatorStrlen,
		"match":   evaluatorMatch,
		"replace": evaluatorReplace,
		"tvg_id":  evaluatorToTvgId,
		"title":   evaluatorTitle,
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
	return (bool)(re.MatchString(subject)), nil
}
func evaluatorReplace(args ...interface{}) (interface{}, error) {
	subject := args[0].(string)
	refind := args[1].(string)
	replace := args[2].(string)

	re := cache.Regexp(refind)
	return (string)(re.ReplaceAllString(subject, replace)), nil
}
func evaluatorToTvgId(args ...interface{}) (interface{}, error) {
	subject := args[0].(string)
	re := cache.Regexp(`(?i)\b(SD|HD|FHD)\b`)

	subject = re.ReplaceAllString(subject, "")

	subject = strings.Replace(subject, "&", "and", -1)
	subject = strings.TrimSpace(subject)

	re = cache.Regexp(`[^a-zA-Z0-9]`)
	tvgId := re.ReplaceAllString(subject, "")

	return tvgId, nil
}

func evaluatorTitle(args ...interface{}) (interface{}, error) {
	subject := args[0].(string)
	re := cache.Regexp(`(?i)\b(SD|HD|FHD)\b`)

	subject = re.ReplaceAllStringFunc(subject, func(s string) string {
		return strings.ToUpper(s)
	})

	subject = strings.Title(subject)
	return strings.TrimSpace(subject), nil
}
