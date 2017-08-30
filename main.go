package main

import (
	"encoding/json"
	"fmt"
	auth "github.com/swapagarwal/aws-remote/Godeps/_workspace/src/github.com/abbot/go-http-auth"
	"github.com/swapagarwal/aws-remote/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/swapagarwal/aws-remote/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/ec2"
	"github.com/swapagarwal/aws-remote/Godeps/_workspace/src/github.com/gorilla/mux"
	"net/http"
	"os"
)

type Instance struct {
	ID       string
	PublicIP string
}

type Instances []Instance

func Secret(user, realm string) string {
	if user == os.Getenv("login") {
		return os.Getenv("password")
	}
	return ""
}

func welcome(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	fmt.Fprintf(w, "Welcome, %s!", r.Username)
}

func ListEC2(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	svc := ec2.New(&aws.Config{Region: aws.String("us-east-1")})

	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	instances := Instances{}

	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			instance := Instance{}
			instance.ID = *inst.InstanceId
			if inst.PublicIpAddress != nil {
				instance.PublicIP = *inst.PublicIpAddress
			}
			instances = append(instances, instance)
		}
	}

	json.NewEncoder(w).Encode(instances)
}

func StartEC2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceId := vars["id"]

	svc := ec2.New(&aws.Config{Region: aws.String("us-east-1")})
	params := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	resp, err := svc.StartInstances(params)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, resp)
}

func StopEC2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceId := vars["id"]

	svc := ec2.New(&aws.Config{Region: aws.String("us-east-1")})
	params := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	resp, err := svc.StopInstances(params)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, resp)
}

func ListS3(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	svc := s3.New(session.New())

	resp, err := svc.ListBuckets(nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
	json.NewEncoder(w).Encode(resp)
}

func main() {
	authenticator := auth.NewBasicAuthenticator("aws-remote.herokuapp.com", Secret)

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/", authenticator.Wrap(welcome))

	r.HandleFunc("/ec2/", authenticator.Wrap(ListEC2))
	r.HandleFunc("/ec2/start/{id}", StartEC2)
	r.HandleFunc("/ec2/stop/{id}", StopEC2)

	r.HandleFunc("/s3/", authenticator.Wrap(ListS3))

	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}
