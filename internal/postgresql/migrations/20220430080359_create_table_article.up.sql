CREATE TABLE article (
  id SERIAL PRIMARY KEY,
  author text not null,
  title text not null,
  body text not null,
  created timestamp(3) default current_timestamp(3)
);