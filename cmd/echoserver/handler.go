package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/google/uuid"
)

func handler(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	response := NewResponse(requestID)

	defer func() {
		if r := recover(); r != nil {
			log.Warnf("%s: an unrecoverable error occurred:\n'%s' - %s", requestID, r, string(debug.Stack()))
			response := fmt.Sprintf(`{"id":"%s","error":["%s"]]}`, requestID, r)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(response))
		}
	}()
	receivedTimestamp := time.Now()
	log.Infof("received incoming request id %s", requestID)

	if len(r.Host) > 0 {
		response.Request.Hostname = &r.Host
	}
	if len(r.Method) > 0 {
		response.Request.Method = &r.Method
	}
	if len(r.URL.Path) > 0 {
		response.Request.Path = &r.URL.Path
	}
	if len(r.Proto) > 0 {
		response.Request.Protocol = &r.Proto
	}
	referer := r.Referer()
	if len(referer) > 0 {
		response.Request.Referer = &referer
	}
	if len(r.RemoteAddr) > 0 {
		response.Request.RemoteAddr = &r.RemoteAddr
	}
	userAgent := r.UserAgent()
	if len(userAgent) > 0 {
		response.Request.UserAgent = &userAgent
	}
	response.Request.Size = r.ContentLength

	log.Debugf("%s: processing request headers...", requestID)
	for header, value := range r.Header {
		if response.Request.Header == nil {
			response.Request.Header = map[string]interface{}{}
		}
		response.Request.Header[header] = value
	}

	log.Debugf("%s: processing url query parameters...", requestID)
	fullQuery := r.URL.Query().Encode()
	if len(fullQuery) > 0 {
		log.Debugf("%s: ...processing query '%s'...", requestID, fullQuery)
		splitQuery := strings.Split(fullQuery, "&")
		for _, query := range splitQuery {
			queryComponents := strings.Split(query, "=")
			if len(queryComponents) != 2 {
				break
			}
			if response.Request.Query == nil {
				response.Request.Query = map[string]interface{}{}
			}
			response.Request.Query[queryComponents[0]] = queryComponents[1]
			log.Debugf("%s: ...processed query component '%s'", requestID, query)
		}
	} else {
		log.Debugf("%s: ...no url query parameters found", requestID)
	}

	log.Debugf("%s: processing request body...", requestID)
	if r.Body != http.NoBody {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			log.Debugf("%s: (raw request body: '%s')", requestID, string(body))
			var jsonBodyData map[string]interface{}
			if err = json.Unmarshal(body, &jsonBodyData); err != nil {
				errorMessage := fmt.Sprintf("failed to parse request body into json: '%s'", err)
				log.Warnf("%s: ...%s", requestID, errorMessage)
				response.Errors = append(response.Errors, errorMessage)
				var bodyAsString = string(body)
				response.Request.Body = &bodyAsString
				log.Debugf("%s: ...set body to raw text input", requestID)
			} else {
				response.Request.Body = jsonBodyData
				log.Debugf("%s: ...set body to parsed json input", requestID)
			}
		} else {
			errorMessage := fmt.Sprintf("failed to read request body into memory: '%s'", err)
			log.Warnf("%s: ...%s", requestID, errorMessage)
			response.Errors = append(response.Errors, errorMessage)
		}
	} else {
		log.Debugf("%s: ...no request body found", requestID)
	}

	log.Debugf("%s: parsing form data...", requestID)
	if err := r.ParseForm(); err != nil {
		errorMessage := fmt.Sprintf("error parsing form data: '%s'", err)
		log.Warnf("%s: ...%s", requestID, errorMessage)
		response.Errors = append(response.Errors, errorMessage)
	} else {
		log.Debugf("%s: ...parsed form data", requestID)
		log.Debugf("%s: processing form data...", requestID)

		formsCount := 0
		for formField, formValue := range r.Form {
			if r.PostForm[formField] == nil {
				if response.Request.Query == nil {
					response.Request.Query = map[string]interface{}{}
				}
				response.Request.Query[formField] = formValue
				log.Debugf("%s: ...processed form field '%s'", requestID, formField)
				formsCount++
			}
		}
		log.Debugf("%s: ...processed %v form fields", requestID, formsCount)

		postFormsCount := 0
		for formField, formValue := range r.PostForm {
			if response.Request.Form == nil {
				response.Request.Form = map[string]interface{}{}
			}
			response.Request.Form[formField] = formValue
			log.Debugf("%s: ...processed post form field '%s'", requestID, formField)
			postFormsCount++
		}
		log.Debugf("%s: ...processed %v post form fields", requestID, postFormsCount)
	}

	log.Debugf("%s: processing basic authentication...", requestID)
	if username, password, ok := r.BasicAuth(); ok {
		response.Request.Username = &username
		response.Request.Password = &password
	} else {
		log.Debugf("%s: ...no basic auth found", requestID)
	}

	log.Debugf("%s: processing request cookies...", requestID)
	cookies := r.Cookies()
	if len(cookies) > 0 {
		for _, cookie := range cookies {
			cookieData := Cookie{
				Name:     cookie.Name,
				Value:    cookie.Value,
				Domain:   cookie.Domain,
				Expires:  cookie.Expires.String(),
				Secure:   cookie.Secure,
				HTTPOnly: cookie.HttpOnly,
			}
			response.Request.Cookies = append(response.Request.Cookies, cookieData)
			log.Debugf("%s: ...processed cookie with name '%s'", requestID, cookie.Name)
		}
	} else {
		log.Debugf("%s: ...no cookies found", requestID)
	}

	log.Debugf("%s: processing metadata...", requestID)
	response.Metadata.RequestProcessingCompleted()

	responseBody, err := json.Marshal(response)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to parse response object into json: '%s'", err)
		log.Warnf("%s: ...%s", requestID, errorMessage)
		response.Errors = append(response.Errors, errorMessage)
	}

	log.Debugf("%s: sending response back to client...", requestID)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if len(response.Errors) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
	if responseSize, err := w.Write(responseBody); err != nil {
		log.Warnf("%s: failed to send response back to client...", requestID)
	} else {
		unit := []string{"bytes", "kb", "Mb", "Gb", "Tb"}
		unitIndex := 0
		simpleResponseSize := float64(responseSize)
		for simpleResponseSize > 1000 {
			simpleResponseSize = simpleResponseSize / 1000
			unitIndex++
		}
		log.Infof("%s: sent response of size %v %s back to client after %s", requestID, simpleResponseSize, unit[unitIndex], time.Since(receivedTimestamp))
	}
}
