package pathlib_test

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/skalt/pathlib.go"
)

func ExampleUserHomeDir() {
	homeDir := expect(pathlib.UserHomeDir())
	fmt.Println(os.Getenv("HOME") == string(homeDir))
	// Output:
	// true
}

func ExampleUserCacheDir() {
	stash := os.Getenv("XDG_CACHE_HOME")
	defer func() {
		_ = os.Setenv("XDG_CACHE_HOME", stash)
	}()
	fmt.Println("On Unix:")

	{ // UserCacheDir() returns $XDG_CACHE_HOME if it's set
		expected := "/example/.cache"
		if err := os.Setenv("XDG_CACHE_HOME", expected); err != nil {
			panic(err)
		}
		actual := expect(pathlib.UserCacheDir())
		fmt.Printf("$XDG_CACHE_HOME:       : %q\n", expected)
		fmt.Printf("UserCacheDir()): %q\n", actual)
	}

	{ // if $XDG_CACHE_HOME is unset, return the OS-specific default.
		if err := os.Unsetenv("XDG_CACHE_HOME"); err != nil {
			panic(err)
		}
		home := expect(pathlib.UserHomeDir())
		actual := strings.Replace(
			expect(pathlib.UserCacheDir()).String(),
			home.String(),
			"$HOME",
			1,
		)
		fmt.Printf("$XDG_CACHE_HOME        : %q\n", os.Getenv("XDG_CACHE_HOME"))
		fmt.Printf("UserCacheDir()): %q\n", actual)

	}

	// Output:
	// On Unix:
	// $XDG_CACHE_HOME:       : "/example/.cache"
	// UserCacheDir()): "/example/.cache"
	// $XDG_CACHE_HOME        : ""
	// UserCacheDir()): "$HOME/.cache"
}

func ExampleUserConfigDir() {
	stash := os.Getenv("XDG_CONFIG_HOME")
	home := expect(pathlib.UserHomeDir())
	defer func() {
		_ = os.Setenv("XDG_CONFIG_HOME", stash)
	}()

	checkOutput := func() {
		xdg_config_home := os.Getenv("XDG_CONFIG_HOME")
		val, err := pathlib.UserConfigDir()
		if err == nil {
			fmt.Printf(
				"%q => %q\n",
				xdg_config_home,
				strings.Replace(val.String(), home.String(), "$HOME", 1),
			)
		} else {
			fmt.Printf("%q => Err(%q)\n", xdg_config_home, err.Error())
		}
	}
	fmt.Println("On Unix")
	_ = os.Setenv("XDG_CONFIG_HOME", "/foo/bar")
	checkOutput()
	_ = os.Unsetenv("XDG_CONFIG_HOME")
	checkOutput()
	_ = os.Setenv("XDG_CONFIG_HOME", "./my_config")
	checkOutput()
	// Output:
	// On Unix
	// "/foo/bar" => "/foo/bar"
	// "" => "$HOME/.config"
	// "./my_config" => Err("path in $XDG_CONFIG_HOME is relative")
}

func TestUserHomeDir(t *testing.T) {
	homeVar := "HOME"
	switch runtime.GOOS {
	case "windows":
		homeVar = "USERPROFILE"
	case "plan9":
		homeVar = "home"
	case "js", "wasip1":
		t.Skip("Skipping on " + runtime.GOOS)
	}
	t.Setenv(homeVar, "")
	_, err := pathlib.UserHomeDir()
	if err == nil {
		t.Fatal("expected error when $HOME is unset")
	}

}
