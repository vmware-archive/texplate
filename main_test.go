package main_test

import (
	"io/ioutil"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {

	var (
		pathToMain    string
		tempDir       string
		writeTempFile func(string) string
	)

	BeforeSuite(func() {
		var err error
		pathToMain, err = gexec.Build("github.com/ljfranklin/texplate-cli")
		Expect(err).To(BeNil())

		tempDir, err = ioutil.TempDir("", "texplate")
		Expect(err).To(BeNil())
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
		_ = os.RemoveAll(tempDir)
	})

	It("interpolates the input yaml into the template", func() {
		inputFile := writeTempFile(`
input_key: input_value
`)
		templateFile := writeTempFile(`
template_key: {{ .input_key }}
`)

		cmd := exec.Command(pathToMain, "execute", "-f", inputFile, templateFile)
		output, err := cmd.Output()
		Expect(err).To(BeNil(), string(output))

		Expect(string(output)).To(MatchYAML(`
template_key: input_value
`))
	})

	It("interpolates multiple input files", func() {
		inputFile1 := writeTempFile(`
key1: value1
`)
		inputFile2 := writeTempFile(`
key2: value2
`)
		templateFile := writeTempFile(`
first: {{.key1}}
second: {{.key2}}
`)

		cmd := exec.Command(pathToMain, "execute", "-f", inputFile1, "-f", inputFile2, templateFile)
		output, err := cmd.Output()
		Expect(err).To(BeNil(), string(output))

		Expect(string(output)).To(MatchYAML(`
first: value1
second: value2
`))
	})

	It("includes the sprig helpers", func() {
		inputFile := writeTempFile(`
whitespace: "   value   "
`)
		templateFile := writeTempFile(`
trimmed: {{ trim .whitespace }}
`)

		cmd := exec.Command(pathToMain, "execute", "-f", inputFile, templateFile)
		output, err := cmd.Output()
		Expect(err).To(BeNil(), string(output))

		Expect(string(output)).To(MatchYAML(`
trimmed: value
`))
	})

	It("converts the output to JSON", func() {
		inputFile := writeTempFile(`
input_key: input_value
`)
		templateFile := writeTempFile(`
template_key: {{ .input_key }}
`)

		cmd := exec.Command(pathToMain, "execute", "-f", inputFile, "--output-format", "json", templateFile)
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).To(BeNil())
		Eventually(session).Should(gexec.Exit(0))

		Expect(string(session.Out.Contents())).To(MatchJSON(`{
"template_key": "input_value"
}`))
	})

	It("converts the output to YAML", func() {
		inputFile := writeTempFile(`
input_key: input_value
`)
		templateFile := writeTempFile(`{
"template_key": "{{ .input_key }}"
}`)

		cmd := exec.Command(pathToMain, "execute", "-f", inputFile, "--output-format", "yaml", templateFile)
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).To(BeNil())
		Eventually(session).Should(gexec.Exit(0))

		Expect(string(session.Out.Contents())).To(MatchYAML(`
template_key: input_value
`))
	})

	Describe("failure cases", func() {
		It("exits 1 if template arg is not provided", func() {
			cmd := exec.Command(pathToMain, "execute")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			Eventually(session).Should(gexec.Exit(1))
			Expect(session.Err.Contents()).To(ContainSubstring("template path"))
		})

		It("exits 1 if template does not exist", func() {
			cmd := exec.Command(pathToMain, "execute", "path/that/does/not/exist")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			Eventually(session).Should(gexec.Exit(1))
			Expect(session.Err.Contents()).To(ContainSubstring("path/that/does/not/exist"))
		})

		It("exits 1 if input file does not exist", func() {
			templateFile := writeTempFile(`{}`)

			cmd := exec.Command(pathToMain, "execute", templateFile, "-f", "path/that/does/not/exist")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			Eventually(session).Should(gexec.Exit(1))
			Expect(session.Err.Contents()).To(ContainSubstring("path/that/does/not/exist"))
		})

		It("exits 1 if input file is not valid YAML/JSON", func() {
			templateFile := writeTempFile(`{}`)
			invalidInputFile := writeTempFile(`{{{{{`)

			cmd := exec.Command(pathToMain, "execute", templateFile, "-f", invalidInputFile)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			Eventually(session).Should(gexec.Exit(1))
			Expect(session.Err.Contents()).To(ContainSubstring(invalidInputFile))
		})

		It("exits 1 if template is not valid for text/template", func() {
			templateFile := writeTempFile(`{{`)

			cmd := exec.Command(pathToMain, "execute", templateFile)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			Eventually(session).Should(gexec.Exit(1))
			Expect(session.Err.Contents()).To(ContainSubstring(templateFile))
		})

		It("exits 1 if template references an unknown key", func() {
			templateFile := writeTempFile(`{{ .foo }}`)

			cmd := exec.Command(pathToMain, "execute", templateFile)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			Eventually(session).Should(gexec.Exit(1))
			Expect(session.Err.Contents()).To(ContainSubstring(templateFile))
		})

		It("exits 1 if output format is JSON and template is not valid YAML/JSON", func() {
			templateFile := writeTempFile(`:`)

			cmd := exec.Command(pathToMain, "execute", templateFile, "--output-format", "json")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			Eventually(session).Should(gexec.Exit(1))
			Expect(session.Err.Contents()).To(ContainSubstring(templateFile))
		})

		It("exits 1 if output format is YAML and template is not valid YAML/JSON", func() {
			templateFile := writeTempFile(`:`)

			cmd := exec.Command(pathToMain, "execute", templateFile, "--output-format", "yaml")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			Eventually(session).Should(gexec.Exit(1))
			Expect(session.Err.Contents()).To(ContainSubstring(templateFile))
		})

		It("exits 1 if output format is not supported", func() {
			templateFile := writeTempFile(`{}`)

			cmd := exec.Command(pathToMain, "execute", templateFile, "--output-format", "not-valid")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			Eventually(session).Should(gexec.Exit(1))
			Expect(session.Err.Contents()).To(ContainSubstring("not-valid"))
		})
	})

	writeTempFile = func(fileContents string) string {
		tmpFile, err := ioutil.TempFile(tempDir, "")
		Expect(err).To(BeNil())
		defer tmpFile.Close()

		_, err = tmpFile.WriteString(fileContents)
		Expect(err).To(BeNil())

		return tmpFile.Name()
	}
})
