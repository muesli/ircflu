package catserver

import "testing"

func Test_ParseFirstLine(t *testing.T) {

    to, msg := ParseFirstLine("a")
    if len(to) != 0 || msg != "a" {
        t.Error("Single char")
    }

    to, msg = ParseFirstLine("#foo hi there")
    if len(to) != 1 || to[0] != "#foo" || msg != "hi there" {
        t.Error("One channel")
    }

    to, msg = ParseFirstLine("#foo,@bar hi there")
    if len(to) != 2 || to[0] != "#foo" || to[1] != "bar" || msg != "hi there" {
        t.Error("Multiple recipients")
    }
}
