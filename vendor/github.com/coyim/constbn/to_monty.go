package constbn

/*
 * Convert a modular integer to Montgomery representation. The integer x[]
 * MUST be lower than m[], but with the same announced bit length.
 */
func toMonty(x, m []base) {
	for k := baseLen(m); k > zero; k-- {
		muladdSmall(x, zero, m)
	}
}
