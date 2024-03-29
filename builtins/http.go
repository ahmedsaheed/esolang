package builtins

import (
	"esolang/lang-esolang/object"
	"fmt"
	"io"
	"net/http"
	"strings"
)

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
		resp, err := http.Get(url)
		if err != nil {
			return newError("HTTP Error: Error making GET request %s", url)
		}
		body, err := io.ReadAll(resp.Body)
		stringifiedResp := string(body)

		// TODO: We can check if it [Object.Hash] or [Object.Array] and convert it to respective object
		return &object.String{Value: stringifiedResp}
	}

	outGoingBody := args[2]

	if requestType == "POST" {
		if outGoingBody.Type() != object.HASH_OBJ {
			return newError("Body must be of type Hash, got %s", args[2].Type())
		}
		body := outGoingBody.(*object.Hash).Inspect()

		resp, err := http.Post(url, "application/json", strings.NewReader(body))
		if err != nil {
			return newError("HTTP Error: Error making POST request %s", url)
		}
		respBody, err := io.ReadAll(resp.Body)
		statusCode := fmt.Sprintf(" Status Code: %d", resp.StatusCode)
		stringifiedResp := string(respBody) + statusCode
		// TODO: We can check if it [Object.Hash] or [Object.Array] and convert it to respective object
		return &object.String{Value: stringifiedResp}

	}

	if requestType == "PATCH" {
		if outGoingBody.Type() != object.HASH_OBJ {
			return newError("Body must be of type Hash, got %s", args[2].Type())
		}
		body := outGoingBody.(*object.Hash).Inspect()
		client := &http.Client{}
		req, err := http.NewRequest("PATCH", url, strings.NewReader(body))
		if err != nil {
			return newError("HTTP Error: Error making PATCH request %s", url)
		}
		resp, err := client.Do(req)

		if err != nil {
			return newError("HTTP Error: Error making PATCH request %s", url)
		}
		respBody, err := io.ReadAll(resp.Body)
		stringifiedResp := string(respBody)
		// TODO: We can check if it [Object.Hash] or [Object.Array] and convert it to respective object
		return &object.String{Value: stringifiedResp}
	}

	return NULL
}

func checkBrackets(text string) bool {
	if len(text) < 2 {
		return false // String is too short
	}
	return text[0] == '[' && text[len(text)-1] == ']'
}

// covert string to object.Hash
func convertToHash(text string) *object.Hash {
	individualObjs := strings.Split(text, "},")
	finalHash := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
	text = text[1 : len(text)-1]

	for _, individualObj := range individualObjs {
		// get each key value pair
		replacer := strings.NewReplacer("{", "", "}", "").Replace(individualObj)
		keyValuePairs := strings.Split(replacer, ",")
		fmt.Println("Key Value Pairs are", keyValuePairs)

		// create a new hash
		hash := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
		for _, pair := range keyValuePairs {
			// get key value pair
			keyValue := strings.Split(pair, ":")
			if len(keyValue) != 2 {
				continue
			}
			key := object.String{
				Value: strings.Trim(keyValue[0], " "),
			}
			hashKey := key.HashKey()

			value := strings.Trim(keyValue[1], " ")
			fmt.Println("Key is ", key, " Value is ", value)
			hash.Pairs[hashKey] = object.HashPair{Key: &object.String{Value: key.Value}, Value: &object.String{Value: value}}

			finalHash.Pairs[hashKey] = object.HashPair{Key: &object.String{Value: key.Value}, Value: hash}

		}

	}

	fmt.Println("Converted Hash is ", finalHash.Inspect())
	return finalHash
}
