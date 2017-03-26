package main

import (
	"fmt"
	"testing"
)

func TestSaveAndLoadPage(t *testing.T) {
	title := "TestPage"
	p1 := &Page{Title: title, Body: []byte("This is a sample page.")}
	p1.save()
	p2, err := loadPage(title)
	if err != nil {
		t.Log("something goes wrong")
		t.Fail()
	}
	fmt.Println(string(p2.Body))
}
