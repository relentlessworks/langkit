package model

import (
	"testing"
)

func TestDetectEnglish(t *testing.T) {
	result := Detect("The quick brown fox jumps over the lazy dog and runs to the river")
	if result.Code != "en" {
		t.Errorf("expected en, got %s", result.Code)
	}
	if result.Confidence <= 0 {
		t.Errorf("expected positive confidence, got %.2f", result.Confidence)
	}
	if result.Script != "Latin" {
		t.Errorf("expected Latin script, got %s", result.Script)
	}
}

func TestDetectFrench(t *testing.T) {
	result := Detect("Le chat est sur le table et le chien est dans le jardin")
	if result.Code != "fr" {
		t.Errorf("expected fr, got %s", result.Code)
	}
}

func TestDetectGerman(t *testing.T) {
	result := Detect("Der Hund ist groß und die Katze ist klein aber beide sind süß")
	if result.Code != "de" {
		t.Errorf("expected de, got %s", result.Code)
	}
}

func TestDetectSpanish(t *testing.T) {
	result := Detect("El gato está en la mesa y el perro está en el suelo")
	if result.Code != "es" {
		t.Errorf("expected es, got %s", result.Code)
	}
}

func TestDetectRussian(t *testing.T) {
	result := Detect("Привет мир, как дела? Я не знаю что сказать в этом тесте")
	if result.Code != "ru" {
		t.Errorf("expected ru, got %s", result.Code)
	}
	if result.Script != "Cyrillic" {
		t.Errorf("expected Cyrillic script, got %s", result.Script)
	}
}

func TestDetectJapanese(t *testing.T) {
	result := Detect("これは日本語のテキストです。こんにちは世界。")
	if result.Code != "ja" {
		t.Errorf("expected ja, got %s", result.Code)
	}
}

func TestDetectKorean(t *testing.T) {
	result := Detect("안녕하세요, 만나서 반갑습니다. 저는 학생입니다.")
	if result.Code != "ko" {
		t.Errorf("expected ko, got %s", result.Code)
	}
	if result.Script != "Hangul" {
		t.Errorf("expected Hangul script, got %s", result.Script)
	}
}

func TestDetectChinese(t *testing.T) {
	result := Detect("你好世界，这是一个中文测试文本。")
	if result.Code != "zh" {
		t.Errorf("expected zh, got %s", result.Code)
	}
	if result.Script != "Han" {
		t.Errorf("expected Han script, got %s", result.Script)
	}
}

func TestDetectArabic(t *testing.T) {
	result := Detect("مرحبا بالعالم هذا نص تجريبي باللغة العربية")
	if result.Code != "ar" {
		t.Errorf("expected ar, got %s", result.Code)
	}
	if result.Script != "Arabic" {
		t.Errorf("expected Arabic script, got %s", result.Script)
	}
}

func TestDetectEmpty(t *testing.T) {
	result := Detect("")
	if result.Code != "und" {
		t.Errorf("expected und, got %s", result.Code)
	}
}

func TestDetectWhitespace(t *testing.T) {
	result := Detect("   \n\t  ")
	if result.Code != "und" {
		t.Errorf("expected und, got %s", result.Code)
	}
}

func TestDetectBatch(t *testing.T) {
	texts := []string{
		"Hello world this is English text",
		"Bonjour le monde ceci est du texte français",
		"Hola mundo este es un texto en español",
	}
	results := DetectBatch(texts)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Result.Code != "en" {
		t.Errorf("expected en for text 0, got %s", results[0].Result.Code)
	}
	if results[1].Result.Code != "fr" {
		t.Errorf("expected fr for text 1, got %s", results[1].Result.Code)
	}
	if results[2].Result.Code != "es" {
		t.Errorf("expected es for text 2, got %s", results[2].Result.Code)
	}
}

func TestLanguageByCode(t *testing.T) {
	lang, ok := LanguageByCode("en")
	if !ok {
		t.Fatal("expected to find en")
	}
	if lang.Name != "English" {
		t.Errorf("expected English, got %s", lang.Name)
	}

	_, ok = LanguageByCode("xx")
	if ok {
		t.Error("expected not to find xx")
	}
}

func TestListLanguageCodes(t *testing.T) {
	codes := ListLanguageCodes()
	if len(codes) != len(SupportedLanguages) {
		t.Errorf("expected %d codes, got %d", len(SupportedLanguages), len(codes))
	}
}

func TestDetectScript(t *testing.T) {
	tests := []struct {
		text   string
		script string
	}{
		{"Hello World", "Latin"},
		{"Привет мир", "Cyrillic"},
		{"こんにちは", "Hiragana"},
		{"안녕하세요", "Hangul"},
		{"你好世界", "Han"},
		{"مرحبا", "Arabic"},
		{"שלום עולם", "Hebrew"},
		{"नमस्ते दुनिया", "Devanagari"},
		{"สวัสดีชาวโลก", "Thai"},
	}

	for _, tt := range tests {
		script, _ := DetectScript(tt.text)
		if script != tt.script {
			t.Errorf("for %q: expected %s, got %s", tt.text, tt.script, script)
		}
	}
}

func TestSupportedScripts(t *testing.T) {
	scripts := SupportedScripts()
	if len(scripts) < 20 {
		t.Errorf("expected at least 20 scripts, got %d", len(scripts))
	}
}
