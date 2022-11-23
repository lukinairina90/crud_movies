# CRUD Movie List Management Application


### To run an application, in the terminal run the command `make run`. This command will start database, migrations and application containers.
## GET Movies list
```bash
curl --location --request GET 'localhost:8080/movies'
```
## GET Movie
```bash
curl --location --request GET 'http://localhost:8080/movie/1'
```
## CREATE Movie
```bash
curl --location --request POST 'http://localhost:8080/movie' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Побег из Шоушенка",
    "description": "Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.",
    "production_year": 1994,
    "genre": "Drama",
    "actors": "Tim Robbins, Morgan Freeman, Bob Gunton, William Sadler",
    "poster": "https://www.imdb.com/title/tt0111161/mediaviewer/rm10105600/"
}'
```
## UPDATE Movie
```bash
curl --location --request PUT 'http://localhost:8080/movie/1' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Побег из Шоушенка",
    "description": "Some updated description",
    "production_year": 2004,
    "genre": "Comedy",
    "actors": "Tim Robbins, Morgan Freeman, Bob Gunton, William Sadler",
    "poster": "https://www.imdb.com/title/tt0111161/mediaviewer/rm10105600/"
}'
```
## DELETE Movie
```bash
curl --location --request DELETE 'http://localhost:8080/movie/1'
```