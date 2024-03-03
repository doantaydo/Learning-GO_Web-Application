package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestFormValid(t *testing.T) {
	r := httptest.NewRequest("POST", "/anything", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("Got invalid when should have been valid")
	}
}

func TestFormRequired(t *testing.T) {
	r := httptest.NewRequest("POST", "/anything", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/anything", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Form shows does not have required fields when it does")
	}
}

func TestFormHas(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	if form.Has("anything") {
		t.Error("Form shows field is not blank when field is blank")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	if !form.Has("a") {
		t.Error("Form shows field is blank when field is not blank")
	}
}

func TestFormMinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	if form.MinLength("a", 3) {
		t.Error("Form shows min length for non-existent field")
	}

	isErr := form.Errors.Get("a")
	if isErr == "" {
		t.Error("Should have an error, but did not get one")
	}

	postedData = url.Values{}
	postedData.Add("a", "abc")
	form = New(postedData)

	if !form.MinLength("a", 3) {
		t.Error("Form shows not enough length for enough length field")
	}

	isErr = form.Errors.Get("a")
	if isErr != "" {
		t.Error("Should not have an error, but got one")
	}

	if form.MinLength("a", 4) {
		t.Error("Form does not show min length for not enough length field")
	}
}

func TestIsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("abc")

	if form.Valid() {
		t.Error("Form shows valid for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("a", "mai@gmail.com")
	form = New(postedData)
	form.IsEmail("a")

	if !form.Valid() {
		t.Error("Form shows invalid email for valid email")
	}

	postedData = url.Values{}
	postedData.Add("a", "abc")
	form = New(postedData)
	form.IsEmail("a")

	if form.Valid() {
		t.Error("Form shows valid email for invalid email")
	}
}
