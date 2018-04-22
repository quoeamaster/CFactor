package test

import (
	"testing"
	"reflect"
	"fmt"
	"strings"
)

/*
 *	struct to describe a user profile
 */
type UserProfile struct {
	Name Name
	Age int `age=int; explanation=exp@34`
	AddressGeneral string
}
/*
 *	struct to describe a Name
 */
type Name struct {
	FirstName string
	LastName string
}

/* -------- */
/*	 ctor	*/
/* -------- */

func NewUserProfile(firstName, lastName, addressGeneral string, age int) UserProfile {
	n := Name{
		FirstName: firstName,
		LastName: lastName,
	}
	up := UserProfile{
		Name: n,
		AddressGeneral: addressGeneral,
		Age: age,
	}
	return up
}

/* -------- */
/*	helper	*/
/* -------- */




/* ---------------- */
/*		TEST		*/
/* ---------------- */

func TestReflectionTypeAccess(t *testing.T) {
	up := NewUserProfile("Alex", "Beanie", "Germany", 35)

	rType := reflect.TypeOf(up)
	//rElemType := rType.Elem()

	// basic acccess
	fmt.Println("* name of type =>", rType.Name());
	fmt.Println("* number of fields =>", rType.NumField())

	for i:=0; i<rType.NumField(); i++ {
		rField := rType.Field(i)
		fmt.Printf("\t# field: [name - %v] [type = %v] [tag ^ %v] \n", rField.Name, rField.Type, rField.Tag)
	}

	// reflect.New is a Ptr... must call Ptr.Elem() to get back the object (non pointer)
	ptrReflectedTypeInstance := reflect.New(rType)
	objReflectedTypeInstance := ptrReflectedTypeInstance.Elem()
	objReflectedTypeInstance.Field(1).SetInt(19)
	objReflectedTypeInstance.Field(2).SetString("Belgium instead of Germany")
	fmt.Println(ptrReflectedTypeInstance)
	fmt.Println(objReflectedTypeInstance)

	objReflectedTypeInstance.FieldByName("AddressGeneral").SetString("whatever address available...")

	field := objReflectedTypeInstance.FieldByName("addressGeneral")
	fmt.Println("have field 'addressGeneral' instead?", field)
	// ** TIP => the "field" String() method contains "invalid reflect.Value" means it is not a VALID field
	fmt.Println(strings.Contains(field.String(), "invalid"))
	fmt.Println(objReflectedTypeInstance)

	// adding another Object / struct level
	rTypeForName := reflect.TypeOf(Name{"", ""})
	ptrNameStruct := reflect.New(rTypeForName)
	objNameStruct := ptrNameStruct.Elem()
	// set values to fields
	objNameStruct.Field(0).SetString("Jonathan")
	objNameStruct.Field(1).SetString("Joestar")
	fmt.Println("\n", objNameStruct)

	// set Name back to UserProfile
	objReflectedTypeInstance.FieldByName("Name").Set(objNameStruct)
	fmt.Println("\nUserProfile completed =>", objReflectedTypeInstance)
}
