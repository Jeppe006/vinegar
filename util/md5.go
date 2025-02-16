package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// VerifyFileMD5 checks the named file with the named MD5 checksum.
func VerifyFileMD5(name string, sum string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	fsum := hex.EncodeToString(h.Sum(nil))

	if sum != fsum {
		return fmt.Errorf("file %s has checksum mismatch", name)
	}

	return nil
}
