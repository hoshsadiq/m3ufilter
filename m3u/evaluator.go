package m3u

import (
	"errors"
	"fmt"
	"github.com/grafov/m3u8"
	"github.com/hoshsadiq/m3ufilter/regex"
	"github.com/maja42/goval"
	"github.com/maja42/no-comment"
	"strings"
)

var evaluator = goval.NewEvaluator()

func evaluate(ms *m3u8.MediaSegment, expr string) (result interface{}, err error) {
	if expr[0] == '=' {
		return strings.TrimSpace(expr[1:]), nil
	}

	attrs := make(map[string]string)
	for _, attr := range ms.Attributes {
		attrs[attr.Key] = attr.Value
	}

	variables := map[string]interface{}{
		"Name":     ms.Title,
		"Uri":      ms.URI,
		"Duration": ms.Duration,
		"Attr":     attrs,
	}

	expr = nocomment.StripCStyleComments(expr)

	//fmt.Printf("Evaluating `%s` using vars %v\n", expr, variables)
	return evaluator.Evaluate(expr, variables, getEvaluatorFunctions())
}

func evaluateBool(ms *m3u8.MediaSegment, expr string) (result bool, err error) {
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

func evaluateStr(ms *m3u8.MediaSegment, expr string) (result string, err error) {
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

	re := regex.GetCache(regexString)
	return (bool)(re.MatchString(subject)), nil
}
func evaluatorReplace(args ...interface{}) (interface{}, error) {
	subject := args[0].(string)
	refind := args[1].(string)
	replace := args[2].(string)

	re := regex.GetCache(refind)
	return (string)(re.ReplaceAllString(subject, replace)), nil
}
func evaluatorToTvgId(args ...interface{}) (interface{}, error) {
	subject := args[0].(string)
	re := regex.GetCache(`(?i)\b(SD|HD|FHD)\b`)

	subject = re.ReplaceAllString(subject, "")

	subject = strings.Replace(subject, "&", "and", -1)
	subject = strings.TrimSpace(subject)

	re = regex.GetCache(`[^a-zA-Z0-9]`)
	tvgId := re.ReplaceAllString(subject, "")

	return tvgId, nil
}

func evaluatorTitle(args ...interface{}) (interface{}, error) {
	subject := args[0].(string)
	re := regex.GetCache(`(?i)\b(SD|HD|FHD)\b`)

	subject = re.ReplaceAllStringFunc(subject, func(s string) string {
		return strings.ToUpper(s)
	})

	subject = strings.Title(subject)
	return strings.TrimSpace(subject), nil
}
