// Copyright 2015 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.
//
// Author: Marc Berhault (marc@cockroachlabs.com)

package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/cockroachdb/cockroach-prod/drivers"
	"github.com/cockroachdb/cockroach/util"
	"github.com/cockroachdb/cockroach/util/log"
	"github.com/ghemawat/stream"
)

const (
	dockerMachineVersionStringPrefix = "docker-machine version "
	dockerMachineBinary              = "docker-machine"
	dockerMachineStoragePath         = "${HOME}/.docker/machine"
	cockroachNodeName                = `cockroach-%d`
)

var (
	cockroachNodeRegexp = regexp.MustCompile(`^cockroach-([0-9]+)$`)
)

// MakeNodeName generates a cockroach node name for the given ID.
func MakeNodeName(id int) string {
	return fmt.Sprintf(cockroachNodeName, id)
}

// CheckDockerMachine verifies that docker-machine is installed and
// runnable.
func CheckDockerMachine() error {
	cmd := exec.Command(dockerMachineBinary, "-v")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	if !strings.HasPrefix(string(out), dockerMachineVersionStringPrefix) {
		return util.Errorf("bad output %s for docker-machine -v, expected string prefix %q",
			out, dockerMachineVersionStringPrefix)
	}
	return nil
}

// ListMachines returns a list of machine names.
func ListMachines() ([]string, error) {
	return stream.Contents(stream.Command(dockerMachineBinary, "ls", "-q"))
}

// ListCockroachNodes returns a list of machines that are cockroach nodes.
// We could use stream with grep, but let's minimize our dependencies.
// docker-machine is also terrible at proper exit codes.
func ListCockroachNodes() ([]string, error) {
	machines, err := ListMachines()
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for _, mach := range machines {
		if cockroachNodeRegexp.MatchString(mach) {
			ret = append(ret, mach)
		}
	}
	return ret, nil
}

// GetLargestNodeIndex takes a list of node names and returns the largest
// node index seen. Returns 0 if no nodes are passed. Fails on parsing errors.
func GetLargestNodeIndex(nodes []string) (int, error) {
	var largest int
	for _, nodeName := range nodes {
		match := cockroachNodeRegexp.FindStringSubmatch(nodeName)
		if match == nil || len(match) != 2 {
			return -1, util.Errorf("invalid cockroach node name: %s", nodeName)
		}
		index, err := strconv.Atoi(match[1])
		if err != nil {
			return -1, util.Errorf("invalid cockroach node name: %s", nodeName)
		}
		if index > largest {
			largest = index
		}
	}
	return largest, nil
}

// GetMachineConfig gets the machine config from docker-machine.
// It returns unmarshalled json.
func GetMachineConfig(name string) (interface{}, error) {
	contents, err := stream.Contents(stream.Command(dockerMachineBinary, "inspect", name))
	if err != nil {
		return nil, err
	}

	var m interface{}
	err = json.Unmarshal([]byte(strings.Join(contents, "\n")), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// GetDockerFlags returns the list of flags we need to pass to docker to
// talk to the given machine's docker daemon.
// We expect a single line, but then split that line into individual flags.
func GetDockerFlags(name string) ([]string, error) {
	contents, err := stream.Contents(stream.Command(dockerMachineBinary, "config", name))
	if err != nil {
		return nil, err
	}

	if len(contents) != 1 {
		return nil, util.Errorf("expected a single output line, got: %v", contents)
	}
	return strings.Split(contents[0], " "), nil
}

// CreateMachine creates a new docker machine using the passed-in driver
// and name.
func CreateMachine(driver drivers.Driver, name string) error {
	log.Infof("creating docker-machine %s", name)

	args := []string{
		"create",
		"--driver", driver.DockerMachineDriver(),
	}
	args = append(args, driver.DockerMachineCreateArgs()...)
	args = append(args, name)

	log.Infof("running: %s %s", dockerMachineBinary, strings.Join(args, " "))
	cmd := exec.Command(dockerMachineBinary, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// StartMachine invokes "docker-machine start" on the given machine name.
func StartMachine(name string) error {
	log.Infof("starting docker machine %s", name)
	cmd := exec.Command(dockerMachineBinary, "start", name)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// StopMachine invokes "docker-machine stop" on the given machine name.
func StopMachine(name string) error {
	log.Infof("stopping docker machine %s", name)
	cmd := exec.Command(dockerMachineBinary, "stop", name)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
