package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

type swaggerDocument struct {
	Swagger             string                                `json:"swagger"`
	BasePath            string                                `json:"basePath"`
	Paths               map[string]map[string]json.RawMessage `json:"paths"`
	Definitions         map[string]json.RawMessage            `json:"definitions"`
	SecurityDefinitions map[string]map[string]interface{}     `json:"securityDefinitions"`
}

func TestSwaggerDocumentIsValidAndDocumentsRoutes(t *testing.T) {
	specBytes, err := os.ReadFile("swagger.json")
	if err != nil {
		t.Fatalf("read swagger spec: %v", err)
	}

	var spec swaggerDocument
	if err := json.Unmarshal(specBytes, &spec); err != nil {
		t.Fatalf("swagger spec is invalid JSON: %v", err)
	}

	if spec.Swagger != "2.0" {
		t.Fatalf("swagger version = %q, want 2.0", spec.Swagger)
	}
	if spec.BasePath != "/api/v1" {
		t.Fatalf("basePath = %q, want /api/v1", spec.BasePath)
	}

	expectedRoutes := map[string][]string{
		"/health":           {"get"},
		"/auth/register":    {"post"},
		"/auth/login":       {"post"},
		"/expenses":         {"get", "post"},
		"/expenses/summary": {"get"},
		"/expenses/{id}":    {"get", "put", "delete"},
	}

	for path, methods := range expectedRoutes {
		operations, ok := spec.Paths[path]
		if !ok {
			t.Fatalf("swagger path %q is missing", path)
		}

		for _, method := range methods {
			if _, ok := operations[method]; !ok {
				t.Fatalf("swagger path %q is missing %s operation", path, method)
			}
		}
	}

	for _, definition := range []string{"RegisterRequest", "LoginRequest", "ExpenseRequest", "Expense", "ExpenseSummary"} {
		if _, ok := spec.Definitions[definition]; !ok {
			t.Fatalf("swagger definition %q is missing", definition)
		}
	}

	if _, ok := spec.SecurityDefinitions["UserIDHeader"]; !ok {
		t.Fatal("swagger security definition UserIDHeader is missing")
	}
}

func TestSwaggerUIReferencesSpec(t *testing.T) {
	indexBytes, err := os.ReadFile("index.html")
	if err != nil {
		t.Fatalf("read swagger UI: %v", err)
	}

	if !strings.Contains(string(indexBytes), "swagger.json") {
		t.Fatal("swagger UI does not reference swagger.json")
	}
}
