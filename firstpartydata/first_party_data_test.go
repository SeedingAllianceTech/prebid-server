package firstpartydata

import (
	"encoding/json"
	"github.com/mxmCherry/openrtb/v15/openrtb2"
	"github.com/prebid/prebid-server/errortypes"
	"github.com/prebid/prebid-server/openrtb_ext"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestExtractGlobalFPD(t *testing.T) {

	testCases := []struct {
		description string
		input       openrtb_ext.RequestWrapper
		expectedReq openrtb_ext.RequestWrapper
		expectedFpd map[string][]byte
	}{
		{
			description: "Site, app and user data present",
			input: openrtb_ext.RequestWrapper{
				BidRequest: &openrtb2.BidRequest{
					ID: "bid_id",
					Site: &openrtb2.Site{
						ID:   "reqSiteId",
						Page: "http://www.foobar.com/1234.html",
						Publisher: &openrtb2.Publisher{
							ID: "1",
						},
						Ext: json.RawMessage(`{"data": {"somesitefpd": "sitefpdDataTest"}}`),
					},
					User: &openrtb2.User{
						ID:     "reqUserID",
						Yob:    1982,
						Gender: "M",
						Ext:    json.RawMessage(`{"data": {"someuserfpd": "userfpdDataTest"}}`),
					},
					App: &openrtb2.App{
						ID:  "appId",
						Ext: json.RawMessage(`{"data": {"someappfpd": "appfpdDataTest"}}`),
					},
				},
			},
			expectedReq: openrtb_ext.RequestWrapper{BidRequest: &openrtb2.BidRequest{
				ID: "bid_id",
				Site: &openrtb2.Site{
					ID:   "reqSiteId",
					Page: "http://www.foobar.com/1234.html",
					Publisher: &openrtb2.Publisher{
						ID: "1",
					},
				},
				User: &openrtb2.User{
					ID:     "reqUserID",
					Yob:    1982,
					Gender: "M",
				},
				App: &openrtb2.App{
					ID: "appId",
				},
			}},
			expectedFpd: map[string][]byte{
				"site": []byte(`{"somesitefpd": "sitefpdDataTest"}`),
				"user": []byte(`{"someuserfpd": "userfpdDataTest"}`),
				"app":  []byte(`{"someappfpd": "appfpdDataTest"}`),
			},
		},
		{
			description: "App FPD only present",
			input: openrtb_ext.RequestWrapper{
				BidRequest: &openrtb2.BidRequest{
					ID: "bid_id",
					Site: &openrtb2.Site{
						ID:   "reqSiteId",
						Page: "http://www.foobar.com/1234.html",
						Publisher: &openrtb2.Publisher{
							ID: "1",
						},
					},
					App: &openrtb2.App{
						ID:  "appId",
						Ext: json.RawMessage(`{"data": {"someappfpd": "appfpdDataTest"}}`),
					},
				},
			},
			expectedReq: openrtb_ext.RequestWrapper{
				BidRequest: &openrtb2.BidRequest{
					ID: "bid_id",
					Site: &openrtb2.Site{
						ID:   "reqSiteId",
						Page: "http://www.foobar.com/1234.html",
						Publisher: &openrtb2.Publisher{
							ID: "1",
						},
					},
					App: &openrtb2.App{
						ID: "appId",
					},
				},
			},
			expectedFpd: map[string][]byte{
				"app":  []byte(`{"someappfpd": "appfpdDataTest"}`),
				"user": nil,
				"site": nil,
			},
		},
		{
			description: "User FPD only present",
			input: openrtb_ext.RequestWrapper{
				BidRequest: &openrtb2.BidRequest{
					ID: "bid_id",
					Site: &openrtb2.Site{
						ID:   "reqSiteId",
						Page: "http://www.foobar.com/1234.html",
						Publisher: &openrtb2.Publisher{
							ID: "1",
						},
					},
					User: &openrtb2.User{
						ID:     "reqUserID",
						Yob:    1982,
						Gender: "M",
						Ext:    json.RawMessage(`{"data": {"someuserfpd": "userfpdDataTest"}}`),
					},
				},
			},
			expectedReq: openrtb_ext.RequestWrapper{
				BidRequest: &openrtb2.BidRequest{
					ID: "bid_id",
					Site: &openrtb2.Site{
						ID:   "reqSiteId",
						Page: "http://www.foobar.com/1234.html",
						Publisher: &openrtb2.Publisher{
							ID: "1",
						},
					},
					User: &openrtb2.User{
						ID:     "reqUserID",
						Yob:    1982,
						Gender: "M",
					},
				},
			},
			expectedFpd: map[string][]byte{
				"app":  nil,
				"user": []byte(`{"someuserfpd": "userfpdDataTest"}`),
				"site": nil,
			},
		},
		{
			description: "No FPD present in req",
			input: openrtb_ext.RequestWrapper{
				BidRequest: &openrtb2.BidRequest{
					ID: "bid_id",
					Site: &openrtb2.Site{
						ID:   "reqSiteId",
						Page: "http://www.foobar.com/1234.html",
						Publisher: &openrtb2.Publisher{
							ID: "1",
						},
					},
					User: &openrtb2.User{
						ID:     "reqUserID",
						Yob:    1982,
						Gender: "M",
					},
					App: &openrtb2.App{
						ID: "appId",
					},
				},
			},
			expectedReq: openrtb_ext.RequestWrapper{
				BidRequest: &openrtb2.BidRequest{
					ID: "bid_id",
					Site: &openrtb2.Site{
						ID:   "reqSiteId",
						Page: "http://www.foobar.com/1234.html",
						Publisher: &openrtb2.Publisher{
							ID: "1",
						},
					},
					User: &openrtb2.User{
						ID:     "reqUserID",
						Yob:    1982,
						Gender: "M",
					},
					App: &openrtb2.App{
						ID: "appId",
					},
				},
			},
			expectedFpd: map[string][]byte{
				"app":  nil,
				"user": nil,
				"site": nil,
			},
		},
		{
			description: "Site FPD only present",
			input: openrtb_ext.RequestWrapper{
				BidRequest: &openrtb2.BidRequest{
					ID: "bid_id",
					Site: &openrtb2.Site{
						ID:   "reqSiteId",
						Page: "http://www.foobar.com/1234.html",
						Publisher: &openrtb2.Publisher{
							ID: "1",
						},
						Ext: json.RawMessage(`{"data": {"someappfpd": true}}`),
					},
					App: &openrtb2.App{
						ID: "appId",
					},
				},
			},
			expectedReq: openrtb_ext.RequestWrapper{
				BidRequest: &openrtb2.BidRequest{
					ID: "bid_id",
					Site: &openrtb2.Site{
						ID:   "reqSiteId",
						Page: "http://www.foobar.com/1234.html",
						Publisher: &openrtb2.Publisher{
							ID: "1",
						},
					},
					App: &openrtb2.App{
						ID: "appId",
					},
				},
			},
			expectedFpd: map[string][]byte{
				"app":  nil,
				"user": nil,
				"site": []byte(`{"someappfpd": true}`),
			},
		},
	}
	for _, test := range testCases {

		inputReq := &test.input
		fpd, err := ExtractGlobalFPD(inputReq)
		assert.NoError(t, err, "Error should be nil")
		err = inputReq.RebuildRequest()
		assert.NoError(t, err, "Error should be nil")

		assert.Equal(t, test.expectedReq.BidRequest, inputReq.BidRequest, "Incorrect input request after global fpd extraction")

		assert.Equal(t, test.expectedFpd[userKey], fpd[userKey], "Incorrect User FPD")
		assert.Equal(t, test.expectedFpd[appKey], fpd[appKey], "Incorrect App FPD")
		assert.Equal(t, test.expectedFpd[siteKey], fpd[siteKey], "Incorrect Site FPDt")
	}
}

func TestExtractOpenRtbGlobalFPD(t *testing.T) {

	testCases := []struct {
		description     string
		input           openrtb2.BidRequest
		output          openrtb2.BidRequest
		expectedFpdData map[string][]openrtb2.Data
	}{
		{
			description: "Site, app and user data present",
			input: openrtb2.BidRequest{
				ID: "bid_id",
				Imp: []openrtb2.Imp{
					{ID: "impid"},
				},
				Site: &openrtb2.Site{
					ID: "reqSiteId",
					Content: &openrtb2.Content{
						Data: []openrtb2.Data{
							{ID: "siteDataId1", Name: "siteDataName1"},
							{ID: "siteDataId2", Name: "siteDataName2"},
						},
					},
				},
				User: &openrtb2.User{
					ID:     "reqUserID",
					Yob:    1982,
					Gender: "M",
					Data: []openrtb2.Data{
						{ID: "userDataId1", Name: "userDataName1"},
					},
				},
				App: &openrtb2.App{
					ID: "appId",
					Content: &openrtb2.Content{
						Data: []openrtb2.Data{
							{ID: "appDataId1", Name: "appDataName1"},
						},
					},
				},
			},
			output: openrtb2.BidRequest{
				ID: "bid_id",
				Imp: []openrtb2.Imp{
					{ID: "impid"},
				},
				Site: &openrtb2.Site{
					ID:      "reqSiteId",
					Content: &openrtb2.Content{},
				},
				User: &openrtb2.User{
					ID:     "reqUserID",
					Yob:    1982,
					Gender: "M",
				},
				App: &openrtb2.App{
					ID:      "appId",
					Content: &openrtb2.Content{},
				},
			},
			expectedFpdData: map[string][]openrtb2.Data{
				siteContentDataKey: {{ID: "siteDataId1", Name: "siteDataName1"}, {ID: "siteDataId2", Name: "siteDataName2"}},
				userDataKey:        {{ID: "userDataId1", Name: "userDataName1"}},
				appContentDataKey:  {{ID: "appDataId1", Name: "appDataName1"}},
			},
		},
		{
			description: "No Site, app or user data present",
			input: openrtb2.BidRequest{
				ID: "bid_id",
				Imp: []openrtb2.Imp{
					{ID: "impid"},
				},
			},
			output: openrtb2.BidRequest{
				ID: "bid_id",
				Imp: []openrtb2.Imp{
					{ID: "impid"},
				},
			},
			expectedFpdData: map[string][]openrtb2.Data{
				siteContentDataKey: nil,
				userDataKey:        nil,
				appContentDataKey:  nil,
			},
		},
		{
			description: "Site only data present",
			input: openrtb2.BidRequest{
				ID: "bid_id",
				Imp: []openrtb2.Imp{
					{ID: "impid"},
				},
				Site: &openrtb2.Site{
					ID:   "reqSiteId",
					Page: "test/page",
					Content: &openrtb2.Content{
						Data: []openrtb2.Data{
							{ID: "siteDataId1", Name: "siteDataName1"},
						},
					},
				},
			},
			output: openrtb2.BidRequest{
				ID: "bid_id",
				Imp: []openrtb2.Imp{
					{ID: "impid"},
				},
				Site: &openrtb2.Site{
					ID:      "reqSiteId",
					Page:    "test/page",
					Content: &openrtb2.Content{},
				},
			},
			expectedFpdData: map[string][]openrtb2.Data{
				siteContentDataKey: {{ID: "siteDataId1", Name: "siteDataName1"}},
				userDataKey:        nil,
				appContentDataKey:  nil,
			},
		},
		{
			description: "App only data present",
			input: openrtb2.BidRequest{
				ID: "bid_id",
				Imp: []openrtb2.Imp{
					{ID: "impid"},
				},
				App: &openrtb2.App{
					ID: "reqAppId",
					Content: &openrtb2.Content{
						Data: []openrtb2.Data{
							{ID: "appDataId1", Name: "appDataName1"},
						},
					},
				},
			},
			output: openrtb2.BidRequest{
				ID: "bid_id",
				Imp: []openrtb2.Imp{
					{ID: "impid"},
				},
				App: &openrtb2.App{
					ID:      "reqAppId",
					Content: &openrtb2.Content{},
				},
			},
			expectedFpdData: map[string][]openrtb2.Data{
				siteContentDataKey: nil,
				userDataKey:        nil,
				appContentDataKey:  {{ID: "appDataId1", Name: "appDataName1"}},
			},
		},
		{
			description: "User only data present",
			input: openrtb2.BidRequest{
				ID: "bid_id",
				Imp: []openrtb2.Imp{
					{ID: "impid"},
				},
				Site: &openrtb2.Site{
					ID: "reqSiteId",
				},
				App: &openrtb2.App{
					ID: "reqAppId",
				},
				User: &openrtb2.User{
					ID:     "reqUserId",
					Yob:    1982,
					Gender: "M",
					Data: []openrtb2.Data{
						{ID: "userDataId1", Name: "userDataName1"},
					},
				},
			},
			output: openrtb2.BidRequest{
				ID: "bid_id",
				Imp: []openrtb2.Imp{
					{ID: "impid"},
				},
				Site: &openrtb2.Site{
					ID: "reqSiteId",
				},
				App: &openrtb2.App{
					ID: "reqAppId",
				},
				User: &openrtb2.User{
					ID:     "reqUserId",
					Yob:    1982,
					Gender: "M",
				},
			},
			expectedFpdData: map[string][]openrtb2.Data{
				siteContentDataKey: nil,
				userDataKey:        {{ID: "userDataId1", Name: "userDataName1"}},
				appContentDataKey:  nil,
			},
		},
	}
	for _, test := range testCases {

		inputReq := &test.input

		res := ExtractOpenRtbGlobalFPD(inputReq)

		assert.Equal(t, &test.output, inputReq, "Result request is incorrect")
		assert.Equal(t, test.expectedFpdData[siteContentDataKey], res[siteContentDataKey], "siteContentData data is incorrect")
		assert.Equal(t, test.expectedFpdData[userDataKey], res[userDataKey], "userData is incorrect")
		assert.Equal(t, test.expectedFpdData[appContentDataKey], res[appContentDataKey], "appContentData is incorrect")

	}
}

func TestExtractBidderConfigFPD(t *testing.T) {

	if specFiles, err := ioutil.ReadDir("./tests/extractbidderconfigfpd"); err == nil {
		for _, specFile := range specFiles {
			fileName := "./tests/extractbidderconfigfpd/" + specFile.Name()

			fpdFile, err := loadFpdFile(fileName)
			if err != nil {
				t.Errorf("Unable to load file: %s", fileName)
			}
			var extReq openrtb_ext.ExtRequestPrebid
			err = json.Unmarshal(fpdFile.InputRequestData, &extReq)
			if err != nil {
				t.Errorf("Unable to unmarshal input request: %s", fileName)
			}
			reqExt := openrtb_ext.RequestExt{}
			reqExt.SetPrebid(&extReq)
			fpdData, err := ExtractBidderConfigFPD(&reqExt)

			if len(fpdFile.ValidationErrors) > 0 {
				assert.Equal(t, err.Error(), fpdFile.ValidationErrors[0].Message, "Incorrect first party data error message")
				continue
			}

			assert.Nil(t, reqExt.GetPrebid().BidderConfigs, "Bidder specific FPD config should be removed from request")

			assert.Nil(t, err, "No error should be returned")
			assert.Equal(t, len(fpdFile.BidderConfigFPD), len(fpdData), "Incorrect fpd data")

			for bidderName, bidderFPD := range fpdFile.BidderConfigFPD {

				if bidderFPD.Site != nil {
					resSite := fpdData[bidderName].Site
					for k, v := range bidderFPD.Site {
						assert.NotNil(t, resSite[k], "Property is not found in result site")
						assert.JSONEq(t, string(v), string(resSite[k]), "site is incorrect")
					}
				} else {
					assert.Nil(t, fpdData[bidderName].Site, "Result site should be also nil")
				}

				if bidderFPD.App != nil {
					resApp := fpdData[bidderName].App
					for k, v := range bidderFPD.App {
						assert.NotNil(t, resApp[k], "Property is not found in result app")
						assert.JSONEq(t, string(v), string(resApp[k]), "app is incorrect")
					}
				} else {
					assert.Nil(t, fpdData[bidderName].App, "Result app should be also nil")
				}

				if bidderFPD.User != nil {
					resUser := fpdData[bidderName].User
					for k, v := range bidderFPD.User {
						assert.NotNil(t, resUser[k], "Property is not found in result user")
						assert.JSONEq(t, string(v), string(resUser[k]), "site is incorrect")
					}
				} else {
					assert.Nil(t, fpdData[bidderName].User, "Result user should be also nil")
				}
			}
		}
	}
}

func TestResolveFPD(t *testing.T) {

	if specFiles, err := ioutil.ReadDir("./tests/resolvefpd"); err == nil {
		for _, specFile := range specFiles {
			fileName := "./tests/resolvefpd/" + specFile.Name()

			fpdFile, err := loadFpdFile(fileName)
			if err != nil {
				t.Errorf("Unable to load file: %s", fileName)
			}

			var inputReq openrtb2.BidRequest
			err = json.Unmarshal(fpdFile.InputRequestData, &inputReq)
			if err != nil {
				t.Errorf("Unable to unmarshal input request: %s", fileName)
			}

			var inputReqCopy openrtb2.BidRequest
			err = json.Unmarshal(fpdFile.InputRequestData, &inputReqCopy)
			if err != nil {
				t.Errorf("Unable to unmarshal input request: %s", fileName)
			}

			var outputReq openrtb2.BidRequest
			err = json.Unmarshal(fpdFile.OutputRequestData, &outputReq)
			if err != nil {
				t.Errorf("Unable to unmarshal output request: %s", fileName)
			}

			reqExtFPD := make(map[string][]byte, 3)
			reqExtFPD["site"] = fpdFile.GlobalFPD["site"]
			reqExtFPD["app"] = fpdFile.GlobalFPD["app"]
			reqExtFPD["user"] = fpdFile.GlobalFPD["user"]

			reqFPD := make(map[string][]openrtb2.Data, 3)

			reqFPDSiteContentData := fpdFile.GlobalFPD[siteContentDataKey]
			if len(reqFPDSiteContentData) > 0 {
				var siteConData []openrtb2.Data
				err = json.Unmarshal(reqFPDSiteContentData, &siteConData)
				if err != nil {
					t.Errorf("Unable to unmarshal site.content.data: %s", fileName)
				}
				reqFPD[siteContentDataKey] = siteConData
			}

			reqFPDAppContentData := fpdFile.GlobalFPD[appContentDataKey]
			if len(reqFPDAppContentData) > 0 {
				var appConData []openrtb2.Data
				err = json.Unmarshal(reqFPDAppContentData, &appConData)
				if err != nil {
					t.Errorf("Unable to unmarshal app.content.data: %s", fileName)
				}
				reqFPD[appContentDataKey] = appConData
			}

			reqFPDUserData := fpdFile.GlobalFPD[userDataKey]
			if len(reqFPDUserData) > 0 {
				var userData []openrtb2.Data
				err = json.Unmarshal(reqFPDUserData, &userData)
				if err != nil {
					t.Errorf("Unable to unmarshal app.content.data: %s", fileName)
				}
				reqFPD[userDataKey] = userData
			}
			if fpdFile.BidderConfigFPD == nil {
				fpdFile.BidderConfigFPD = make(map[openrtb_ext.BidderName]*openrtb_ext.ORTB2)
				fpdFile.BidderConfigFPD["appnexus"] = &openrtb_ext.ORTB2{}
			}

			resultFPD, errL := ResolveFPD(&inputReq, fpdFile.BidderConfigFPD, reqExtFPD, reqFPD, []string{"appnexus"})

			if len(errL) == 0 {
				assert.Equal(t, inputReq, inputReqCopy, "Original request should not be modified")

				bidderFPD := resultFPD["appnexus"]

				if outputReq.Site != nil && len(outputReq.Site.Ext) > 0 {
					resSiteExt := bidderFPD.Site.Ext
					expectedSiteExt := outputReq.Site.Ext
					bidderFPD.Site.Ext = nil
					outputReq.Site.Ext = nil
					assert.JSONEq(t, string(expectedSiteExt), string(resSiteExt), "site.ext is incorrect")

					assert.Equal(t, outputReq.Site, bidderFPD.Site, "Site is incorrect")
				}
				if outputReq.App != nil && len(outputReq.App.Ext) > 0 {
					resAppExt := bidderFPD.App.Ext
					expectedAppExt := outputReq.App.Ext
					bidderFPD.App.Ext = nil
					outputReq.App.Ext = nil
					assert.JSONEq(t, string(expectedAppExt), string(resAppExt), "app.ext is incorrect")

					assert.Equal(t, outputReq.App, bidderFPD.App, "App is incorrect")
				}
				if outputReq.User != nil && len(outputReq.User.Ext) > 0 {
					resUserExt := bidderFPD.User.Ext
					expectedUserExt := outputReq.User.Ext
					bidderFPD.User.Ext = nil
					outputReq.User.Ext = nil
					assert.JSONEq(t, string(expectedUserExt), string(resUserExt), "user.ext is incorrect")

					assert.Equal(t, outputReq.User, bidderFPD.User, "User is incorrect")
				}
			} else {
				assert.ElementsMatch(t, errL, fpdFile.ValidationErrors, "Incorrect first party data warning message")
			}

		}
	}
}

func TestExtractFPDForBidders(t *testing.T) {

	if specFiles, err := ioutil.ReadDir("./tests/extractfpdforbidders"); err == nil {
		for _, specFile := range specFiles {
			fileName := "./tests/extractfpdforbidders/" + specFile.Name()
			fpdFile, err := loadFpdFile(fileName)

			if err != nil {
				t.Errorf("Unable to load file: %s", fileName)
			}

			var expectedRequest openrtb2.BidRequest
			err = json.Unmarshal(fpdFile.OutputRequestData, &expectedRequest)
			if err != nil {
				t.Errorf("Unable to unmarshal input request: %s", fileName)
			}

			resultRequest := &openrtb_ext.RequestWrapper{}
			resultRequest.BidRequest = &openrtb2.BidRequest{}
			err = json.Unmarshal(fpdFile.InputRequestData, resultRequest.BidRequest)
			assert.NoError(t, err, "Error should be nil")

			resultFPD, errL := ExtractFPDForBidders(resultRequest)

			if len(fpdFile.ValidationErrors) > 0 {
				assert.Equal(t, len(fpdFile.ValidationErrors), len(errL), "Incorrect number of errors was returned")
				assert.ElementsMatch(t, errL, fpdFile.ValidationErrors, "Incorrect errors were returned")
				//in case or error no further assertions needed
				continue
			}
			assert.Empty(t, errL, "Error should be empty")
			assert.Equal(t, len(resultFPD), len(fpdFile.BiddersFPDResolved))

			for bidderName, expectedValue := range fpdFile.BiddersFPDResolved {
				actualValue := resultFPD[bidderName]
				if expectedValue.Site != nil {
					if len(expectedValue.Site.Ext) > 0 {
						assert.JSONEq(t, string(expectedValue.Site.Ext), string(actualValue.Site.Ext), "Incorrect first party data")
						expectedValue.Site.Ext = nil
						actualValue.Site.Ext = nil
					}
					assert.Equal(t, expectedValue.Site, actualValue.Site, "Incorrect first party data")
				}
				if expectedValue.App != nil {
					if len(expectedValue.App.Ext) > 0 {
						assert.JSONEq(t, string(expectedValue.App.Ext), string(actualValue.App.Ext), "Incorrect first party data")
						expectedValue.App.Ext = nil
						actualValue.App.Ext = nil
					}
					assert.Equal(t, expectedValue.App, actualValue.App, "Incorrect first party data")
				}
				if expectedValue.User != nil {
					if len(expectedValue.User.Ext) > 0 {
						assert.JSONEq(t, string(expectedValue.User.Ext), string(actualValue.User.Ext), "Incorrect first party data")
						expectedValue.User.Ext = nil
						actualValue.User.Ext = nil
					}
					assert.Equal(t, expectedValue.User, actualValue.User, "Incorrect first party data")
				}
			}

			if expectedRequest.Site != nil {
				if len(expectedRequest.Site.Ext) > 0 {
					assert.JSONEq(t, string(expectedRequest.Site.Ext), string(resultRequest.BidRequest.Site.Ext), "Incorrect site in request")
					expectedRequest.Site.Ext = nil
					resultRequest.BidRequest.Site.Ext = nil
				}
				assert.Equal(t, expectedRequest.Site, resultRequest.BidRequest.Site, "Incorrect site in request")
			}
			if expectedRequest.App != nil {
				if len(expectedRequest.App.Ext) > 0 {
					assert.JSONEq(t, string(expectedRequest.App.Ext), string(resultRequest.BidRequest.App.Ext), "Incorrect app in request")
					expectedRequest.App.Ext = nil
					resultRequest.BidRequest.App.Ext = nil
				}
				assert.Equal(t, expectedRequest.App, resultRequest.BidRequest.App, "Incorrect app in request")
			}
			if expectedRequest.User != nil {
				if len(expectedRequest.User.Ext) > 0 {
					assert.JSONEq(t, string(expectedRequest.User.Ext), string(resultRequest.BidRequest.User.Ext), "Incorrect user in request")
					expectedRequest.User.Ext = nil
					resultRequest.BidRequest.User.Ext = nil
				}
				assert.Equal(t, expectedRequest.User, resultRequest.BidRequest.User, "Incorrect user in request")
			}

		}
	}
}

func loadFpdFile(filename string) (fpdFile, error) {
	var fileData fpdFile
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return fileData, err
	}
	err = json.Unmarshal(fileContents, &fileData)
	if err != nil {
		return fileData, err
	}

	return fileData, nil
}

type fpdFile struct {
	InputRequestData   json.RawMessage                                    `json:"inputRequestData,omitempty"`
	OutputRequestData  json.RawMessage                                    `json:"outputRequestData,omitempty"`
	BidderConfigFPD    map[openrtb_ext.BidderName]*openrtb_ext.ORTB2      `json:"bidderConfigFPD,omitempty"`
	BiddersFPDResolved map[openrtb_ext.BidderName]*ResolvedFirstPartyData `json:"biddersFPDResolved,omitempty"`
	GlobalFPD          map[string]json.RawMessage                         `json:"globalFPD,omitempty"`
	ValidationErrors   []*errortypes.BadInput                             `json:"validationErrors,omitempty"`
}
