package gogo

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var jStats jujuStatus

// Spinup will create one cluster
func (j *Juju) Spinup() {
	bootstrap := ""
	tmp := "JUJU_DATA=/tmp/" + j.Name
	fmt.Printf("kind is %s\n", j.Kind)
	if j.Kind == "aws" {
		j.SetAWSCreds()
		bootstrap = j.AwsCl.Region
	} else if j.Kind == "maas" {
		j.SetMAASCloud()
		j.SetMAASCreds()
		bootstrap = j.MaasCl.Type
	}

	cmd := exec.Command("juju", "bootstrap", bootstrap) // with aws this is is expecting region ex - juju bootstrap aws/us-west-2
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	commandResult(out, err, "bootstrap")

	cmd = exec.Command("juju", "deploy", j.Bundle)
	cmd.Env = append(os.Environ(), tmp)
	out, err = cmd.CombinedOutput()
	commandResult(out, err, "deploy")
}

// DisplayStatus will ask juju for status
func (j *Juju) DisplayStatus() {
	tmp := "JUJU_DATA=/tmp/" + j.Name
	cmd := exec.Command("juju", "status")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	commandResult(out, err, "display status")
}

// ClusterReady will check status and return true if cluster is running
func (j *Juju) ClusterReady() bool {
	tmp := "JUJU_DATA=/tmp/" + j.Name
	cmd := exec.Command("juju", "status", "--format=json")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s failed with %s\n", "get cluster deets", err)
	}

	json.Unmarshal([]byte(out), &jStats)

	for k := range jStats.Machines {
		machineStatus := jStats.Machines[k].MachStatus["current"]
		if machineStatus != "started" {
			fmt.Println("Cluster Not Ready")
			return false
		}
	}

	for k := range jStats.ApplicationResults {
		appStatus := jStats.ApplicationResults[k].AppStatus["current"]
		if appStatus != "active" {
			fmt.Println("Cluster Not Ready")
			return false
		}
	}

	fmt.Println("Cluster Ready")
	return true
}

// GetKubeConfig will cat out kubernetes config to stdout
func (j *Juju) GetKubeConfig() {
	tmp := "JUJU_DATA=/tmp/" + j.Name
	cmd := exec.Command("juju", "ssh", "kubernetes-master/0", "cat", "config")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	commandResult(out, err, "get kube config")
}

// DestroyCluster will kill off one cluster
func (j *Juju) DestroyCluster() {
	tmp := "JUJU_DATA=/tmp/" + j.Name
	cmd := exec.Command("juju", "destroy-controller", "--destroy-all-models", "lab", "-y")
	cmd.Env = append(os.Environ(), tmp)
	out, err := cmd.CombinedOutput()
	commandResult(out, err, "destroy-controller")
}

func commandResult(out []byte, err error, command string) {
	fmt.Printf("\n%s\n", string(out))
	if err != nil {
		log.Fatalf("%s failed with %s\n", command, err)
	}
}

// Create is an example of spinning up multiple clusters
func (j *Juju) Create(clusters []string) {
	// clusters := []string{"d8048274-2bc6-49bf-81fd-846aeaddf2fe", "97c19eda-7aeb-4eee-a35c-57dc3755d98f"}

	// for _, cluster := range clusters {
	// 	j.p.wg.Add(1)
	// 	go j.Spinup()
	// }
	// j.p.wg.Wait()
}
