package radixs

import "testing"

func TestGetLongestMatch(t *testing.T) {
	assert := newAssert(t)
	tr, err := FromMap(pairs)
	assert(err == nil, "error creating tree from map", "err:", err)

	key := "smarties"
	expected := "smart"
	prefix, value, err := tr.LongestMatch(key)
	assert(err == nil, "longest match for key:", key, "not found, err:", err)
	assert(prefix == expected, "unmatched longest prefix, expected:", "smart", "got:", prefix)
	assert(value == pairs[expected], "unmatched longest prefix value, expected:", pairs[expected], "got:", value)

	key = "rubberized"
	expected = "rubberized"
	prefix, value, err = tr.LongestMatch(key)
	assert(err == nil, "longest match for key:", key, "not found, err:", err)
	assert(prefix == expected, "unmatched longest prefix, expected:", expected, "got:", prefix)
	assert(value == pairs[expected], "unmatched longest prefix value, expected:", pairs[expected], "got:", value)

	_, _, err = tr.LongestMatch("smallest")
	assert(err != nil, "longest match for key:", key, "should not exist")
}

func TestGetWithParams(t *testing.T) {
	assert := newAssert(t)
	tr := New(WithParams('/', ':'))

	err := tr.SetWithParams("/api/v1/projects/:project", "ProjectHandler")
	assert(err == nil, "error setting key:", "/api/v1/projects/:project", "error:", err)

	err = tr.SetWithParams("/api/v1/projects/:project/instances/:instance", "InstanceHandler")
	assert(err == nil, "error setting key:", "/api/v1/projects/:project/instances/:instance", "error:", err)

	err = tr.SetWithParams("/api/v1/projects/:project/instances/:instance/databases/:database", "DatabaseHandler")
	assert(err == nil, "error setting key:", "/api/v1/projects/:project/instances/:instance/databases/:database", "error:", err)

	err = tr.SetWithParams("/api/v1/projects/:project/instances/:instance/applications/:application", "ApplicationHandler")
	assert(err == nil, "error setting key:", "/api/v1/projects/:project/instances/:instance/applications/:application", "error:", err)

	err = tr.SetWithParams("/api/v1/accounts", "AccountsHandler")
	assert(err == nil, "error setting key:", "/api/v1/accounts", "error:", err)

	params := map[string]string{}
	key := "/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W"
	value, err := tr.GetWithParams(key, params)
	assert(err == nil, "failed to set key:", key, "error:", err, "params", params)
	assert(value == "ProjectHandler", "wrong value for key:", key, "got:", value, "expected:", "ProjectHandler", "params", params)
	assert(len(params) == 1 && params["project"] == "01FW1D5RWNR6MEZDJZZYJX8G2W", "invalid parameters for key:", key, "params", params)

	params = map[string]string{}
	key = "/api/v1/accounts"
	value, err = tr.GetWithParams(key, params)
	assert(err == nil, "failed to get key:", key, "error:", err, "params", params)
	assert(value == "AccountsHandler", "wrong value for key:", key, "got:", value, "expected:", "AccountsHandler", "params", params)
	assert(len(params) == 0, "invalid parameters for key:", key, "params", params)

	params = map[string]string{}
	key = "/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459"
	value, err = tr.GetWithParams(key, params)
	assert(err == nil, "failed to get key:", key, "error:", err, "params", params)
	assert(value == "InstanceHandler", "wrong value for key:", key, "got:", value, "expected:", "InstanceHandler", "params", params)
	assert(
		len(params) == 2 && params["project"] == "01FW1D5RWNR6MEZDJZZYJX8G2W" && params["instance"] == "31459",
		"invalid parameters for key:", key, "params", params,
	)

	params = map[string]string{}
	key = "/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459/databases/ordersdb"
	value, err = tr.GetWithParams(key, params)
	assert(err == nil, "failed to set key:", key, "error:", err, "params", params)
	assert(value == "DatabaseHandler", "wrong value for key:", key, "got:", value, "expected:", "DatabaseHandler", "params", params)
	assert(
		len(params) == 3 && params["project"] == "01FW1D5RWNR6MEZDJZZYJX8G2W" &&
			params["instance"] == "31459" && params["database"] == "ordersdb",
		"invalid parameters for key:", key, "params", params,
	)

	params = map[string]string{}
	key = "/api/v1/projects/01FW1D5RWNR6MEZDJZZYJX8G2W/instances/31459/applications/application1"
	value, err = tr.GetWithParams(key, params)
	assert(err == nil, "failed to get key:", key, "error:", err, "params", params)
	assert(value == "ApplicationHandler", "wrong value for key:", key, "got:", value, "expected:", "ApplicationHandler", "params", params)
	assert(
		len(params) == 3 && params["project"] == "01FW1D5RWNR6MEZDJZZYJX8G2W" &&
			params["instance"] == "31459" && params["application"] == "application1",
		"invalid parameters for key:", key, "params", params,
	)
}

func TestNeighborMatch(t *testing.T) {
	assert := newAssert(t)

	tr, err := FromMap(pairs)
	assert(err == nil, "error creating tree from map", "err:", err)

	_ = tr.Set("small", 67)
	_ = tr.Set("sma", 677)

	expect := map[string]interface{}{
		"sma": 677, "small": 67, "smaller": 81, "smallish": 82, "smart": 83,
	}

	neighboors := make(map[string]interface{})
	err = tr.NeighborMatch("smalle", neighboors)
	assert(err == nil, "error in neighbor match", "err:", err)

	for k, v := range expect {
		assert(neighboors[k] == v, "invalid result for key:", k, "expected:", v, "got:", neighboors[k])
	}
}

func TestRegressionGetSearchKeyExhaustion(t *testing.T) {
	assert := newAssert(t)

	tr, err := FromMap(pairs)
	assert(err == nil, "error creating tree from map", "err:", err)
	_ = tr.Set("small", 67)

	_, err = tr.Get("smalle")
	assert(err != nil, "get: key should not exist", "err:", err)

	match, _, err := tr.LongestMatch("smalle")
	assert(err == nil && match == "small", "longest match: invalid match", match, "err:", err)

	neighboors := make(map[string]interface{})
	err = tr.NeighborMatch("smalle", neighboors)
	assert(err == nil && len(neighboors) == 4, "neighbor match: invalid matches:", neighboors, "err:", err)
}

func TestRegressionGetWithParamsFirstParamElement(t *testing.T) {
	assert := newAssert(t)
	tr := New(WithParams(':', '@'))
	err := tr.SetWithParams("@namespace:documents:accounts:@accountId:@subscriptionId:@resourceType:@resourceId", "value")
	assert(err == nil, "error setting with params:", err)

	err = tr.SetWithParams("@namespace:files:accounts:@accountId:@subscriptionId:@resourceType:@resourceId", "value")
	assert(err == nil, "error setting with params:", err)

	params := map[string]string{}
	value, err := tr.GetWithParams("my-company:documents:accounts:E7B4320A06A1:DBCAB1AD:document:46D05077510E", params)
	assert(err == nil, "error getting with params:", err, "params", params)
	assert(value != nil && value.(string) == "value", "wrong value:", value, "params", params)
	assert(
		params["namespace"] == "my-company" &&
			params["accountId"] == "E7B4320A06A1" &&
			params["subscriptionId"] == "DBCAB1AD" &&
			params["resourceType"] == "document" &&
			params["resourceId"] == "46D05077510E",
		"wrong parameters", params)

	params = map[string]string{}
	value, err = tr.GetWithParams("my-company:files:accounts:E7B4320A06A1:DBCAB1AD:file:46D05077510E", params)
	assert(err == nil, "error getting with params:", err, "params", params)
	assert(value != nil && value.(string) == "value", "wrong value:", value, "params", params)
	assert(
		params["namespace"] == "my-company" &&
			params["accountId"] == "E7B4320A06A1" &&
			params["subscriptionId"] == "DBCAB1AD" &&
			params["resourceType"] == "file" &&
			params["resourceId"] == "46D05077510E",
		"wrong parameters", params)
}

func TestRegressionGetMultipleParamsSameKey(t *testing.T) {
	assert := newAssert(t)
	tr := New(WithParams(':', '@'))
	err := tr.SetWithParams("urn:documents:accounts:@accountId:@subscriptionId:@resourceType:@resourceId", "value")
	assert(err == nil, "error setting with params:", err)

	params := map[string]string{}
	value, err := tr.GetWithParams("urn:documents:accounts:E7B4320A06A1:DBCAB1AD:document:46D05077510E", params)
	assert(err == nil, "error getting with params:", err, "params", params)
	assert(value != nil && value.(string) == "value", "wrong value:", value, "params", params)
	assert(
		params["accountId"] == "E7B4320A06A1" &&
			params["subscriptionId"] == "DBCAB1AD" &&
			params["resourceType"] == "document" &&
			params["resourceId"] == "46D05077510E",
		"wrong parameters", params)

	err = tr.SetWithParams("urn:documents:accounts:@accountId:@subscriptionId:@resourceType:@resourceId:admin", "admin")
	assert(err == nil, "error setting with params:", err)

	params = map[string]string{}
	value, err = tr.GetWithParams("urn:documents:accounts:E7B4320A06A1:DBCAB1AD:document:46D05077510E:admin", params)
	assert(err == nil, "error getting with params:", err, "params", params)
	assert(value != nil && value.(string) == "admin", "wrong value:", value, "params", params)
	assert(
		params["accountId"] == "E7B4320A06A1" &&
			params["subscriptionId"] == "DBCAB1AD" &&
			params["resourceType"] == "document" &&
			params["resourceId"] == "46D05077510E",
		"wrong parameters", params)

	err = tr.SetWithParams("urn:documents:accounts:@accountId:@subscriptionId:@resourceType:@resourceId:admin:@adminId", "admin")
	assert(err == nil, "error setting with params:", err)

	params = map[string]string{}
	value, err = tr.GetWithParams("urn:documents:accounts:E7B4320A06A1:DBCAB1AD:document:46D05077510E:admin:XYZ", params)
	assert(err == nil, "error getting with params:", err, "params", params)
	assert(value != nil && value.(string) == "admin", "wrong value:", value, "params", params)
	assert(
		params["accountId"] == "E7B4320A06A1" &&
			params["subscriptionId"] == "DBCAB1AD" &&
			params["resourceType"] == "document" &&
			params["resourceId"] == "46D05077510E",
		params["adminId"] == "XYZ",
		"wrong parameters", params)

	params = map[string]string{}
	_, err = tr.GetWithParams("urn:documents:accounts:E7B4320A06A1", params)
	assert(err != nil, "expected key not found:", err, "params", params)
}
