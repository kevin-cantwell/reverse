package reverse_test

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kevin-cantwell/reverse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader", func() {
	It("Should read in reverse", func() {
		file, info, err := createFileFromLines(
			"foo",
			"bar",
			"baz",
		)
		defer file.Close()

		rev := reverse.NewReader(file)
		offset, err := rev.SeekToEnd()

		Expect(err).To(BeNil())
		Expect(offset).To(Equal(info.Size()))

		// Now we can use a regular scanner to scan lines in reverse!
		scanner := bufio.NewScanner(rev)
		scanner.Scan()
		Expect(scanner.Text()).To(Equal("zab"))
		scanner.Scan()
		Expect(scanner.Text()).To(Equal("rab"))
		scanner.Scan()
		Expect(scanner.Text()).To(Equal("oof"))

		Expect(scanner.Err()).To(BeNil())
	})
})

func createFileFromLines(data ...string) (*os.File, os.FileInfo, error) {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, nil, err
	}

	_, err = tmp.Write([]byte(strings.Join(data, "\n")))

	info, err := tmp.Stat()
	if err != nil {
		tmp.Close()
		return nil, nil, err
	}

	return tmp, info, err
}
