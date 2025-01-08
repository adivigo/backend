create table users (
	id serial primary key,
	email VARCHAR(255) unique,
	password VARCHAR(255),
	first_name VARCHAR(255),
	last_name VARCHAR(255),
	phone_number VARCHAR(255) default '111111111111',
	image VARCHAR(255),
    role VARCHAR(255),
    point INT,
	created_at timestamp default now(),
	updated_at timestamp
);

INSERT INTO users (email, password) VALUES
('tes@mail.com', '1234'),
('tes1@mail.com', '1234'),
('tes2@mail.com', '1234'),
('tes3@mail.com', '1234'),
('tes4@mail.com', '1234');

create table movies (
    id serial primary key,
    title VARCHAR(255),
    image VARCHAR(255),
    banner VARCHAR(255),
    tag VARCHAR(255),
    release_date DATE,
    duration TIME,
    synopsis TEXT,
    created_at timestamp default now(),
	updated_at timestamp
);

INSERT INTO movies (title, image, banner, tag, release_date, duration, synopsis) VALUES
('Movie 1', 'movie1.jpg', 'banner1.jpg', '-', '2023-01-01', '01:30:00', 'Synopsis of Movie 1'),
('Movie 2', 'movie2.jpg', 'banner2.jpg', 'recommended', '2023-02-01', '02:00:00', 'Synopsis of Movie 2'),
('Movie 3', 'movie3.jpg', 'banner3.jpg', 'recommended', '2023-03-01', '01:45:00', 'Synopsis of Movie 3'),
('Movie 4', 'movie4.jpg', 'banner4.jpg', '-', '2023-04-01', '01:20:00', 'Synopsis of Movie 4'),
('Movie 5', 'movie5.jpg', 'banner5.jpg', '-', '2023-05-01', '02:10:00', 'Synopsis of Movie 5');

create table genres (
    id serial primary key,
    genre_name VARCHAR(255),
    created_at timestamp default now(),
	updated_at timestamp
);

INSERT INTO genres (genre_name) VALUES
('Action'),
('Drama'),
('Comedy'),
('Horror'),
('Sci-Fi');

create table movie_genres (
    id serial primary key,
    movie_id int REFERENCES movies(id) ON DELETE CASCADE,
    genre_id int REFERENCES genres(id) ON DELETE CASCADE,
    created_at timestamp default now(),
	updated_at timestamp
);

INSERT INTO movie_genres (movie_id, genre_id) VALUES
(1, 1),  
(2, 2),  
(3, 3),  
(4, 4),
(5, 5);

create table directors (
    id serial primary key,
    director_name VARCHAR(255),
    created_at timestamp default now(),
	updated_at timestamp
);

INSERT INTO directors (director_name) VALUES
('Director 1'),
('Director 2'),
('Director 3'),
('Director 4'),
('Director 5');

create table movie_directors (
    id serial primary key,
    movie_id int REFERENCES movies(id) on delete cascade,
    director_id int REFERENCES genres(id) on delete cascade,
    created_at timestamp default now(),
	updated_at timestamp
);

INSERT INTO movie_directors (movie_id, director_id) VALUES
(1, 1),
(2, 2),
(3, 3),
(4, 4),
(5, 5);

create table casts (
    id serial primary key,
    cast_name VARCHAR(255),
    created_at timestamp default now(),
	updated_at timestamp
);

INSERT INTO casts (cast_name) VALUES
('Cast 1'),
('Cast 2'),
('Cast 3'),
('Cast 4'),
('Cast 5');

create table movie_casts (
    id serial primary key,
    movie_id int REFERENCES movies(id),
    cast_id int REFERENCES casts(id),
    created_at timestamp default now(),
	updated_at timestamp
);

INSERT INTO movie_casts (movie_id, cast_id) VALUES
(1, 1),
(2, 2),
(3, 3),
(4, 4),
(5, 5);

create table cinemas (
    id serial primary key,
    name VARCHAR(255),
    image VARCHAR(255),
    date DATE,
    time TIME,
    list_city  VARCHAR(255),
    created_at timestamp default now(),
	updated_at timestamp
);

INSERT INTO cinemas (name, image, date, time, list_city) VALUES
('Cinema 1', 'cinema1.jpg', '2023-01-01', '18:00:00', 'City A'),
('Cinema 2', 'cinema2.jpg', '2023-01-02', '19:00:00', 'City B'),
('Cinema 3', 'cinema3.jpg', '2023-01-03', '20:00:00', 'City C'),
('Cinema 4', 'cinema4.jpg', '2023-01-04', '21:00:00', 'City D'),
('Cinema 5', 'cinema5.jpg', '2023-01-05', '22:00:00', 'City E');

create table payments (
    id serial primary key,
    name VARCHAR(255),
    image VARCHAR(255),
    created_at timestamp default now(),
	updated_at timestamp
);

INSERT INTO payments (name, image) VALUES
('Credit Card', 'credit_card.jpg'),
('Debit Card', 'debit_card.jpg'),
('PayPal', 'paypal.jpg'),
('Cash', 'cash.jpg'),
('Gift Card', 'gift_card.jpg');

create table seats (
    id serial primary key,
    name VARCHAR(255),
    created_at timestamp default now(),
	updated_at timestamp
);

INSERT INTO seats (name) VALUES
('A1'),
('A2'),
('A3'),
('A4'),
('A5');

CREATE TABLE orders (
    id serial PRIMARY KEY,
    cinema_id int REFERENCES cinemas(id) ON DELETE CASCADE,
    movie_id int REFERENCES movies(id) ON DELETE CASCADE,
    seats VARCHAR[],
    payment_id int REFERENCES payments(id) ON DELETE CASCADE,
    user_id int REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(255) default 'pending',
    virtual_id int8,
    total_price int,
    expiry_date DATE,
    barcode VARCHAR(255),
    created_at timestamp DEFAULT now(),
    updated_at timestamp
);

INSERT INTO orders (cinema_id, movie_id, seats, payment_id, user_id, status, virtual_id, total_price, expiry_date, barcode)
VALUES 
(1, 1, ARRAY['A1', 'A2', 'A3'], 1, 1, 'confirmed', 30001, 150000, '2023-12-31', '1234567890'),
(1, 2, ARRAY['B1', 'B2'], 2, 2, 'confirmed', 30002, 100000, '2023-12-31', '1234567891'),
(2, 3, ARRAY['C1', 'C2', 'C3', 'C4'], 3, 3, 'pending', 30003, 200000, '2023-12-31', '1234567892'),
(2, 4, ARRAY['D1'], 4, 4, 'confirmed', 30004, 50000, '2023-12-31', '1234567893'),
(3, 5, ARRAY['E1', 'E2', 'E3', 'E4', 'E5'], 5, 5, 'canceled', 30005, 250000, '2023-12-31', '1234567894');