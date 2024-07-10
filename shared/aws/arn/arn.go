package arn

import "fmt"

func RoleArn(account string, role string) string {
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", account, role)
}
