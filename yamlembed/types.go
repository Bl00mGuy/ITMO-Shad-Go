package yamlembed

import (
	"strings"
)

type Foo struct {
	A string `yaml:"aa"`
	p int64  `yaml:"-"`
}

type Bar struct {
	I      int64    `yaml:"-"`
	B      string   `yaml:"b"`
	UpperB string   `yaml:"-"`
	OI     []string `yaml:"oi,omitempty"`
	F      []any    `yaml:"f,flow"`
}

type Baz struct {
	Foo `yaml:",inline"`
	Bar `yaml:",inline"`
}

func toUpperCase(input string) string {
	return strings.ToUpper(input)
}

func (b *Bar) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Bar
	if err := unmarshal((*plain)(b)); err != nil {
		return err
	}
	b.UpperB = toUpperCase(b.B)
	return nil
}

func (b *Baz) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var foo Foo
	if err := unmarshal(&foo); err != nil {
		return err
	}
	var bar Bar
	if err := unmarshal(&bar); err != nil {
		return err
	}
	b.Foo = foo
	b.Bar = bar
	return nil
}

func initFoo(a string, p int64) Foo {
	return Foo{
		A: a,
		p: p,
	}
}

func initBar(i int64, b string, oi []string, f []any) Bar {
	return Bar{
		I:  i,
		B:  b,
		OI: oi,
		F:  f,
	}
}
