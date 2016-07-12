package Documentation

import (

	"fmt"
	"os"
	"log"
)

func Test(){
	fmt.Printf("Jei")
}

func ExtractDocumentation(path string)  {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 1024)
	_ , err = file.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", b)

	err = file.Close()
	if  err != nil {
		log.Fatal(err)
	}

	
}

