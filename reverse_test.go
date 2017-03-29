package reverse_test

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kevin-cantwell/reverse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader", func() {
	Context("#Read", func() {
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
	Context("#ReadForward", func() {
		It("Should read forward", func() {
			file, _, _ := createFileFromLines(
				"foo",
				"bar",
				"baz",
			)
			defer file.Close()

			rev := reverse.NewReader(file)

			foo := make([]byte, 4)
			n, err := rev.ReadForward(foo)
			Expect(err).To(BeNil())
			Expect(n).To(Equal(4))
			Expect(string(foo)).To(Equal("foo\n"))

			bar := make([]byte, 4)
			n, err = rev.ReadForward(bar)
			Expect(err).To(BeNil())
			Expect(n).To(Equal(4))
			Expect(string(bar)).To(Equal("bar\n"))

			baz := make([]byte, 4)
			n, err = rev.ReadForward(baz)
			Expect(err).To(Equal(io.EOF))
			Expect(n).To(Equal(3))
			Expect(string(baz[0:3])).To(Equal("baz"))
		})
	})
})

func createFileFromLines(data ...string) (*os.File, os.FileInfo, error) {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, nil, err
	}

	_, err = tmp.Write([]byte(strings.Join(data, "\n")))
	tmp.Close()
	tmp, err = os.Open(tmp.Name()) // Re-open it so the seeker is not at "end" after writing

	info, err := tmp.Stat()
	if err != nil {
		tmp.Close()
		return nil, nil, err
	}

	return tmp, info, err
}
