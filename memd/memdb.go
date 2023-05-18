package main

import (
	"database/sql"
	"fmt"
	"github.com/hashicorp/go-memdb"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

type Employee struct {
	ID         int
	Name       string
	Department string
	Age        int
}

var sqldb *sql.DB

const DSN = "host=localhost user=postgres password=Simform@123 dbname=Memdb port=5432 sslmode=disable"

func main() {
	var err error
	//CONNECTING TO DATABASE
	sqldb, err = sql.Open("postgres", DSN)
	if err != nil {
		log.Fatal("Error connecting to database: Reason: ", err)
	}
	fmt.Println("Successfully connected to database")
	//Defining the DB schema
	fmt.Println("Defining DB schema")
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"employee": &memdb.TableSchema{
				Name: "employee",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "ID"},
					},
					"name": &memdb.IndexSchema{
						Name:    "name",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Name"},
					},
					"department": &memdb.IndexSchema{
						Name:    "department",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Department"},
					},
					"age": &memdb.IndexSchema{
						Name:    "age",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Age"},
					},
				},
			},
		},
	}
	//Create a new database
	fmt.Println("Creating new database")
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
	//Creating  write transaction
	txn := db.Txn(true)

	//Insert some employees
	employee := []*Employee{
		&Employee{ID: 1, Name: "Heema", Department: "Golang", Age: 20},
		&Employee{ID: 2, Name: "Khushi", Department: "AX", Age: 30},
		&Employee{ID: 3, Name: "Riya", Department: "Dotnet", Age: 38},
		&Employee{ID: 2, Name: "Dhatri", Department: "Python", Age: 35},
	}

	//insert into db
	fmt.Println("Insert into db")
	for _, val := range employee {
		if err := txn.Insert("employee", val); err != nil {
			log.Fatal("Error creating records", err)
		}
	}

	fmt.Println("Successfully created records")
	//Committing the transaction
	fmt.Println("Committing the transaction")
	txn.Commit()

	//Creating read only transaction
	txn = db.Txn(false)
	defer txn.Abort()

	//Lookup by department using First
	//Here inline condition is used similar as gorm
	raw, err := txn.First("employee", "department", "Dotnet")
	if err != nil {
		log.Fatal("Error finding records. Reason: ", err)
	}
	fmt.Println(strings.Repeat("-", 20))
	fmt.Printf("Employee name of %s department is: %s Age: %d\n", raw.(*Employee).Department, raw.(*Employee).Name, raw.(*Employee).Age)

	//Lookup by id using first-->
	//It does not throw error like gorm if the record is not found. It will just give empty interface error if we try to access columns using the return tx instance
	//This first method is similar to first method in gorm.
	first, err := txn.First("employee", "id")
	if err != nil {
		log.Fatal("Error finding records. Reason: ", err)
	}
	fmt.Println(strings.Repeat("-", 20))
	fmt.Printf("Employee details of %d id is: \n", first.(*Employee).ID)
	fmt.Printf("Name: %s, Department: %s, Age: %d\n", first.(*Employee).Name, first.(*Employee).Department, first.(*Employee).Age)

	//Get all records use get method which is similar to find method in gorm.
	//Here it will not give duplicate record error like gorm even if we try to give same id twice and the unique constraint is given true

	get, err := txn.Get("employee", "id")
	if err != nil {
		log.Fatal("Error finding records. Reason: ", err)
	}
	for val := get.Next(); val != nil; val = get.Next() {
		obj := val.(*Employee)
		fmt.Println(strings.Repeat("-", 20))
		fmt.Printf("Employee details of %d id is: \n", obj.ID)
		fmt.Printf("Name: %s, Department: %s Age:%d\n", obj.Name, obj.Department, obj.Age)
	}

	// Range scan over people with ages between 25 and 35 inclusive
	fmt.Println(strings.Repeat("-", 20))
	fmt.Println("People aged after 25:")
	scan, err := txn.LowerBound("employee", "age", 25)
	for val := scan.Next(); val != nil; val = scan.Next() {
		p := val.(*Employee)
		fmt.Println(strings.Repeat("-", 20))
		fmt.Printf("Employee details of %d id is: \n", p.ID)
		fmt.Printf("Name: %s, Department: %s Age:%d\n", p.Name, p.Department, p.Age)
	}
}
