package builtins

import (
	"encoding/json"
	"esolang/lang-esolang/object"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type Response struct {
	Body       string
	StatusCode int
	Status     string
	Header     http.Header
}

func _http(args ...object.Object) object.Object {
	requestType := args[0].(*object.String).Value
	if len(args) < 2 || len(args) > 3 {
		return newError("Invalid Arguement got=%d, expected=3", len(args))
	}
	if args[0].Type() != object.STRING_OBJ {
		if strings.ToUpper(requestType) != "GET" && strings.ToUpper(requestType) != "POST" && strings.ToUpper(requestType) != "PATCH" {
			return newError("Request type can either be a `GET` or `POST` or `PATCH`")
		}
		return newError("Request type must be of type String, got %s", args[0].Type())
	}

	if args[1].Type() != object.STRING_OBJ {
		return newError("URL must be of type String, got %s", args[1].Type())
	}
	url := args[1].(*object.String).Value
	if requestType == "GET" {
		// make sure body is empty
		if len(args) == 3 {
			return newError("GET request should not have a body")
		}
		response, err := http.Get(url)
		if err != nil {
			return newError("HTTP Error: Error making GET request %s", url)
		}
		return outputResp(response)
	}

	outGoingBody := args[2]

	if requestType == "POST" {
		if outGoingBody.Type() != object.HASH_OBJ {
			return newError("Body must be of type Hash, got %s", args[2].Type())
		}
		body := outGoingBody.(*object.Hash).Inspect()

		response, err := http.Post(url, "application/json", strings.NewReader(body))
		if err != nil {
			return newError("HTTP Error: Error making POST request %s", url)
		}

		if err != nil {
			return newError("HTTP Error: could read response body %s", err)
		}
		return outputResp(response)
	}

	if requestType == "PATCH" {
		if outGoingBody.Type() != object.HASH_OBJ {
			return newError("Body must be of type Hash, got %s", args[2].Type())
		}
		body := outGoingBody.(*object.Hash).Inspect()
		client := &http.Client{}
		req, err := http.NewRequest(requestType, url, strings.NewReader(body))
		if err != nil {
			return newError("HTTP Error: Error making PATCH request %s", url)
		}
		response, err := client.Do(req)

		if err != nil {
			return newError("HTTP Error: Error making PATCH request %s", url)
		}
		return outputResp(response)
	}

	return NULL
}

func checkBrackets(text string) bool {
	if len(text) < 2 {
		return false // String is too short
	}
	return text[0] == '[' && text[len(text)-1] == ']'
}

func outputResp(res *http.Response) object.Object {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		newError("HTTP Error: could read response body %s", err)
	}
	stringifiedBody := string(body)

	resp := Response{
		Body:       stringifiedBody,
		StatusCode: res.StatusCode,
		Status:     res.Status,
		Header:     res.Header,
	}
	return generateHashFromResponse(resp)
}

func generateHashFromResponse(res Response) *object.Hash {
	hash := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
	headerHash := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
	// bodyHash := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}

	for key, value := range res.Header {
		keyValue := (&object.String{Value: key}).HashKey()
		hashKey := object.HashKey{Type: object.STRING_OBJ, Value: keyValue.Value}
		headerHash.Pairs[hashKey] = object.HashPair{Key: &object.String{Value: key}, Value: &object.String{Value: value[0]}}
	}

	headerKey := (&object.String{Value: "Header"}).HashKey()
	hashKey := object.HashKey{Type: object.STRING_OBJ, Value: headerKey.Value}
	hash.Pairs[hashKey] = object.HashPair{Key: &object.String{Value: "Header"}, Value: headerHash}

	statusCodeKey := (&object.String{Value: "StatusCode"}).HashKey()
	hashKey = object.HashKey{Type: object.STRING_OBJ, Value: statusCodeKey.Value}
	hash.Pairs[hashKey] = object.HashPair{Key: &object.String{Value: "StatusCode"}, Value: &object.Integer{Value: int64(res.StatusCode)}}

	statusKey := (&object.String{Value: "Status"}).HashKey()
	hashKey = object.HashKey{Type: object.STRING_OBJ, Value: statusKey.Value}
	hash.Pairs[hashKey] = object.HashPair{Key: &object.String{Value: "Status"}, Value: &object.String{Value: res.Status}}

	// populate bodyHash

	bodyKey := (&object.String{Value: "Body"}).HashKey()
	bodyHash := generateResponseBodyHash(res.Body)
	bodyHashKey := object.HashKey{Type: object.STRING_OBJ, Value: bodyKey.Value}
	hash.Pairs[bodyHashKey] = object.HashPair{Key: &object.String{Value: "Body"}, Value: bodyHash}

	return hash
}

func generateResponseBodyHash(body string) *object.Hash {
	// VERY MUCH WIP
	hash := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}

	if checkBrackets(body) {
		// parse the string as a json array
		var arr []interface{}
		err := json.Unmarshal([]byte(body), &arr)
		if err != nil {
			fmt.Println("error parsing json object")
			return &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
		}
		for i, v := range arr {
			hashValueKey := (&object.Integer{Value: int64(i)}).HashKey()
			hashKey := object.HashKey{Type: object.INTEGER_OBJ, Value: hashValueKey.Value}
			hash.Pairs[hashKey] = object.HashPair{Key: &object.Integer{Value: int64(i)}, Value: &object.String{Value: fmt.Sprintf("%v", v)}}
		}
		// fmt.Println("value of hash arrayy ", hash.Inspect())
		return hash
	}

	// parse the string as a json object
	var obj map[string]interface{}
	err := json.Unmarshal([]byte(body), &obj)
	if err != nil {
		fmt.Println("error parsing json object")
		return &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
	}
	for k, v := range obj {
		hashValueKey := (&object.String{Value: k}).HashKey()
		hashKey := object.HashKey{Type: object.STRING_OBJ, Value: hashValueKey.Value}
		// fmt.Println("value of k1 ", k)
		// nested hash are currently converted to string - let fix this

		if reflect.TypeOf(v).Kind() == reflect.Map {
			for k1, v1 := range v.(map[string]interface{}) {
				innerHash := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
				hashValueKey := (&object.String{Value: k1}).HashKey()
				hashKey := object.HashKey{Type: object.STRING_OBJ, Value: hashValueKey.Value}
				innerHash.Pairs[hashKey] = object.HashPair{Key: &object.String{Value: k1}, Value: &object.String{Value: fmt.Sprintf("%v", v1)}}
				hash.Pairs[hashKey] = object.HashPair{Key: &object.String{Value: k}, Value: innerHash}

			}
		} else {
			hash.Pairs[hashKey] = object.HashPair{Key: &object.String{Value: k}, Value: &object.String{Value: fmt.Sprintf("%v", v)}}
		}
	}
	// fmt.Println("value of hash hashable ", hash.Inspect())
	return hash
}
