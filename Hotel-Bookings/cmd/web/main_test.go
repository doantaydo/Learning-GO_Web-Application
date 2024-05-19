package main

import "testing"

func TestSetUpAppConfig(t *testing.T) {
	_, err := SetUpAppConfig()
	if err != nil {
		t.Error("Failed SetUpAppConfig()!!")
	}
}
