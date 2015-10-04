
package main

import (
"math"
"bytes"
  "net/http"
  "net/rpc"
  "encoding/json"
"strings"
 "fmt"
 "io/ioutil"
"strconv"
	
  )

type Call struct{
List flist `json:"list"`
}

type flist  struct {

Meta fmeta `json:"-"`
Resources []fresources `json:"resources"`
}

type fmeta struct {
Type string `json:"-"` 
Start int32 `json:"-"`
Count  int32 `json:"-"`
}

type fresources struct {
Resource fresource `json:"resource"`

}

type fresource struct {
Classname string `json:"classname"`
Fields ffields `json:"fields"`
}

type ffields struct{
Price string `json:"price"`
Symbol string `json:"symbol"`

}

type Args struct{
StockSymbolAndPercentage string
UserBudget float64
}

type Quote struct {
Stocks string `json:"stocksymbol"`
UnvestedAmount float64	`json:"stockprice"`
TradeId int `json:"id"`
}

type Id struct{
TradeId int `json:"id"`
}

type UpdQuote struct {
Stocks string `json:"stocksymbol"`
UnvestedAmount float64	`json:"stockprice"`
}


type StockCalc int

var M map[int]Quote


   func (t *StockCalc) StockPrice(args *Args, quote *Quote) error {
   
     

    a := string(args.StockSymbolAndPercentage[:])

	a = strings.Replace(a,":",",",-1)
	a = strings.Replace(a,"%",",",-1)
	a = strings.Replace(a,",,",",",-1)
	a = strings.Trim(a," ")
	a = strings.Replace(a,"\"","",-1)
	a = strings.TrimSpace(a)
	a = strings.TrimSuffix(a,",")
	Stockstmp:= strings.Split(a,",")
	

	var ReqUrl string
	


	for  i :=0; i < len(Stockstmp) ; i++ {
	i = i+1

	temp,_ := strconv.ParseFloat(Stockstmp[i],64)
	temp= (temp * args.UserBudget * 0.01) 

	ReqUrl= ReqUrl + (Stockstmp[i-1] + ",")
	
	
		}
		ReqUrl = strings.TrimSuffix(ReqUrl,",")
		
  UrlStr := "http://finance.yahoo.com/webservice/v1/symbols/" + ReqUrl +  "/quote?format=json"
 
 
 client := &http.Client{}
 
resp, _  := client.Get(UrlStr)
req, _ := http.NewRequest("GET", UrlStr, nil)

req.Header.Add("If-None-Match", "application/json")
req.Header.Add("Content-Type", "application/x-www-form-urlencoded")


// make request
resp, _  = client.Do(req)
if( resp.StatusCode >= 200 && resp.StatusCode < 300 ) {
  var C Call
  body, _ := ioutil.ReadAll(resp.Body) 
  
 err := json.Unmarshal(body, &C)
 
 n:= len(Stockstmp)
 
 
Quo:= make ( []float64,n,n)

 
 for  i :=0; i < n ; i++ {
	i = i + 1
	TempFloat,_ := strconv.ParseFloat(Stockstmp[i],64)
	Quo[i] = (TempFloat * args.UserBudget * 0.01)

	}
	
	
	var buffer bytes.Buffer
	q:=0
  for _ ,Sample := range  C.List.Resources {
  

temp1:= Sample.Resource.Fields.Symbol
temp2,_:=strconv.ParseFloat(Sample.Resource.Fields.Price,64) 
temp3:= (int)(Quo[q+1]/temp2)
temp4:= math.Mod(Quo[q+1],temp2)
q = q + 2

quote.Stocks = fmt.Sprintf("%s:%g:%d",temp1,temp2,temp3)
quote.UnvestedAmount = quote.UnvestedAmount + temp4
buffer.WriteString(quote.Stocks)
buffer.WriteString(",")
}

M = map[int]Quote{
quote.TradeId : {quote.Stocks,quote.UnvestedAmount,quote.TradeId},
}

quote.TradeId = len(M) 
quote.Stocks = (buffer.String())
quote.Stocks = strings.TrimSuffix(quote.Stocks, ",")





 if( err == nil ) {
      fmt.Println("Completed")
    }
  } else {
    fmt.Println(resp.Status);
  
  }
return nil
 }

  
  func (t *StockCalc) UpdStockPrice(id *Id, updquote *UpdQuote) error {
  
    var Tmp2,ttmp string
    var Sampletmp1 []string
	

		
	var  tmp = M[id.TradeId]
	ttmp = string(tmp.Stocks[:])
	ttmp = strings.Replace(ttmp,",", ":",-1)
	ttmp = strings.Trim(ttmp," ")
	ttmp = strings.TrimSpace(ttmp)

	Sampletmp1 = strings.Split(ttmp,":")

	
	for i:=0 ; i < len(Sampletmp1) ; i ++ {
	Tmp2 = Tmp2 +","+ Sampletmp1[i]
	i = i + 2
	
	}
	

	 Tmp2 = strings.TrimLeft(Tmp2,",")

        

 
  UrlStr := "http://finance.yahoo.com/webservice/v1/symbols/" + Tmp2 +  "/quote?format=json"
   
   
   client := &http.Client{}
   
  resp, _  := client.Get(UrlStr)
  req, _ := http.NewRequest("GET", UrlStr, nil)
  
  req.Header.Add("If-None-Match", "application/json")
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  
  
  // make request
  resp, _  = client.Do(req)
  if( resp.StatusCode >= 200 && resp.StatusCode < 300 ) {
    var C Call

    body, _ := ioutil.ReadAll(resp.Body) 
    
   _ = json.Unmarshal(body, &C)
   


 var buf bytes.Buffer

 k:= 1
  for _ ,Sample := range  C.List.Resources {

  temp1 := Sample.Resource.Fields.Symbol
temp2,_ := strconv.ParseFloat(Sample.Resource.Fields.Price,64) 
temp3,_:= strconv.ParseFloat(Sampletmp1[k],64)

if (temp3 > temp2) {
updquote.Stocks = fmt.Sprintf("%s:%s%v:%v",temp1,"-",temp2,Sampletmp1[k +1 ])
buf.WriteString(updquote.Stocks)
buf.WriteString(",")

} else if (temp3 < temp2){

updquote.Stocks = fmt.Sprintf("%s:%s%v:%v",temp1,"+",temp2,Sampletmp1[k +1 ])
buf.WriteString(updquote.Stocks)
buf.WriteString(",")

}else {

updquote.Stocks = fmt.Sprintf("%s:%v:%v",temp1,temp2,Sampletmp1[k +1 ])
buf.WriteString(updquote.Stocks)
buf.WriteString(",")
}

updquote.Stocks = (buf.String())
updquote.Stocks = strings.TrimSuffix(updquote.Stocks, ",")
updquote.UnvestedAmount = tmp.UnvestedAmount
k = k + 3
 }
      }
        return nil
 }

  func main() {
stockcalc:= new(StockCalc)
rpc.Register(stockcalc)
rpc.HandleHTTP()

err := http.ListenAndServe(":2081", nil)
if err != nil {
fmt.Println(err.Error())
}
}


  

