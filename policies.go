package core

import "net/http"

func GetPolicyResult(name string, data interface{}, version string, predefinedId string) (any, error) {
	result, err := getPolicyProvider(http.DefaultClient).Evaluate(name, data, version, predefinedId)
	if result != nil {
		return result.Result, err
	} else {
		return nil, err
	}
}
