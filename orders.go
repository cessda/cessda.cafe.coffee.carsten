// orders.go

package main

import (
	"errors"
)

type order struct {
  ID        string `json:"id"`
  OrderedAt string `json:"OrderedAt"`
  Type      string `json:"type"`
  ReadyAt   string `json:"readyAt"`
  Retrieved bool   `json:"retrieved"`
}

var orderList = []order{
  order{ID: "c1be03bf-d9cc-486b-92af-3d91c27d3ba5", Type: "COFFEE", OrderedAt: "2019-02-16T11:31:47+0000", ReadyAt: "2019-02-16T11:32:34+0000"},
  order{ID: "90fcb5bd-0f08-4656-8306-4e1efaaea2b0", Type: "CAPPUCCINO", OrderedAt: "2019-02-16T10:31:47+0000", ReadyAt: "2019-02-16T10:32:34+0000", Retrieved: true},
}

func getAllOrders() []order {
  return orderList
}

func getOrderbyID(id string) (*order, error) {
  for _, o := range orderList {
    if o.ID == id {
      return &o, nil
    }
  }
  return nil, errors.New("Order not found")
}



