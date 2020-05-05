## Objective
This is a sample guideline to pull the tags attached across a GCP orgs using [CFT scorecard](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/cli/docs/scorecard.md). The most of the steps here are referred to [README](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/cli/docs/scorecard.md) in CFT scorecard.

## Prerequisite
Below are the tools/binary required to get installed beforehand.
* [Git](https://git-scm.com/)
* [Golang](https://golang.org/)
* [GCP CLI](https://cloud.google.com/sdk/docs/quickstarts/)

Run the following commands to install gcp cli.
```
# CFT cli for OS X
curl -o cft https://storage.googleapis.com/cft-cli/latest/cft-darwin-amd64

# CFT cli for Linux
curl -o cft https://storage.googleapis.com/cft-cli/latest/cft-linux-amd64

# executable
chmod +x cft

# Clone the policy library
git clone --branch feature/gcp-labels https://github.com/ericyz/policy-library.git
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
gcloud asset export --organization $ORG_ID --output-path gs://$BUCKET/resource_inventory.json --content-type resource --billing-project $PROJECT
```

## Run Scorecard
### Commands
```
./cft scorecard --policy-path policy-library --bucket=$BUCKET --target-project=$PROJECT --output-format csv --output-metadata-fields key,value --output-path inventory-reports
```

### Output File
#### scorecard.csv
The analysis report in a csv formate generated by cft cli with the following columns

- Category: Always "Others" in this report
- Constraint: Always "report-labels" indicating it's a label report
- Resource: A fully-qualified name of url to represent a resource
- Message: Resource Type
- key: The key of the label attached to the resource
- value: The value of the label attached to the resource

## Parse Result
Running the following command will generate four 6 csv files by parsing the scorecard.csv report

### Command
```
go run scorecard-result-parser.go
```

### Output Files
#### key-counts.csv
A csv report to count the occurances key of labels in a org. The columns are:
- Identifier: The key of the label
- Counts: The number of usage of the key

#### value-counts.csv
A csv report to count the occurances value of labels in a org
- Identifier: The value of the label
- Counts: The number of usage of the value

#### keyvalue-counts.csv
A csv report to count the occurances key-value pair of labels in a org, with a format of "key:value"
- Identifier: The key-value pair with a format of "key:value"
- Counts: The number of usage of the key-value pair

#### key-counts-by-resource.csv
A csv report to count the occurances key of labels in a org grouping by resource. The columns are:
- Resource: The name of GCP resource
- Identifier: The key of the label
- Counts: The number of usage of the key

#### value-counts-by-resource.csv
A csv report to count the occurances value of labels in a org grouping by resource
- Resource: The name of GCP resource
- Identifier: The value of the label
- Counts: The number of usage of the value

#### keyvalue-counts-by-resource.csv
A csv report to count the occurances key-value pair of labels in a org, with a format of "key:value", grouping by resource
- Resource: The name of GCP resource
- Identifier: The key-value pair with a format of "key:value"
- Counts: The number of usage of the key-value pair

## License
Apache 2.0 - See [LICENSE](./LICENSE) for more information.



