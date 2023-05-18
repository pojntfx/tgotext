package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/spf13/cobra"
)

func main() {
	POTHeader := `# SOME DESCRIPTIVE TITLE.
# Copyright (C) YEAR THE PACKAGE'S COPYRIGHT HOLDER
# This file is distributed under the same license as the PACKAGE package.
# FIRST AUTHOR <EMAIL@ADDRESS>, YEAR.
#
#, fuzzy
msgid ""
msgstr ""
"Project-Id-Version: PACKAGE VERSION\n"
"Report-Msgid-Bugs-To: \n"
"POT-Creation-Date: 2022-12-22 23:13+0100\n"
"PO-Revision-Date: YEAR-MO-DA HO:MI+ZONE\n"
"Last-Translator: FULL NAME <EMAIL@ADDRESS>\n"
"Language-Team: LANGUAGE <LL@li.org>\n"
"Language: \n"
"MIME-Version: 1.0\n"
"Content-Type: text/plain; charset=CHARSET\n"
"Content-Transfer-Encoding: 8bit\n"
`
	objName := "Lang"
	rootCmd := &cobra.Command{Use: "tgotext"}

	cmdParse := &cobra.Command{
		Use:   "parse [template file to parse]",
		Short: "Parse a template file for translatable strings",
		Long:  "The given template file is checked for all instances of \"{{ ." + objName + ".Get \"<text>\" }}\".",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			f, err := os.Open(args[0])
			if err != nil {
				panic(err)
			}
			defer f.Close()

			re := regexp.MustCompile(`\{\{\s*.` + objName + `.Get "(.*)"\s*\}\}`)

			fScanner := bufio.NewScanner(f)
			fScanner.Split(bufio.ScanLines)

			printHeader, err := rootCmd.Flags().GetBool("header")
			if printHeader {
				fmt.Println(POTHeader)
			}

			lineNumber := 1
			for fScanner.Scan() {
				stringsFound := re.FindAllStringSubmatch(fScanner.Text(), -1)
				for i := 0; i < len(stringsFound); i++ {
					fmt.Printf("#: %s:%d\nmsgid %q\nmsgstr \"\"\n\n", path.Base(args[0]), lineNumber, stringsFound[i][1])
				}
				lineNumber++
			}
		},
	}
	cmdParse.Flags().StringVarP(&objName, "object", "o", objName, "The name of the Locale object used in the template (without dot prefix!)")
	rootCmd.PersistentFlags().Bool("header", false, "Print POT header")
	rootCmd.AddCommand(cmdParse)

	rootCmd.Execute()
}
