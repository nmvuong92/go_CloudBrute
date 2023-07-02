go run main.go -d 812c-222-253-42-176.ngrok-free.app -k xxx -m storage -t 80 -T 10 -w "./data/storage_small.txt"
go run main.go -d github.com -k github -m storage -t 80 -T 10 -w "./data/storage_small.txt"
go run main.go -d google.com -k github -m storage -t 80 -T 10 -w "./data/storage_small.txt"


cloudbrute -C /opt/external-go-modules/CloudBrute/config -d xyz.com -k xyz -c amazon -t 10 -w /opt/external-go-modules/CloudBrute/data/storage_small.txt -o /tmp/awss3-storagebrute-8e1159cb-13b3-4367-909d-ba891df0816c/awss3-bruteforce-.txt

go run main.go -C /Users/vuong/projects/google/CloudBrute/config  -d watchtowr.com -k watchtowr -m storage -t 80 -T 10 -w "./data/storage_small.txt" -c amazon -o "out.txt"