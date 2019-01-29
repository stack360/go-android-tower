package client

import (
    "bytes"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "encoding/json"
)

type Client struct {
    BaseURL    *url.URL
    httpClient *http.Client
}

func NewClient(path string) *Client {
    httpClient := http.DefaultClient
    u, _ := url.Parse(path)
    c := &Client {BaseURL: u, httpClient: httpClient}
    return c
}

func (c *Client) RunFunction(funcName string, params []byte) (string, error) {
    body := make(map[string]interface{})
    json.Unmarshal(params, &body)
    req, err := c.newRequest("POST", "/api/run_function", body)
    if err != nil {
        return "Error encountered.", err
    }
    resp, err := c.do(req)
    data, _ := ioutil.ReadAll(resp.Body)
    return string(data), nil
}

func (c *Client) ListDevices() (string, error){
    req, err := c.newRequest("GET", "/api/devices", nil)
    if err != nil {
        return "Error encountered.", err
    }
    resp, err := c.do(req)
    data, _ := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()
    return string(data), nil
}

func (c *Client) newRequest(method string, path string, body map[string]interface{}) (*http.Request, error) {
    rel := &url.URL{Path: path}
    u := c.BaseURL.ResolveReference(rel)
    var buf io.ReadWriter
    if body != nil {
        buf = new(bytes.Buffer)
        err := json.NewEncoder(buf).Encode(body)
        if err != nil {
            return nil, err
        }
    }
    req, err := http.NewRequest(method, u.String(), buf)
    if err != nil {
        return nil, err
    }
    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }
    req.Header.Set("Accept", "application/json")
    return req, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}
