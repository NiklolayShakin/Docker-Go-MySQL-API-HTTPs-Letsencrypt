package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/acme/autocert"
)

type coordinates struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type addr struct {
	Address     string      `json:"address"`
	Coordinates coordinates `json:"coordinates"`
	Comment     string      `json:"comment"`
}

type contactInfo struct {
	Name      []string `json:"name"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	Addresses []addr   `json:"addresses"`
}

//DB scheme
// id
// first_name
// full_name
// email
// phone_number
// address_city
// address_street
// address_house
// address_entrance
// address_floor
// address_office
// address_comment
// location_latitude
// location_longitude
// amount_charged
// user_id
// user_agent
// created_at
// address_doorcode

type dbPerson struct {
	id                 sql.NullString
	first_name         sql.NullString
	full_name          sql.NullString
	email              sql.NullString
	phone_number       sql.NullString
	address_city       sql.NullString
	address_street     sql.NullString
	address_house      sql.NullString
	address_entrance   sql.NullString
	address_floor      sql.NullString
	address_office     sql.NullString
	address_comment    sql.NullString
	location_latitude  sql.NullString
	location_longitude sql.NullString
	amount_charged     sql.NullString
	user_id            sql.NullString
	user_agent         sql.NullString
	created_at         sql.NullString
	address_doorcode   sql.NullString
}

// implimenting stringer inteface for contactInfo stuct
func (c contactInfo) String() string {
	result := fmt.Sprintf("Name: %v\nEmail: %v\nPhone: %v\nAddress:", c.Name, c.Email, c.Phone)

	for i, a := range c.Addresses {
		result = result + fmt.Sprintf("\n\taddress_%d: \n\t\t%v\n\t\tCoordinates: %v\n\t\tComment: %v", i+1, a.Address, a.Coordinates, a.Comment)
	}

	return result
}

// geting data from DB
func getPerson(num int) []*dbPerson {
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", dbUser, dbPassword, dbHost, dbPort, dbName))

	if err != nil {
		log.Print(err.Error())
		return []*dbPerson{}
	}
	defer db.Close()

	results, err := db.Query(fmt.Sprintf("SELECT * FROM person WHERE phone_number='+%v';", num))
	if err != nil {
		log.Print(err.Error())
	}

	var persons []*dbPerson

	for results.Next() {
		var p dbPerson
		err = results.Scan(
			&p.id,
			&p.first_name,
			&p.full_name,
			&p.email,
			&p.phone_number,
			&p.address_city,
			&p.address_street,
			&p.address_house,
			&p.address_entrance,
			&p.address_floor,
			&p.address_office,
			&p.address_comment,
			&p.location_latitude,
			&p.location_longitude,
			&p.amount_charged,
			&p.user_id,
			&p.user_agent,
			&p.created_at,
			&p.address_doorcode)

		if err != nil {
			log.Print(err.Error())
		}

		persons = append(persons, &p)

	}

	return persons
}

//fill the contact out of DB response
func populateContact(dbresponse []*dbPerson) (contact contactInfo) {

	if len(dbresponse) == 0 {
		log.Print("------------------------")
		return
	}

	for i, person := range dbresponse {

		//first cyclce , populate initial data

		if i == 0 {
			//check that first_name is not empty before adding to the slice
			if len(person.first_name.String) != 0 {
				contact.Name = []string{person.first_name.String, person.full_name.String}
			} else {
				contact.Name = []string{person.full_name.String}
			}
			contact.Phone = person.phone_number.String
			contact.Email = person.email.String

			//if addr exist
			if person.address_city.Valid == true {
				//chek and populate if coordinates exist
				var (
					coord   coordinates
					address string
					comment string
				)
				if person.location_latitude.Valid && person.location_longitude.Valid {
					coord = coordinates{person.location_latitude.String, person.location_longitude.String}
				}
				//addr
				address = person.address_city.String + ", " + person.address_street.String + ", д." + person.address_house.String + ", кв." + person.address_office.String
				comment = person.address_comment.String

				if len(person.address_entrance.String) > 0 {
					comment = comment + " подъезд:" + person.address_entrance.String
				}

				if len(person.address_doorcode.String) > 0 {
					comment = comment + " код:" + person.address_doorcode.String
				}

				address2add := addr{address, coord, comment}

				contact.Addresses = append(contact.Addresses, address2add)
			}

		}
		// if we have more than one result, check are there a new names
		// if we found new name, adding them to the addresses slice
		if i > 0 {
			var (
				newFirstName = true
				newFullName  = true
			)
			//add new and not emtpy names to slice
			for _, cn := range contact.Name {
				if cn == person.first_name.String || len(person.first_name.String) == 0 {
					newFirstName = false
					break
				}
			}

			if newFirstName {
				contact.Name = append(contact.Name, person.first_name.String)
			}

			for _, cn := range contact.Name {
				if cn == person.full_name.String {
					newFullName = false
					break
				}
			}

			if newFullName {
				contact.Name = append(contact.Name, person.full_name.String)
			}

		}

		// if we have more than one result from DB , check addrresses
		// any new address we will add it to the addr slice
		if i > 0 {
			var (
				newAddress = true
			)
			currentPersonAddr := person.address_city.String + ", " + person.address_street.String + ", д." + person.address_house.String + ", кв." + person.address_office.String

			for _, a := range contact.Addresses {
				if a.Address == currentPersonAddr {
					newAddress = false
					break
				}
			}

			if newAddress && len(person.address_city.String) != 0 {

				var coord coordinates

				if person.location_latitude.Valid && person.location_longitude.Valid {
					coord = coordinates{person.location_latitude.String, person.location_longitude.String}
				}
				comment := person.address_comment.String

				if len(person.address_entrance.String) > 0 {
					comment = comment + " подъезд:" + person.address_entrance.String
				}

				if len(person.address_doorcode.String) > 0 {
					comment = comment + " код:" + person.address_doorcode.String
				}

				address2add := addr{currentPersonAddr, coord, comment}

				contact.Addresses = append(contact.Addresses, address2add)
			}

		}

	}

	return
}

func handler(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["phone"]

	if len(keys) < 1 {
		fmt.Fprintf(w, "%v\n", "{ERROR:\"You should use at least one allowed key \"}")
		return
	}

	// check number format remove non digits symbols
	phone := keys[0]
	if len(phone) > 20 || len(phone) < 11 {
		fmt.Fprintf(w, "%v\n", "{ERROR:\"Phone number has to be 11...20 characters including ( ) - + \"}")
		return
	}
	//remove non-digits
	re := regexp.MustCompile("[^0-9]*")
	phone = string(re.ReplaceAll([]byte(phone), []byte("")))
	re = regexp.MustCompile("^8")
	phone = string(re.ReplaceAll([]byte(phone), []byte("7")))

	if len(phone) != 11 {
		fmt.Fprintf(w, "%v\n", "{ERROR:\"Number is not 11-digits after removing non-digit symbols\"}")
		return
	}

	if ok {
		num, err := strconv.Atoi(phone)
		if err != nil {
			fmt.Println(err)
			return
		}

		persons := getPerson(num)

		if len(persons) == 0 {
			fmt.Fprintf(w, "%v\n", "{The number not found}")
			return
		}
		result := populateContact(persons)

		//maping struct to json
		contact, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Fprintf(w, "%v\n", string(contact))
		return
	}
	fmt.Fprintf(w, "%v\n", "{ERROR:\"Wrong request\"}")
	return
}

// // https redirect, use if  use need redirect http reuqests to https servers
// func redirectHTTP(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "GET" {
// 		http.Error(w, "Bad method", http.StatusBadRequest)
// 		return
// 	}

// 	host, _, err := net.SplitHostPort(r.Host)
// 	if err != nil {
// 		host = r.Host
// 	}
// 	host = net.JoinHostPort(host, "443")

// 	target := "https://" + host + r.URL.RequestURI()

// 	http.Redirect(w, r, target, http.StatusFound)
// }

func main() {

	http.HandleFunc(securePath, handler)

	if !useHttps {
		//starting http service
		fmt.Printf("Starting server at port %v\n", webPort)
		log.Fatal(http.ListenAndServe(":"+webPort, nil))
	} else {

		m := &autocert.Manager{
			Cache:      autocert.DirCache("secret-dir"),
			Prompt:     autocert.AcceptTOS,
			Email:      certEmail,
			HostPolicy: autocert.HostWhitelist(certDomain, certDomain),
		}
		s := &http.Server{
			Addr:      ":https",
			TLSConfig: m.TLSConfig(),
		}

		//starting https service
		log.Fatal(s.ListenAndServeTLS("", ""))

	}

}
