// Copyright © 2019 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package builder

import (
	"github.com/homeport/pina-golada/internal/golada/logger"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/homeport/pina-golada/pkg/annotation"
	"github.com/homeport/pina-golada/pkg/files"
	"github.com/homeport/pina-golada/pkg/inspector"
)

func TestBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pgl internal golada builder")
}

// AssetInjector is the injector of the asset injector instance
var AssetInjector AssetProvider

// @pgl(package=builder&injector=AssetInjector)
type AssetProvider interface {
	// @pgl(asset=../../../assets/tests/fileTestFolder/test.txt&compressor=tar&type=file)
	GetMainGoFile() (d files.Directory, err error)

	// @pgl(asset=../../../assets/tests/issue-27/directory-1&compressor=tar)
	GetInfoFileNo1() (d files.Directory, err error)

	// @pgl(asset=../../../assets/tests/issue-27/directory-2&compressor=tar)
	GetInfoFileNo2() (d files.Directory, err error)
}

var _ = Describe("should generate files correctly", func() {

	var (
		l logger.Logger
	)

	_ = BeforeEach(func() {
		l = logger.NewDefaultLogger(&DevNullWriter{}, logger.Debug)
	})

	_ = It("should create a compilable file", func() {
		stream, e := inspector.NewFileStream("./")
		Expect(e).To(BeNil())

		astStream := inspector.NewAstStream(stream.Filter(func(file inspector.File) bool {
			return strings.Contains(file.FileInfo.Name(), "builder_test.go")
		}))
		interfaces := astStream.Find()
		Expect(len(interfaces)).To(BeEquivalentTo(1))

		builder := NewBuilder(interfaces[0], &PinaGoladaInterface{
			Injector: "AssetInjector",
		}, annotation.NewPropertyParser(), l)
		b, e := builder.BuildFile()

		Expect(e).To(BeNil())
		Expect(b).To(Not(BeNil()))

		// Uncomment the following line to print the generated source for for manual inspection:
		// fmt.Println(string(b))
	})

	_ = It("should not create duplicate import entries", func() {
		stream, e := inspector.NewFileStream("./")
		Expect(e).To(BeNil())

		astStream := inspector.NewAstStream(stream.Filter(func(file inspector.File) bool {
			return strings.Contains(file.FileInfo.Name(), "builder_test.go")
		}))
		interfaces := astStream.Find()
		Expect(len(interfaces)).To(BeEquivalentTo(1))

		builder := NewBuilder(
			interfaces[0],
			&PinaGoladaInterface{
				Injector: "AssetInjector",
			},
			annotation.NewPropertyParser(),
			l,
		)

		b, err := builder.BuildFile()
		Expect(err).ToNot(HaveOccurred())
		Expect(b).ToNot(BeNil())

		Expect(strings.Count(string(b), "github.com/homeport/pina-golada/pkg/files/paths")).To(BeEquivalentTo(1))
	})

	_ = It("should not write starting and closing parenthesis for a function with no return type", func() {
		stream, e := inspector.NewFileStream("./")
		Expect(e).To(BeNil())

		astStream := inspector.NewAstStream(stream.Filter(func(file inspector.File) bool {
			return strings.Contains(file.FileInfo.Name(), "builder_test.go")
		}))
		interfaces := astStream.Find()
		Expect(len(interfaces)).To(BeEquivalentTo(1))

		builder := NewBuilder(interfaces[0], &PinaGoladaInterface{
			Injector: "AssetInjector",
		}, annotation.NewPropertyParser(), l)
		b, e := builder.BuildFile()

		Expect(e).To(BeNil())
		Expect(b).To(Not(BeNil()))
		Expect(strings.Count(string(b), "func init() ()")).To(BeEquivalentTo(0))
	})
})
