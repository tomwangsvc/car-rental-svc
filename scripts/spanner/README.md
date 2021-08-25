# Usage

```shell
# gcloud config configurations activate svc-dev

export GCP_PROJECT_ID=data-fabric-323905
export SPANNER_INSTANCE_ID=tom-wang-dev
export SPANNER_DATABASE_ID=car-svc
gcloud spanner databases create ${SPANNER_DATABASE_ID} --instance=${SPANNER_INSTANCE_ID}

# ../shell/setenv.sh dev

export ENV=dev
go run migratex.go -env_id=${ENV} -gcp_project_id=${GCP_PROJECT_ID} -spanner_instance_id=${SPANNER_INSTANCE_ID} -spanner_database_id=${SPANNER_DATABASE_ID}
```
