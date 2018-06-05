/*
 *  Copyright Project - CFactor, Author - quoeamaster, (C) 2018
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

// testing Struct for having hierarchical Structures.
package TOML

import (
	"time"
	"bytes"
	"fmt"
	"CFactor/common"
	"strconv"
	"reflect"
)

/*
 *	a struct to encapsulate a "transaction"
 */

// Struct wrapping up a "transaction"
type TransactionRecord struct {
	// struct to describe the "client" involved
	Client Client `toml:"client" additional:"parent"`

	Amount float32 `toml:"amount"`

	// struct to describe the "broker" involved
	Broker Broker `toml:"broker" additional:"parent"`
}

/*
 *	a struct to describe a "broker"
 */

// Struct wrapping up a "Broker"
type Broker struct {
	FullName string `toml:"broker.fullname"`
	Id string `toml:"broker.id"`
	Licences []string `toml:"broker.licences"`

	LicenceExpiryDate time.Time `toml:"broker.licenceExpiryDate"`
}

/*
 *	a struct to describe a "client / user"
 */

// Struct wrapping up a "client"
type Client struct {
	FullName string `toml:"client.fullname"`
	Id string `toml:"client.id"`

	// struct to describe the client address
	Address ClientAddress `toml:"client.address" additional:"parent"`
// TODO: testing ptr of struct instead of Address value
	//AddressPtr *ClientAddress `toml:"client.addressPtr" additional:"parent"`
}

/*
 *	a struct to describe an "address" for the client
 */

// Struct wrapping up a client's "address"
type ClientAddress struct {
	StreetNum int `toml:"client.address.streetnum"`
	StreetName string `toml:"client.address.streetname"`
	City string `toml:"client.address.city"`
	Country string `toml:"client.address.country"`

	// struct to describe (Lat, Lon) pair
	GeoPoint GeoPoint `toml:"client.address.geopoint" additional:"parent"`
}

/*
 *	a struct to describe a "geopoint" for the client address
 */

// Struct wrapping up a "geo-point" for a location on the address
type GeoPoint struct {
	Lat float64 `toml:"client.address.geopoint.Lat"`
	Lon float64 `toml:"client.address.geopoint.Lon"`

	LatLonArr []float64 `toml:"client.address.geopoint.LatLonArr"`
}


/* ------------------------ */
/*	String() declaration	*/
/* ------------------------ */

// return a string representation of a TransactionRecord
func (o *TransactionRecord) String() string {
	var bBuffer bytes.Buffer

	bBuffer.WriteString("Client => {")
	bBuffer.WriteString(o.Client.String())
	bBuffer.WriteString("}, \n\nBroker => {")
	bBuffer.WriteString(o.Broker.String())
	bBuffer.WriteString("}, \n\nAmount => {")
	bBuffer.WriteString(fmt.Sprintf("%v", o.Amount))
	bBuffer.WriteString("}\n")

	return bBuffer.String()
}
// return a string representation of a Client
func (o *Client) String() string {
	var bBuffer bytes.Buffer

	bBuffer.WriteString("\tfullname = ")
	bBuffer.WriteString(o.FullName)
	bBuffer.WriteString(", id = ")
	bBuffer.WriteString(o.Id)
	bBuffer.WriteString(", address = ")
	bBuffer.WriteString(o.Address.String())
	bBuffer.WriteString("$\n")

	return bBuffer.String()
}
// return a string representation of a Broker
func (o *Broker) String() string {
	var bBuffer bytes.Buffer

	bBuffer.WriteString("\n\tfullname = ")
	bBuffer.WriteString(o.FullName)
	bBuffer.WriteString(", id = ")
	bBuffer.WriteString(o.Id)
	bBuffer.WriteString(", licenceExpiryDate = ")
	bBuffer.WriteString(common.FormatTimeToString("", o.LicenceExpiryDate))
	bBuffer.WriteString(", licences = ")
	bBuffer.WriteString(fmt.Sprintf("%v", o.Licences))
	bBuffer.WriteString(">\n")

	return bBuffer.String()
}
// return a string representation of a ClientAddress
func (o *ClientAddress) String() string {
	var bBuffer bytes.Buffer

	bBuffer.WriteString("\n\t\tStreetNum = ")
	bBuffer.WriteString(strconv.Itoa(o.StreetNum))
	bBuffer.WriteString(", StreetName = ")
	bBuffer.WriteString(o.StreetName)
	bBuffer.WriteString(", city = ")
	bBuffer.WriteString(o.City)
	bBuffer.WriteString(", country = ")
	bBuffer.WriteString(o.Country)
	bBuffer.WriteString(", geopoint = ")
	bBuffer.WriteString(o.GeoPoint.String())
	bBuffer.WriteString("}")

	return bBuffer.String()
}
// return a string representation of a GeoPoint
func (o *GeoPoint) String() string {
	var bBuffer bytes.Buffer

	bBuffer.WriteString("\n\t\t\tLat = ")
	bBuffer.WriteString(fmt.Sprintf("%v", o.Lat))
	bBuffer.WriteString(", Lon = ")
	bBuffer.WriteString(fmt.Sprintf("%v", o.Lon))
	bBuffer.WriteString(", LatLonArr = ")
	bBuffer.WriteString(fmt.Sprintf("%v", o.LatLonArr))
	//bBuffer.WriteString("\n")

	return bBuffer.String()
}


/* -------------------- */
/*	lifecycle hooks     */
/* -------------------- */

// the lifeCycle Hook method implementation (check IConfig.go)
func (o *TransactionRecord) SetStructsReferences(structRefMap *map[string]interface{}) (err error) {
	structRefMapVal := *structRefMap
	if len(structRefMapVal)==0 {
		return nil
	}
	for key, structRef := range structRefMapVal {
		switch key {
		case "TOML.Client":
			o.Client = reflect.Indirect(reflect.ValueOf(structRef)).Interface().(Client)
		case "TOML.Broker":
			o.Broker = reflect.Indirect(reflect.ValueOf(structRef)).Interface().(Broker)
		case "TOML.ClientAddress":
			o.Client.Address = reflect.Indirect(reflect.ValueOf(structRef)).Interface().(ClientAddress)
		case "TOML.GeoPoint":
			o.Client.Address.GeoPoint = reflect.Indirect(reflect.ValueOf(structRef)).Interface().(GeoPoint)
		default:
			return fmt.Errorf("unknown struct type! [%v]", key)
		}
	}	// end -- for (structRef)

	// recovery if necessary
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	return nil
}
