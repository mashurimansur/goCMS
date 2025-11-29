-- +goose Up
CREATE TABLE persons(
    id INT AUTO_INCREMENT,
    name VARCHAR(50),
    PRIMARY KEY (id)
);

-- +goose StatementBegin
INSERT INTO persons(name) VALUES ('Huri');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE persons;
-- +goose StatementEnd
