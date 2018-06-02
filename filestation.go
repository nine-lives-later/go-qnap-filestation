package filestation

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	netUrl "net/url"
	"reflect"
	"strings"
	"time"
)

var defaultConfigOptions = &ConfigOptions{
	APICallTimeout: 60 * time.Second,
}

// ConfigOptions contains some advanced settings on server communication.
type ConfigOptions struct {
	APICallTimeout time.Duration
}

// FileStationSession is a container for our session state.
type FileStationSession struct {
	Host          string
	SessionID     string
	Transport     *http.Transport
	ConfigOptions *ConfigOptions
}

// String returns the session's hostname.
func (s *FileStationSession) String() string {
	return s.Host
}

// APIRequest builds our request before sending it to the server.
type APIRequest struct {
	Method      string
	URL         string
	Body        string
	ContentType string
	QueryParams QueryParameters
}

// RequestError contains information about any error we get from a request.
type RequestError struct {
	Status int `json:"status,omitempty"`
}

// QueryParameters contains the parameters provided to the script.
type QueryParameters map[string]string

// Error returns the error message.
func (r RequestError) Error() string {
	msg := "unknown"

	switch r.Status {
	case 8: // disabled
		msg = "API is disabled"
	case 4: // permission denied
		msg = "Permission denied"
	case 3: // session expired
		msg = "Session expired"
	}

	return fmt.Sprintf("Status Code: %v (%v)", r.Status, msg)
}

// Connect sets up our connection to the Zevenet system.
func Connect(host, username, password string, configOptions *ConfigOptions) (*FileStationSession, error) {
	var url string
	if !strings.HasPrefix(host, "http") {
		url = fmt.Sprintf("https://%s", host)
	} else {
		url = host
	}
	if configOptions == nil {
		configOptions = defaultConfigOptions
	}

	// create the session
	session := &FileStationSession{
		Host: url,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		ConfigOptions: configOptions,
	}

	// initialize the session
	err := session.initialize()
	if err != nil {
		return nil, err
	}

	// perform login
	err = session.Login(username, password)
	if err != nil {
		return nil, err
	}

	// done
	return session, nil
}

func (s *FileStationSession) initialize() error {
	return nil
}

// apiCall is used to query the ZAPI.
func (s *FileStationSession) apiCall(options *APIRequest) ([]byte, error) {
	var req *http.Request
	client := &http.Client{
		Transport: s.Transport,
		Timeout:   s.ConfigOptions.APICallTimeout,
	}

	// build url
	var url strings.Builder

	url.WriteString(s.Host)
	url.WriteString("/")
	url.WriteString(options.URL)
	url.WriteString("?")

	if s.SessionID != "" {
		url.WriteString("sid=")
		url.WriteString(netUrl.QueryEscape(s.SessionID))
		url.WriteString("&")
	}

	for k, v := range options.QueryParams {
		url.WriteString(k)
		url.WriteString("=")
		url.WriteString(netUrl.QueryEscape(v))
		url.WriteString("&")
	}

	body := bytes.NewReader([]byte(options.Body))
	req, _ = http.NewRequest(strings.ToUpper(options.Method), url.String(), body)

	//fmt.Println("REQ -- ", options.Method, " ", url.String(), " -- ", options.Body)

	if len(options.ContentType) > 0 {
		req.Header.Set("Content-Type", options.ContentType)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, _ := ioutil.ReadAll(res.Body)

	//fmt.Println("RES --", res.StatusCode, " -- ", string(data))

	return data, nil
}

//Generic delete
func (s *FileStationSession) delete(scriptPath string, params QueryParameters) error {
	req := &APIRequest{
		Method:      "delete",
		URL:         scriptPath,
		QueryParams: params,
	}

	_, callErr := s.apiCall(req)
	return callErr
}

func (s *FileStationSession) post(body interface{}, scriptPath string, params QueryParameters) error {
	marshalJSON, err := jsonMarshal(body)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "post",
		URL:         scriptPath,
		Body:        strings.TrimRight(string(marshalJSON), "\n"),
		ContentType: "application/json",
		QueryParams: params,
	}

	_, callErr := s.apiCall(req)
	return callErr
}

func (s *FileStationSession) put(body interface{}, scriptPath string, params QueryParameters) error {
	marshalJSON, err := jsonMarshal(body)
	if err != nil {
		return err
	}

	return s.putRaw(marshalJSON, scriptPath, params)
}

func (s *FileStationSession) putRaw(body []byte, scriptPath string, params QueryParameters) error {
	req := &APIRequest{
		Method:      "put",
		URL:         scriptPath,
		Body:        strings.TrimRight(string(body), "\n"),
		ContentType: "application/json",
		QueryParams: params,
	}

	_, callErr := s.apiCall(req)
	return callErr
}

//Get a url and populate an entity. If the entity does not exist (404) then the
//passed entity will be untouched and false will be returned as the second parameter.
//You can use this to distinguish between a missing entity or an actual error.
func (s *FileStationSession) getForEntity(e interface{}, scriptPath string, params QueryParameters) error {
	resp, err := s.getRaw(scriptPath, params)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, e)
	if err != nil {
		return s.checkError(resp)
	}

	return nil
}

func (s *FileStationSession) getRaw(scriptPath string, params QueryParameters) ([]byte, error) {
	req := &APIRequest{
		Method:      "get",
		URL:         scriptPath,
		ContentType: "application/json",
		QueryParams: params,
	}

	return s.apiCall(req)
}

// checkError handles any errors we get from our API requests. It returns either the
// message of the error, if any, or nil.
func (s *FileStationSession) checkError(resp []byte) error {
	if len(resp) == 0 {
		return nil
	}

	var reqError RequestError

	err := json.Unmarshal(resp, &reqError)
	if err != nil {
		return fmt.Errorf("%s\n%s", err.Error(), string(resp[:]))
	}

	return reqError
}

// jsonMarshal specifies an encoder with 'SetEscapeHTML' set to 'false' so that <, >, and & are not escaped. https://golang.org/pkg/encoding/json/#Marshal
// https://stackoverflow.com/questions/28595664/how-to-stop-json-marshal-from-escaping-and
func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// Helper to copy between transfer objects and model objects to hide the myriad of boolean representations
// in the iControlREST api. DTO fields can be tagged with bool:"yes|enabled|true" to set what true and false
// marshal to.
func marshal(to, from interface{}) error {
	toVal := reflect.ValueOf(to).Elem()
	fromVal := reflect.ValueOf(from).Elem()
	toType := toVal.Type()
	for i := 0; i < toVal.NumField(); i++ {
		toField := toVal.Field(i)
		toFieldType := toType.Field(i)
		fromField := fromVal.FieldByName(toFieldType.Name)
		if fromField.Interface() != nil && fromField.Kind() == toField.Kind() {
			toField.Set(fromField)
		} else if toField.Kind() == reflect.Bool && fromField.Kind() == reflect.String {
			switch fromField.Interface() {
			case "yes", "enabled", "enable", "true":
				toField.SetBool(true)
				break
			case "no", "disabled", "disable", "false", "":
				toField.SetBool(false)
				break
			default:
				return fmt.Errorf("Unknown boolean conversion for %s: %s", toFieldType.Name, fromField.Interface())
			}
		} else if fromField.Kind() == reflect.Bool && toField.Kind() == reflect.String {
			tag := toFieldType.Tag.Get("bool")
			switch tag {
			case "yes":
				toField.SetString(toBoolString(fromField.Interface().(bool), "yes", "no"))
				break
			case "enabled":
				toField.SetString(toBoolString(fromField.Interface().(bool), "enabled", "disabled"))
				break
			case "enable":
				toField.SetString(toBoolString(fromField.Interface().(bool), "enable", "disable"))
				break
			case "true":
				toField.SetString(toBoolString(fromField.Interface().(bool), "true", "false"))
				break
			}
		} else {
			return fmt.Errorf("Unknown type conversion %s -> %s", fromField.Kind(), toField.Kind())
		}
	}
	return nil
}

func toBoolString(b bool, trueStr, falseStr string) string {
	if b {
		return trueStr
	}
	return falseStr
}
