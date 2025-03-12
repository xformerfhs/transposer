//
// SPDX-FileCopyrightText: Copyright 2024 Frank Schwab
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileType: SOURCE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Frank Schwab
//
// Version: 1.0.0
//
// Change history:
//    2024-12-29: V1.0.0: Created.
//

package main

import (
	"transposer/logger"
)

// ******** Private constants ********

const (
	rcOK              = 0
	rcParameterError  = 1
	rcProcessingError = 2
)

// ******** Private functions ********

// printUsageError prints an error message and the usage information.
func printUsageError(msgNum byte, msgText string) int {
	logger.PrintError(msgNum, msgText)
	return printUsage()
}

// printUsageErrorf print a formatted error message and the usage information.
func printUsageErrorf(msgNum byte, msgFormat string, a ...any) int {
	logger.PrintErrorf(msgNum, msgFormat, a...)

	return printUsage()
}

// printUsage prints the usage.
func printUsage() int {
	encryptFlagSet.Usage()

	return rcParameterError
}

func printProcessingError(msgNum byte, msgText string) int {
	logger.PrintError(msgNum, msgText)
	return rcProcessingError
}

func printProcessingErrorf(msgNum byte, msgFormat string, a ...any) int {
	logger.PrintErrorf(msgNum, msgFormat, a...)
	return rcProcessingError
}
