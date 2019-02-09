
- creating vendor using go modules

```
export GO111MODULE=on
go mod init
go mod tidy
go mod vendor
```

- deploying to appengine

`gcloud app deploy`


- explore gcloud

`gcloud topic gcloudignore`

You can stream logs from the command line by running:
`gcloud app logs tail -s default`

To view your application in the web browser run:
`gcloud app browse`