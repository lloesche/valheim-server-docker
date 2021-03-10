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
	"os"
	"regexp"
	"strings"

	"github.com/golang/glog"
)

func main() {
	envMatch := flag.String("env-match", "VALHEIM_LOG_FILTER_MATCH", "Valheim match filter env varname prefix")
	envPrefix := flag.String("env-startswith", "VALHEIM_LOG_FILTER_STARTSWITH", "Valheim starts-with filter env varname prefix")
	envSuffix := flag.String("env-endswith", "VALHEIM_LOG_FILTER_ENDSWITH", "Valheim ends-with filter env varname prefix")
	envContains := flag.String("env-contains", "VALHEIM_LOG_FILTER_CONTAINS", "Valheim contains filter varname prefix")
	envRegexp := flag.String("env-regexp", "VALHEIM_LOG_FILTER_REGEXP", "Valheim regexp filter varname prefix")
	envFilterEmpty := flag.String("env-empty", "VALHEIM_LOG_FILTER_EMPTY", "Valheim empty-line filter varname")
	flag.Parse()

	if *envMatch == "" || *envPrefix == "" || *envSuffix == "" || *envContains == "" || *envRegexp == "" || *envFilterEmpty == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	var matchFilters []string
	var prefixFilters []string
	var suffixFilters []string
	var containsFilters []string
	var regexpFilters []*regexp.Regexp
	filterEmpty := false

	glog.Info("Configuring Valheim server log filter")
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		envVar := pair[0]
		varValue := pair[1]
		if len(varValue) == 0 {
			continue
		}
		if strings.HasPrefix(envVar, *envMatch) {
			glog.Infof("Removing log lines matching %s", varValue)
			matchFilters = append(matchFilters, varValue)
		} else if strings.HasPrefix(envVar, *envPrefix) {
			glog.Infof("Removing log lines starting with %s", varValue)
			prefixFilters = append(prefixFilters, varValue)
		} else if strings.HasPrefix(envVar, *envSuffix) {
			glog.Infof("Removing log lines ending with %s", varValue)
			suffixFilters = append(suffixFilters, varValue)
		} else if strings.HasPrefix(envVar, *envContains) {
			glog.Infof("Removing log lines containing %s", varValue)
			containsFilters = append(containsFilters, varValue)
		} else if strings.HasPrefix(envVar, *envRegexp) {
			glog.Infof("Removing log lines matching regexp %s", varValue)
			regexpFilters = append(regexpFilters, regexp.MustCompile(varValue))
		} else if envVar == *envFilterEmpty {
			filterEmpty = varValue == "true"
			glog.Infof("Removing empty log lines: %t", filterEmpty)
		}
	}
	glog.Flush()

	scanner := bufio.NewScanner(os.Stdin)
Input:
	for scanner.Scan() {
		if glog.V(10) {
			glog.Infof("Processing line '%s'", scanner.Text())
		}
		if filterEmpty && len(scanner.Text()) == 0 {
			glog.V(5).Info("Skipping empty line")
			continue
		}
		for _, filter := range matchFilters {
			if scanner.Text() == filter {
				if glog.V(5) {
					glog.Infof("Line matched '%s'", filter)
				}
				continue Input
			}
		}
		for _, filter := range prefixFilters {
			if strings.HasPrefix(scanner.Text(), filter) {
				if glog.V(5) {
					glog.Infof("Line matched prefix filter '%s'", filter)
				}
				continue Input
			}
		}
		for _, filter := range suffixFilters {
			if strings.HasSuffix(scanner.Text(), filter) {
				if glog.V(5) {
					glog.Infof("Line matched suffix filter '%s'", filter)
				}
				continue Input
			}
		}
		for _, filter := range containsFilters {
			if strings.Contains(scanner.Text(), filter) {
				if glog.V(5) {
					glog.Infof("Line contains filter %s", filter)
				}
				continue Input
			}
		}
		for _, filter := range regexpFilters {
			if filter.MatchString(scanner.Text()) {
				if glog.V(5) {
					glog.Infof("Line matched regexp filter '%s'", filter)
				}
				continue Input
			}
		}
		if glog.V(8) {
			glog.Info("Line matched no filters")
		}
		glog.Flush()
		fmt.Println(scanner.Text())
	}

	if scanner.Err() != nil {
		glog.Error(scanner.Err())
	}
	glog.Flush()
}
