/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package cmd

import (
	"fmt"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var exportAppName string
var exportAppOwner string
var exportAppWithKeys bool

//var flagExportAPICmdToken string
// ExportApp command related usage info
const exportAppCmdLiteral = "export-app"
const exportAppCmdShortDesc = "Export App"

const exportAppCmdLongDesc = "Export an Application from a specified  environment"

const exportAppCmdExamples = utils.ProjectName + ` ` + exportAppCmdLiteral + ` -n SampleApp -o admin -e dev
` + utils.ProjectName + ` ` + exportAppCmdLiteral + ` -n SampleApp -o admin -e prod
NOTE: All the 3 flags (--name (-n), --owner (-o) and --environment (-e)) are mandatory`

// exportAppCmd represents the exportApp command
var ExportAppCmd = &cobra.Command{
	Use: exportAppCmdLiteral + " (--name <name-of-the-application> --owner <owner-of-the-application> --environment " +
		"<environment-from-which-the-app-should-be-exported>)",
	Short:   exportAppCmdShortDesc,
	Long:    exportAppCmdLongDesc,
	Example: exportAppCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + exportAppCmdLiteral + " called")
		var appsExportDirectoryPath = filepath.Join(utils.ExportDirectory, utils.ExportedAppsDirName, cmdExportEnvironment)

		cred, err := getCredentials(cmdExportEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeExportAppCmd(cred, appsExportDirectoryPath)
	},
}

func executeExportAppCmd(credential credentials.Credential, appsExportDirectoryPath string) {
	accessToken, preCommandErr := credentials.GetOAuthAccessToken(credential, cmdExportEnvironment)

	if preCommandErr == nil {
		adminEndpiont := utils.GetAdminEndpointOfEnv(cmdExportEnvironment, utils.MainConfigFilePath)
		resp, err := getExportAppResponse(exportAppName, exportAppOwner, adminEndpiont, accessToken)
		if err != nil {
			utils.HandleErrorAndExit("Error exporting Application: "+exportAppName, err)
		}

		// Print info on response
		utils.Logf(utils.LogPrefixInfo+"ResponseStatus: %v\n", resp.Status())
		if resp.StatusCode() == http.StatusOK {
			WriteApplicationToZip(exportAppName, exportAppOwner, appsExportDirectoryPath, resp)
		} else {
			fmt.Println("Error " + string(resp.Body()))
		}
	} else {
		// error exporting Application
		fmt.Println("Error exporting Application:" + preCommandErr.Error())
	}
}

// WriteApplicationToZip
// @param exportAppName : Name of the Application to be exported
// @param exportAppOwner : Owner of the Application to be exported
// @param resp : Response returned from making the HTTP request (only pass a 200 OK)
// Exported Application will be written to a zip file
func WriteApplicationToZip(exportAppName, exportAppOwner, zipLocationPath string,
	resp *resty.Response) {
	zipFilename := replaceUserStoreDomainDelimiter(exportAppOwner) + "_" + exportAppName + ".zip" // admin_testApp.zip
	// Writes the REST API response to a temporary zip file
	tempZipFile, err := utils.WriteResponseToTempZip(zipFilename, resp)
	if err != nil {
		utils.HandleErrorAndExit("Error creating the temporary zip file to store the exported application", err)
	}

	err = utils.CreateDirIfNotExist(zipLocationPath)
	if err != nil {
		utils.HandleErrorAndExit("Error creating dir to store zip archive: "+zipLocationPath, err)
	}

	exportedFinalZip := filepath.Join(zipLocationPath, zipFilename)
	// Add application_params.yaml file inside the zip and create a new zip file in exportedFinalZip location
	err = impl.IncludeParamsFileToZip(tempZipFile, exportedFinalZip, utils.ParamFileApplication)
	if err != nil {
		utils.HandleErrorAndExit("Error creating the final zip archive", err)
	}
	fmt.Println("Successfully exported Application!")
	fmt.Println("Find the exported Application at " + exportedFinalZip)
}

// The Application owner name is used to construct a unique name for the app export zip.
// When an app belonging to a user from a secondary user store is exported, the owner name will have
// the format '<Userstore_domain>/<Username>'. The '/' character will be mistakenly considerd as a
// file separator character, resulting in an invalid path being constructed.
// Therefore this function overcomes this issue by replacing the '/' character.
func replaceUserStoreDomainDelimiter(username string) string {
	return strings.ReplaceAll(username, "/", "#")
}

// ExportApp
// @param name : Name of the Application to be exported
// @param apimEndpoint : API Manager Endpoint for the environment
// @param accessToken : Access Token for the resource
// @return response Response in the form of *resty.Response
func getExportAppResponse(name, owner, adminEndpoint, accessToken string) (*resty.Response, error) {
	adminEndpoint = utils.AppendSlashToString(adminEndpoint)
	query := "export/applications?appName=" + name + utils.SearchAndTag + "appOwner=" + owner

	if exportAppWithKeys {
		query += "&withKeys=true"
	}

	url := adminEndpoint + query
	utils.Logln(utils.LogPrefixInfo+"ExportApp: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	headers[utils.HeaderAccept] = utils.HeaderValueApplicationZip

	resp, err := utils.InvokeGETRequest(url, headers)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//init using Cobra
func init() {
	RootCmd.AddCommand(ExportAppCmd)
	ExportAppCmd.Flags().StringVarP(&exportAppName, "name", "n", "",
		"Name of the Application to be exported")
	ExportAppCmd.Flags().StringVarP(&exportAppOwner, "owner", "o", "",
		"Owner of the Application to be exported")
	ExportAppCmd.Flags().StringVarP(&cmdExportEnvironment, "environment", "e",
		"", "Environment to which the Application should be exported")
	ExportAppCmd.Flags().BoolVarP(&exportAppWithKeys, "withKeys", "",
		false, "Export keys for the application ")
	_ = ExportAppCmd.MarkFlagRequired("environment")
	_ = ExportAppCmd.MarkFlagRequired("owner")
	_ = ExportAppCmd.MarkFlagRequired("name")
}
