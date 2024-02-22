package wallet

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-ca/api"
	id "github.com/hyperledger/fabric-ca/lib"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type Wallet struct {
	contractapi.Contract
}

var (
	mspClient *id.Identity
	client    *id.Client
)

// aff: applicant, employer, admin
func (c *Wallet) Register(ctx contractapi.TransactionContextInterface, username string, secret string, aff string) pb.Response {
	regReq := &api.RegistrationRequest{
		Name:           username,
		Secret:         secret,
		Type:           "idemix",
		MaxEnrollments: -1,
		Affiliation:    aff,
	}

	exuser, err := checkClientIdentity(ctx)
	if err != nil {
		return pb.Response{
			Status:  500,
			Message: "failed to check user existence",
		}
	}

	if exuser == username {
		return pb.Response{
			Status:  400,
			Message: "user already exist",
		}
	}

	_, err = mspClient.RegisterAndEnroll(regReq)
	if err != nil {
		log.Fatalf("Failed to register user: %s", err)
		return pb.Response{
			Status:  500,
			Message: "failed to register user",
		}
	}

	return pb.Response{
		Status:  200,
		Message: "registration success",
	}
}

func (c *Wallet) Login(ctx contractapi.TransactionContextInterface, username string, secret string) pb.Response {
	enrollReq := &api.EnrollmentRequest{
		Name:   username,
		Secret: secret,
	}

	exuser, err := checkClientIdentity(ctx)
	if err != nil {
		return pb.Response{
			Status:  500,
			Message: "failed to check user existence",
		}
	}

	if exuser != username {
		return pb.Response{
			Status:  400,
			Message: "user not exist",
		}
	}

	_, err = client.Enroll(enrollReq)
	if err != nil {
		log.Fatalf("Failed to enroll user: %s", err)
		return pb.Response{
			Status:  500,
			Message: "failed to enroll user",
		}
	}

	return pb.Response{
		Status:  200,
		Message: "enrollment success",
	}
}

func checkClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}
