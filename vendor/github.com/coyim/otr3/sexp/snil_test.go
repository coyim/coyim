package sexp

import "testing"

func Test_Snil_First_ReturnsItself(t *testing.T) {
	assertEquals(t, Snil{}, Snil{}.First())
}

func Test_Snil_Second_ReturnsItself(t *testing.T) {
	assertEquals(t, Snil{}, Snil{}.Second())
}

func Test_Snil_Value_ReturnsNil(t *testing.T) {
	assertEquals(t, nil, Snil{}.Value())
}

func Test_Snil_String_ReturnsAStringRepresentation(t *testing.T) {
	assertEquals(t, "()", Snil{}.String())
}
