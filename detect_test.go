package cpython_test

import (
	"errors"
	"testing"

	"github.com/paketo-buildpacks/packit"
	cpython "github.com/paketo-community/cpython"
	"github.com/paketo-community/cpython/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		buildpackYMLParser *fakes.VersionParser

		detect packit.DetectFunc
	)

	it.Before(func() {

		buildpackYMLParser = &fakes.VersionParser{}

		detect = cpython.Detect(buildpackYMLParser)
	})

	it("returns a plan that provides cpython", func() {
		result, err := detect(packit.DetectContext{
			WorkingDir: "/working-dir",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Plan).To(Equal(packit.BuildPlan{
			Provides: []packit.BuildPlanProvision{
				{Name: cpython.Cpython},
			},
			Or: []packit.BuildPlan{
				{
					Provides: []packit.BuildPlanProvision{
						{Name: cpython.Python},
					},
				},
			},
		}))

	})

	context("when the source code contains a buildpack.yml file", func() {
		it.Before(func() {
			buildpackYMLParser.ParseVersionCall.Returns.Version = "some-version"
		})

		it("returns a plan that provides and requires that version of cpython", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: "/working-dir",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Plan).To(Equal(packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: cpython.Cpython},
				},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: cpython.Cpython,
						Metadata: cpython.BuildPlanMetadata{
							Version:       "some-version",
							VersionSource: "buildpack.yml",
						},
					},
				},
				Or: []packit.BuildPlan{
					{
						Provides: []packit.BuildPlanProvision{
							{Name: cpython.Python},
						},
						Requires: []packit.BuildPlanRequirement{
							{
								Name: cpython.Python,
								Metadata: cpython.BuildPlanMetadata{
									Version:       "some-version",
									VersionSource: "buildpack.yml",
								},
							},
						},
					},
				},
			}))
		})
	})

	context("failure cases", func() {
		context("when the buildpack.yml parser fails", func() {
			it.Before(func() {
				buildpackYMLParser.ParseVersionCall.Returns.Err = errors.New("failed to parse buildpack.yml")
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: "/working-dir",
				})
				Expect(err).To(MatchError("failed to parse buildpack.yml"))
			})
		})
	})
}
