package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"encoding/json"
	"log"
	"io"
	"net/http"
	"strings"
	"regexp"
)

func main() {
	type Movie struct {
		Url string
		Title string
		Year string
		Rated string
		Released string
		Runtime string
		Genre string
		Director string
		Writer string
		Actors string
		Plot string
		Language string
		Country string
		Awards string
		Poster string
		Ratings []struct {
			Source string
			Value string
		}
		Metascore string
		imdbRating string
		imdbVotes string
		imdbID string
		Type string
		DVD string
		BoxOffice string
		Production string
		Website string
		Response string
	}

	type ImdbMovie struct {
		Year string
		Rated string
		Released string
		Runtime string
		Genre string
		Director string
		Writer string
		Actors string
		Plot string
		Language string
		Country string
		Awards string
		Poster string
		Ratings []struct {
			Source string
			Value string
		}
		Metascore string
		imdbRating string
	}

	var movies []Movie

	kpCollector := colly.NewCollector()

	kpCollector.OnRequest(func(kpr *colly.Request) {
		fmt.Println("Visiting", kpr.URL)
	})

	imdbVisited := false
	
	kpCollector.OnHTML(".list-with-image-section-list ul li", func(kpe *colly.HTMLElement) {
		var movie Movie
		url := kpe.ChildAttr("a", "href")
		title := kpe.ChildText("h4")
		
		if !imdbVisited && strings.Contains(url, "imdb.com") {
			imdbVisited = true

			compiledUrl := regexp.MustCompile(`tt\d+`)
			imdbId := compiledUrl.FindString(url)

			imdbUrl := "https://www.omdbapi.com/?i=" + imdbId + "&apikey=13d6e6df"
			imdbResponse, err := http.Get(imdbUrl)
			
			if err != nil {
				fmt.Println("IMDB Error: ", err)
			}
			
			defer imdbResponse.Body.Close()
			imdbData, err := io.ReadAll(imdbResponse.Body)
			
			if err != nil {
				fmt.Println("IMDB Error: ", err)
			}
			
			var imdbMovie ImdbMovie
			jsonErr := json.Unmarshal(imdbData, &imdbMovie)

			if jsonErr != nil {
				log.Fatal(jsonErr)
			}
			fmt.Printf("%+v\n", url)
			fmt.Printf("%+v\n", title)
			fmt.Printf("%+v\n", imdbMovie)
			movie = Movie{
				Url: url,
				Title: title,
				Year: imdbMovie.Year,
				Rated: imdbMovie.Rated,
				Released: imdbMovie.Released,
				Runtime: imdbMovie.Runtime,
				Genre: imdbMovie.Genre,
				Director: imdbMovie.Director,
				Writer: imdbMovie.Writer,
				Actors: imdbMovie.Actors,
				Plot: imdbMovie.Plot,
				Language: imdbMovie.Language,
				Country: imdbMovie.Country,
				Awards: imdbMovie.Awards,
				Poster: imdbMovie.Poster,
				Ratings: imdbMovie.Ratings,
				Metascore: imdbMovie.Metascore,
				imdbRating: imdbMovie.imdbRating,
			}
			movieJson, err := json.MarshalIndent(movie, "", "    ")
			if err != nil {
				log.Fatal(err)
			}
			
			fmt.Println(string(movieJson))
			movies = append(movies, movie)
		}
	})

	kpCollector.Visit("https://kinopavasaris.lt/en/line-up/")
}
