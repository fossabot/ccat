package mutators

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

func init() {
	simpleRegister("unzstd", unzstd, withDescription("decompress zstd data"))
	simpleRegister("zstd", czstd, withDescription("compress zstd data"), withExpectingBinary(true))
}

func unzstd(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	d, err := zstd.NewReader(in)
	if err != nil {
		return 0, err
	}
	defer d.Close()

	n, err := io.Copy(out, d)
	return n, err
}

func czstd(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	e, err := zstd.NewWriter(out)
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(e, in)
	e.Close()
	return n, err
}
