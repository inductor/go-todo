package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGetTodos(t *testing.T) {
	tests := []struct {
		testData   *Todo
		wantBody   string
		wantStatus int
	}{
		{
			testData:   nil,
			wantBody:   "[]",
			wantStatus: 200,
		},
		{
			testData:   &Todo{Todo: "gorilla", Done: true},
			wantBody:   `[{"id":1,"todo":"gorilla","done":true}]`,
			wantStatus: 200,
		},
	}

	for _, tt := range tests {
		// insert data
		if tt.testData != nil {
			if _, err := db.Exec("insert into TODO (todo, done) values (?, ?)", tt.testData.Todo, tt.testData.Done); err != nil {
				t.Fatalf("failed to insert testdata: %s", err)
			}
		}

		r := httptest.NewRecorder()
		getTodos(r)

		gotBody := strings.TrimRight(r.Body.String(), "\n")
		if gotBody != tt.wantBody {
			t.Fatalf("unexpected response body got : %s, want: %s", gotBody, tt.wantBody)
		}
		gotStatus := r.Result().StatusCode
		if gotStatus != tt.wantStatus {
			t.Fatalf("unexpected response status got : %d, want: %d", gotStatus, tt.wantStatus)
		}
	}

	t.Cleanup(func() {
		db.Exec("delete from TODO")
	})
}
func TestCreateTodo(t *testing.T) {
	tests := []struct {
		requestBody string
		wantStatus  int
	}{
		{
			requestBody: `{"todo":"gorilla","done":true}`,
			wantStatus:  200,
		},
	}
	for _, tt := range tests {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(tt.requestBody))
		createTodo(rec, req)
		gotStatus := rec.Result().StatusCode
		if gotStatus != tt.wantStatus {
			t.Fatalf("unexpected response status got : %d, want: %d, error: %s", gotStatus, tt.wantStatus, rec.Body.String())
		}
	}
	t.Cleanup(func() {
		if err := os.Remove("todo.db"); err != nil {
			t.Fatal(err)
		}
	})
}
