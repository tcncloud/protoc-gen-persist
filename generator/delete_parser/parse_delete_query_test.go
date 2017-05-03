package delete_parser_test

import (
	"fmt"
	"testing"
	d "github.com/tcncloud/protoc-gen-persist/generator/delete_parser"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDeleteParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "delete parser suite")
}

var _ = Describe("Delete Parser", func() {
	It("can parse", func() {
		p := d.NewParser("DELETE FROM example_table START(?) END(?) KIND(CO)")
		kr, err := p.Expr()

		Expect(err).To(BeNil())
		fmt.Printf("kr: %s\n", *kr)
	})
	It("fails when table is not there", func() {
		p := d.NewParser("DELETE FROM START(?) END(?) KIND(CO)")
		_, err := p.Expr()

		Expect(err).ToNot(BeNil())
		fmt.Printf("err: %s\n", err)
	})
	It("fails when table is int", func() {
		p := d.NewParser("DELETE FROM 12345 START(?) END(?) KIND(CO)")
		_, err := p.Expr()

		Expect(err).ToNot(BeNil())
		fmt.Printf("err: %s\n", err)
	})
	It("fails when table is float", func() {
		p := d.NewParser("DELETE FROM 123.45 START(?) END(?) KIND(CO)")
		_, err := p.Expr()

		Expect(err).ToNot(BeNil())
		fmt.Printf("err: %s\n", err)
	})
	It("fails when table contains a .", func() {
		p := d.NewParser("DELETE FROM hello.world START(?) END(?) KIND(CO)")
		_, err := p.Expr()

		Expect(err).ToNot(BeNil())
		fmt.Printf("err: %s\n", err)
	})

	It("fails when table is a .", func() {
		p := d.NewParser("DELETE FROM . START(?) END(?) KIND(CO)")
		_, err := p.Expr()

		Expect(err).ToNot(BeNil())
		fmt.Printf("err: %s\n", err)
	})

	It("can parse start and end keys of different sized lists", func() {
		p := d.NewParser("DELETE FROM example_table START(?, ?, ?) END(?, ?) KIND(OO)")
		kr, err := p.Expr()

		Expect(err).To(BeNil())
		Expect(len(kr.Start)).To(Equal(3))
		Expect(len(kr.End)).To(Equal(2))
		fmt.Printf("kr: %s\n", *kr)
	})

	It("fails with multiple start declarations", func() {
		p := d.NewParser("DELETE FROM example_table START(?) END(?) START(?) KIND(CO)")
		_, err := p.Expr()

		Expect(err).ToNot(BeNil())
		fmt.Printf("err: %s\n", err)
	})

	It("fails with multiple end declarations", func() {
		p := d.NewParser("DELETE FROM example_table START(?) END(?) END(?) KIND(CO)")
		_, err := p.Expr()

		Expect(err).ToNot(BeNil())
		fmt.Printf("err: %s\n", err)
	})

	It("fails with multiple kind declarations", func() {
		p := d.NewParser("DELETE FROM example_table KIND(OC) START(?) KIND(CO) END(?)")
		_, err := p.Expr()

		Expect(err).ToNot(BeNil())
		fmt.Printf("err: %s\n", err)
	})

	It("parses strings and floats and ints and question marks, and empty keys", func() {
		p := d.NewParser("DELETE FROM example_table START(?, 'test 11 1.2', 1234, 1.234) END() KIND(CC)")
		kr, err := p.Expr()

		Expect(err).To(BeNil())

		Expect(len(kr.Start)).To(Equal(4))
		Expect(kr.Start[0].Type).To(Equal("?"))
		Expect(kr.Start[1].Type).To(Equal("STRING"))
		Expect(kr.Start[2].Type).To(Equal("INT"))
		Expect(kr.Start[3].Type).To(Equal("FLOAT"))

		Expect(len(kr.End)).To(Equal(0))

		fmt.Printf("kr: %s\n", *kr)
	})

	It("parses CO, OC, CC, OO kinds", func() {
		kinds := []string{"CO", "OC", "CC", "OO"}
		expected := map[string]string{
			"CO": "ClosedOpen",
			"OC": "OpenClosed",
			"CC": "ClosedClosed",
			"OO": "OpenOpen",
		}
		for _, k := range kinds {
			p := d.NewParser(fmt.Sprintf("DELETE FROM e START() END() KIND(%s)", k))
			kr, err := p.Expr()
			Expect(err).To(BeNil())
			Expect(kr.Kind).To(Equal(expected[k]))
		}
	})

})

