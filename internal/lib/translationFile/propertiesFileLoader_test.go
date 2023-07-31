package translationFile

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
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

func TestPropertiesFileLoader_parsePropertiesFile(t *testing.T) {

	t.Run("Regular input", func(t *testing.T) {
		result, err := (&PropertiesFileLoader{}).parsePropertiesFile([]byte(`
foo.bar.baz= flerp
blip= bazer=bing#comment
# This entire line is a comment
spaces =    "magic"  `))

		if err != nil {
			t.Fatalf("Failed to parse properties file with error %s", err.Error())
		}

		if len(result) != 3 {
			t.Fatalf(fmt.Sprintf("Expected 3 translations got %d", len(result)))
		}

		spaceIndex := slices.IndexFunc(result, func(res yaml.MapItem) bool { return res.Key.(string) == "spaces" })

		if spaceIndex == -1 || result[spaceIndex].Value != "\"magic\"" {
			t.Fatalf("Spaces in path should be stripped %s", result[spaceIndex].Value)
		}
	})

	t.Run("Empty input", func(t *testing.T) {
		result, err := (&PropertiesFileLoader{}).parsePropertiesFile([]byte(""))
		if err != nil {
			t.Fatalf("Failed to parse properties file with error %s", err.Error())
		}

		if len(result) != 0 {
			t.Fatalf(fmt.Sprintf("Expected 0 translations got %d", len(result)))
		}
	})

	// input is contains a "string" node. Aka a value that is both a node in the tree and a string.
	t.Run("String node input", func(t *testing.T) {
		result, err := (&PropertiesFileLoader{}).parsePropertiesFile([]byte(`
login=Login
login.bang=Login!`))

		if err != nil {
			t.Fatalf("Failed to parse properties file with error %s", err.Error())
		}

		loginIndex := slices.IndexFunc(result, func(res yaml.MapItem) bool { return res.Key.(string) == "login" })

		if loginIndex == -1 {
			t.Fatalf("Expected to find login key")
		}

		loginStringkeyIdx := slices.IndexFunc(result[loginIndex].Value.(yaml.MapSlice), func(res yaml.MapItem) bool { return res.Key.(string) == StringNodePlaceholder })
		loginBangIdx := slices.IndexFunc(result[loginIndex].Value.(yaml.MapSlice), func(res yaml.MapItem) bool { return res.Key.(string) == "bang" })

		if loginStringkeyIdx == -1 || loginBangIdx == -1 {
			t.Fatalf("Expected to find login keys")
		}

		if result[loginIndex].Value.(yaml.MapSlice)[loginStringkeyIdx].Value != "Login" {
			t.Fatalf("Expected string node value to be Login")
		}

		if result[loginIndex].Value.(yaml.MapSlice)[loginBangIdx].Value != "Login!" {
			t.Fatalf("Expected login.bang value to be Login!")
		}
	})
}
