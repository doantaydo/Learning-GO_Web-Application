package main

import "testing"

func TestSetUpAppConfig(t *testing.T) {
	err := SetUpAppConfig()
	if err != nil {
		t.Error("Failed SetUpAppConfig()!!")
	}
}
