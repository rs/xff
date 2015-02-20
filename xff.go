package xff

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

var privateMasks = func() []net.IPNet {
	masks := []net.IPNet{}
	for _, cidr := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "fc00::/7"} {
		_, net, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}
		masks = append(masks, *net)
	}
	return masks
}()

// IsPublicIP returns true if the given IP can be routed on the Internet
func IsPublicIP(ip net.IP) bool {
	if !ip.IsGlobalUnicast() {
		return false
	}
	for _, mask := range privateMasks {
		if mask.Contains(ip) {
			return false
		}
	}
	return true
}

// Parse parses the X-Forwarded-For Header and returns the IP address.
func Parse(ipList string) string {
	for _, ip := range strings.Split(ipList, ",") {
		ip = strings.TrimSpace(ip)
		if IP := net.ParseIP(ip); IP != nil && IsPublicIP(IP) {
			return ip
		}
	}
	return ""
}

func parseXFP(port string) string {
	return port
}

// XFF is a middleware to update RemoteAdd from X-Fowarded-* Headers.
func XFF(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		xff := r.Header.Get("X-Forwarded-For")
		xfp := r.Header.Get("X-Forwarded-Port")
		var ip string
		if xff != "" {
			ip = Parse(xff)
		}
		var port string
		if xfp != "" {
			port = Parse(xfp)
		}
		if ip != "" && port != "" {
			r.RemoteAddr = fmt.Sprintf("%s:%s", ip, port)
		} else {
			oip, oport, err := net.SplitHostPort(r.RemoteAddr)
			if err == nil {
				if ip != "" {
					r.RemoteAddr = fmt.Sprintf("%s:%s", ip, oport)

				} else if port != "" {
					r.RemoteAddr = fmt.Sprintf("%s:%s", oip, port)

				}
			}
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}