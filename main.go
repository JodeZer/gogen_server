package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type H5PayRequest struct {
	Attach      string `json:"attach,omitempty" url:"attach,omitempty"`       //透传参数
	Busicd      string `json:"busicd" url:"busicd"`                           //交易类型
	BackUrl     string `json:"backUrl,omitempty" url:"backUrl,omitempty"`     //异步通知接收地址
	Chcd        string `json:"chcd" url:"chcd"`                               //支付渠道代码   WXP:ALP
	FrontUrl    string `json:"frontUrl" url:"frontUrl"`                       //支付成功或失败后跳转到的url
	GoodsInfo   string `json:"goodsInfo,omitempty" url:"goodsInfo,omitempty"` //订单商品详情
	Mchntid     string `json:"mchntid" url:"mchntid"`                         //商户号
	OrderNum    string `json:"orderNum" url:"orderNum"`                       //订单号
	Txamt       string `json:"txamt" url:"txamt"`                             //订单金额
	Sign        string `json:"sign" url:"-"`                                  //签名
	Terminalid  string `json:"terminalid" url:"terminalid"`                   //终端号
	Version     string `json:"version" url:"version"`                         //版本号
	OutOrderNum string `json:"outOrderNum,omitempty" url:"outOrderNum,omitempty"`
	SignType    string `json:"signType,omitempty" url:"signType,omitempty"`
	Charset     string `json:"charset,omitempty" url:"charset,omitempty"`
}

var Domain string

const signKey = "zsdfyreuoyamdphhaweyrjbvzkgfdycs"

func main() {
	filepath := flag.String("f", "example.conf", "gen url by conf")
	flag.Parse()

	h5 := encapConfigData(*filepath)
	genUrl(h5)

}
func encapConfigData(path string) (h5 *H5PayRequest) {
	node := "h5"
	c := Config{}
	c.InitConfig(path)
	//for k,v:=range c.Mymap{
	//	fmt.Printf("k:%s v:%s\n",k,v)
	//}
	h5 = &H5PayRequest{}
	h5.OrderNum = strconv.Itoa(int(time.Now().UnixNano()))
	h5.Attach = c.Read(node, "attach")
	h5.BackUrl = c.Read(node, "backUrl")
	h5.Busicd = c.Read(node, "busicd")
	h5.Chcd = c.Read(node, "chcd")
	h5.FrontUrl = c.Read(node, "frontUrl")
	h5.GoodsInfo = c.Read(node, "goodsInfo")
	h5.Mchntid = c.Read(node, "mchntid")
	h5.Txamt = c.Read(node, "txamt")
	h5.Terminalid = c.Read(node, "terminalid")
	h5.Version = c.Read(node, "version")
	h5.Charset = c.Read(node, "charset")
	h5.OutOrderNum = c.Read(node, "outOrderNum")
	h5.SignType = c.Read(node, "signType")
	theSignKey := c.Read(node, "signKey")
	Domain = c.Read(node, "domain")
	h5.Sign = signWithSha(h5, theSignKey)
	return

}
func genUrl(h5 *H5PayRequest) {
	bytes, _ := json.Marshal(h5)
	fmt.Printf("json: ==>%s\n", bytes)
	b64str := base64.StdEncoding.EncodeToString(bytes)
	url := Domain + "?data=" + b64str
	fmt.Println(url + "\n")
}
func signWithSha(s interface{}, signKey string) string {
	signBuffer, _ := Query(s)
	signString := signBuffer.String() + signKey
	fmt.Printf("sign string==> %s\n", signString)
	//h:=sha1.New()
	//h.Write([]byte(signString))
	//signBytes:=h.Sum(nil)
	if s.(*H5PayRequest).Version == "2.0" {
		signString = fmt.Sprintf("%x", sha256.Sum256([]byte(signString)))
	} else {
		signString = fmt.Sprintf("%x", sha1.Sum([]byte(signString)))
	}

	//fmt.Printf("sign\n  %v\n\n",signBytes)
	return signString
}
func main1() {
	chcd := flag.String("c", "WXP", "channel code")
	env := flag.String("e", "test", "local | test | product address")

	merid := flag.String("m", "000000000000012", "merchant id")
	flag.Parse()
	cmd := flag.Arg(0)
	fmt.Printf("%s %s %s %s\n", *chcd, *merid, *env, cmd)

	reqData := &H5PayRequest{
		Busicd:   "WPAY",
		BackUrl:  "http://www.baidu.com",
		Chcd:     *chcd,
		FrontUrl: "http://www.baidu.com",
		//Mchntid:"999118888881312",
		//Mchntid:"100000000010001",
		Mchntid:  *merid,
		OrderNum: strconv.Itoa(int(time.Now().UnixNano())),
		Txamt:    "000000000001",
		//		Txdir:"Q",
		Terminalid: "00000001",
		Version:    "1.0",
	}

	reqData.Sign = signWithSha(reqData, signKey)

	bytes, _ := json.Marshal(reqData)
	fmt.Printf("json: ==>%s\n", bytes)
	b64str := base64.StdEncoding.EncodeToString(bytes)
	var url string
	//fmt.Printf("encode ==>%s\n",b64str)
	switch *env {
	case "test":
		url = "http://test.quick.ipay.so/scanpay/unified?data=" + b64str
	case "product":
		url = "http://showmoney.cn/scanpay/unified?data=" + b64str
	case "local":
		url = "http://10.30.1.195:6800/scanpay/unified?data=" + b64str

	}
	//url:="http://showmoney.cn/scanpay/unified?data=" + b64str
	//url:="http://10.30.1.195:6800/scanpay/unified?data=" + b64str
	//fmt.Sprintf(url,b64str)
	fmt.Println(url + "\n")
	//req,_:=http.NewRequest("GET",url,nil)
	//client:=&http.Client{}
	//response, _:=client.Do(req)
	//defer response.Body.Close()
	//rspByte,_:=ioutil.ReadAll(response.Body)
	//fmt.Printf("response string  %s",string(rspByte))
	//http://test.quick.ipay.so/scanpay/unified?data=eyJidXNpY2QiOiJXUEFZIiwiYmFja1VybCI6Imh0dHA6Ly93d3cuYmFpZHUuY29tIiwiY2hjZCI6IldYUCIsImZyb250VXJsIjoiaHR0cDovL3d3dy5iYWlkdS5jb20iLCJtY2hudGlkIjoiMDAwMDAwMDAwMDAwMDEyIiwib3JkZXJOdW0iOiIxNDczODM0NDUxNjUwNTg5MDAwIiwidHhhbXQiOiIwMDAwMDAwMDAwMDEiLCJ0eGRpciI6IlEiLCJzaWduIjoiYmE2Yzg2NmI4ODA3MGM2ZTg1MTMyOTY1MWFjMzAwN2M4Mjk2MWI0MCIsInRlcm1pbmFsaWQiOiIwMDAwMDAwMSIsInZlcnNpb24iOiIxLjAifQ==
	//fmt.Println(signWithSha1(url,""))

}

func Query(s interface{}, excludes ...string) (buf bytes.Buffer, err error) {
	if s == nil {
		return
	}

	v, err := Values(s)
	if err != nil {
		return buf, err
	}

	return QueryValues(v), nil
}

// QueryValues implements encoding of values into URL query parameters without escape
func QueryValues(v url.Values, excludes ...string) (buf bytes.Buffer) {
	if v == nil {
		return
	}

	keys := make([]string, 0, len(v))
	for k := range v {
		if StringInSlice(k, excludes) {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		vs := v[k]
		prefix := k + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf

}
func StringInSlice(a string, list []string) bool {
	if len(list) == 0 {
		return false
	}
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

var timeType = reflect.TypeOf(time.Time{})

var encoderType = reflect.TypeOf(new(Encoder)).Elem()

type Encoder interface {
	EncodeValues(key string, v *url.Values) error
}

func Values(v interface{}) (url.Values, error) {
	values := make(url.Values)
	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return values, nil
		}
		val = val.Elem()
	}

	if v == nil {
		return values, nil
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("query: Values() expects struct input. Got %v", val.Kind())
	}

	err := reflectValue(values, val, "")
	return values, err
}

func reflectValue(values url.Values, val reflect.Value, scope string) error {
	var embedded []reflect.Value

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous { // unexported
			continue
		}

		sv := val.Field(i)
		tag := sf.Tag.Get("url")
		if tag == "-" {
			continue
		}
		name, opts := parseTag(tag)
		if name == "" {
			if sf.Anonymous && sv.Kind() == reflect.Struct {
				// save embedded struct for later processing
				embedded = append(embedded, sv)
				continue
			}

			name = sf.Name
		}

		if scope != "" {
			name = scope + "[" + name + "]"
		}

		if opts.Contains("omitempty") && isEmptyValue(sv) {
			continue
		}

		if sv.Type().Implements(encoderType) {
			if !reflect.Indirect(sv).IsValid() {
				sv = reflect.New(sv.Type().Elem())
			}

			m := sv.Interface().(Encoder)
			if err := m.EncodeValues(name, &values); err != nil {
				return err
			}
			continue
		}

		if sv.Kind() == reflect.Slice || sv.Kind() == reflect.Array {
			var del byte
			if opts.Contains("comma") {
				del = ','
			} else if opts.Contains("space") {
				del = ' '
			} else if opts.Contains("semicolon") {
				del = ';'
			} else if opts.Contains("brackets") {
				name = name + "[]"
			}

			if del != 0 {
				s := new(bytes.Buffer)
				first := true
				for i := 0; i < sv.Len(); i++ {
					if first {
						first = false
					} else {
						s.WriteByte(del)
					}
					s.WriteString(valueString(sv.Index(i), opts))
				}
				values.Add(name, s.String())
			} else {
				for i := 0; i < sv.Len(); i++ {
					k := name
					if opts.Contains("numbered") {
						k = fmt.Sprintf("%s%d", name, i)
					}
					values.Add(k, valueString(sv.Index(i), opts))
				}
			}
			continue
		}

		if sv.Type() == timeType {
			values.Add(name, valueString(sv, opts))
			continue
		}

		for sv.Kind() == reflect.Ptr {
			if sv.IsNil() {
				break
			}
			sv = sv.Elem()
		}

		if sv.Kind() == reflect.Struct {
			reflectValue(values, sv, name)
			continue
		}

		values.Add(name, valueString(sv, opts))
	}

	for _, f := range embedded {
		if err := reflectValue(values, f, scope); err != nil {
			return err
		}
	}

	return nil
}

// valueString returns the string representation of a value.
func valueString(v reflect.Value, opts tagOptions) string {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Bool && opts.Contains("int") {
		if v.Bool() {
			return "1"
		}
		return "0"
	}

	if v.Type() == timeType {
		t := v.Interface().(time.Time)
		if opts.Contains("unix") {
			return strconv.FormatInt(t.Unix(), 10)
		}
		return t.Format(time.RFC3339)
	}

	return fmt.Sprint(v.Interface())
}

// isEmptyValue checks if a value should be considered empty for the purposes
// of omitting fields with the "omitempty" option.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	if v.Type() == timeType {
		return v.Interface().(time.Time).IsZero()
	}

	return false
}

// tagOptions is the string following a comma in a struct field's "url" tag, or
// the empty string. It does not include the leading comma.
type tagOptions []string

// parseTag splits a struct field's url tag into its name and comma-separated
// options.
func parseTag(tag string) (string, tagOptions) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}

// Contains checks whether the tagOptions contains the specified option.
func (o tagOptions) Contains(option string) bool {
	for _, s := range o {
		if s == option {
			return true
		}
	}
	return false
}

const middle = "========="

type Config struct {
	Mymap  map[string]string
	strcet string
}

func (c *Config) InitConfig(path string) {
	c.Mymap = make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		s := strings.TrimSpace(string(b))
		//fmt.Println(s)
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			c.strcet = strings.TrimSpace(s[n1+1 : n2])
			continue
		}

		if len(c.strcet) == 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}

		key := c.strcet + middle + frist
		c.Mymap[key] = strings.TrimSpace(second)
	}
}

func (c Config) Read(node, key string) string {
	key = node + middle + key
	v, found := c.Mymap[key]
	if !found {
		return ""
	}
	return v
}
