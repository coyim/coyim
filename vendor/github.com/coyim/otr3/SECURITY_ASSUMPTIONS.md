# Security Assumptions and Caveats

Here we list all the assumptions and caveats from a security perspective for the code in OTR3.

1. This code has not been audited, and there are no guarantees that it will fulfill the security properties of the OTR protocol. (NOTE: this is not true anymore - see link to audit at TODO-add-here)
2. Zeroing `byte` slices wipes the value from memory in the Golang VM.
3. `byte` slices and `big.Int` instances are not likely to be copied to other places in memory by the Golang GC.
4. Assigning 0 to a `big.Int` wipes the previous value from memory. (NOTE: this is not true anymore - we do a stronger kind of wiping)
5. Modular exponentiation and other similar `big.Int` operations don't leak enough timing information to be useful for side channel attacks. (Or OTR provides enough blinding to counter act this). The libotr implementation uses MPIs from libgcrypt, that seem to be implemented in a similar manner to `big.Int` operations. (NOTE: this is not true anymore - the current implementation uses a constant time modular exponentiation operation).
6. Locking of sensitive memory is sufficient to stop that memory from being swapped.
