package translationFile

import (
	"os"
	"testing"
)

func TestYamlLoader_Load(t *testing.T) {

	tmpFile, err := os.CreateTemp(os.TempDir(), "TranslationBot-")
	if err != nil {
		t.Fatalf("Failed to create test input file %s", err.Error())
	}
	defer tmpFile.Close()

	written, err := tmpFile.Write([]byte("" +
		"zip: zap\n" +
		"foo:\n" +
		"  bar: baz\n" +
		"  biz: bip\n" +
		"alpha: \n" +
		"  omega: two\n" +
		"ben: bob\n"))
	if err != nil {
		t.Fatalf("Failed to write test yaml file %s", err.Error())
	} else if written == 0 {
		t.Fatalf("Failed to write test yaml file. No bytes written")
	}

	translations := make(map[string]*Translation)
	yamlLoader := YamlLoader{}

	_, err = yamlLoader.Load(
		"english",
		[]string{"english"},
		"english",
		tmpFile.Name(),
		translations)
	if err != nil {
		t.Fatalf("Failed to load yaml translation file with error %s", err.Error())
	}

	if len(translations) != 5 {
		t.Fatalf("Expected 5 translations")
	}
}
