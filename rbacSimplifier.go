package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type Rule struct {
	ApiGroups []string `yaml:"apiGroups"`
	Resources []string `yaml:"resources"`
	Verbs     []string `yaml:"verbs"`
}

type Rules struct {
	Rules []Rule `yaml:"rules"`
}

func showHelp() {
	fmt.Println("Usage: go run rbacSimplifier.go [OPTIONS]")
	fmt.Println("Options:")
	fmt.Println("  --input-file FILE    input file")
	fmt.Println("  --help               display this help and exit")
}

func main() {

	inputFile := flag.String("input-file", "", "Input yaml file")
	help := flag.Bool("help", false, "Display help")

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *inputFile == "" {
		log.Println("No input file provided")
		showHelp()
		return
	}

	f, err := os.ReadFile(*inputFile)

	if err != nil {
		log.Fatal(err)
	}

	var rules Rules

	if err := yaml.Unmarshal(f, &rules); err != nil {
		log.Fatal(err)
	}

	normalizedRules := make(map[string]map[string]struct{})
	//let's find all the resources first and initialize the inner map
	for _, rule := range rules.Rules {

		if len(rule.ApiGroups) != 1 {
			log.Fatal("More than one apiGroups per rule is not supported")
		}
		if rule.ApiGroups[0] == string('*') {
			rule.ApiGroups[0] = ""
		}
		for _, resource := range rule.Resources {

			key := fmt.Sprintf("%s.%s", rule.ApiGroups, resource)

			_, ok := normalizedRules[key]

			if !ok {
				normalizedRules[key] = make(map[string]struct{})
			}
		}
	}

	// Let's get the verbs now
	for _, rule := range rules.Rules {

		for _, resource := range rule.Resources {

			key := fmt.Sprintf("%s.%s", rule.ApiGroups, resource)
			for _, verb := range rule.Verbs {

				_, ok := normalizedRules[key][verb]

				if !ok {
					normalizedRules[key][verb] = struct{}{}
				}
			}
		}
	}

	// Short the apiGroups to produce a standard and comparable output
	keys := make([]string, 0, len(normalizedRules))
	for key := range normalizedRules {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var outputRules Rules
	for _, key := range keys {

		parts := strings.Split(key, "[")
		parts = strings.Split(parts[1], "]")

		apiGroup := strings.TrimSpace(parts[0])
		//Remove the . symbol
		resource := strings.TrimSpace(parts[1][1:])

		var verbs []string
		for verb, _ := range normalizedRules[key] {
			verbs = append(verbs, verb)
		}
		// Short the verbs to produce a standard and comparable output
		sort.Strings(verbs)
		outputRules.Rules = append(outputRules.Rules, Rule{ApiGroups: []string{apiGroup}, Resources: []string{resource}, Verbs: verbs})
	}
	out, err := yaml.Marshal(outputRules)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(out))
}

func printMap(toPrintMap map[string]map[string]struct{}, orderedKey []string) {
	for _, key := range orderedKey {
		fmt.Printf("%s %v\n", key, toPrintMap[key])
	}
}
