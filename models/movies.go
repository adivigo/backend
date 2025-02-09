package models

import (
	"context"
	"fmt"
	"latihan_gin/lib"
	"time"

	"github.com/jackc/pgx/v5"
)

type Movie struct {
	Id int `json:"id"`
	Title string `json:"title" form:"title"`
	// Image *multipart.FileHeader `json:"image" form:"image" binding:"required"`
	// Banner *multipart.FileHeader `json:"banner" form:"banner" binding:"required"`
	Tag string `json:"tag" form:"tag"`
	GenreName *string `json:"genreName" form:"genre_name"`
	DirectorName *string `json:"directorName" form:"director_name"`
	CastName *string `json:"castName" form:"cast_name"`
	// ReleaseDate time.Time `json:"releaseDate" form:"release_date"`
	Duration string `json:"duration" form:"duration"`
	Synopsis string `json:"synopsis" form:"synopsis"`
	// CreatedAt time.Time `json:"createdAt" db:"created_at"`
	// UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}

type MovieBody struct {
	Movie
	Image string `json:"image"`
	Banner string `json:"banner"`
	ReleaseDate string `json:"releaseDate" form:"release_date"`
}

type MovieData struct {
	Movie
	Image string `json:"image" form:"image" binding:"required"`
	Banner string `json:"banner" form:"banner" binding:"required"`
	ReleaseDate time.Time `json:"releaseDate" form:"release_date"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}

type Cinemas struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
	Date time.Time `json:"date"`
	Time time.Time `json:"time"`
	ListCity string `json:"listCity"`
}

type Seats struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type ListMovies []MovieData


func FindAllSeats() []Seats {
	conn, erro := lib.DB()
	if erro != nil {
		fmt.Println(erro)
	}
	defer conn.Close(context.Background())

	var rows pgx.Rows
	var err error

	rows, _ = conn.Query(context.Background(), `
			select seats.id, seats.name from seats
		`)

	seats, err := pgx.CollectRows(rows, pgx.RowToStructByName[Seats])
	if err != nil {
		fmt.Println(err)
	}

	return seats
}

func FindAllCinemas() []Cinemas {
	conn, erro := lib.DB()
	if erro != nil {
		fmt.Println(erro)
	}
	defer conn.Close(context.Background())

	var rows pgx.Rows
	var err error

	rows, _ = conn.Query(context.Background(), `
			select cinemas.id, cinemas.name, cinemas.image, cinemas.date, cinemas.time, cinemas.list_city from cinemas
		`)

	cinemas, err := pgx.CollectRows(rows, pgx.RowToStructByName[Cinemas])
	if err != nil {
		fmt.Println(err)
	}

	return cinemas
}

func FindOneCinema(paramId int) Cinemas {
	var cinema Cinemas

	conn, _ := lib.DB()
	defer conn.Close(context.Background())

	err := conn.QueryRow(context.Background(), `
		select cinemas.id, cinemas.name, cinemas.image, cinemas.date, cinemas.time, cinemas.list_city from cinemas
		WHERE cinemas.id = $1
	`, paramId).Scan(
		&cinema.Id,
		&cinema.Name,
		&cinema.Image,
		&cinema.Date,
		&cinema.Time,
		&cinema.ListCity,
	)
	if err != nil {
		fmt.Println(err)
	}
	return cinema
}

func FindAllMovies(search, sortBy, sortOrder string, page, pageLimit int) ListMovies {
	conn, erro := lib.DB()
	if erro != nil {
		fmt.Println(erro)
	}
	defer conn.Close(context.Background())

	var rows pgx.Rows
	var err error

	offset := (page - 1) * pageLimit

	if search != "" {
		rows, err = conn.Query(context.Background(), `
			select movies.id, title, image, banner, tag,  genres.genre_name, directors.director_name, casts.cast_name, release_date, duration, synopsis, movies.created_at, movies.updated_at from movies
			left join movie_genres ON movies.id = movie_genres.movie_id
			left join genres on movie_genres.genre_id = genres.id
			left join movie_directors on movies.id = movie_directors.movie_id
			left join directors on directors.id = movie_directors.director_id 
			left join movie_casts on movies.id = movie_casts.movie_id 
			left join casts on casts.id = movie_casts.cast_id
			WHERE title ILIKE $1
			ORDER BY ` + sortBy + ` ` + sortOrder + `
			LIMIT $2 OFFSET $3
		`, "%"+search+"%", pageLimit, offset)
	} else {
		rows, err = conn.Query(context.Background(), `
			select movies.id, title, image, banner, tag, genres.genre_name, directors.director_name, casts.cast_name, release_date, duration, synopsis, movies.created_at, movies.updated_at from movies
			left join movie_genres ON movies.id = movie_genres.movie_id
			left join genres on movie_genres.genre_id = genres.id
			left join movie_directors on movies.id = movie_directors.movie_id
			left join directors on directors.id = movie_directors.director_id 
			left join movie_casts on movies.id = movie_casts.movie_id 
			left join casts on casts.id = movie_casts.cast_id
			ORDER BY ` + sortBy + ` ` + sortOrder + `
			LIMIT $1 OFFSET $2
		`, pageLimit, offset)
		// rows, err = conn.Query(context.Background(), `
		// SELECT title from movies
		// `)
		// log.Println(err.Error())
	}

	if err != nil {
		fmt.Println(err)
	}

	movies, err := pgx.CollectRows(rows, pgx.RowToStructByName[MovieData])
	if err != nil {
		fmt.Println(err)
	}

	return movies
}

func FindOneMovie(paramId int) MovieData {
	var movie MovieData

	conn, _ := lib.DB()
	defer conn.Close(context.Background())

	err := conn.QueryRow(context.Background(), `
		select movies.id, title, image, banner, tag,  genres.genre_name, directors.director_name, casts.cast_name, release_date, duration, synopsis, movies.created_at, movies.updated_at 
		from movies
		left join movie_genres ON movies.id = movie_genres.movie_id
		left join genres on movie_genres.genre_id = genres.id
		left join movie_directors on movies.id = movie_directors.movie_id
		left join directors on directors.id = movie_directors.director_id 
		left join movie_casts on movies.id = movie_casts.movie_id 
		left join casts on casts.id = movie_casts.cast_id
		WHERE movies.id = $1
	`, paramId).Scan(
		&movie.Id, 
		&movie.Title, 
		&movie.Image, 
		&movie.Banner, 
		&movie.Tag, 
		&movie.GenreName, 
		&movie.DirectorName, 
		&movie.CastName, 
		&movie.ReleaseDate, 
		&movie.Duration,  
		&movie.Synopsis,  
		&movie.CreatedAt, 
		&movie.UpdatedAt)
	if err != nil {
		fmt.Println(err)
	}
	return movie
}

// func FindOneMovieByTitle(title string) Movie {
// 	var movie Movie

// 	conn := lib.DB()
// 	defer conn.Close(context.Background())

// 	conn.QueryRow(context.Background(), `
// 		SELECT id, title, description, created_at, updated_at
// 		FROM movies WHERE title = $1
// 	`, title).Scan(&movie.Id, &movie.Title, &movie.Description, &movie.CreatedAt, &movie.UpdatedAt)
	
// 	return movie
// }

func InsertMovie(movie MovieBody, releaseDate time.Time) MovieData {
	conn, _ := lib.DB()
	defer conn.Close(context.Background())
	
	var newlyCreated MovieData

	conn.QueryRow(context.Background(), `
		INSERT INTO "movies" (title, image, banner, tag, release_date, duration, synopsis) VALUES
		($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, title, image, banner, tag, release_date, duration, synopsis, created_at, updated_at
	`, movie.Title, movie.Image, movie.Banner, movie.Tag, releaseDate, movie.Duration, movie.Synopsis).Scan(
		&newlyCreated.Id, 
		&newlyCreated.Title, 
		&newlyCreated.Image, 
		&newlyCreated.Banner, 
		&newlyCreated.Tag, 
		&newlyCreated.ReleaseDate, 
		&newlyCreated.Duration, 
		&newlyCreated.Synopsis, 
		&newlyCreated.CreatedAt, 
		&newlyCreated.UpdatedAt,
	)

	return newlyCreated
}

func UpdateMovie(movie MovieData) MovieData {
	conn, _ := lib.DB()
	defer conn.Close(context.Background())
	
	var updatedMovie MovieData

	err := conn.QueryRow(context.Background(), `
		UPDATE movies SET title=$1, image=$2, banner=$3, tag=$4, release_date=$5, duration=$6, synopsis=$7 WHERE id = $8
		RETURNING id, title, image, banner, tag, release_date, duration, synopsis, created_at, updated_at
	`, movie.Title, movie.Image, movie.Banner, movie.Tag, movie.ReleaseDate, movie.Duration, movie.Synopsis, movie.Id).Scan(
		&updatedMovie.Id, 
		&updatedMovie.Title, 
		&updatedMovie.Image, 
		&updatedMovie.Banner, 
		&updatedMovie.Tag, 
		&updatedMovie.ReleaseDate, 
		&updatedMovie.Duration, 
		&updatedMovie.Synopsis, 
		&updatedMovie.CreatedAt, 
		&updatedMovie.UpdatedAt,
	)
	fmt.Println(err)
	return updatedMovie
}

func RemoveMovie(id int) (MovieData) {
	conn, _ := lib.DB()
	defer conn.Close(context.Background())

	deletedMovie := MovieData{}

	err := conn.QueryRow(context.Background(),`
		DELETE FROM movies
		WHERE id = $1
		RETURNING id, title, image, banner, tag, release_date, duration, synopsis, created_at, updated_at
	`, id).Scan(
		&deletedMovie.Id,
		&deletedMovie.Title,
		&deletedMovie.Image,
		&deletedMovie.Banner,
		&deletedMovie.Tag,
		&deletedMovie.ReleaseDate,
		&deletedMovie.Duration,
		&deletedMovie.Synopsis,
		&deletedMovie.CreatedAt,
		&deletedMovie.UpdatedAt,
	)
	fmt.Println(err)
	return deletedMovie
}

func CountMovies(search string) int {
	conn, _ := lib.DB()
	defer conn.Close(context.Background())

	var total int

	search = fmt.Sprintf("%%%s%%", search)

	conn.QueryRow(context.Background(), `
		SELECT COUNT(id) FROM movies
		WHERE title ILIKE $1
	`, search).Scan(&total)

	return total
}