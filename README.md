# T212 to Digrin CLI
Golang CLI tool for fetching T212 reports via API call and transforming them to be used in Digrin portfolio tracker. Stores the reports in AWS S3.

```bash
echo "T212_API_KEY=xxx" >> .env
echo "AWS_ACCESS_KEY_ID=xxx" >> .env # or use aws configure
echo "AWS_SECRET_ACCESS_KEY=xxx" >> .env # or use aws configure
echo "AWS_REGION=xxx" >> .env # or use aws configure
echo "BUCKET_NAME=xxx" >> .env
```

```bash
go mod download
```

```bash
go run main.py
```

# TODO

- [ ] investigate option of go routines

- [ ] add [log](https://pkg.go.dev/log) / [logrus](https://github.com/sirupsen/logrus)

- [ ] add decorators or decorators like functionality

- [ ] yield fetchReports ?

- [ ] manage secrets via aws secrets
