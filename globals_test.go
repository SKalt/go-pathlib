package pathlib_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/skalt/pathlib.go"
)

func ExampleUserHomeDir() {
	homeDir := pathlib.UserHomeDir().Unwrap()
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
		actual := pathlib.UserCacheDir().Unwrap()
		fmt.Printf("$XDG_CACHE_HOME:       : %q\n", expected)
		fmt.Printf("UserCacheDir().Unwrap(): %q\n", actual)
	}

	{ // if $XDG_CACHE_HOME is unset, return the OS-specific default.
		if err := os.Unsetenv("XDG_CACHE_HOME"); err != nil {
			panic(err)
		}
		home := pathlib.UserHomeDir().Unwrap()
		actual := strings.Replace(
			pathlib.UserCacheDir().Unwrap().String(),
			home.String(),
			"$HOME",
			1,
		)
		fmt.Printf("$XDG_CACHE_HOME        : %q\n", os.Getenv("XDG_CACHE_HOME"))
		fmt.Printf("UserCacheDir().Unwrap(): %q\n", actual)

	}

	// Output:
	// On Unix:
	// $XDG_CACHE_HOME:       : "/example/.cache"
	// UserCacheDir().Unwrap(): "/example/.cache"
	// $XDG_CACHE_HOME        : ""
	// UserCacheDir().Unwrap(): "$HOME/.cache"
}

func ExampleUserConfigDir() {
	stash := os.Getenv("XDG_CONFIG_HOME")
	home := pathlib.UserHomeDir().Unwrap()
	defer func() {
		_ = os.Setenv("XDG_CONFIG_HOME", stash)
	}()

	checkOutput := func() {
		xdg_config_home := os.Getenv("XDG_CONFIG_HOME")
		val, err := pathlib.UserConfigDir().Unpack()
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
