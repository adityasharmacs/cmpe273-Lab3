package main
import  ("github.com/julienschmidt/httprouter"
    "fmt"
    "net/http"
    "strconv"
    "encoding/json"
    "strings"
    "sort")


type KeyValue struct{
  Key int `json:"key,omitempty"`
  Value string  `json:"value,omitempty"`
} 


var keyVal1,keyVal2,keyVal3 [] KeyValue
var indexer1,indexer2,indexer3 int
type ByKey []KeyValue
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }


func GetKeys(rw http.ResponseWriter, request *http.Request,p httprouter.Params){
  portArray := strings.Split(request.Host,":")
  if(portArray[1]=="3000"){
    sort.Sort(ByKey(keyVal1))
    result,_:= json.Marshal(keyVal1)
    fmt.Fprintln(rw,string(result))
  }else if(portArray[1]=="3001"){
    sort.Sort(ByKey(keyVal2))
    result,_:= json.Marshal(keyVal2)
    fmt.Fprintln(rw,string(result))
  }else{
    sort.Sort(ByKey(keyVal3))
    result,_:= json.Marshal(keyVal3)
    fmt.Fprintln(rw,string(result))
  }
}

func PutKeys(rw http.ResponseWriter, request *http.Request,p httprouter.Params){
  portArray := strings.Split(request.Host,":")
  key,_ := strconv.Atoi(p.ByName("key_id"))
  if(portArray[1]=="3000"){
    keyVal1 = append(keyVal1,KeyValue{key,p.ByName("value")})
    indexer1++
  }else if(portArray[1]=="3001"){
    keyVal2 = append(keyVal2,KeyValue{key,p.ByName("value")})
    indexer2++
  }else{
    keyVal3 = append(keyVal3,KeyValue{key,p.ByName("value")})
    indexer3++
  } 
}

func GetKey(rw http.ResponseWriter, request *http.Request,p httprouter.Params){ 
  out := keyVal1
  ind := indexer1
  portArray := strings.Split(request.Host,":")
  if(portArray[1]=="3001"){
    out = keyVal2 
    ind = indexer2
  }else if(portArray[1]=="3002"){
    out = keyVal3
    ind = indexer3
  } 
  key,_ := strconv.Atoi(p.ByName("key_id"))
  for i:=0 ; i< ind ;i++{
    if(out[i].Key==key){
      result,_:= json.Marshal(out[i])
      fmt.Fprintln(rw,string(result))
    }
  }
}



func main(){
  indexer1 = 0
  indexer2 = 0
  indexer3 = 0
  mux := httprouter.New()
    mux.GET("/keys",GetKeys)
    mux.GET("/keys/:key_id",GetKey)
    mux.PUT("/keys/:key_id/:value",PutKeys)
    go http.ListenAndServe(":3000",mux)
    go http.ListenAndServe(":3001",mux)
    go http.ListenAndServe(":3002",mux)
    select {}
}
