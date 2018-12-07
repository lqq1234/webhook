package main

import(
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"bytes"
	"strings"
	"os"
	"log"
)

type HarborEvent struct {
	Events []struct{
		Action  string `json:"action"`
		Target  struct{
				Repository    string   `json:"repository"`
				Tag           string   `json:"tag"`
				Url           string   `json:"url"`
			}
	}   `json:"events"`
}


type ImageHook struct {
	Pushdata struct{Tag string `json:"tag"`} `json:"push_data"`
	Respository struct{RepoName string `json:"repo_name"`} `json:"repository"`
	Host	string	`json:"host"`
}


func main() {

	http.HandleFunc("/image/harbor/webhook", func(w http.ResponseWriter, r *http.Request) {
		bts,_ := ioutil.ReadAll(r.Body)
		var harborEvent HarborEvent
		json.Unmarshal(bts,&harborEvent)
		for _,v := range harborEvent.Events{
			var imagehook ImageHook
			if v.Action == "push" && len(v.Target.Tag) > 0 {
				rs := strings.Split(v.Target.Url,"/")
				imagehook.Pushdata.Tag = v.Target.Tag
				imagehook.Respository.RepoName = v.Target.Repository
				imagehook.Host = rs[2]
				SendImagePushEvent(imagehook)
			}
		}
	})
	log.Fatal(http.ListenAndServe(":8089", nil))
}

func SendImagePushEvent(imageHook ImageHook)(err error){

	hookurl := os.Getenv("hookurl")
	//hookurl := "http://127.0.0.1:8080/v1/image/webhook"
	bts,_ := json.Marshal(imageHook)

	client := http.Client{}
	req,_ := http.NewRequest("POST",hookurl,bytes.NewReader(bts))
	resp,err := client.Do(req)

	respbts,_ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(respbts))

	return
}
