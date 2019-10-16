package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"time"
	"crypto/tls"
	"os"
	"github.com/pborman/getopt/v2"
	"strings"
	"encoding/json"
//	"log"
//	"regexp"
//	"strconv"
)


/*############
## Global ##
############ */
var url string
var urldelete string
var payload = strings.NewReader("test")
var req, _ = http.NewRequest("GET", url, nil)


type Message []struct {
	Ref       string `json:"_ref"`
	Ipv4Addrs []struct {
		Ref              string `json:"_ref"`
		ConfigureForDhcp bool   `json:"configure_for_dhcp"`
		Host             string `json:"host"`
		Ipv4Addr         string `json:"ipv4addr"`
	} `json:"ipv4addrs"`
	Name string `json:"name"`
	View string `json:"view"`
}
var m Message  






/*#############
## Funcs      ##
############## */

func input_from_command_line() string {
    var u string
    var ud string
    var p string
    var action string

    optName := getopt.StringLong("hostname", 'h', "", "Hostname (FQDN ex: bvinfsd01.ii-corpnet.com)(needed for add)")
    optIP := getopt.StringLong("ipaddr", 'i', "", "IP Address (needed for add)")
    optAlias := getopt.StringLong("alias", 's', "", "alias (needed for create_alias)")
    optType := getopt.StringLong("recordtype", 'r', "", "Record Type (A,Alias,Cname)")
    optAction := getopt.StringLong("action", 'a', "", "Action (add,create_alias,delete_alias,delete,query)")
    optHelp := getopt.BoolLong("help", 0, "Help")
    getopt.Parse()

    if *optHelp {
        getopt.Usage()
        os.Exit(0)
    }


    switch *optAction {
    	case "add":
                u = "https://172.25.35.103/wapi/v2.7/record:host"
		p="{ \"name\""+":\""+*optName+"\""+","+"\"ipv4addrs\""+":[ {"+"\"ipv4addr\":"+"\""+*optIP+"\""+"} ] }"
    	case "create_alias":
		*optAction = "query"
		u = "https://172.25.35.103/wapi/v2.7/record:host?name~="+*optName
    		ud = "https://172.25.35.103/wapi/v2.7/"
                action = "create_alias"
		p="{ \"name\""+":\""+*optName+"\""+","+"\"ipv4addrs\""+":[ {"+"\"ipv4addr\":"+"\""+*optIP+"\""+"} ], "+ "\"aliases\": [ \""+*optAlias+"\" ]   }"
    	case "delete_alias":
		*optAction = "query"
		u = "https://172.25.35.103/wapi/v2.7/record:host?name~="+*optName
    		ud = "https://172.25.35.103/wapi/v2.7/"
                action = "delete_alias"
		p = "{ \"name\":\"win10.ii-corpnet.com\",\"ipv4addrs\":[ { \"ipv4addr\":\"172.25.35.205\"} ],\"aliases\": [ ] }"
    	case "delete":
		*optAction = "query"
		u = "https://172.25.35.103/wapi/v2.7/record:host?name~="+*optName
    		ud = "https://172.25.35.103/wapi/v2.7/"
                action = "delete"
    	case "query":
		u = "https://172.25.35.103/wapi/v2.7/record:host?_return_fields%2B=aliases&name="+*optName
    	}

    	return(*optName+"!"+*optIP+"!"+*optType+"!"+*optAction+"!"+u+"!"+p+"!"+ud+"!"+action)
}



func main() {


	return_slice := strings.Split(input_from_command_line(), "!")

	
	timeout := time.Duration(30 * time.Second)
	tr := &http.Transport{
		MaxIdleConnsPerHost: 10,
	}	
	client := http.Client{
		Transport: tr,
		Timeout:   timeout,
	}


	trSkipVerify := &http.Transport{
	MaxIdleConnsPerHost: 10,
	TLSClientConfig: &tls.Config {
		InsecureSkipVerify: true,
	},
	}
	client.Transport = trSkipVerify

	url = return_slice[4]

	if return_slice[3] == "add" {
		payload = strings.NewReader(return_slice[5])
        	req, _ = http.NewRequest("POST", url,  payload)
        }
	if return_slice[3] == "query" {
        	req, _ = http.NewRequest("GET", url,  nil)
        }



	req.Header.Add("cookie", "ibapauth=%22ip%3D172.22.12.75%2Cclient%3DAPI%2Cgroup%3Dadmin-group%2Cctime%3D1570019370%2Ctimeout%3D6000%2Cmtime%3D1570019370%2Csu%3D1%2Cauth%3DLOCAL%2Cuser%3Dansible%2CoZ7O5Rtv7nJSPRji%2BYUs2DmebRLfL3gnHVw%22")
        req.Header.Add("authorization", "Basic YW5zaWJsZTppbmZvYmxveA==")

	resp, err := client.Do(req)

	 if err != nil {
                 fmt.Println(err)
                 os.Exit(1)
         }

        body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &m)



        fmt.Println(resp)
	fmt.Printf(string(body))



	urldelete = return_slice[6]

	if return_slice[7] == "delete" {
		fmt.Println("Deleting host ....")
        	req, err = http.NewRequest("DELETE", urldelete+m[0].Ref,  nil)
	 	if err != nil {
                 	fmt.Println(err)
                 	os.Exit(1)
         	}
        
	req.Header.Add("cookie", "ibapauth=%22ip%3D172.22.12.75%2Cclient%3DAPI%2Cgroup%3Dadmin-group%2Cctime%3D1570019370%2Ctimeout%3D6000%2Cmtime%3D1570019370%2Csu%3D1%2Cauth%3DLOCAL%2Cuser%3Dansible%2CoZ7O5Rtv7nJSPRji%2BYUs2DmebRLfL3gnHVw%22")
        req.Header.Add("authorization", "Basic YW5zaWJsZTppbmZvYmxveA==")
	respd, errd := client.Do(req)
	if errd != nil {
                 fmt.Println(err)
                 os.Exit(1)
        }

        fmt.Println(respd)
	}




/*#############
#   Create Alias
###############*/

	if return_slice[7] == "create_alias" {
		fmt.Println("Creating alais ....", return_slice[5])
		payload = strings.NewReader(return_slice[5])
        	req, err = http.NewRequest("PUT", urldelete+m[0].Ref,  payload)
	 	if err != nil {
                 	fmt.Println(err)
                 	os.Exit(1)
         	}
        
	req.Header.Add("cookie", "ibapauth=%22ip%3D172.22.12.75%2Cclient%3DAPI%2Cgroup%3Dadmin-group%2Cctime%3D1570019370%2Ctimeout%3D6000%2Cmtime%3D1570019370%2Csu%3D1%2Cauth%3DLOCAL%2Cuser%3Dansible%2CoZ7O5Rtv7nJSPRji%2BYUs2DmebRLfL3gnHVw%22")
        req.Header.Add("authorization", "Basic YW5zaWJsZTppbmZvYmxveA==")
	respd, errd := client.Do(req)
	if errd != nil {
                 fmt.Println(err)
                 os.Exit(1)
        }

        fmt.Println(respd)
	}



/*#############
#  Delete Alias
###############*/

	if return_slice[7] == "delete_alias" {
		fmt.Println("hello there delete alias")
		payload = strings.NewReader(return_slice[5])
        	req, err = http.NewRequest("PUT", urldelete+m[0].Ref,  payload)
	 	if err != nil {
                 	fmt.Println(err)
                 	os.Exit(1)
         	}
        
	req.Header.Add("cookie", "ibapauth=%22ip%3D172.22.12.75%2Cclient%3DAPI%2Cgroup%3Dadmin-group%2Cctime%3D1570019370%2Ctimeout%3D6000%2Cmtime%3D1570019370%2Csu%3D1%2Cauth%3DLOCAL%2Cuser%3Dansible%2CoZ7O5Rtv7nJSPRji%2BYUs2DmebRLfL3gnHVw%22")
        req.Header.Add("authorization", "Basic YW5zaWJsZTppbmZvYmxveA==")
	respd, errd := client.Do(req)
	if errd != nil {
                 fmt.Println(err)
                 os.Exit(1)
        }

        fmt.Println(respd)
	}
}
