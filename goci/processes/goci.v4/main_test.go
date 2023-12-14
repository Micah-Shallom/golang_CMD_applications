package main

import (
	"bytes"
	"errors"
	"os/exec"
	"testing"
)

func TestRun(t *testing.T) {
	var testCases = []struct {
		name   string
		proj   string
		out    string
		expErr error
	}{
		{name: "success", proj: "./testdata/tool/",out: "Go Build: SUCCESS\nGofmt: SUCCESS\nGit Push: SUCCESS\n",expErr: nil},
		{name: "fail", proj: "./testdata/toolErr", out: "", expErr: &stepErr{step: "go build"}},
		{name: "failFormat", proj: "./testdata/toolFmtErr/", out: "", expErr: &stepErr{step: "go fmt"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			err := run(tc.proj, &out)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Expected error: %q. Got %q", tc.expErr, err)
					return
				}
				if !errors.Is(err,tc.expErr){
					t.Errorf("Expected error: %q. Got %q.\n", tc.expErr, err)
				}
				return 
			}
			if err != nil {
				t.Errorf("Unexpected error: %q\n", err)
			}
			if out.String() != tc.out{
				t.Errorf("Expected output: %q. Got %q\n", tc.out, out.String())
			}
		})
	}
}


func setupGit(t *testing.T, proj string) func () {
	t.Helper()

	//use the LookPath() function to enquire if git is an installed package
	gitExec, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}

	var gitCMDList = [] struct {
		name	string
		args 	[]string
		dir		string
		env 	[]string
	}{
		{"staging",[]string{"add"}, "./testdata/tool", nil},
		{"commit",[]string{"commit", "-m", "update: testCommit"}, "./testdata/tool", nil},
	}

	for _, g := range gitCMDList {
		g.args = append(g.args, g.dir)
		gitCmd := exec.Command(gitExec, g.args...)
		if err := gitCmd.Run(); err != nil {
			t.Fatal(err)
		}
	}
}