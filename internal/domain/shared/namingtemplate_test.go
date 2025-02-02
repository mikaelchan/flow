package shared_test

import (
	"testing"

	"github.com/mikaelchan/hamster/internal/domain/shared"
)

func TestNamingTemplate_NewParser(t *testing.T) {
	tmpl := shared.NamingTemplate("x{title}y{year}z")

	parser, err := shared.NewNamingTemplateParser(tmpl)
	if err != nil {
		t.Fatal(err)
	}
	if parser == nil {
		t.Fatal("invalid naming template")
	}
}

func TestNamingTemplate_Generate(t *testing.T) {
	tmpl := shared.NamingTemplate("TITLE-{title}-year-{year}-z")
	parser, err := shared.NewNamingTemplateParser(tmpl)
	if err != nil {
		t.Fatal(err)
	}
	result := parser.Generate(map[string]string{"title": "test", "year": "2024"})
	if result != "TITLE-test-year-2024-z" {
		t.Fatal("invalid result")
	}

}
