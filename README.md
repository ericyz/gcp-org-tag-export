## Prerequisite
Below are the tools/binary required to get installed beforehand
* [Git](https://git-scm.com/)
* [Golang](https://golang.org/)
* [GCP CLI](https://cloud.google.com/sdk/docs/quickstarts/)

Git Repos for cft cli and policy library.
* [CFT cli](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit)
* [Policy Library](https://github.com/ericyz/policy-library.git)

Run the following commands to clone above repos and build the cft binary.
```
git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
git clone --branch feature/gcp-lables https://github.com/ericyz/policy-library.git
(cd cloud-foundation-toolkit/cli && make build)
```

## Create folders
The following commands will create the folders to store the GCP label analysis report.
```
mkdir -p inventory-reports
```

## Set the environment variables
The following commands will set the environments variables used in the gcloud commands.
```
PROJECT=YOUR_PROJECT_ID
BUCKET=YOUR_CAI_BUCKET
REGION=YOUR_DEFAULT_REGION
ORG_ID=YOUR_ORG_ID
USER_EMAIL=$(gcloud config list account --format "value(core.account)")
```

## Bucket for CAI
```
gsutil mb -l $REGION -p $PROJECT gs://$BUCKET
```

## API and Permissions
The following commands will enable the Cloud Asset API and grant the permission needed.
```
gcloud services enable cloudasset.googleapis.com --project $PROJECT

gcloud organizations add-iam-policy-binding $ORG_ID --member=user:$USER_EMAIL --role roles/cloudasset.viewer

gsutil iam ch user:$USER_EMAIL:objectViewer gs://$BUCKET 
```

## CAI Export
```
gcloud asset export --organization $ORG_ID --output-path gs://$BUCKET/iam_inventory.json --content-type iam-policy --billing-project $PROJECT

gcloud asset export --organization $ORG_ID --output-path gs://$BUCKET/resource_inventory.json --content-type resource --billing-project $PROJECT
```

## Run Scorecard
```
./cloud-foundation-toolkit/cli/bin/cft scorecard --policy-path policy-library --bucket=$BUCKET --target-project=$PROJECT --output-format csv --output-metadata-fields key,value --output-path inventory-reports
```

## Parse Result
```
go run csv-extract.go
```