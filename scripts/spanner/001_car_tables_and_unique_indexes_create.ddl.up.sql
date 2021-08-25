CREATE TABLE car (
  brand_name STRING(1024) NOT NULL,
  car_id STRING(1024) NOT NULL,
  date_created TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
  date_updated TIMESTAMP OPTIONS (allow_commit_timestamp = true),
  model_name STRING(1024),
  test BOOL NOT NULL
) PRIMARY KEY (car_id);


CREATE TABLE customer (
  age INT64 NOT NULL,
  customer_id STRING(1024) NOT NULL,
  date_created TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
  date_updated TIMESTAMP OPTIONS (allow_commit_timestamp = true),
  ethnicity STRING(1024) NOT NULL,
  gender STRING(1024) NOT NULL,
  name STRING(1024) NOT NULL,
  test BOOL NOT NULL
) PRIMARY KEY (customer_id);

CREATE TABLE car_customer_association (
  car_id STRING(1024) NOT NULL,
  customer_id STRING(1024) NOT NULL,
  date_created TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
  date_rental_end TIMESTAMP,
  date_rental_start TIMESTAMP,
  date_updated TIMESTAMP OPTIONS (allow_commit_timestamp = true),
  id STRING(1024) NOT NULL,
  test BOOL NOT NULL
) PRIMARY KEY (id);

CREATE INDEX car_by_brand_name ON car(brand_name);
CREATE INDEX car_customer_association_by_date_rental_end_and_date_rental_start ON car_customer_association(date_rental_end, date_rental_start);
