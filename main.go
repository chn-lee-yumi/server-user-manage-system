package main

import (
    "fmt"
    "net/http"
    "log"
    "strings"
    //"io"
    "io/ioutil"
    "encoding/hex"
    //"encoding/json"
    "crypto/sha1"
    "time"
    "encoding/binary"
    "os/exec"
    "net/url"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) { // 首页处理
    log.Println("req.URL.Path: ",r.URL.Path)
    s,_:=ioutil.ReadFile("index.html")
    w.Write(s)
}

/*func Res(w http.ResponseWriter, r *http.Request) { // 静态资源文件处理
    log.Println("req.URL.Path: ",r.URL.Path)
    s,_:=ioutil.ReadFile(strings.Trim(r.URL.Path,"/"))
    w.Write(s)
}*/

type AddUserPost struct {
    Ticket string //新建用户必须的凭证
    Username string //用户名
    PublicKey string //公钥
}

func AddUser(w http.ResponseWriter, r *http.Request) { // API，新建用户
    //log.Println("req: ",r)
    log.Println("req.URL.Path: ",r.URL.Path)
    log.Println("req.Method: ",r.Method)
    if r.Method != "POST" {
        fmt.Fprintln(w,`{"code":-1,"msg":"请使用POST"}`)
        return
    }
    body,_:=ioutil.ReadAll(r.Body)
    //log.Println("Body: ",string(body))
    url_parsed, err:=url.ParseQuery(string(body))
    if err!=nil {
        fmt.Printf("URL参数解析错误：ParseQuery", err.Error())
        fmt.Fprintln(w,`{"code":-2,"msg":"URL参数解析错误：ParseQuery"}`)
        return
    }
    //没有处理key不存在的情况
    log.Println(url_parsed["Ticket"][0])
    log.Println(url_parsed["Username"][0])
    log.Println(url_parsed["PublicKey"][0])
    /*jsonStr, err := json.Marshal(url_parsed)
    if err!=nil {
        fmt.Printf("URL参数解析错误：json.Marshal", err.Error())
        fmt.Fprintln(w,`{"code":-2,"msg":"URL参数解析错误：json.Marshal"}`)
        return
    }
    log.Println(jsonStr)
    d:=AddUserPost{}
    err = json.Unmarshal(jsonStr, &d)
    if err!=nil {
        fmt.Printf("URL参数解析错误：json.Unmarshal", err.Error())
        fmt.Fprintln(w,`{"code":-2,"msg":"URL参数解析错误：json.Unmarshal"}`)
        return
    }*/
    d:=&AddUserPost{Ticket: url_parsed["Ticket"][0],
        Username: url_parsed["Username"][0],
        PublicKey: url_parsed["PublicKey"][0]}

    //log.Println("req.Body: ",d)
    if d.Ticket=="" || d.Username=="" || d.PublicKey=="" {
        fmt.Fprintln(w,`{"code":-3,"msg":"信息不完整"}`)
        return
    }
    if _,ok:=tickets[d.Ticket]; !ok{
        fmt.Fprintln(w,`{"code":-4,"msg":"凭据错误"}`)
        return
    }
    //校验公钥格式
    tmp_ticket:=strings.Split(d.PublicKey, " ")
    log.Println(tmp_ticket)
    if len(tmp_ticket)!=3{
        fmt.Fprintln(w,`{"code":-6,"msg":"公钥格式错误"}`)
        return
    }
    if strings.IndexAny(tmp_ticket[0], "ssh-")!=0{
        fmt.Fprintln(w,`{"code":-6,"msg":"公钥格式错误"}`)
        return
    }

    //下面开始新建linux用户
    cmd := exec.Command("useradd", "-m", d.Username, "-p", "*")
    output, err := cmd.Output()
    if err != nil {
        fmt.Printf("创建用户错误", err.Error())
        fmt.Fprintln(w,`{"code":-5,"msg":"创建用户错误"}`)
        return
    }
    fmt.Printf(string(output))
    //新建.ssh目录
    cmd = exec.Command("mkdir", "/home/"+d.Username+"/.ssh")
    output, err = cmd.Output()
    if err != nil {
        fmt.Printf("创建.ssh目录错误:%s", err.Error())
        fmt.Fprintln(w,`{"code":-5,"msg":"创建.ssh目录错误"}`)
        return
    }
    fmt.Printf(string(output))
    //修改.ssh目录权限
    cmd = exec.Command("chmod", "700", "/home/"+d.Username+"/.ssh")
    output, err = cmd.Output()
    if err != nil {
        fmt.Printf("修改.ssh目录权限错误:%s", err.Error())
        fmt.Fprintln(w,`{"code":-5,"msg":"修改.ssh目录权限错误"}`)
        return
    }
    fmt.Printf(string(output))
    //创建authorized_keys
    err = ioutil.WriteFile("/home/"+d.Username+"/.ssh/authorized_keys",[]byte(d.PublicKey),0600)
    if err != nil {
        fmt.Printf("创建authorized_keys错误:%s", err.Error())
        fmt.Fprintln(w,`{"code":-5,"msg":"创建authorized_keys错误"}`)
        return
    }
    fmt.Println("创建authorized_keys成功")
    //修改所属用户和组
    cmd = exec.Command("chown", "-R", d.Username+":"+d.Username, "/home/"+d.Username+"/.ssh")
    output, err = cmd.Output()
    if err != nil {
        fmt.Printf("修改所属用户和组错误:%s", err.Error())
        fmt.Fprintln(w,`{"code":-5,"msg":"修改所属用户和组错误"}`)
        return
    }
    fmt.Printf(string(output))

    delete(tickets,d.Ticket)
    log.Println("Del Ticket: ",d.Ticket)
    fmt.Fprintln(w,`{"code":0,"msg":"用户创建成功"}`)
}


func NewTicket(w http.ResponseWriter, r *http.Request) { // API，新建用户
    //log.Println("req: ",r)
    log.Println("req.URL.Path: ",r.URL.Path)
    timestamp:=time.Now().UnixNano()
    buf := make([]byte, 8)
    binary.BigEndian.PutUint64(buf, uint64(timestamp))
    h := sha1.New()
    h.Write(buf)
    ticket:=hex.EncodeToString(h.Sum(nil))
    log.Println("Add Ticket: ",ticket)
    tickets[ticket]=1
    fmt.Fprintln(w,ticket)
}

var tickets map[string]int //存放可用Ticket的字典

func main() {
    tickets=make(map[string]int)
    http.HandleFunc("/", IndexHandler)
    //http.Handle("/res/", http.HandlerFunc(Res))
    //http.Handle("/res/", http.FileServer(http.Dir("."))) // 静态资源文件处理
    http.HandleFunc("/api/adduser", AddUser)
    http.HandleFunc("/api/newticket", NewTicket)
    err := http.ListenAndServe(":8000", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
