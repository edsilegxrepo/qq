package html

import (
	"testing"
)

func TestHtmlCodec(t *testing.T) {
	codec := &Codec{}

	rawHtml := `
<!DOCTYPE html>
<html>
<head><title>Test Page</title></head>
<body>
	<h1>Main Header</h1>
	<p class="intro">Welcome to test page</p>
	<!-- This is a comment -->
	<ul>
		<li>Item 1</li>
		<li>Item 2</li>
	</ul>
</body>
</html>
`
	var result map[string]any
	err := codec.Unmarshal([]byte(rawHtml), &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	htmlMap, ok := result["html"].(map[string]any)
	if !ok {
		t.Fatalf("expected html key, got %v", result)
	}
	body, ok := htmlMap["body"].(map[string]any)
	if !ok {
		t.Fatalf("expected body tag in html, got %v", htmlMap)
	}
	if body == nil {
		t.Errorf("expected body to not be nil")
	}

	// Empty document
	var emptyResult map[string]any
	if err := codec.Unmarshal([]byte(""), &emptyResult); err != nil {
		t.Fatalf("unexpected error on empty html: %v", err)
	}
}
