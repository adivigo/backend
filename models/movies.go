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
	GenreName string `json:"genreName" form:"genre_name"`
	DirectorName string `json:"directorName" form:"director_name"`
	CastName string `json:"castName" form:"cast_name"`
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

type ListMovies []Movie

func FindAllMovies(search, sortBy, sortOrder string, page, pageLimit int) []MovieData {
	conn, _ := lib.DB()
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
			select movies.id, title, image, banner, tag,  genres.genre_name, directors.director_name, casts.cast_name, release_date, duration, synopsis, movies.created_at, movies.updated_at from movies
			left join movie_genres ON movies.id = movie_genres.movie_id
			left join genres on movie_genres.genre_id = genres.id
			left join movie_directors on movies.id = movie_directors.movie_id
			left join directors on directors.id = movie_directors.director_id 
			left join movie_casts on movies.id = movie_casts.movie_id 
			left join casts on casts.id = movie_casts.cast_id
			ORDER BY ` + sortBy + ` ` + sortOrder + `
			LIMIT $1 OFFSET $2
		`, pageLimit, offset)
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

// func UpdateMovie(movie Movie) Movie {
// 	conn := lib.DB()
// 	defer conn.Close(context.Background())
	
// 	var updatedMovie Movie

// 	conn.QueryRow(context.Background(), `
// 		UPDATE movies SET title=$1, description=$2 WHERE id = $3
// 		RETURNING id, title, description, created_at, updated_at
// 	`, movie.Title, movie.Description, movie.Id).Scan(
// 		&updatedMovie.Id, 
// 		&updatedMovie.Title, 
// 		&updatedMovie.Description, 
// 		&updatedMovie.CreatedAt, 
// 		&updatedMovie.UpdatedAt,
// 	)
// 	return updatedMovie
// }

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