// go:build ignore
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/net/hyperhttp"
)

var xmlSampleBody = `<JSONObject><Header><Security><UsernameToken><Id>UsernameToken1234</Id><Username>username123</Username>` +
	`<Password Type="http://foo.bar.com">Password123</Password></UsernameToken></Security></Header></JSONObject>`

func main() {
	cfg := config.Load()
	cfg.ServiceName = config.String("http-server-tim")
	// Switch to OTLP since default reporting endpoint is ZIPKIN
	cfg.Reporting.Endpoint = config.String("localhost:5442")
	cfg.Reporting.TraceReporterType = config.TraceReporterType_OTLP
	// Add "xml" to allowed content types
	cfg.DataCapture.AllowedContentTypes = append(cfg.DataCapture.AllowedContentTypes,
		config.String("xml"))

	// Set the cert file and key file to launch a TLS server.
	// certFile := "./traceable/mitmproxy_self_gen_certs/mitm_domain.crt"
	// keyFile := "./traceable/mitmproxy_self_gen_certs/mitm_domain.key"

	flusher := hypertrace.Init(cfg)
	defer flusher()

	r := mux.NewRouter()
	// curl -v -XPOST http://localhost:8081/foo -H 'Accept: application/json, */*;q=0.5' -H 'Content-type: application/json'
	//  -H 'x-forwarded-for: 73.67.10.215' -d '{"name":"John Q Public"}'
	r.Handle("/foo", hyperhttp.NewHandler(
		http.HandlerFunc(fooHandler),
		"/foo",
	))

	r.Handle("/payload", hyperhttp.NewHandler(
		http.HandlerFunc(payloadHandler),
		"/payload",
	))

	// GET /outgoing makes an exit call
	r.Handle("/outgoing", hyperhttp.NewHandler(
		http.HandlerFunc(outgoingCallHandler),
		"/outgoing",
	))

	// GET request with a query string
	r.Handle("/bar", hyperhttp.NewHandler(
		http.HandlerFunc(barHandler),
		"/bar",
	))
	// POST request with header "Content-Type: application/x-www-form-urlencoded"
	r.Handle("/barbody", hyperhttp.NewHandler(
		http.HandlerFunc(barBodyHandler),
		"/barbody",
	))

	r.Handle("/multipart", hyperhttp.NewHandler(
		http.HandlerFunc(multipartFormHandler),
		"/multipart",
	))

	r.Handle("/echouppercase", hyperhttp.NewHandler(
		http.HandlerFunc(echoUpperCaseHandler),
		"/echouppercase",
	))

	r.Handle("/echo", hyperhttp.NewHandler(
		http.HandlerFunc(echoHandler),
		"/echo",
	))

	// curl -v -XPOST http://localhost:8081/getxml -H 'Accept: application/json, */*;q=0.5' -H 'Content-type: application/json'
	//  -H 'x-forwarded-for: 73.67.10.215' -d '{"name":"John Q Public"}'
	r.Handle("/getxml", hyperhttp.NewHandler(
		http.HandlerFunc(getXmlHandler),
		"/getxml",
	))
	// Uncomment to enable TLS server
	// log.Fatal(http.ListenAndServeTLS(":8081", certFile, keyFile, r))
	// G114 (CWE-676): Use of net/http serve function that has no support for setting timeouts (Confidence: HIGH, Severity: MEDIUM)
	// #nosec G114
	log.Fatal(http.ListenAndServe(":8081", r))
}

type person struct {
	Name          string `json:"name"`
	Password      string `json:password`
	Authorization string `json:authorization`
	Age           int    `json:"age"`
	CurrentCity   city   `json:"current_city"`
	Cities        []city `json:"previous_cities"`
}

type city struct {
	Name       string `json:"name"`
	Population int    `json:"population"`
}

type bigJsonField struct {
	Name         string `json:"Name"`
	Telephone    int64  `json:"Telephone"`
	City         string `json:"City"`
	CreditCardNo string `json:"credit card no"`
	Field2       string `json:"field2"`
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
	sBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	p := &person{}
	err = json.Unmarshal(sBody, p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	<-time.After(5 * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(fmt.Sprintf("{\"param0\": \"Hello %s\", \"authorization\": \"Bearer some-jwt-token\", \"param2\":{\"param3\":\"param4\",\"value\":\"00000\"}}", p.Name)))
}

var responseBody string = `{"serviceCount":4,"personalDetailChildren":[{"id":"41d29bf9-3d7b-4e7b-90f1-e405693b592f","firstName":"Bruce","lastName":"Meyer","dateOfBirth":"1951-4-20",` +
	`"personalDetailChildren":null,"phoneNumber":235932802,"hobbies":["badminton"],"country":"Italy","serviceCount":0},{"firstName":"Sharon","lastName":"Alonso","country":"Vietnam",` +
	`"hobbies":["writing","volleyball","football"],"serviceCount":0,"personalDetailChildren":null,"id":"00ead215-4695-4cbf-817d-6a822751810f","phoneNumber":351877815,` +
	`"dateOfBirth":"1986-11-17"},{"firstName":"Joseph","lastName":"Peterson","phoneNumber":307642061,"serviceCount":0,"id":"20cbf408-ccfd-4cd5-b71c-ec994d5f9eeb","hobbies":["basketball"],` +
	`"country":"Cuba","personalDetailChildren":null,"dateOfBirth":"1957-7-10"},{"id":"6400dce5-d7b8-415d-975c-37a90c3a47a8","dateOfBirth":"1985-8-2","country":"Chile",` +
	`"personalDetailChildren":null,"firstName":"Harry","lastName":"Meyer","phoneNumber":878252759,"hobbies":["volleyball","table tennis","bowling"],"serviceCount":0},{"firstName":"Joseph",` +
	`"country":"India","serviceCount":0,"personalDetailChildren":null,"id":"0cec233b-3795-4078-a553-fd4ae7b9eaf3","lastName":"Austin","phoneNumber":314419610,"dateOfBirth":"1994-5-9",` +
	`"hobbies":["scuba diving","singing","acting"]},{"serviceCount":0,"id":"161484ce-5ec6-445e-8f27-f59eba14b7fe","firstName":"Michael","phoneNumber":567866599,"country":"Cuba",` +
	`"personalDetailChildren":null,"lastName":"Gonzalez","dateOfBirth":"1953-8-16","hobbies":["video gaming","walking","photography"]},{"phoneNumber":654287521,"dateOfBirth":"1981-11-22",` +
	`"hobbies":["scuba diving","table tennis","acting"],"country":"Guatemala","personalDetailChildren":null,"id":"f2134ee3-a6ad-4ee5-858e-3384303a9c3e","firstName":"Jack","lastName":"Austin",` +
	`"serviceCount":0},{"firstName":"Cassandra","phoneNumber":861162853,"dateOfBirth":"1956-12-8","personalDetailChildren":null,"id":"f7797a75-0317-46fd-a7de-6c12fc1b881d","lastName":"Johnson",` +
	`"hobbies":["blogging","writing","badminton"],"country":"United Kingdom","serviceCount":0},{"serviceCount":0,"personalDetailChildren":null,"lastName":"Matthews","phoneNumber":261799483,` +
	`"hobbies":["table tennis"],"country":"Israel","id":"42e7dd1b-fee7-4c38-9d90-c16bda78401c","firstName":"Martin","dateOfBirth":"1993-8-31"},{"phoneNumber":146698856,"dateOfBirth":"1995-12-3",` +
	`"country":"Portugal","id":"6b83a18d-0db3-48ed-8c7a-fcf5a4ffa2ff","lastName":"Hudson","serviceCount":0,"personalDetailChildren":null,"firstName":"Nina","hobbies":["football"]},{"serviceCount":0,` +
	`"personalDetailChildren":null,"id":"c2df72de-8c8e-4090-bf5b-ca9e5d4eee99","phoneNumber":642291654,"dateOfBirth":"1944-2-20","hobbies":["basketball","football","acting"],"country":"Jamaica",` +
	`"firstName":"Molly","lastName":"Davis"},{"lastName":"Elliot","dateOfBirth":"1947-2-7","country":"Jamaica","id":"aca0c316-f71f-496d-8c62-aa977a4111c5","phoneNumber":132619301,"hobbies":["writing"],` +
	`"serviceCount":0,"personalDetailChildren":null,"firstName":"Rose"},{"id":"d4329f55-855d-4c79-8371-a6574cbd42fa","firstName":"Molly","lastName":"Washington","dateOfBirth":"1993-2-3",` +
	`"serviceCount":0,"phoneNumber":815893907,"hobbies":["video gaming","yoga"],"country":"Portugal","personalDetailChildren":null},{"firstName":"Cindy","phoneNumber":958956657,"dateOfBirth":"1921-2-22",` +
	`"country":"Canada","serviceCount":0,"personalDetailChildren":null,"id":"3809db0c-97f3-40c9-b0d1-e8d97dbe4015","lastName":"Peterson","hobbies":["archery"]},{"personalDetailChildren":null,` +
	`"id":"5409ff19-cda8-4741-8a42-53e42a211ab3","firstName":"Nick","phoneNumber":758113006,"dateOfBirth":"1995-9-20","hobbies":["hockey","basketball","archery"],"country":"Netherlands","lastName":"Richardson",` +
	`"serviceCount":0},{"lastName":"Wilson","phoneNumber":386712898,"dateOfBirth":"1972-2-15","hobbies":["hockey","video gaming","badminton"],"serviceCount":0,"id":"a498a1c2-900e-4cde-89a3-83ed91dc1444",` +
	`"firstName":"David","country":"Germany","personalDetailChildren":null},{"lastName":"Hernandez","dateOfBirth":"1974-8-23","personalDetailChildren":null,` +
	`"country":"China","serviceCount":0,"id":"6fc67f4f-cd3d-4691-a8a1-22c7e2ae1289","firstName":"Julie","phoneNumber":505837958,` +
	`"hobbies":["basketball","blogging","reading"]},{"serviceCount":0,"id":"926516ef-4059-40dc-bfd8-bcfe4ff50207","lastName":"Washington",` +
	`"phoneNumber":481016364,"dateOfBirth":"1935-7-9","hobbies":["running","badminton","bowling"],"country":"Bahrain","firstName":"Jack",` +
	`"personalDetailChildren":null},{"dateOfBirth":"1922-11-13","serviceCount":0,"personalDetailChildren":null,` +
	`"id":"3cc3d3de-691b-4a5e-9041-f379ba45361d","firstName":"Adam","hobbies":["scuba diving","bowling"],"country":"Spain",` +
	`"lastName":"Olson","phoneNumber":932252939},{"id":"f3a11e19-390f-4271-9ed1-b9a918b4ae09","firstName":"Nick",` +
	`"phoneNumber":892624868,"serviceCount":0,"lastName":"Jackson","dateOfBirth":"1951-4-6",` +
	`"hobbies":["badminton","volleyball"],"country":"Colombia","personalDetailChildren":null}],` +
	`"id":"1b156b63-e965-41b3-a0d4-2ffb7f79bae6","lastName":"Austin","phoneNumber":865542508,` +
	`"dateOfBirth":"1992-3-8","hobbies":["writing"],"country":"Bahrain","firstName":"Mark"}`

func payloadHandler(w http.ResponseWriter, r *http.Request) {
	sBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	p := &bigJsonField{}
	err = json.Unmarshal(sBody, p)
	if err != nil {
		fmt.Printf("unmarshal error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	<-time.After(3 * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(responseBody))
}

// Should be a GET request
func barHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	author := r.FormValue("author")

	<-time.After(30 * time.Millisecond)
	// echo custom headers. must be called before w.WriteHeader() which adds the status code
	w.Header().Set("Content-Type", "application/json")
	for k, v := range r.Header {
		//fmt.Printf("header k: %s, val: %v\n", k, v)
		if strings.HasPrefix(strings.ToLower(k), "x-") {
			w.Header().Set(k, strings.Join(v, ","))
		}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf("{\"message\": \"Hello %s\"}", author)))
}

// Should be a POST request.
func barBodyHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	author := r.PostFormValue("author")

	<-time.After(300 * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", "timothercookie23=cval1;secure")
	w.Header().Add("Set-Cookie", "timothercookie789=cval1;secure")
	w.Header().Add("Set-Cookie", "timothercookie=cval2;secure")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf("{\"message\": \"Hello %s\"}", author)))
}

const multipartMaxSize int64 = 1024 * 1024

func multipartFormHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(multipartMaxSize)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	author := r.PostFormValue("author")
	f, fh, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("error while reading file object from form: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("file size: %d\n", fh.Size)
	buf := make([]byte, fh.Size)
	_, err = f.Read(buf)
	if err != nil {
		fmt.Printf("error while reading file into buf: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// G303 (CWE-377): File creation in shared tmp directory without using ioutil.Tempfile (Confidence: HIGH, Severity: MEDIUM)
	// #nosec G303
	err = os.WriteFile("/tmp/gadat1.png", buf, 0600)
	if err != nil {
		fmt.Printf("error while writing file: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	<-time.After(300 * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", "cookie23=cval1;secure")
	w.Header().Add("Set-Cookie", "othercookie789=cval1;secure")
	w.Header().Add("Set-Cookie", "othercookie=cval2;secure")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf("{\"message\": \"Hello %s\"}", author)))
}

func outgoingCallHandler(w http.ResponseWriter, r *http.Request) {
	client := http.Client{
		Transport: hyperhttp.NewTransport(
			http.DefaultTransport,
		),
	}
	// "In order to correlate the server call to this client call we need to pass the server request's context."
	req, _ := http.NewRequestWithContext(r.Context(), "GET", "https://httpbin.org/get", nil)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	for key, val := range resp.Header {
		if len(val) > 0 {
			w.Header().Set(key, val[len(val)-1])
		}
	}

	// For HTTP client, the span.End() is not invoked until the response body closer is invoked. See opentelemetry-go-contrib/instrumentation/net/http/otelhttp/transport.go.
	// In transport.go the response body is of type io.ReadCloser and it is wrapped by a ReadCloser whose Close() function calls span.End().
	// So the span will not be ended until either res.Body.Close() is invoked or the whole body is read. So in order to capture the client span, you need to read the body or
	// invoke res.Body.Close().
	// It is advisable because there’s a chance you might have a memory leak if you do not close the response body.
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	sb := string(body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, sb)
}

func echoUpperCaseHandler(w http.ResponseWriter, r *http.Request) {
	sBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	<-time.After(2 * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(strings.ToUpper(string(sBody))))
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	sBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	<-time.After(2 * time.Millisecond)

	reqContentType := r.Header.Get("Content-Type")

	if len(reqContentType) != 0 {
		w.Header().Set("Content-Type", reqContentType)
	}
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(string(sBody)))
}

func getXmlHandler(w http.ResponseWriter, r *http.Request) {
	sBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	p := &person{}
	err = json.Unmarshal(sBody, p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	<-time.After(30 * time.Millisecond)

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(xmlSampleBody))
}
