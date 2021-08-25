INSERT INTO car (
  brand_name,
  car_id,
  date_created,
  model_name,
  test
)
values (
  "honda",
  "@CAR_ID_HONDA_ACCORD@",
  PENDING_COMMIT_TIMESTAMP(),
  "accord",
  True
),
(
  "honda",
  "@CAR_ID_HONDA_CIVIC@",
  PENDING_COMMIT_TIMESTAMP(),
  "civic",
  True
);
