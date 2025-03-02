package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/crowdsecurity/crowdsec/pkg/csplugin"
	"github.com/crowdsecurity/crowdsec/pkg/models"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type LAPI struct {
	router     *gin.Engine
	loginResp  models.WatcherAuthResponse
	bouncerKey string
	t          *testing.T
}

func SetupLAPITest(t *testing.T) LAPI {
	t.Helper()
	router, loginResp, err := InitMachineTest()
	if err != nil {
		t.Fatal(err.Error())
	}

	APIKey, err := CreateTestBouncer()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	return LAPI{
		router:     router,
		loginResp:  loginResp,
		bouncerKey: APIKey,
	}
}

func (l *LAPI) InsertAlertFromFile(path string) *httptest.ResponseRecorder {
	alertReader := GetAlertReaderFromFile(path)
	return l.RecordResponse("POST", "/v1/alerts", alertReader)
}

func (l *LAPI) RecordResponse(verb string, url string, body *strings.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(verb, url, body)
	if err != nil {
		l.t.Fatal(err)
	}
	req.Header.Add("X-Api-Key", l.bouncerKey)
	AddAuthHeaders(req, l.loginResp)
	l.router.ServeHTTP(w, req)
	return w
}

func InitMachineTest() (*gin.Engine, models.WatcherAuthResponse, error) {
	router, err := NewAPITest()
	if err != nil {
		return nil, models.WatcherAuthResponse{}, fmt.Errorf("unable to run local API: %s", err)
	}

	loginResp, err := LoginToTestAPI(router)
	if err != nil {
		return nil, models.WatcherAuthResponse{}, fmt.Errorf("%s", err.Error())
	}
	return router, loginResp, nil
}

func LoginToTestAPI(router *gin.Engine) (models.WatcherAuthResponse, error) {
	body, err := CreateTestMachine(router)
	if err != nil {
		return models.WatcherAuthResponse{}, fmt.Errorf("%s", err.Error())
	}
	err = ValidateMachine("test")
	if err != nil {
		log.Fatalln(err.Error())
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/watchers/login", strings.NewReader(body))
	req.Header.Add("User-Agent", UserAgent)
	router.ServeHTTP(w, req)

	loginResp := models.WatcherAuthResponse{}
	err = json.NewDecoder(w.Body).Decode(&loginResp)
	if err != nil {
		return models.WatcherAuthResponse{}, fmt.Errorf("%s", err.Error())
	}
	return loginResp, nil
}

func AddAuthHeaders(request *http.Request, authResponse models.WatcherAuthResponse) {
	request.Header.Add("User-Agent", UserAgent)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authResponse.Token))
}

func TestSimulatedAlert(t *testing.T) {
	lapi := SetupLAPITest(t)
	lapi.InsertAlertFromFile("./tests/alert_minibulk+simul.json")
	alertContent := GetAlertReaderFromFile("./tests/alert_minibulk+simul.json")
	//exclude decision in simulation mode

	w := lapi.RecordResponse("GET", "/v1/alerts?simulated=false", alertContent)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"Ip 91.121.79.178 performed crowdsecurity/ssh-bf (6 events over `)
	assert.NotContains(t, w.Body.String(), `"message":"Ip 91.121.79.179 performed crowdsecurity/ssh-bf (6 events over `)
	//include decision in simulation mode

	w = lapi.RecordResponse("GET", "/v1/alerts?simulated=true", alertContent)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"Ip 91.121.79.178 performed crowdsecurity/ssh-bf (6 events over `)
	assert.Contains(t, w.Body.String(), `"message":"Ip 91.121.79.179 performed crowdsecurity/ssh-bf (6 events over `)
}

func TestCreateAlert(t *testing.T) {
	lapi := SetupLAPITest(t)
	// Create Alert with invalid format

	w := lapi.RecordResponse("POST", "/v1/alerts", strings.NewReader("test"))
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"message\":\"invalid character 'e' in literal true (expecting 'r')\"}", w.Body.String())

	// Create Alert with invalid input
	alertContent := GetAlertReaderFromFile("./tests/invalidAlert_sample.json")

	w = lapi.RecordResponse("POST", "/v1/alerts", alertContent)
	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "{\"message\":\"validation failure list:\\n0.scenario in body is required\\n0.scenario_hash in body is required\\n0.scenario_version in body is required\\n0.simulated in body is required\\n0.source in body is required\"}", w.Body.String())

	// Create Valid Alert
	w = lapi.InsertAlertFromFile("./tests/alert_sample.json")
	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "[\"1\"]", w.Body.String())
}

func TestCreateAlertChannels(t *testing.T) {

	apiServer, err := NewAPIServer()
	if err != nil {
		log.Fatalln(err.Error())
	}
	apiServer.controller.PluginChannel = make(chan csplugin.ProfileAlert)
	apiServer.InitController()

	loginResp, err := LoginToTestAPI(apiServer.router)
	if err != nil {
		log.Fatalln(err.Error())
	}
	lapi := LAPI{router: apiServer.router, loginResp: loginResp}

	var pd csplugin.ProfileAlert
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		pd = <-apiServer.controller.PluginChannel
		wg.Done()
	}()

	go lapi.InsertAlertFromFile("./tests/alert_ssh-bf.json")
	wg.Wait()
	assert.Equal(t, len(pd.Alert.Decisions), 1)
	apiServer.Close()
}

func TestAlertListFilters(t *testing.T) {
	lapi := SetupLAPITest(t)
	lapi.InsertAlertFromFile("./tests/alert_ssh-bf.json")
	alertContent := GetAlertReaderFromFile("./tests/alert_ssh-bf.json")

	//bad filter

	w := lapi.RecordResponse("GET", "/v1/alerts?test=test", alertContent)
	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "{\"message\":\"Filter parameter 'test' is unknown (=test): invalid filter\"}", w.Body.String())

	//get without filters

	w = lapi.RecordResponse("GET", "/v1/alerts", emptyBody)
	assert.Equal(t, 200, w.Code)
	//check alert and decision
	assert.Contains(t, w.Body.String(), "Ip 91.121.79.195 performed 'crowdsecurity/ssh-bf' (6 events over ")
	assert.Contains(t, w.Body.String(), `scope":"Ip","simulated":false,"type":"ban","value":"91.121.79.195"`)

	//test decision_type filter (ok)

	w = lapi.RecordResponse("GET", "/v1/alerts?decision_type=ban", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Ip 91.121.79.195 performed 'crowdsecurity/ssh-bf' (6 events over ")
	assert.Contains(t, w.Body.String(), `scope":"Ip","simulated":false,"type":"ban","value":"91.121.79.195"`)

	//test decision_type filter (bad value)

	w = lapi.RecordResponse("GET", "/v1/alerts?decision_type=ratata", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "null", w.Body.String())

	//test scope (ok)

	w = lapi.RecordResponse("GET", "/v1/alerts?scope=Ip", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Ip 91.121.79.195 performed 'crowdsecurity/ssh-bf' (6 events over ")
	assert.Contains(t, w.Body.String(), `scope":"Ip","simulated":false,"type":"ban","value":"91.121.79.195"`)

	//test scope (bad value)

	w = lapi.RecordResponse("GET", "/v1/alerts?scope=rarara", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "null", w.Body.String())

	//test scenario (ok)

	w = lapi.RecordResponse("GET", "/v1/alerts?scenario=crowdsecurity/ssh-bf", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Ip 91.121.79.195 performed 'crowdsecurity/ssh-bf' (6 events over ")
	assert.Contains(t, w.Body.String(), `scope":"Ip","simulated":false,"type":"ban","value":"91.121.79.195"`)

	//test scenario (bad value)

	w = lapi.RecordResponse("GET", "/v1/alerts?scenario=crowdsecurity/nope", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "null", w.Body.String())

	//test ip (ok)

	w = lapi.RecordResponse("GET", "/v1/alerts?ip=91.121.79.195", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Ip 91.121.79.195 performed 'crowdsecurity/ssh-bf' (6 events over ")
	assert.Contains(t, w.Body.String(), `scope":"Ip","simulated":false,"type":"ban","value":"91.121.79.195"`)

	//test ip (bad value)

	w = lapi.RecordResponse("GET", "/v1/alerts?ip=99.122.77.195", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "null", w.Body.String())

	//test ip (invalid value)

	w = lapi.RecordResponse("GET", "/v1/alerts?ip=gruueq", emptyBody)
	assert.Equal(t, 500, w.Code)
	assert.Equal(t, `{"message":"unable to convert 'gruueq' to int: invalid address: invalid ip address / range"}`, w.Body.String())

	//test range (ok)

	w = lapi.RecordResponse("GET", "/v1/alerts?range=91.121.79.0/24&contains=false", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Ip 91.121.79.195 performed 'crowdsecurity/ssh-bf' (6 events over ")
	assert.Contains(t, w.Body.String(), `scope":"Ip","simulated":false,"type":"ban","value":"91.121.79.195"`)

	//test range

	w = lapi.RecordResponse("GET", "/v1/alerts?range=99.122.77.0/24&contains=false", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "null", w.Body.String())

	//test range (invalid value)

	w = lapi.RecordResponse("GET", "/v1/alerts?range=ratata", emptyBody)
	assert.Equal(t, 500, w.Code)
	assert.Equal(t, `{"message":"unable to convert 'ratata' to int: invalid address: invalid ip address / range"}`, w.Body.String())

	//test since (ok)

	w = lapi.RecordResponse("GET", "/v1/alerts?since=1h", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Ip 91.121.79.195 performed 'crowdsecurity/ssh-bf' (6 events over ")
	assert.Contains(t, w.Body.String(), `scope":"Ip","simulated":false,"type":"ban","value":"91.121.79.195"`)

	//test since (ok but yields no results)

	w = lapi.RecordResponse("GET", "/v1/alerts?since=1ns", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "null", w.Body.String())

	//test since (invalid value)

	w = lapi.RecordResponse("GET", "/v1/alerts?since=1zuzu", emptyBody)
	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), `{"message":"while parsing duration: time: unknown unit`)

	//test until (ok)

	w = lapi.RecordResponse("GET", "/v1/alerts?until=1ns", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Ip 91.121.79.195 performed 'crowdsecurity/ssh-bf' (6 events over ")
	assert.Contains(t, w.Body.String(), `scope":"Ip","simulated":false,"type":"ban","value":"91.121.79.195"`)

	//test until (ok but no return)

	w = lapi.RecordResponse("GET", "/v1/alerts?until=1m", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "null", w.Body.String())

	//test until (invalid value)

	w = lapi.RecordResponse("GET", "/v1/alerts?until=1zuzu", emptyBody)
	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), `{"message":"while parsing duration: time: unknown unit`)

	//test simulated (ok)

	w = lapi.RecordResponse("GET", "/v1/alerts?simulated=true", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Ip 91.121.79.195 performed 'crowdsecurity/ssh-bf' (6 events over ")
	assert.Contains(t, w.Body.String(), `scope":"Ip","simulated":false,"type":"ban","value":"91.121.79.195"`)

	//test simulated (ok)

	w = lapi.RecordResponse("GET", "/v1/alerts?simulated=false", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Ip 91.121.79.195 performed 'crowdsecurity/ssh-bf' (6 events over ")
	assert.Contains(t, w.Body.String(), `scope":"Ip","simulated":false,"type":"ban","value":"91.121.79.195"`)

	//test has active decision

	w = lapi.RecordResponse("GET", "/v1/alerts?has_active_decision=true", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Ip 91.121.79.195 performed 'crowdsecurity/ssh-bf' (6 events over ")
	assert.Contains(t, w.Body.String(), `scope":"Ip","simulated":false,"type":"ban","value":"91.121.79.195"`)

	//test has active decision

	w = lapi.RecordResponse("GET", "/v1/alerts?has_active_decision=false", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "null", w.Body.String())

	//test has active decision (invalid value)

	w = lapi.RecordResponse("GET", "/v1/alerts?has_active_decision=ratatqata", emptyBody)
	assert.Equal(t, 500, w.Code)
	assert.Equal(t, `{"message":"'ratatqata' is not a boolean: strconv.ParseBool: parsing \"ratatqata\": invalid syntax: unable to parse type"}`, w.Body.String())

}

func TestAlertBulkInsert(t *testing.T) {
	lapi := SetupLAPITest(t)
	//insert a bulk of 20 alerts to trigger bulk insert
	lapi.InsertAlertFromFile("./tests/alert_bulk.json")
	alertContent := GetAlertReaderFromFile("./tests/alert_bulk.json")

	w := lapi.RecordResponse("GET", "/v1/alerts", alertContent)
	assert.Equal(t, 200, w.Code)
}

func TestListAlert(t *testing.T) {
	lapi := SetupLAPITest(t)
	lapi.InsertAlertFromFile("./tests/alert_sample.json")
	// List Alert with invalid filter

	w := lapi.RecordResponse("GET", "/v1/alerts?test=test", emptyBody)
	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "{\"message\":\"Filter parameter 'test' is unknown (=test): invalid filter\"}", w.Body.String())

	// List Alert

	w = lapi.RecordResponse("GET", "/v1/alerts", emptyBody)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "crowdsecurity/test")
}

func TestCreateAlertErrors(t *testing.T) {
	lapi := SetupLAPITest(t)
	alertContent := GetAlertReaderFromFile("./tests/alert_sample.json")

	//test invalid bearer
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/alerts", alertContent)
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", "ratata"))
	lapi.router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	//test invalid bearer
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/v1/alerts", alertContent)
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", lapi.loginResp.Token+"s"))
	lapi.router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

}

func TestDeleteAlert(t *testing.T) {
	lapi := SetupLAPITest(t)
	lapi.InsertAlertFromFile("./tests/alert_sample.json")

	// Fail Delete Alert
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/alerts", strings.NewReader(""))
	AddAuthHeaders(req, lapi.loginResp)
	req.RemoteAddr = "127.0.0.2:4242"
	lapi.router.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code)
	assert.Equal(t, `{"message":"access forbidden from this IP (127.0.0.2)"}`, w.Body.String())

	// Delete Alert
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/v1/alerts", strings.NewReader(""))
	AddAuthHeaders(req, lapi.loginResp)
	req.RemoteAddr = "127.0.0.1:4242"
	lapi.router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"nbDeleted":"1"}`, w.Body.String())
}

func TestDeleteAlertTrustedIPS(t *testing.T) {
	cfg := LoadTestConfig()
	// IPv6 mocking doesn't seem to work.
	// cfg.API.Server.TrustedIPs = []string{"1.2.3.4", "1.2.4.0/24", "::"}
	cfg.API.Server.TrustedIPs = []string{"1.2.3.4", "1.2.4.0/24"}
	cfg.API.Server.ListenURI = "::8080"
	server, err := NewServer(cfg.API.Server)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = server.InitController()
	if err != nil {
		log.Fatal(err.Error())
	}
	router, err := server.Router()
	if err != nil {
		log.Fatal(err.Error())
	}
	loginResp, err := LoginToTestAPI(router)
	if err != nil {
		log.Fatal(err.Error())
	}
	lapi := LAPI{
		router:    router,
		loginResp: loginResp,
		t:         t,
	}

	assertAlertDeleteFailedFromIP := func(ip string) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/v1/alerts", strings.NewReader(""))

		AddAuthHeaders(req, loginResp)
		req.RemoteAddr = ip + ":1234"

		router.ServeHTTP(w, req)
		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), fmt.Sprintf(`{"message":"access forbidden from this IP (%s)"}`, ip))
	}

	assertAlertDeletedFromIP := func(ip string) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/v1/alerts", strings.NewReader(""))
		AddAuthHeaders(req, loginResp)
		req.RemoteAddr = ip + ":1234"

		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, `{"nbDeleted":"1"}`, w.Body.String())
	}

	lapi.InsertAlertFromFile("./tests/alert_sample.json")
	assertAlertDeleteFailedFromIP("4.3.2.1")
	assertAlertDeletedFromIP("1.2.3.4")

	lapi.InsertAlertFromFile("./tests/alert_sample.json")
	assertAlertDeletedFromIP("1.2.4.0")
	lapi.InsertAlertFromFile("./tests/alert_sample.json")
	assertAlertDeletedFromIP("1.2.4.1")
	lapi.InsertAlertFromFile("./tests/alert_sample.json")
	assertAlertDeletedFromIP("1.2.4.255")

	lapi.InsertAlertFromFile("./tests/alert_sample.json")
	assertAlertDeletedFromIP("127.0.0.1")

}
