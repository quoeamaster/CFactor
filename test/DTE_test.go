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

package test

import (
	"testing"
	"fmt"

	"github.com/BurntSushi/toml"
	"bytes"
)

// DTE => Down.To.Earth testing

func TestDummyTest(t *testing.T) {
	fmt.Println("\n*** [TestDummyTest]")
	if true {
		t.Logf("[log] basic testing on TOML access [%v]", "FAILED")
		//t.Errorf("basic testing on TOML access [%v]", "FAILED")
	}
	fmt.Println("** test done **")
}

/*
 *	simple struct to declare an address
 */
type AddressBook struct {
	Name string
	City string
	Street string
	Country string
}

func NewAddressBook(name, city, street, country string) AddressBook {
	addr := AddressBook{}
	addr.Name = name
	addr.City = city
	addr.Street = street
	addr.Country = country

	return addr
}

func (a *AddressBook) String() string {
	if a != nil {
		return fmt.Sprintf("[%v] addr => %v, %v, %v", a.Name, a.Street, a.City, a.Country)
	}
	return ""
}

// ` to start a multi-line string
// * a sample TOML
const ADDRESS_BOOK_SAMPLE_TOML_1  = `
	Name = "Simon Beistar"
	City = "Amsterdam"
	Street = "Stadhouderskade 7"
	Country = "Netherlands"
`

/*
 *	Decode a string into the given AddressBook struct; once decoded,
 *	access / read the fields of the struct is like AddressBookInstance.{field}
 */
func TestTOMLPackageRead(t *testing.T) {
	fmt.Println("\n*** [TestTOMLPackageRead]")
	addressBookInstance := AddressBook{}
	metaData, e := toml.Decode(ADDRESS_BOOK_SAMPLE_TOML_1, &addressBookInstance)
	if e != nil {
		panic(e)
	}
	// the actual contents are set to the instance (pointer to be provided)
	fmt.Println(addressBookInstance.String())
	// meta data after "reflection"
	fmt.Println(metaData.Keys())

	fmt.Println("\ntry again using io.Reader => ")
	convertStringToReader()
}

func convertStringToReader() {
	addressBookInstance := AddressBook{}
	bBytes := []byte(ADDRESS_BOOK_SAMPLE_TOML_1)
	// using bytes package create a new Reader
	bReader := bytes.NewReader(bBytes)

	_, e := toml.DecodeReader(bReader, &addressBookInstance)
	if e!=nil {
		panic(e)
	}
	fmt.Println(addressBookInstance.String())
}

/* ----------- encode ------------- */


