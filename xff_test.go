package xff

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse_none(t *testing.T) {
	res := Parse("")
	assert.Equal(t, "", res)
}

func TestParse_localhost(t *testing.T) {
	res := Parse("127.0.0.1")
	assert.Equal(t, "", res)
}

func TestParse_localnet(t *testing.T) {
	res := Parse("fe80::2acf:e9ff:fe16:b5cd")
	assert.Equal(t, "", res)
}

func TestParse_invalid(t *testing.T) {
	res := Parse("invalid")
	assert.Equal(t, "", res)
}

func TestParse_invalid_sioux(t *testing.T) {
	res := Parse("123#1#2#3")
	assert.Equal(t, "", res)
}

func TestParse_invalid_private_lookalike(t *testing.T) {
	res := Parse("102.3.2.1")
	assert.Equal(t, "102.3.2.1", res)
}

func TestParse_valid(t *testing.T) {
	res := Parse("68.45.152.220")
	assert.Equal(t, "68.45.152.220", res)
}

func TestParse_multi_first(t *testing.T) {
	res := Parse("12.13.14.15, 68.45.152.220")
	assert.Equal(t, "12.13.14.15", res)
}

func TestParse_multi_last(t *testing.T) {
	res := Parse("192.168.110.162, 190.57.149.90")
	assert.Equal(t, "190.57.149.90", res)
}

func TestParse_multi_with_invalid(t *testing.T) {
	res := Parse("192.168.110.162, invalid, 190.57.149.90")
	assert.Equal(t, "190.57.149.90", res)
}

func TestParse_multi_with_invalid2(t *testing.T) {
	res := Parse("192.168.110.162, 190.57.149.90, invalid")
	assert.Equal(t, "190.57.149.90", res)
}

func TestParse_multi_with_invalid_sioux(t *testing.T) {
	res := Parse("192.168.110.162, 190.57.149.90, 123#1#2#3")
	assert.Equal(t, "190.57.149.90", res)
}

func TestGetRemoteAddr(t *testing.T) {
	assert.Equal(t, "1.2.3.4:1234", GetRemoteAddr(&http.Request{RemoteAddr: "1.2.3.4:1234"}))
	assert.Equal(t, "[2001:db8:0:1:1:1:1:1]:1234", GetRemoteAddr(&http.Request{RemoteAddr: "[2001:db8:0:1:1:1:1:1]:1234"}))
	assert.Equal(t, "[2001:db8:0:1:1:1:1:1]:1234", GetRemoteAddr(&http.Request{
		RemoteAddr: "1.2.3.4:1234",
		Header:     http.Header{"X-Forwarded-For": []string{"2001:db8:0:1:1:1:1:1"}},
	}))
	assert.Equal(t, "1.2.3.4:4321", GetRemoteAddr(&http.Request{
		RemoteAddr: "1.2.3.4:1234",
		Header: http.Header{
			"X-Forwarded-Port": []string{"4321"},
		},
	}))

}
