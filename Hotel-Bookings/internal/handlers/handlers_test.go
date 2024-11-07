package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/doantaydo/Learning-GO_Web-Application/Hotel-Bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"new reservations", "/admin/reservations-new", "GET", http.StatusOK},
	{"all reservations", "/admin/reservations-all", "GET", http.StatusOK},
	{"show reservations", "/admin/reservations/new/1/show", "GET", http.StatusOK},
	{"non-existent", "/non/existent", "GET", http.StatusNotFound},
	{"show reservations calendar", "/admin/reservations-calendar", "GET", http.StatusOK},
	{"show reservations calendar with year", "/admin/reservations-calendar?y=2024&month=8", "GET", http.StatusOK},
	{"favicon", "/favicon.ico", "GET", http.StatusOK},
	//-------------------------------------------------------------------------------//
	// {"post-search-availability", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-01"},
	// }, http.StatusOK},
	// {"post-search-availability-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-01"},
	// }, http.StatusOK},
	// {"post-make-reservation", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "John"},
	// 	{key: "last_name", value: "Smith"},
	// 	{key: "email", value: "John@Smith.com"},
	// 	{key: "last_name", value: "1334-113-12323"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, e := range theTests {
		resp, err := testServer.Client().Get(testServer.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("For %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.MakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case when reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := url.Values{}
	reqBody.Add("first_name", "Akihito")
	reqBody.Add("last_name", "Shu")
	reqBody.Add("email", "doantayd@gmail.com")
	reqBody.Add("phone", "2748476277")

	startDate, err := time.Parse("2006-01-02", "2050-01-01")
	if err != nil {
		log.Println(err)
		return
	}
	endDate, err := time.Parse("2006-01-02", "2050-01-02")
	if err != nil {
		log.Println(err)
		return
	}
	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    1,
	}

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case when reservation is not in session (reset everything)
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for empty reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case when form is empty
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for empty form: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case when invalid data
	reqBody = url.Values{}
	reqBody.Add("first_name", "A")
	reqBody.Add("last_name", "Shu")
	reqBody.Add("email", "doantayd@gmail.com")
	reqBody.Add("phone", "2748476277")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code for invalid form: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case fail when insert reservation
	reqBody = url.Values{}
	reqBody.Add("first_name", "Akihito")
	reqBody.Add("last_name", "Shu")
	reqBody.Add("email", "doantayd@gmail.com")
	reqBody.Add("phone", "2748476277")
	reservation = models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    2,
	}

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler failed when trying to fail inserting reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case fail when insert room restriction
	reqBody = url.Values{}
	reqBody.Add("first_name", "Akihito")
	reqBody.Add("last_name", "Shu")
	reqBody.Add("email", "doantayd@gmail.com")
	reqBody.Add("phone", "2748476277")
	reservation = models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    1000,
	}

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler failed when trying to fail inserting room restriction: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestReservationSummary(t *testing.T) {
	startDate, err := time.Parse("2006-01-02", "2050-01-01")
	if err != nil {
		log.Println(err)
		return
	}
	endDate, err := time.Parse("2006-01-02", "2050-01-02")
	if err != nil {
		log.Println(err)
		return
	}
	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case when reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/reservation-summary", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationSummary handler returned wrong response code for empty reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostAvailability(t *testing.T) {
	reqBody := url.Values{}
	reqBody.Add("start", "2050-01-01")
	reqBody.Add("end", "2050-01-02")

	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("PostAvailability handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case can not parse form
	req, _ = http.NewRequest("POST", "/search-availability", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned wrong response code when can not parse form: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case can not parse start date
	reqBody = url.Values{}
	reqBody.Add("start", "invalid")
	reqBody.Add("end", "2050-01-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned wrong response code when can not parse start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case can not parse end date
	reqBody = url.Values{}
	reqBody.Add("start", "2050-01-02")
	reqBody.Add("end", "invalid")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned wrong response code when can not parse end date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case can not connect to database
	reqBody = url.Values{}
	reqBody.Add("start", "2050-01-30")
	reqBody.Add("end", "2050-02-01")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned wrong response code when can not connect to database: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case there is no availability room
	reqBody = url.Values{}
	reqBody.Add("start", "2050-01-02")
	reqBody.Add("end", "2050-02-01")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler returned wrong response code when no availability room: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// test case: rooms are not available
	reqBody := url.Values{}
	reqBody.Add("start", "2050-01-01")
	reqBody.Add("end", "2050-01-01")
	reqBody.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("Failed to parse json")
	}
	if j.Message != "Available!" {
		t.Errorf("AvailabilityJSON handler can not get availability: got %s, wanted %s", j.Message, "Available!")
	} else if j.OK != false {
		t.Errorf("AvailabilityJSON handler gets available room when rooms are not available")
	}

	// test case: rooms are available
	reqBody = url.Values{}
	reqBody.Add("start", "2050-01-01")
	reqBody.Add("end", "2050-01-01")
	reqBody.Add("room_id", "2")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("Failed to parse json")
	}
	if j.Message != "Available!" {
		t.Errorf("AvailabilityJSON handler can not get availability: got %s, wanted %s", j.Message, "Available!")
	} else if j.OK != true {
		t.Errorf("AvailabilityJSON handler gets no available room when rooms are available")
	}

	// test case: SearchAvailabilityByDatesByRoomID get errors
	reqBody = url.Values{}
	reqBody.Add("start", "2050-01-01")
	reqBody.Add("end", "2050-01-01")
	reqBody.Add("room_id", "1000")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("Failed to parse json")
	}
	if j.Message != "Error connecting to database" {
		t.Errorf("AvailabilityJSON handler can not get error when connecting database: got %s, wanted %s", j.Message, "Error connecting to database")
	}

	// test case: can not parse form
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("Failed to parse json")
	}
	if j.Message != "Internal server error" {
		t.Errorf("AvailabilityJSON handler can not get error when parsing form: got %s, wanted %s", j.Message, "Internal server error")
	}
}

func TestChooseRoom(t *testing.T) {
	// test case: normal
	reservation := models.Reservation{
		RoomID: 1,
	}

	req, _ := http.NewRequest("GET", "/choose-room", nil)
	req.RequestURI = "/choose-room/1"
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case: reservation is empty
	req, _ = http.NewRequest("GET", "/choose-room", nil)
	req.RequestURI = "/choose-room/1"
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code when reservation is empty: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case: url param is wrong
	req, _ = http.NewRequest("GET", "/choose-room", nil)
	req.RequestURI = "/choose-room/a"
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code when url param is wrong: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestBookRoom(t *testing.T) {
	req, _ := http.NewRequest("GET", "/book-room?id=1&s=2050-01-01&e=2050-01-02", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case can not get room by ID
	req, _ = http.NewRequest("GET", "/book-room?id=3&s=2050-01-01&e=2050-01-02", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong when can not get room by ID: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid-credentials",
		"admin@admin.com",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"jack@user.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-data",
		"not-email",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

func TestPostShowLogin(t *testing.T) {
	// range through all tests
	for _, e := range loginTests {
		postData := url.Values{}
		postData.Add("email", e.email)
		postData.Add("password", "password")

		// create request
		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("Failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("Failed %s: expected location %s, but got %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("Failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var updateReservationTests = []struct {
	name             string
	src              string
	year             string
	expectedLocation string
}{
	{
		"SRC IS ALL",
		"all",
		"",
		"/admin/reservations-all",
	},
	{
		"SRC IS NEW",
		"new",
		"",
		"/admin/reservations-new",
	},
	{
		"SRC IS CALENDAR",
		"calendar",
		"",
		"/admin/reservations-calendar",
	},
	{
		"HAVE YEAR",
		"",
		"2024",
		"/admin/reservations-calendar?y=2024&m=01",
	},
}

func TestPostShowReservation(t *testing.T) {
	for _, e := range updateReservationTests {
		postData := url.Values{}
		postData.Add("year", e.year)
		if e.year != "" {
			postData.Add("month", "01")
		} else {
			postData.Add("month", "")
		}

		// create request
		req, _ := http.NewRequest("POST", "/admin/reservations/"+e.src+"/1", strings.NewReader(postData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.AdminPostShowReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("Failed %s: expected code %d, but got %d", e.name, http.StatusSeeOther, rr.Code)
		}

		// get the URL from test
		actualLoc, _ := rr.Result().Location()
		if actualLoc.String() != e.expectedLocation {
			t.Errorf("Failed %s: expected location %s, but got %s", e.name, e.expectedLocation, actualLoc.String())
		}
	}

}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
