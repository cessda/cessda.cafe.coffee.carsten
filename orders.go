// orders.go

package main

import (
  "time"
  "github.com/satori/go.uuid"
)

type order struct {
  ID        string `json:"id"`
  OrderedAt string `json:"OrderedAt"`
  Type      string `json:"type"`
  ReadyAt   string `json:"readyAt"`
  Retrieved bool   `json:"retrieved"`
}

// Initial coffe orders
var orderList = []order{
  order{ID: "c1be03bf-d9cc-486b-92af-3d91c27d3ba5", Type: "COFFEE", OrderedAt: "2019-02-16T11:31:47+0000", ReadyAt: "2019-02-16T11:32:34+0000"},
  order{ID: "90fcb5bd-0f08-4656-8306-4e1efaaea2b0", Type: "CAPPUCCINO", OrderedAt: "2019-02-16T10:31:47+0000", ReadyAt: "2019-02-16T10:32:34+0000", Retrieved: true},
}

// return entire order history
func getAllOrders() []order {
  return orderList
}

// check whether a coffee is still brewing
func systemBrewing() bool {
  return false
}

// check whether a coffee needs picking up
func orderWaiting() bool {
  for _, o := range orderList {
    if ! o.Retrieved {
      return true
    }
  }
  return false
}

func retrieveOrder(id string) (*order, bool) {
  for index, o := range orderList {
    if o.ID == id {
      if ! o.Retrieved {
        orderList[index].Retrieved = true
        o.Retrieved = true
        return &o, true
      } else {
        return &o, false
      }
    }
  }
  return nil, false
}

func getOrderbyID(id string) (*order, bool) {
  for _, o := range orderList {
    if o.ID == id {
      return &o, true
    }
  }
  return nil, false
}

func newOrder(Type string) (*order, bool) {

  myorderid, _ := uuid.NewV4()

  var newOrder order
  newOrder.ID = myorderid.String()
  newOrder.Type = Type
  newOrder.OrderedAt = time.Now().Format(time.RFC3339)
  newOrder.ReadyAt = time.Now().Add(time.Minute*1).Format(time.RFC3339)

  orderList = append(orderList, newOrder)

  theNewOrder, success := getOrderbyID(newOrder.ID)
  return theNewOrder, success

}


