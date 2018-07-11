package main

import (
         "fmt"
         "strings"
         "strconv"
         "io"
         "bufio"
         "os"
       )  

func main() {
	if len(os.Args) == 1 { 
		fmt.Println("Please provide one or more words to search.")
	} else {
		data, err := os.Open("UnicodeData.txt")
		if err != nil {
			panic(err)
		}
		for _, c := range Filter(data, strings.Join(os.Args[1:], " ")) {
			fmt.Println(format(c))

		}
	
	}

		
}

func parseLine(line string) (rune, string) {
	fields := strings.Split(line, ";") 
	code := fields[0]
	name := fields[1]

	codeInt, _ := strconv.ParseInt(code, 16, 32)

	return rune(codeInt), name
}

type CharName struct {
	Char rune 
	Name string
}

func Filter(data io.Reader, query string) []CharName {
	scanner := bufio.NewScanner(data)
	results := []CharName{}

	query = strings.ToUpper(query)

	for scanner.Scan(){
		line := scanner.Text()
	    char, name :=  parseLine(line)

	    if strings.Contains(line, query) {
	    	charName := CharName{char, name}
	    	results = append(results, charName)
	    }
	}
	
	return results
}

func format(c CharName) string {
	return fmt.Sprintf("U+%04X\t%c\t%s", c.Char, c.Char, c.Name)
}





