# DCK version

This directory contains the construction-kit version of go-cuddlymenu. The original Go sources are preserved at their original paths (revision `75f1c185be4d999e38708bc724ee9ad4a6e4d828`), with small asset accessors so both versions use the same embedded resources.

Run the original with `go run ./cmd/cuddlymenu` and this version with `go run ./dck/cmd/cuddlymenu` from the repository root.

The choreography and assets stay local; reusable rendering and effects live in `../../lib/democonstructionkit`. Second Reality retains its original ST3 music synchronization.
