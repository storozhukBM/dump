package dump_test

import (
	"github.com/storozhukBM/dump"
)

type txOffset struct {
	TxName   string
	idx      uint64
	deadline uint64
}

func Example() {
	offset := txOffset{
		TxName:   "Final",
		idx:      34,
		deadline: 16000000000,
	}
	body := "txBody%1"
	hashCode := uint64(9487746)
	codeIsValid := false

	dump.Dump("Tx commit. ", offset, body, hashCode, codeIsValid)
	// Tx commit. offset: `{TxName:Final idx:34 deadline:16000000000}`; body: `txBody%1`; hashCode: `9487746`; codeIsValid: `false`
}
