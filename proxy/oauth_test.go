package proxy

import (
  "fmt"
  "testing"
)

func TestPercentEncode(t *testing.T) {
  var tests = []struct {
    input  string
    expect string
  }{
    {input: "Ladies + Gentlemen", expect: "Ladies%20%2B%20Gentlemen"},
    {input: "An encoded string!", expect: "An%20encoded%20string%21"},
    {input: "Dogs, Cats & Mice", expect: "Dogs%2C%20Cats%20%26%20Mice"},
    {input: "☃", expect: "%E2%98%83"},
  }

  for i, tt := range tests {
    name := fmt.Sprintf("Encode_%d", i)
    t.Run(name, func(t *testing.T) {
      result := PercentEncode(tt.input)
      if result != tt.expect {
        t.Errorf("got \"%s\", expected \"%s\"", result, tt.expect)
      }
    })
  }
}
