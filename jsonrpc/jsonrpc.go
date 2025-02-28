//go:build !solution

package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

func getReflectMethod(service interface{}, methodName string) reflect.Value {
	return reflect.ValueOf(service).MethodByName(methodName)
}

func checkMethodValidity(method reflect.Value, methodName string) error {
	if !method.IsValid() {
		return fmt.Errorf("method %s not found", methodName)
	}
	return nil
}

func prepareRequestBody(r *http.Request, reqType reflect.Type) (reflect.Value, error) {
	reqValue := reflect.New(reqType.Elem())
	if err := json.NewDecoder(r.Body).Decode(reqValue.Interface()); err != nil {
		return reflect.Value{}, fmt.Errorf("failed to deserialize request body: %v", err)
	}
	return reqValue, nil
}

func checkMethodReturnValues(result []reflect.Value) error {
	if len(result) != 2 {
		return fmt.Errorf("method must return two values")
	}
	err := result[1].Interface()
	if err != nil {
		return err.(error)
	}
	return nil
}

func sendJSONResponse(w http.ResponseWriter, rsp reflect.Value) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(rsp.Interface())
}

func MakeHandler(service interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		methodName := r.URL.Path[1:]

		method := getReflectMethod(service, methodName)

		if err := checkMethodValidity(method, methodName); err != nil {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		reqType := method.Type().In(1)
		reqValue, err := prepareRequestBody(r, reqType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := method.Call([]reflect.Value{reflect.ValueOf(ctx), reqValue})

		if err := checkMethodReturnValues(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := sendJSONResponse(w, result[0]); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func encodeRequest(req interface{}) ([]byte, error) {
	return json.Marshal(req)
}

func createPostRequest(ctx context.Context, url string, reqBody []byte) (*http.Request, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	return httpReq, nil
}

func executeRequest(httpReq *http.Request) (*http.Response, error) {
	httpClient := http.DefaultClient
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	return resp, nil
}

func processResponse(resp *http.Response, rsp interface{}) error {
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read error response body: %v", err)
		}
		return errors.New(string(bodyBytes))
	}

	return json.NewDecoder(resp.Body).Decode(rsp)
}

func Call(ctx context.Context, endpoint string, method string, req, rsp interface{}) error {
	reqBody, err := encodeRequest(req)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s", endpoint, method)

	httpReq, err := createPostRequest(ctx, url, reqBody)
	if err != nil {
		return err
	}

	resp, err := executeRequest(httpReq)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
	}(resp.Body)

	return processResponse(resp, rsp)
}
