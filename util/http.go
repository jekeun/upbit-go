package util

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type RequestOptions struct {
	Url     string
	Method  string
	Body    io.Reader
	Query   map[string]string
	Headers map[string]string
}

func Request(options *RequestOptions, result interface{}) (
	err error,
) {
	client := &http.Client{}

	req, err := http.NewRequest(options.Method, options.Url, options.Body)
	if err != nil {
		return
	}

	if options.Query != nil {
		q := req.URL.Query()
		for index, value := range options.Query {
			q.Add(index, value)
		}

		req.URL.RawQuery = q.Encode()
	}

	if options.Headers != nil {
		for prop, value := range options.Headers {
			req.Header.Add(prop, value)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	Body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}


	////fmt.Println(res.Header)
	//fmt.Println(res.StatusCode)
	//if res.StatusCode != 200 {
	//	bodyString := string(Body)
	//	fmt.Println(bodyString)
	//}
	////remaingReq := res.Header.Get("Remaining-Req")
	////fmt.Println(remaingReq)
	err = json.Unmarshal(Body, result)
	if err != nil {
		//bodyString := string(Body)
		//fmt.Println(bodyString)
		//
		////fmt.Println(Body)
		////fmt.Println(err)
		return
	}
	return
}
