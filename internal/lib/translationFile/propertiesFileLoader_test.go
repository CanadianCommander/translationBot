package translationFile

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"os"
	"path/filepath"
	"testing"
)

func TestPropertiesFileLoader_Load(t *testing.T) {

	tmpFile, err := os.CreateTemp(os.TempDir(), "TranslationBot-")
	if err != nil {
		t.Fatalf("Failed to create test input file %s", err.Error())
	}
	defer tmpFile.Close()

	err = os.Chdir("/tmp/")
	if err != nil {
		t.Fatalf("Failed to change working directory %s", err.Error())
	}

	dummyPack := git.LanguagePack{
		Name: "dummy",
		TranslationFiles: map[string]string{
			"english": filepath.Base(tmpFile.Name()),
		},
		Project: &git.Project{
			BaseDir: "/",
		},
	}

	written, err := tmpFile.Write([]byte("" +
		"foo.bar.baz= flerp\n" +
		"blip= bazer=bing#comment\n" +
		"# This entire line is a comment \n" +
		"   spaces =\"magic\"" +
		"\n"))
	if err != nil {
		t.Fatalf("Failed to write test properties file %s", err.Error())
	} else if written == 0 {
		t.Fatalf("Failed to write test properties file. No bytes written")
	}

	translations := make(map[string]*Translation)
	propertiesFileLoader := PropertiesFileLoader{}

	_, err = propertiesFileLoader.Load(
		"english",
		[]string{"english"},
		"english",
		&dummyPack,
		translations)
	if err != nil {
		t.Fatalf("Failed to load properties translation file with error %s", err.Error())
	}

	if len(translations) != 3 {
		t.Fatalf(fmt.Sprintf("Expected 3 translations got %d", len(translations)))
	}
	if translations["dummy.spaces"] == nil {
		t.Fatalf("Spaces in path should be stripped")
	}
	if translations["dummy.blip"].SourceValue == " bazer=bing" {
		t.Fatalf("Comments should be stripped")
	}
}
