package mutators

import (
	"io"

	"github.com/batmac/ccat/log"

	"github.com/klauspost/compress/s2"
)

func init() {
	simpleRegister("uns2", uns2, withDescription("decompress s2 data"),
		withCategory("decompress"),
	)
	simpleRegister("s2", cs2, withDescription("compress to s2 data"),
		withCategory("compress"),
	)

	simpleRegister("unsnap", uns2, withDescription("decompress snappy data"),
		withCategory("decompress"),
	)
	simpleRegister("snap", csnappy, withDescription("compress to snappy data"),
		withCategory("compress"),
	)
}

func uns2(out io.WriteCloser, in io.ReadCloser) (int64, error) {
	d := s2.NewReader(in)
	if d == nil {
		log.Fatal("s2 decompressor failed to init")
	}
	n, err := io.Copy(out, d)
	return n, err
}

func cs2(dst io.WriteCloser, src io.ReadCloser) (int64, error) {
	return _cs2(dst, src)
}

func csnappy(dst io.WriteCloser, src io.ReadCloser) (int64, error) {
	return _cs2(dst, src, s2.WriterSnappyCompat())
}

func _cs2(dst io.WriteCloser, src io.ReadCloser, opts ...s2.WriterOption) (int64, error) {
	enc := s2.NewWriter(dst, opts...)
	n, err := io.Copy(enc, src)
	if err != nil {
		enc.Close()
		return 0, err
	}
	// Blocks until compression is done.
	enc.Close()
	return n, nil
}
