package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Help example to passing distances as a program argument
const HELP_PASING_DISTANCES_ARG_EXAMPLE = "-distances=100,200.65,-300.47"

// Help message to passing distances as a program argument
const HELP_PASING_DISTANCES_ARG = "Required ordered list of distances to each satelite Kenobi,Skywalker,Sato.\n\t\tPlease use keyword 'distances' with '=' and coma ',' as list separator values.\n\t\texample: cmd " + HELP_PASING_DISTANCES_ARG_EXAMPLE

// Help example to passing messages as a program argument
const HELP_PASING_MESSAGES_ARG_EXAMPLE = "-message=this..the.complete.message,.is.the..message,.is...message"

// Help message for passing messages as a program argument
const HELP_PASING_MESSAGES_ARG = "Required list of messages transmited to each satelite Kenobi,Skywalker,Sato.\n\t\tPlease use keyword 'messages' with '=' and coma ',' as list separator values.\n\t\tAlso use '.' to word separator (don't use empty spaces just '.' instead)\n\t\texample: cmd " + HELP_PASING_MESSAGES_ARG_EXAMPLE

func AskForHelp() (askedForHelp bool) {
	cmdArgs := os.Args
	helpArgRegex := regexp.MustCompile(`(-h)|(help)`)
	askedForHelp = false
	for _, arg := range cmdArgs {
		//checks if asking for help
		if helpArgRegex.MatchString(arg) {
			log.Println("Hello to Operation Fire Quasar.")
			log.Print("\nUsage:\n")
			log.Print("\n\toperation-fire-quasar <arguments>\n")
			log.Print("\nThe arguments are:\n")
			log.Print("\n\t-distances\n")
			log.Print("\t\t" + HELP_PASING_DISTANCES_ARG + "\n")
			log.Print("\n\t-messages\n")
			log.Print("\t\t" + HELP_PASING_MESSAGES_ARG + "\n")
			log.Print("\nexamples:\n")
			log.Printf("\n\toperation-fire-quasar %s %s\n", HELP_PASING_DISTANCES_ARG_EXAMPLE, HELP_PASING_MESSAGES_ARG_EXAMPLE)
			log.Println()
			askedForHelp = true
		}
	}
	return askedForHelp
}

// Searchs the command args to detect if web server profile is present
func IsProfileServerArgPresent() (isPresent bool) {
	cmdArgs := os.Args
	serverArgRegex := regexp.MustCompile(`-profile=server`)
	isPresent = false
	for _, arg := range cmdArgs {
		if serverArgRegex.MatchString(arg) {
			isPresent = true
			return
		}
	}
	return
}

// Parses command args to get distances and messages list
func ParseArgs() (distances []float32, messages [][]string, err error) {
	cmdArgs := os.Args
	distancesArgRegex := regexp.MustCompile(`\bdistances\b`)
	messagesArgRegex := regexp.MustCompile(`\bmessages\b`)
	for i, arg := range cmdArgs {
		if i == 0 {
			continue
		}
		if distancesArgRegex.MatchString(arg) {
			var errParsing error
			distances, errParsing = ParseDistances(arg)
			if errParsing != nil {
				return distances, messages, errParsing
			}
		} else if messagesArgRegex.MatchString(arg) {
			var errParsing error
			messages, errParsing = ParseMessages(arg)
			if errParsing != nil {
				return distances, messages, errParsing
			}
		}
	}

	// checks if distances was populated
	if len(distances) == 0 {
		// no parseable values for distances detected
		return distances, messages, errors.New("no parseable values for distances argument detected.\n\t\t" + HELP_PASING_DISTANCES_ARG)
	}

	if len(messages) == 0 {
		// no parseable values for messages detected
		return distances, messages, errors.New("no parseable values for messages arg detected.\n\t\t" + HELP_PASING_MESSAGES_ARG)
	}
	return distances, messages, nil
}

// Parse messages from arg string
// input: the string argument
// output: the messsages list in []string for each message
// erorr1: if error parsing list is detected
// erorr2: if error parsing message word list is detected
func ParseMessages(arg string) (messages [][]string, err error) {
	//gets separator '=' idx
	separatorIdx := strings.Index(arg, "=")
	if separatorIdx > 0 {
		// split distances list values strings
		values := strings.Split(arg[separatorIdx+1:], ",")

		// checks if at least have 3 values
		if len(values) < 3 {
			return messages, fmt.Errorf("error parsing distances, list values count: %d. Need 3 at least", len(values))
		}

		// intializes messages list
		messages = make([][]string, len(values))

		// parse values
		for i, value := range values {
			// parse value to message []string format
			message := strings.Split(value, ".")
			if len(message) == 0 {
				return messages, fmt.Errorf("error parsing messages, list value: '%s'", value)
			}
			// convert to float32
			messages[i] = message
		}
	}
	return messages, nil
}

// Parse distances from arg string
// input: the string argument
// output: the distances list in float32
// erorr1: if error parsing list is detected
// erorr2: if error parsing floats values is detected
func ParseDistances(arg string) (distances []float32, err error) {
	//gets separator '=' idx
	separatorIdx := strings.Index(arg, "=")
	if separatorIdx > 0 {
		// split distances list values strings
		values := strings.Split(arg[separatorIdx+1:], ",")

		// checks if at least have 3 values
		if len(values) < 3 {
			return distances, fmt.Errorf("error parsing distances, list values count: %d. Need 3 at least", len(values))
		}

		// intializes distances
		distances = make([]float32, len(values))

		// parse values
		for i, value := range values {
			// parse value to float
			valueFloat, parseErr := strconv.ParseFloat(value, 32)
			if parseErr != nil {
				return distances, fmt.Errorf("error parsing distances, list value: %s . %e", values, parseErr)
			}
			// convert to float32
			distances[i] = float32(valueFloat)
		}
	}
	return distances, nil
}
