package treesitter_test

import (
	"ahmedash95/php-lsp-server/pkg/treesitter"
	"testing"
)

func TestGetSymbols(t *testing.T) {

	tests := map[string]struct {
		code     string
		expected []treesitter.Symbol
	}{
		"single variable": {
			code: `<?php
			$foo = 'bar';`,
			expected: []treesitter.Symbol{
				{Name: "foo", Kind: treesitter.Kind_Variable, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 4, OffsetEnd: 7}},
			},
		},
		"multiple variables": {
			code: `<?php
			$foo = 'bar';
			$bar = 'foo';`,
			expected: []treesitter.Symbol{
				{Name: "foo", Kind: treesitter.Kind_Variable, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 4, OffsetEnd: 7}},
				{Name: "bar", Kind: treesitter.Kind_Variable, Position: treesitter.Position{LineStart: 2, LineEnd: 2, OffsetStart: 4, OffsetEnd: 7}},
			},
		},
		"function": {
			code: `<?php
			function foo() {
				return 'bar';
			}`,
			expected: []treesitter.Symbol{
				{Name: "foo", Kind: treesitter.Kind_Function, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 12, OffsetEnd: 15}},
			},
		},
		"multiple functions": {
			code: `<?php
			function foo() {
				return 'bar';
			}
			function bar() {
				return 'foo';
			}`,
			expected: []treesitter.Symbol{
				{Name: "foo", Kind: treesitter.Kind_Function, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 12, OffsetEnd: 15}},
				{Name: "bar", Kind: treesitter.Kind_Function, Position: treesitter.Position{LineStart: 4, LineEnd: 4, OffsetStart: 12, OffsetEnd: 15}},
			},
		},
		"nested functions": {
			code: `<?php
			function foo() {
				function bar() {
					return 'foo';
				}
			},
		}`,
			expected: []treesitter.Symbol{
				{Name: "foo", Kind: treesitter.Kind_Function, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 12, OffsetEnd: 15}},
				{Name: "bar", Kind: treesitter.Kind_Function, Position: treesitter.Position{LineStart: 2, LineEnd: 2, OffsetStart: 13, OffsetEnd: 16}},
			},
		},
		"class": {
			code: `<?php 
			class Foo {}
			`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 9, OffsetEnd: 12}},
			},
		},
		"multiple classes": {
			code: `<?php
			class Foo {}
			class Bar {}
			`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 9, OffsetEnd: 12}},
				{Name: "Bar", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 2, LineEnd: 2, OffsetStart: 9, OffsetEnd: 12}},
			},
		},
		"class with inheritance": {
			code: `<?php
			class Foo extends Bar {}
			`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 9, OffsetEnd: 12}},
				{Name: "Bar", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 21, OffsetEnd: 24}},
			},
		},
		"class implement interface": {
			code: `<?php
			class Foo implements Bar {}
			`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 9, OffsetEnd: 12}},
				{Name: "Bar", Kind: treesitter.Kind_Interface, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 24, OffsetEnd: 27}},
			},
		},
		"class with multiple interfaces": {
			code: `<?php
			class Foo implements Bar, Baz {}
			`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 9, OffsetEnd: 12}},
				{Name: "Bar", Kind: treesitter.Kind_Interface, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 24, OffsetEnd: 27}},
				{Name: "Baz", Kind: treesitter.Kind_Interface, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 29, OffsetEnd: 32}},
			},
		},
		"interface": {
			code: `<?php
			interface Foo {}
			`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Interface, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 13, OffsetEnd: 16}},
			},
		},
		"multiple interfaces": {
			code: `<?php
			interface Foo {}
			interface Bar {}
			`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Interface, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 13, OffsetEnd: 16}},
				{Name: "Bar", Kind: treesitter.Kind_Interface, Position: treesitter.Position{LineStart: 2, LineEnd: 2, OffsetStart: 13, OffsetEnd: 16}},
			},
		},
		"interface with inheritance": {
			code: `<?php
			interface Foo extends Bar {}
			`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Interface, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 13, OffsetEnd: 16}},
				{Name: "Bar", Kind: treesitter.Kind_Interface, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 25, OffsetEnd: 28}},
			},
		},
		"interface with multiple inheritance": {
			code: `<?php
			interface Foo extends Bar, Baz {}
			`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Interface, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 13, OffsetEnd: 16}},
				{Name: "Bar", Kind: treesitter.Kind_Interface, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 25, OffsetEnd: 28}},
				{Name: "Baz", Kind: treesitter.Kind_Interface, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 30, OffsetEnd: 33}},
			},
		},
		"trait": {
			code: `<?php
			trait Foo {}
			`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 9, OffsetEnd: 12}},
			},
		},
		"multiple traits": {
			code: `<?php
			trait Foo {}
			trait Bar {}
			`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 9, OffsetEnd: 12}},
				{Name: "Bar", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 2, LineEnd: 2, OffsetStart: 9, OffsetEnd: 12}},
			},
		},
		"trait with inheritance": {
			code: `<?php
			trait Foo {
				use Bar;
			},
			`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 9, OffsetEnd: 12}},
				{Name: "Bar", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 2, LineEnd: 2, OffsetStart: 8, OffsetEnd: 11}},
			},
		},
		"trait with multiple inheritance": {
			code: `<?php
			trait Foo {
				use Bar, Baz;
			},
			`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 9, OffsetEnd: 12}},
				{Name: "Bar", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 2, LineEnd: 2, OffsetStart: 8, OffsetEnd: 11}},
				{Name: "Baz", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 2, LineEnd: 2, OffsetStart: 13, OffsetEnd: 16}},
			},
		},
		"class method": {
			code: `<?php
			class Foo {
				public function bar() {
					return 'foo';
				}
				public static function baz() {
					return 'bar';
				}
			}
			`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 9, OffsetEnd: 12}},
				{Name: "bar", Kind: treesitter.Kind_Method, Position: treesitter.Position{LineStart: 2, LineEnd: 2, OffsetStart: 20, OffsetEnd: 23}},
				{Name: "baz", Kind: treesitter.Kind_Method, Position: treesitter.Position{LineStart: 5, LineEnd: 5, OffsetStart: 27, OffsetEnd: 30}},
			},
		},
		"method with inheritance and interface and trait": {
			code: `<?php
			class Foo extends Bar implements Baz {
				use Qux;
			}
			`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 9, OffsetEnd: 12}},
				{Name: "Bar", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 21, OffsetEnd: 24}},
				{Name: "Baz", Kind: treesitter.Kind_Interface, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 36, OffsetEnd: 39}},
				{Name: "Qux", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 2, LineEnd: 2, OffsetStart: 8, OffsetEnd: 11}},
			},
		},
		"class property": {
			code: `<?php
			class Foo {
				public $bar = 'foo';
				public static $baz = 'bar';
				public $bax;
				public const QUX = 'qux';
			}
			`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 9, OffsetEnd: 12}},
				{Name: "bar", Kind: treesitter.Kind_Property, Position: treesitter.Position{LineStart: 2, LineEnd: 2, OffsetStart: 12, OffsetEnd: 15}},
				{Name: "baz", Kind: treesitter.Kind_Property, Position: treesitter.Position{LineStart: 3, LineEnd: 3, OffsetStart: 19, OffsetEnd: 22}},
				{Name: "bax", Kind: treesitter.Kind_Property, Position: treesitter.Position{LineStart: 4, LineEnd: 4, OffsetStart: 12, OffsetEnd: 15}},
				{Name: "QUX", Kind: treesitter.Kind_Constant, Position: treesitter.Position{LineStart: 5, LineEnd: 5, OffsetStart: 17, OffsetEnd: 20}},
			},
		},

		"constant definition": {
			code: `<?php
			define('FOO', 'foo');
			const BAR = 'bar';
			`,
			expected: []treesitter.Symbol{
				{Name: "'FOO'", Kind: treesitter.Kind_Constant, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 10, OffsetEnd: 15}},
				{Name: "BAR", Kind: treesitter.Kind_Constant, Position: treesitter.Position{LineStart: 2, LineEnd: 2, OffsetStart: 9, OffsetEnd: 12}},
			},
		},
		"class with method and method body": {
			code: `<?php
			class Foo {
				public function bar() {
					$name = 'Ahmed';
					$github = 'ahmedash95';
				}
			}`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 9, OffsetEnd: 12}},
				{Name: "bar", Kind: treesitter.Kind_Method, Position: treesitter.Position{LineStart: 2, LineEnd: 2, OffsetStart: 20, OffsetEnd: 23}},
			},
		},
		"abstract class": {
			code: `<?php
			abstract class Foo {
				abstract public function bar();
			}`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 18, OffsetEnd: 21}},
				{Name: "bar", Kind: treesitter.Kind_Method, Position: treesitter.Position{LineStart: 2, LineEnd: 2, OffsetStart: 29, OffsetEnd: 32}},
			},
		},
		"abstract function": {
			code: `<?php
			class Foo {
				abstract public function bar();
				abstract protected function doRead(#[\SensitiveParameter] string $sessionId): string;
			}`,
			expected: []treesitter.Symbol{
				{Name: "Foo", Kind: treesitter.Kind_Class, Position: treesitter.Position{LineStart: 1, LineEnd: 1, OffsetStart: 9, OffsetEnd: 12}},
				{Name: "bar", Kind: treesitter.Kind_Method, Position: treesitter.Position{LineStart: 2, LineEnd: 2, OffsetStart: 29, OffsetEnd: 32}},
				{Name: "doRead", Kind: treesitter.Kind_Method, Position: treesitter.Position{LineStart: 3, LineEnd: 3, OffsetStart: 32, OffsetEnd: 38}},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := treesitter.GetDocumentSymbols(tc.code)
			if len(actual) == 0 {
				t.Errorf("Expected GetSymbols to return %d symbols, got %d", len(tc.expected), len(actual))
			}
			if len(actual) != len(tc.expected) {
				t.Errorf("Expected %v, got %v \n %v\n %v\n", len(tc.expected), len(actual), tc.expected, actual)
			}
			for i, _ := range tc.expected {
				if tc.expected[i] != actual[i] {
					t.Errorf("Expected \n%v,\n Got %v,\n code: %v", tc.expected, actual, tc.code)
				}
			}
		})
	}
}
