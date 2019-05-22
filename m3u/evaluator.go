package m3u

import (
	"errors"
	"fmt"
	"github.com/grafov/m3u8"
	"github.com/hoshsadiq/m3ufilter/regex"
	"github.com/maja42/goval"
	"strings"
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
		"strlen":  evaluatorStrlen,
		"match":   evaluatorMatch,
		"replace": evaluatorReplace,
		"tvg_id":  evaluatorToTvgId,
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
	subject = strings.Replace(subject, "FHD", "", -1)
	subject = strings.Replace(subject, "HD", "", -1)
	subject = strings.Replace(subject, "SD", "", -1)
	subject = strings.TrimSpace(subject)

	re := regex.GetCache("[^a-zA-Z0-9]")
	return re.ReplaceAllString(subject, ""), nil
}

// the below might come in handy.
//func createNewTvgId(title string) string {
//	title = strings.Replace(title, "FHD", "", -1)
//	title = strings.Replace(title, "HD", "", -1)
//	title = strings.Replace(title, "SD", "", -1)
//	title = strings.TrimSpace(title)
//
//	// todo this regex needs to be configurable
//	countryRe := regex.GetCache("(?i)(^(USA?|UK|NL)|\\.(uk|us|nl))\b")
//	country := countryRe.ReplaceAllString(title, "$2|$3")
//
//	re := regex.GetCache(fmt.Sprintf("(^(%s)|[^a-zA-Z0-9])", country))
//	parsedTitle := re.ReplaceAllString(title, "")
//
//	countryMatches := strings.Split(country, "|")
//	if countryMatches[1] != "" {
//		country = countryMatches[1]
//	} else {
//		country = countryMatches[0]
//		if country == "USA" {
//			country = "us"
//		}
//	}
//
//	country = strings.ToLower(country)
//
//	if country != "" {
//		return parsedTitle + "." + country
//	}
//
//	log.Warnf("Tried to guess new tvg-id, but country was not found for %s", title)
//	return parsedTitle
//}
