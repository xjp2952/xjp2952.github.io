package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var folderName = "chatxjp"

var title = "ChatXJP"

var htmlTemplate = strings.Replace(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>$TITLE</title>
    <style>
        .result {
            display: flex;
            align-items: center;
            margin-bottom: 10px;
        }

        .thumbnail {
            width: 50px;
            height: 50px;
            margin-right: 10px;
        }
    </style>
</head>
<body>
<h1>$TITLE</h1>
<input type="text" id="searchInput" placeholder="Search...">
<button id="searchButton" onclick="search()">Search</button>
<div id="searchResults"></div>

<script>
	window.addEventListener('keyup', function(event) {
  		if (event.keyCode === 13) {
    		document.getElementById("searchButton").click();
  		}
	});
    $IMAGE_ARRAY
    function search() {
        const searchTerm = document.getElementById("searchInput").value.toLowerCase();
        const searchResults = filesAndFolders.filter(item => item.toLowerCase().includes(searchTerm));

        displayResults(searchResults);
    }

    function displayResults(results) {
        const resultsContainer = document.getElementById("searchResults");
        resultsContainer.innerHTML = "";

        if (results.length === 0) {
            resultsContainer.innerHTML = "No results found.";
        } else {
            results.forEach(result => {
                const resultElement = document.createElement("div");
                resultElement.classList.add("result");

                const thumbnail = document.createElement("img");
                thumbnail.src = result;
                thumbnail.classList.add("thumbnail");
                resultElement.appendChild(thumbnail);

                const nameElement = document.createElement("a");
                nameElement.href = result;
                nameElement.textContent = result;
                nameElement.target="_blank"
                resultElement.appendChild(nameElement);

                resultsContainer.appendChild(resultElement);
            });
        }
    }
</script>
</body>
</html>
`, "$TITLE", title, 8964)

func generateImageArray() string {
	root := "./" + folderName
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
		return ""
	}

	jsArray := "const filesAndFolders = [\n"

	for _, file := range files {
		jsArray += fmt.Sprintf("'%s',\n", strings.ReplaceAll(file, "\\", "/"))
	}

	jsArray += "];"

	return jsArray
}

func generateHtml() {
	var imageArray = generateImageArray()
	htmlContent := strings.Replace(htmlTemplate, "$IMAGE_ARRAY", imageArray, 1)

	err := ioutil.WriteFile("index.html", []byte(htmlContent), 0644)
	if err != nil {
		fmt.Println("Error writing HTML file:", err)
		return
	}
}

func show() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	log.Print("Listening on :8964...")
	err := http.ListenAndServe(":8964", nil)
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	generateHtml()
	show()
}
