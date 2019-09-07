package m3u

import (
	"errors"
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/hoshsadiq/m3ufilter/cache"
	"strings"
)

func evaluate(ms *Stream, expression string) (result interface{}, err error) {
	if expression[0] == '=' {
		return strings.TrimSpace(expression[1:]), nil
	}

	variables := map[string]interface{}{
		"ChNo":     ms.ChNo,
		"Id":       ms.Id,
		"Name":     ms.Name,
		"Uri":      ms.Uri,
		"Duration": ms.Duration,
		"Logo":     ms.Logo,
		"Group":    ms.Group,

		"strlen":  evaluatorStrlen,
		"match":   evaluatorMatch,
		"replace": evaluatorReplace,
		"tvg_id":  evaluatorToTvgId,
		"title":   evaluatorTitle,
	}

	program, err := cache.Expr(expression, variables)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, variables)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}


	//fmt.Printf("Evaluating `%s` using vars %v\n", expr, variables)
	//res, err := evaluator.Evaluate(expr, variables, getEvaluatorFunctions())
	log.Debugf("Debugging expr %s, res = %s, vars = %v", expression, output, variables)
	return output, err
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

func evaluatorStrlen(str string) int {
	return len(str)
}

func evaluatorMatch(subject string, regexString string) bool {
	re := cache.Regexp(regexString)
	return re.MatchString(subject)
}
func evaluatorReplace(subject string, refind string, replace string) string {
	re := cache.Regexp(refind)
	return re.ReplaceAllString(subject, replace)
}
func evaluatorToTvgId(subject string) string {
	re := cache.Regexp(`(?i)\b(SD|HD|FHD)\b`)

	subject = re.ReplaceAllString(subject, "")

	subject = strings.Replace(subject, "&", "and", -1)
	subject = strings.TrimSpace(subject)

	re = cache.Regexp(`[^a-zA-Z0-9]`)
	tvgId := re.ReplaceAllString(subject, "")

	return tvgId
}

func evaluatorTitle(subject string) string {
	re := cache.Regexp(`(?i)\b(SD|HD|FHD)\b`)

	subject = re.ReplaceAllStringFunc(subject, func(s string) string {
		return strings.ToUpper(s)
	})

	subject = strings.Title(subject)
	return strings.TrimSpace(subject)
}
