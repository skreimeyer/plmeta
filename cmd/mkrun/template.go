package main

import (
	"os"
	"path/filepath"
	"text/template"
)

func mkRunScript(l lang) error {
	fname := filepath.Join(l.Path, "run.sh")
	out, err := os.OpenFile(fname,os.O_RDWR|os.O_CREATE,0777)
	if err != nil {
		return (err)
	}
	var runScript string
	if l.Compiled == true {
		runScript = `#!/usr/env bash
for OUTPUT in $(ls | grep -v "\(run.sh\|build.sh\)")
do
echo $OUTPUT [$(date)]>> RESULTS
{ time ./bin/$OUTPUT ; } 2>> RESULTS
done
`
	} else {
		runScript = `#!/usr/env bash
for OUTPUT in $(ls | grep -v "\(run.sh\|build.sh\)")
do
echo $OUTPUT [$(date)]>> RESULTS
{ time {{.RunCmd}} $OUTPUT ; } 2>> RESULTS
done
`
	}
	t := template.Must(template.New("runscript").Parse(runScript))
	err = t.Execute(out, l)
	if err != nil {
		return (err)
	}
	return nil
}

func mkBuild(l lang) error {
	fname := filepath.Join(l.Path, "build.sh")
	out, err := os.OpenFile(fname,os.O_RDWR|os.O_CREATE,0777)
	if err != nil {
		return (err)
	}
	buildScript := `#!/usr/env bash
mkdir bin
for OUTPUT in $(ls | grep -v "\(run.sh\|build.sh\)")
do
{{.BuildCmd}} $OUTPUT
done
`
	t := template.Must(template.New("build").Parse(buildScript))
	err = t.Execute(out, l)
	if err != nil {
		return (err)
	}
	return nil
}
