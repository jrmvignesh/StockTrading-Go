package main

import (
"bufio"
"os"
"log"
 "fmt"
"net/rpc"
  )
  
  
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

  // Create Client
  func main() {
  
  if len(os.Args) != 2 {
fmt.Println("Usage: ", os.Args[0], "server")
os.Exit(1)
}
serverAddress := os.Args[1]

client, err := rpc.DialHTTP("tcp", serverAddress+":2081")
if err != nil {
log.Fatal("dialing:", err)
}

fmt.Print("Enter 1 for Buying Stocks,2 for Portfolio check, 3 for exit: ")
	var Userchoice int
	fmt.Scan(&Userchoice)
switch Userchoice {
case 1 : {

reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Stock Symbol: ")
	StockSymbolAndPercentage, _ := reader.ReadString('\n')
	fmt.Println(StockSymbolAndPercentage)
	fmt.Print("Enter budget: ")
	var UserBudget float64
	fmt.Scan(&UserBudget)
	
	args:= Args{StockSymbolAndPercentage,UserBudget}

var quo Quote 

err = client.Call("StockCalc.StockPrice",args,&quo)

if err != nil {
log.Fatal("arith error:", err)
}

fmt.Println("Stock Response: ", quo.Stocks)
fmt.Println("Unvested Amount: " , quo.UnvestedAmount)
fmt.Println("Trade Id: " , quo.TradeId)

}

case 2 :  {

fmt.Print("Enter TradeID to check the portfolio: ")
	var TradeID int
	fmt.Scan(&TradeID)

id:= Id{TradeID}

var  uquote   UpdQuote

err = client.Call("StockCalc.UpdStockPrice",id,&uquote)

if err != nil {
log.Fatal("Update error:", err)
}

fmt.Println("Updated Stock Response: ",uquote.Stocks)
fmt.Println("Unvested Amount: ", uquote.UnvestedAmount)

}

case 3 : break

}


}



