package main  

  
import (  
    "fmt"  
    "hash/crc32"  
    "sort"     
    "net/http"
    "encoding/json" 
    "io/ioutil"
    "os"
    "strings"
)  
   
type HashingFunction []uint32  

type KeyValue struct{
    Key int `json:"key,omitempty"`
    Value string `json:"value,omitempty"`
}



func (hashCircle HashingFunction) Len() int {  
    return len(hashCircle)  
}  
  
func (hashCircle HashingFunction) Less(i, j int) bool {  
    return hashCircle[i] < hashCircle[j]  
}  
  
func (hashCircle HashingFunction) Swap(i, j int) {  
    hashCircle[i], hashCircle[j] = hashCircle[j], hashCircle[i]  
}  
  
type Node struct {  
    Id       int  
    IP       string    
}  
  
func AddNode(id int, ip string) *Node {  
    return &Node{  
        Id:       id,  
        IP:       ip,  
    }  
}  
  
type HashCircleConsistent struct {  
    Nodes       map[uint32]Node  
    IsPresent   map[int]bool  
    Circle      HashingFunction  
    
}  
  
func NewHash() *HashCircleConsistent {  
    return &HashCircleConsistent{  
        Nodes:     make(map[uint32]Node),   
        IsPresent: make(map[int]bool),  
        Circle:      HashingFunction{},  
    }  
}  
  
func (hashCircle *HashCircleConsistent) AddNode(node *Node) bool {  
 
    if _, ok := hashCircle.IsPresent[node.Id]; ok {  
        return false  
    }  
    str := hashCircle.NodeIP(node)  
    hashCircle.Nodes[hashCircle.HashValue(str)] = *(node)
    hashCircle.IsPresent[node.Id] = true  
    hashCircle.SortHashCircle()  
    return true  
}  
  
func (hashCircle *HashCircleConsistent) SortHashCircle() {  
    hashCircle.Circle = HashingFunction{}  
    for k := range hashCircle.Nodes {  
        hashCircle.Circle = append(hashCircle.Circle, k)  
    }  
    sort.Sort(hashCircle.Circle)  
}  
  
func (hashCircle *HashCircleConsistent) NodeIP(node *Node) string {  
    return node.IP 
}  
  
func (hashCircle *HashCircleConsistent) HashValue(key string) uint32 {  
    return crc32.ChecksumIEEE([]byte(key))  
}  
  
func (hashCircle *HashCircleConsistent) Get(key string) Node {  
    hash := hashCircle.HashValue(key)  
    i := hashCircle.SearchNode(hash)  
    return hashCircle.Nodes[hashCircle.Circle[i]]  
}  

func (hashCircle *HashCircleConsistent) SearchNode(hash uint32) int {  
    i := sort.Search(len(hashCircle.Circle), func(i int) bool {return hashCircle.Circle[i] >= hash })  
    if i < len(hashCircle.Circle) {  
        if i == len(hashCircle.Circle)-1 {  
            return 0  
        } else {  
            return i  
        }  
    } else {  
        return len(hashCircle.Circle) - 1  
    }  
}  
  
func PutKeyValue(circle *HashCircleConsistent, str string, input string){
        ipAddress := circle.Get(str)  
        address := "http://"+ipAddress.IP+"/keys/"+str+"/"+input
		fmt.Println(address)
        req,err := http.NewRequest("PUT",address,nil)
        client := &http.Client{}
        resp, err := client.Do(req)
        if err!=nil{
            fmt.Println("Error:",err)
        }else{
            defer resp.Body.Close()
            fmt.Println("PUT Request successfully completed")
        }  
}  

func GetKeyValue(key string,circle *HashCircleConsistent){
    var out KeyValue 
    ipAddress:= circle.Get(key)
	address := "http://"+ipAddress.IP+"/keys/"+key
	fmt.Println(address)
    response,err:= http.Get(address)
    if err!=nil{
        fmt.Println("Error:",err)
    }else{
        defer response.Body.Close()
        contents,err:= ioutil.ReadAll(response.Body)
        if(err!=nil){
            fmt.Println(err)
        }
        json.Unmarshal(contents,&out)
        result,_:= json.Marshal(out)
        fmt.Println(string(result))
    }
}

func GetAllKeys(address string){
     
    var out []KeyValue
    response,err:= http.Get(address)
    if err!=nil{
        fmt.Println("Error:",err)
    }else{
        defer response.Body.Close()
        contents,err:= ioutil.ReadAll(response.Body)
        if(err!=nil){
            fmt.Println(err)
        }
        json.Unmarshal(contents,&out)
        result,_:= json.Marshal(out)
        fmt.Println(string(result))
    }
}
func main() {   
    circle := NewHash()      
    circle.AddNode(AddNode(0, "127.0.0.1:3000"))
	circle.AddNode(AddNode(1, "127.0.0.1:3001"))
	circle.AddNode(AddNode(2, "127.0.0.1:3002")) 
	
	if(os.Args[1]=="PUT"){
		key := strings.Split(os.Args[2],"/")
        PutKeyValue(circle,key[0],key[1])
    } else if ((os.Args[1]=="GET") && len(os.Args)==3){
    	GetKeyValue(os.Args[2],circle)
    } else {
		GetAllKeys("http://127.0.0.1:3000/keys")
	    GetAllKeys("http://127.0.0.1:3001/keys")
	    GetAllKeys("http://127.0.0.1:3002/keys")
	}
}  
