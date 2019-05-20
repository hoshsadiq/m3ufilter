package m3u

import (
	"errors"
	"fmt"
	"github.com/grafov/m3u8"
	"github.com/hoshsadiq/m3ufilter/regex"
	"github.com/maja42/goval"
)

var evaluator = goval.NewEvaluator()

func evaluate(ms *m3u8.MediaSegment, expr string) (result interface{}, err error) {
	attrs := make(map[string]string)
	for _, attr := range ms.Attributes {
		attrs[attr.Key] = attr.Value
	}

	variables := map[string]interface{}{
		"Title":    ms.Title,
		"Uri":      ms.URI,
		"Duration": ms.Duration,
		"Attr":     attrs,
	}

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
		"strlen": evaluatorStrlen,
		"match":  evaluatorMatch,
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

//func evaluator_match(args ...interface{}) (interface{}, error) {
//	subject := args[0].(*m3u8.MediaSegment)
//	attrKey := args[1].(string)
//
//	attr, err := GetAttr(ms, attrKey)
//	if err != nil {
//		return nil, err
//	}
//
//	return (string)(attr.Value), nil
//}
