package tea

import (
	"bytes"
	"strings"
	"testing"
)

func equals(a, b *CLI) bool {
	domain := a.Domain == b.Domain
	port := a.Port == b.Port
	ssl := a.SSL == b.SSL
	p403 := a.Pages.Page403 == b.Pages.Page403
	p404 := a.Pages.Page404 == b.Pages.Page404
	p500 := a.Pages.Page500 == b.Pages.Page500
	md := a.MD == b.MD
	dev := a.Dev == b.Dev
	ui := a.UI == b.UI

	return domain && port && ssl && p403 && p404 && p500 && md && dev && ui
}

func TestPagesStringDefault(t *testing.T) {
	page := &pages{}
	str := page.String()

	if str != "{403: 404: 500:}" {
		t.Error("invalid defaults", str)
	}
}

func TestPageStringSetValid(t *testing.T) {
	page := &pages{}
	err := page.Set("403=403.html,404=./docs/404.html,500=../500.html")

	if err != nil {
		t.Error("valid parsing failure", err)
	}

	str := page.String()
	if str != "{403:403.html 404:./docs/404.html 500:../500.html}" {
		t.Error("invalid expected value", str)
	}
}

func TestNewCLIDefault(t *testing.T) {
	var b bytes.Buffer
	expected := &CLI{
		Domain: "",
		Port:   8080,
		SSL:    "",
		Pages:  &pages{},
		MD:     false,
		HTML5:  false,
		Dev:    false,
		UI:     false,
	}

	values := NewCLI("test", &b)
	if !equals(values, expected) {
		t.Error("invalid expected value", expected, values)
	}
}

func TestNewCLIParse(t *testing.T) {
	var b bytes.Buffer
	expected := &CLI{
		Domain: "example.com",
		Port:   80,
		SSL:    "",
		Pages: &pages{
			Page403: "403.html",
			Page404: "./docs/404.html",
			Page500: "../500.html",
		},
		MD:    true,
		HTML5: true,
		Dev:   true,
		UI:    true,
	}

	values := NewCLI("test", &b)

	args := strings.Split("-domain example.com -port 80 -pages 403=403.html,404=./docs/404.html,500=../500.html -markdown -html5 -dev -ui", " ")
	values.Parse(args)

	if !equals(values, expected) {
		t.Error("invalid expected value", expected, values)
	}
}
