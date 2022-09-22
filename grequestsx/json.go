package grequestsx

import "github.com/levigross/grequests"

func JsonGet(url string, ro *grequests.RequestOptions, out interface{}, flags ...Flags) error {
	resp, err := DoRegularRequest("GET", url, ro, flags...)
	if err != nil {
		return err
	}
	return resp.JSON(out)
}

func JsonPut(url string, ro *grequests.RequestOptions, out interface{}, flags ...Flags) error {
	resp, err := DoRegularRequest("PUT", url, ro, flags...)
	if err != nil {
		return err
	}
	return resp.JSON(out)
}

func JsonPatch(url string, ro *grequests.RequestOptions, out interface{}, flags ...Flags) error {
	resp, err := DoRegularRequest("PATCH", url, ro, flags...)
	if err != nil {
		return err
	}
	return resp.JSON(out)
}

func JsonDelete(url string, ro *grequests.RequestOptions, out interface{}, flags ...Flags) error {
	resp, err := DoRegularRequest("DELETE", url, ro, flags...)
	if err != nil {
		return err
	}
	return resp.JSON(out)
}

func JsonPost(url string, ro *grequests.RequestOptions, out interface{}, flags ...Flags) error {
	resp, err := DoRegularRequest("POST", url, ro, flags...)
	if err != nil {
		return err
	}
	return resp.JSON(out)
}

func JsonHead(url string, ro *grequests.RequestOptions, out interface{}, flags ...Flags) error {
	resp, err := DoRegularRequest("HEAD", url, ro, flags...)
	if err != nil {
		return err
	}
	return resp.JSON(out)
}

func JsonOptions(url string, ro *grequests.RequestOptions, out interface{}, flags ...Flags) error {
	resp, err := DoRegularRequest("OPTIONS", url, ro, flags...)
	if err != nil {
		return err
	}
	return resp.JSON(out)
}

func JsonReq(verb string, url string, ro *grequests.RequestOptions, out interface{}, flags ...Flags) error {
	resp, err := DoRegularRequest(verb, url, ro, flags...)
	if err != nil {
		return err
	}
	return resp.JSON(out)
}
