package cowbull

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("LocalThinker", func() {

	Describe("Think", func() {
		DescribeTable("invalid digit count",
			func(digits int) {
				thinker := LocalThinker(digits)

				_, err := thinker.Think()
				Ω(err).Should(HaveOccurred())
			},
			Entry("negative", -1),
			Entry("zero", 0),
			Entry("greater than ten", 11))

		DescribeTable("valid digit count",
			func(digits, expectedDigits int) {
				thinker := LocalThinker(digits)

				actualDigits, err := thinker.Think()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(actualDigits).Should(Equal(expectedDigits))
			},
			Entry("one", 1, 1),
			Entry("four", 4, 4),
			Entry("ten", 10, 10))

		Context("when perm returns 0 as first element", func() {
			var thinker *AIThinker
			var expectedDigits int
			var actualDigits int
			var err error

			BeforeEach(func() {
				expectedDigits = 4
				perm := func(n int) []int {
					perm := make([]int, n)
					for i := range perm {
						perm[i] = i
					}
					return perm
				}
				thinker = NewLocalThinker(expectedDigits, perm)
				actualDigits, err = thinker.Think()
			})

			It("should not return an error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should return as many digit number as requested", func() {
				Ω(actualDigits).Should(Equal(expectedDigits))
			})

			It("should shuffle the first element", func() {
				Ω(strings.HasPrefix(thinker.number, "0")).To(BeFalse())
			})

			It("should still contain 0", func() {
				Ω(strings.Contains(thinker.number, "0")).To(BeTrue())
			})
		})
	})

	Describe("Try", func() {
		var number string
		var try string
		var thinker *AIThinker
		var cows, bulls int
		var err error

		JustBeforeEach(func() {
			number = "42"
			// directly inject the number into thinker
			thinker = &AIThinker{number: number}
			cows, bulls, err = thinker.Try(try)
		})

		Context("when the input number has digit count mismatch", func() {
			BeforeEach(func() {
				try = "423"
			})

			It("should return an error", func() {
				Ω(err).Should(HaveOccurred())
				Ω(cows).Should(BeZero())
				Ω(bulls).Should(BeZero())
			})
		})

		Context("when the input number is valid", func() {
			DescribeTable("cows bulls",
				func(number string, expectedCows, expectedBulls int) {
					cows, bulls, err := thinker.Try(number)
					Ω(cows).Should(Equal(expectedCows))
					Ω(bulls).Should(Equal(expectedBulls))
					Ω(err).ShouldNot(HaveOccurred())
				},
				Entry("0, 0", "10", 0, 0),
				Entry("1, 0", "21", 1, 0),
				Entry("2, 0", "24", 2, 0),
				Entry("0, 1", "43", 0, 1),
				Entry("0, 2", "42", 0, 2),
			)
		})
	})
})
