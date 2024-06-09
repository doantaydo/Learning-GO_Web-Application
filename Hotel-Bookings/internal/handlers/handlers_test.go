package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
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
	reqBody := "first_name=Akihito"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Shu")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=doantayd@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=2748476277")

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

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
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
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
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
	reqBody = "first_name=A"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Shu")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=aki.shu@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=2748476277")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
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
	reqBody = "first_name=Akihito"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Shu")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=doantayd@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=2748476277")
	reservation = models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    2,
	}

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
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
	reqBody = "first_name=Akihito"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Shu")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=doantayd@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=2748476277")
	reservation = models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    1000,
	}

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
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
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")

	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
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
	reqBody = "start=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
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
	reqBody = "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=invalid")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
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
	reqBody = "start=2050-01-30"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-02-01")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
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
	reqBody = "start=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-02-01")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
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

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
