package TOML

import (
	"time"
	"bytes"
	"fmt"
	"CFactor/common"
	"strconv"
)

/*
 *	a struct to encapsulate a "transaction"
 */
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
type Broker struct {
	FullName string `toml:"broker.fullname"`
	Id string `toml:"broker.id"`
	Licences []string `toml:"broker.licences"`

	LicenceExpiryDate time.Time `toml:"broker.licenceExpiryDate"`
}

/*
 *	a struct to describe a "client / user"
 */
type Client struct {
	FullName string `toml:"client.fullname"`
	Id string `toml:"client.id"`

	// struct to describe the client address
	Address ClientAddress `toml:"client.address" additional:"parent"`
}

/*
 *	a struct to describe an "address" for the client
 */
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
type GeoPoint struct {
	Lat float32 `toml:"client.address.geopoint.Lat"`
	Lon float32 `toml:"client.address.geopoint.Lon"`

	LatLonArr []float64 `toml:"client.address.geopoint.LatLonArr"`
}


/* ------------------------ */
/*	String() declaration	*/
/* ------------------------ */

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
