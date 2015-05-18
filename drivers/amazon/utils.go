// Copyright 2015 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.
//
// Author: Marc Berhault (marc@cockroachlabs.com)

package amazon

import "github.com/awslabs/aws-sdk-go/aws"

// LoadAWSCredentials loads the credentials using the AWS api. This automatically
// loads from ENV, or from the .aws/credentials file.
// Returns the key-id and secret-key.
func LoadAWSCredentials() (string, string, error) {
	creds, err := aws.DefaultCreds().Credentials()
	if err != nil {
		return "", "", err
	}

	return creds.AccessKeyID, creds.SecretAccessKey, nil
}
