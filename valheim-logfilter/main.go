//   Copyright 2021 Lukas Lösche
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/golang/glog"
)

type PatternAction struct {
	filter string
	cmd    string
}

type RegexpAction struct {
	filter *regexp.Regexp
	cmd    string
}

// valheim-logfilter is a string processor for log lines emitted by Valheim dedicated server.
// This tool is primarily written to remove redundant/broken/uninteresting log lines and it has a
// secondary function running event hooks when something happens in the log. Like notifying on Discord
// when a player logs into the server.
// Match patterns are read from the environment. If a log line matches a pattern the line is either
// removed or a command is executed with the matching log line written to stdin of that command.
// If a command is executed it is run inside a goroutine that is not being waited for. Meaning
// if valheim-logfilter ends (i.e. if Valheim server is stopped) before the command completes
// the command is terminated.
// Also note that this is explicitly written for Valheim server, which emits few log lines per minute.
// If you were to use this for a high throughput log you would want to add rate limiting and pooling
// for command execution. This also means that you might not want to add commands on events that can
// be triggered by unauthenticated users. Like connection attempts for instance.
func main() {
	envMatch := flag.String("env-match", "VALHEIM_LOG_FILTER_MATCH", "Valheim match filter env varname prefix")
	envPrefix := flag.String("env-startswith", "VALHEIM_LOG_FILTER_STARTSWITH", "Valheim starts-with filter env varname prefix")
	envSuffix := flag.String("env-endswith", "VALHEIM_LOG_FILTER_ENDSWITH", "Valheim ends-with filter env varname prefix")
	envContains := flag.String("env-contains", "VALHEIM_LOG_FILTER_CONTAINS", "Valheim contains filter varname prefix")
	envRegexp := flag.String("env-regexp", "VALHEIM_LOG_FILTER_REGEXP", "Valheim regexp filter varname prefix")
	envFilterEmpty := flag.String("env-empty", "VALHEIM_LOG_FILTER_EMPTY", "Valheim empty-line filter varname")
	envFilterUTF8 := flag.String("env-utf8", "VALHEIM_LOG_FILTER_UTF8", "Valheim UTF-8 filter varname")
	flag.Parse()

	if *envMatch == "" || *envPrefix == "" || *envSuffix == "" || *envContains == "" || *envRegexp == "" || *envFilterEmpty == "" || *envFilterUTF8 == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	var matchFilters []PatternAction
	var prefixFilters []PatternAction
	var suffixFilters []PatternAction
	var containsFilters []PatternAction
	var regexpFilters []RegexpAction
	filterEmpty := false
	filterUTF8 := false

	glog.V(1).Info("Configuring Valheim server log filter")
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		envVar := pair[0]
		varValue := pair[1]
		if len(varValue) == 0 {
			continue
		}
		cmd, foundCmdInEnv := os.LookupEnv("ON_" + envVar)
		if strings.HasPrefix(envVar, *envMatch) {
			if foundCmdInEnv {
				glog.V(2).Infof("On log lines matching '%s' running '%s'", varValue, cmd)
			} else {
				glog.V(2).Infof("Removing log lines matching '%s'", varValue)
			}
			matchFilters = append(matchFilters, PatternAction{varValue, cmd})
		} else if strings.HasPrefix(envVar, *envPrefix) {
			if foundCmdInEnv {
				glog.V(2).Infof("On log lines starting with '%s' running '%s'", varValue, cmd)
			} else {
				glog.V(2).Infof("Removing log lines starting with '%s'", varValue)
			}
			prefixFilters = append(prefixFilters, PatternAction{varValue, cmd})
		} else if strings.HasPrefix(envVar, *envSuffix) {
			if foundCmdInEnv {
				glog.V(2).Infof("On log lines ending with '%s' running '%s'", varValue, cmd)
			} else {
				glog.V(2).Infof("Removing log lines ending with '%s'", varValue)
			}
			suffixFilters = append(suffixFilters, PatternAction{varValue, cmd})
		} else if strings.HasPrefix(envVar, *envContains) {
			if foundCmdInEnv {
				glog.V(2).Infof("On log lines containing '%s' running '%s'", varValue, cmd)
			} else {
				glog.V(2).Infof("Removing log lines containing '%s'", varValue)
			}
			containsFilters = append(containsFilters, PatternAction{varValue, cmd})
		} else if strings.HasPrefix(envVar, *envRegexp) {
			if foundCmdInEnv {
				glog.V(2).Infof("On log lines matching regexp '%s' running '%s", varValue, cmd)
			} else {
				glog.V(2).Infof("Removing log lines matching regexp '%s'", varValue)
			}
			regexpFilters = append(regexpFilters, RegexpAction{regexp.MustCompile(varValue), cmd})
		} else if envVar == *envFilterEmpty {
			filterEmpty = varValue == "true"
			glog.V(2).Infof("Removing empty log lines: %t", filterEmpty)
		} else if envVar == *envFilterUTF8 {
			filterUTF8 = varValue == "true"
			glog.V(2).Infof("Removing invalid UTF-8 chars: %t", filterUTF8)
		}
	}
	glog.Flush()

	scanner := bufio.NewScanner(os.Stdin)
Input:
	for scanner.Scan() {
		logLine := scanner.Text()
		if glog.V(10) {
			glog.Infof("Processing line '%s'", logLine)
		}
		if filterEmpty && len(logLine) == 0 {
			if glog.V(5) {
				glog.Info("Skipping empty line")
			}
			continue
		}
		if filterUTF8 && !utf8.ValidString(logLine) {
			if glog.V(5) {
				glog.Info("Removing non UTF-8 character(s)")
			}
			v := make([]rune, 0, len(logLine))
			for i, r := range logLine {
				if r == utf8.RuneError {
					_, size := utf8.DecodeRuneInString(logLine[i:])
					if size == 1 {
						continue
					}
				}
				v = append(v, r)
			}
			logLine = string(v)
		}
		for _, action := range matchFilters {
			if logLine == action.filter {
				if glog.V(5) {
					glog.Infof("Line matched '%s'", action.filter)
				}
				if removeLogLine(action.cmd, logLine) {
					continue Input
				}
			}
		}
		for _, action := range prefixFilters {
			if strings.HasPrefix(logLine, action.filter) {
				if glog.V(5) {
					glog.Infof("Line matched prefix filter '%s'", action.filter)
				}
				if removeLogLine(action.cmd, logLine) {
					continue Input
				}
			}
		}
		for _, action := range suffixFilters {
			if strings.HasSuffix(logLine, action.filter) {
				if glog.V(5) {
					glog.Infof("Line matched suffix filter '%s'", action.filter)
				}
				if removeLogLine(action.cmd, logLine) {
					continue Input
				}
			}
		}
		for _, action := range containsFilters {
			if strings.Contains(logLine, action.filter) {
				if glog.V(5) {
					glog.Infof("Line contains filter '%s'", action.filter)
				}
				if removeLogLine(action.cmd, logLine) {
					continue Input
				}
			}
		}
		for _, action := range regexpFilters {
			if action.filter.MatchString(logLine) {
				if glog.V(5) {
					glog.Infof("Line matched regexp filter '%s'", action.filter)
				}
				if removeLogLine(action.cmd, logLine) {
					continue Input
				}
			}
		}
		if glog.V(8) {
			glog.Info("Line matched no removal filters")
		}
		glog.Flush()
		fmt.Println(logLine)
	}

	if scanner.Err() != nil {
		glog.Error(scanner.Err())
	}
	glog.Flush()
}

// removeLogLine returns true if the cmd arg is an empty string
// or false if it is a non-zero length string. If cmd is not empty
// then the command in it will be executed by runHook() inside of
// a goroutine. We are not waiting for the command to return meaning
// if the Valheim server is being stopped while the command runs
// it will be aborted.
func removeLogLine(cmd string, logLine string) bool {
	if cmd == "" {
		return true
	} else {
		go runHook(cmd, logLine)
	}
	return false
}

// runHook takes a shell command and a log line string as arguments.
// The command will be executed in a bash shell and the log line is
// written to its stdin.
func runHook(cmd string, logLine string) {
	glog.Infof("Running hook %q for %q", cmd, logLine)
	subProcess := exec.Command("/bin/bash", "-c", cmd)
	stdin, err := subProcess.StdinPipe()
	if err != nil {
		glog.Error(err)
	}

	subProcess.Stdout = os.Stdout
	subProcess.Stderr = os.Stderr

	if err = subProcess.Start(); err != nil {
		glog.Error(err)
	}
	glog.Flush()

	io.WriteString(stdin, logLine+"\n")
	stdin.Close()
	subProcess.Wait()
}
