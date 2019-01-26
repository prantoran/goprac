$ gcloud services enable cloudfunctions.googleapis.com


writing a http cloud function

$ gcloud alpha pubsub topics create randomNumbers



setting up go modules

$ export GO11MODULE=on
$ go mod init
$ go mod tidy
    - checks dependencies, downloads them and creates `go.sum`
$ go mod vendor

render your dependencies locally before you deploy the Google Cloud function


deploy to gcloud functions
$ gcloud alpha functions deploy api --entry-point Send --runtime go111 --trigger-http --set-env-vars PROJECT_ID=project-id

response will be a bunch of info including an url, which when called using curl, returns projectid




creating background functions

deploy background function to gcloud

$ gcloud alpha functions deploy consumer --entry-point Consume --runtime go111 --trigger-topic=randomNumbers

$ gcloud alpha functions logs read consumer


deleting functions

$ gcloud alpha functions delete api
$ gcloud alpha functions delete consumer
$ gcloud alpha pubsub topics delete randomNumbers