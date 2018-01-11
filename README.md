# Introduction

this is powerful sql orm library for Golang, it has few dependencies and high performance

### A Simple Example

copy following code to a main.go file and run

```golang
package main

import (
	"encoding/json"
	"fmt"
	"github.com/bigpigeon/toyorm"
	. "unsafe"

	// when database is mysql
	_ "github.com/go-sql-driver/mysql"
	// when database is sqlite3
	//_ "github.com/mattn/go-sqlite3"
)

func JsonEncode(v interface{}) string {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(jsonData)
}

type Product struct {
	toyorm.ModelDefault
	Name  string  `toyorm:"index"`
	Price float64 `toyorm:"index"`
	Count int
	Tag   string `toyorm:"index"`
}

func main() {
	var err error
	var toy *toyorm.Toy
	// when database is mysql, make sure your mysql have toyorm_example schema
	toy, err = toyorm.Open("mysql", "root:@tcp(localhost:3306)/toyorm_example?charset=utf8&parseTime=True")
	// when database is sqlite3
	//toy,err = toyorm.Open("sqlite3", "toyorm_test.db")

	brick := toy.Model(&Product{}).Debug()
	_, err = brick.DropTableIfExist()
	if err != nil {
		panic(err)
	}
	_, err = brick.CreateTableIfNotExist()
	if err != nil {
		panic(err)
	}
	// Insert will set id to source data when primary key is auto_increment
	brick.Insert(&[]Product{
		{Name: "food one", Price: 1, Count: 4, Tag: "food"},
		{Name: "food two", Price: 2, Count: 3, Tag: "food"},
		{Name: "food three", Price: 3, Count: 2, Tag: "food"},
		{Name: "food four", Price: 4, Count: 1, Tag: "food"},
		{Name: "toolkit one", Price: 1, Count: 8, Tag: "toolkit"},
		{Name: "toolkit two", Price: 2, Count: 6, Tag: "toolkit"},
		{Name: "toolkit one", Price: 3, Count: 4, Tag: "toolkit"},
		{Name: "toolkit two", Price: 4, Count: 2, Tag: "toolkit"},
	})
	// find the first food
	{
		var product Product
		_, err := brick.Where(toyorm.ExprEqual, Offsetof(Product{}.Tag), "food").Find(&product)
		if err != nil {
			panic(err)
		}
		fmt.Printf("food %s\n", JsonEncode(product))
	}

	// find all food
	{
		var products []Product
		_, err := brick.Where(toyorm.ExprEqual, Offsetof(Product{}.Tag), "food").Find(&products)
		if err != nil {
			panic(err)
		}
		fmt.Printf("foods %s\n", JsonEncode(products))
	}

	// find count = 2
	{
		var products []Product
		_, err := brick.Where(toyorm.ExprEqual, Offsetof(Product{}.Count), 2).Find(&products)
		if err != nil {
			panic(err)
		}
		fmt.Printf("count > 2 products %s\n", JsonEncode(products))
	}
	// find count = 2 and price > 3
	{
		var products []Product
		_, err := brick.Where(toyorm.ExprEqual, Offsetof(Product{}.Count), 2).And().
			Condition(toyorm.ExprGreater, Offsetof(Product{}.Price), 3).Find(&products)
		if err != nil {
			panic(err)
		}
		fmt.Printf("count = 2 and price > 3 products %s\n", JsonEncode(products))
	}
	// find count = 2 and price > 3 or count = 4
	{
		var products []Product
		_, err := brick.Where(toyorm.ExprEqual, Offsetof(Product{}.Count), 2).And().
			Condition(toyorm.ExprGreater, Offsetof(Product{}.Price), 3).Or().
			Condition(toyorm.ExprEqual, Offsetof(Product{}.Count), 4).Find(&products)
		if err != nil {
			panic(err)
		}
		fmt.Printf("count = 2 and price > 3 or count = 4  products %s\n", JsonEncode(products))
	}
	// find price > 3 and (count = 2 or count = 1)
	{
		var products []Product
		_, err := brick.Where(toyorm.ExprGreater, Offsetof(Product{}.Price), 3).And().Conditions(
			brick.Where(toyorm.ExprEqual, Offsetof(Product{}.Count), 2).Or().
				Condition(toyorm.ExprEqual, Offsetof(Product{}.Count), 1).Search,
		).Find(&products)
		if err != nil {
			panic(err)
		}
		fmt.Printf("price > 3 and (count = 2 or count = 1)  products %s\n", JsonEncode(products))
	}
	// update to count = 4
	{
		_, err := brick.Update(&Product{
			Count: 4,
		})
		if err != nil {
			panic(err)
		}
		var Counters []struct {
			Name  string
			Count int
		}
		_, err = brick.Find(&Counters)
		if err != nil {
			panic(err)
		}
		for _, counter := range Counters {
			fmt.Printf("product name %s, count %d\n", counter.Name, counter.Count)
		}
	}
	// delete with element
	{
		// make a transaction, because I do not really delete a data
		brick := brick.Begin()
		var product Product
		_, err := brick.Where(toyorm.ExprEqual, Offsetof(Product{}.Name), "food four").Find(&product)
		if err != nil {
			panic(err)
		}

		_, err = brick.Delete(&product)
		if err != nil {
			panic(err)
		}

		var disappearProduct Product
		_, err = brick.Where(toyorm.ExprEqual, Offsetof(Product{}.Name), "food four").Find(&disappearProduct)
		fmt.Printf("when I delete name = food four product , I try to find it will have error(%s)\n", err)
		brick.Rollback()
	}
	// delete with condition
	{
		// make a transaction, because I am not really delete a data
		brick := brick.Begin()
		_, err := brick.Where(toyorm.ExprEqual, Offsetof(Product{}.Name), "food four").DeleteWithConditions()
		if err != nil {
			panic(err)
		}
		var product Product
		_, err = brick.Where(toyorm.ExprEqual, Offsetof(Product{}.Name), "food four").Find(&product)
		fmt.Printf("when I delete name = food four product , I try to find it will have error(%s)\n", err)
		brick.Rollback()
	}
}
```

---

### Database connection

import database driver

```golang
// if database is mysql
_ "github.com/go-sql-driver/mysql"
// if database is sqlite3
_ "github.com/mattn/go-sqlite3"
```

create a toy

```golang
// if database is mysql, make sure your mysql have toyorm_example schema
toy, err = toyorm.Open("mysql", "root:@tcp(localhost:3306)/toyorm_example?charset=utf8&parseTime=True")
// if database is sqlite3
toy,err = toyorm.Open("sqlite3", "toyorm_test.db")
```


### Model definition

**example**

```golang
type Extra map[string]interface{}

func (e Extra) Scan(value interface{}) error {
     switch v := value.(type) {
     case string:
          return json.Unmarshal([]byte(v), e)
     case []byte:
          return json.Unmarshal(v, e)
     default:
          return errors.New("not support type")
     }
}

func (e Extra) Value() (driver.Value, error) {
     return json.Marshal(e)
}

type UserDetail struct {
     ID       int  `toyorm:"primary key;auto_increment"`
     UserID   uint `toyorm:"index"`
     MainPage string
     Extra    Extra `toyorm:"type:VARCHAR(1024)"`
}

type Blog struct {
     toyorm.ModelDefault
     UserID  uint   `toyorm:"index"`
     Title   string `toyorm:"index"`
     Content string
}

type User struct {
     toyorm.ModelDefault
     Name    string `toyorm:"unique index"`
     Age     int
     Sex     string
     Detail  *UserDetail
     Friends []*User
     Blog    []Blog
}
```

**special fields**

1. all special fields have some process in handlers, do not try to change it type or set it value


Field Name| Type      | Description
----------|-----------|------------
CreatedAt |time.Time  | generate when element be create
UpdatedAt |time.Time  | generate when element be update/create
DeletedAt |*time.Time | delete mode is soft

**field tags**

1. tag format can be \<key:value\> or \<key\>

2. the following is special tag

Key           |Value                   |Description
--------------|------------------------|-----------
index         | void or string         | use for optimization when search condition have this field,if you want make a combined,just set same index name with fields
unique index  | void or string         | have unique limit index, other same as index
primary key   | void                   | allow multiple primary key,but some operation not support
\-            | void                    | ignore this field in sql
type          | string                  | sql type
column        | string                  | sql column name
auto_increment| void                    | recommend, if your table primary key have auto_increment attribute must add it
autoincrement | void                    | same as auto_increment

other custom TAG will append to end of CREATE TABLE field



#### Bind models

1. model kind must be a struct or a point with struct

2. the model is what information that toyorm know about the table

```golang
brick := toy.Model(&User{})
// or
brick := toy.Model(User{})
```

### Sql Operation

---

#### create table

```golang
var err error
_, err = toy.Model(&User{}).Debug().CreateTable()
// CREATE TABLE user (id BIGINT AUTO_INCREMENT,created_at TIMESTAMP NULL,updated_at TIMESTAMP NULL,deleted_at TIMESTAMP NULL,name VARCHAR(255),age BIGINT ,sex VARCHAR(255) , PRIMARY KEY(id))
// CREATE INDEX idx_user_deletedat ON user(deleted_at)
// CREATE UNIQUE INDEX udx_user_name ON user(name)
_, err =toy.Model(&UserDetail{}).Debug().CreateTable()
// CREATE TABLE user_detail (id BIGINT AUTO_INCREMENT,user_id BIGINT,main_page Text,extra VARCHAR(1024), PRIMARY KEY(id))
// CREATE INDEX idx_user_detail_userid ON user_detail(user_id)
_, err =toy.Model(&Blog{}).Debug().CreateTable()
// CREATE TABLE blog (id BIGINT AUTO_INCREMENT,created_at TIMESTAMP NULL,updated_at TIMESTAMP NULL,deleted_at TIMESTAMP NULL,user_id BIGINT,title VARCHAR(255),content VARCHAR(255) , PRIMARY KEY(id))
// CREATE INDEX idx_blog_deletedat ON blog(deleted_at)
// CREATE INDEX idx_blog_userid ON blog(user_id)
// CREATE INDEX idx_blog_title ON blog(title)
```

#### drop table

```golang
var err error
_, err =toy.Model(&User{}).Debug().DropTable()
// DROP TABLE user
_, err =toy.Model(&UserDetail{}).Debug().DropTable()
// DROP TABLE user_detail
_, err =toy.Model(&Blog{}).Debug().DropTable()
// DROP TABLE blog
```

#### insert/save data

// insert with autoincrement will set id to source data

```golang
user := &User{
    Name: "bigpigeon",
    Age:  18,
    Sex:  "male",
}
_, err = toy.Model(&User{}).Debug().Insert(&user)
// INSERT INTO user(created_at,updated_at,name,age,sex) VALUES(?,?,?,?,?) , args:[]interface {}{time.Time{wall:0xbe8df5112e7f07c8, ext:210013499, loc:(*time.Location)(0x141af80)}, time.Time{wall:0xbe8df5112e7f1768, ext:210017044, loc:(*time.Location)(0x141af80)}, "bigpigeon", 18, "male"}
// print user format with json
/* {
  "ID": 1,
  "CreatedAt": "2018-01-11T20:47:00.780077+08:00",
  "UpdatedAt": "2018-01-11T20:47:00.780081+08:00",
  "DeletedAt": null,
  "Name": "bigpigeon",
  "Age": 18,
  "Sex": "male",
  "Detail": null,
  "Friends": null,
  "Blog": null
}*/
```

// save data use "REPLACE INTO" when primary key exist

```golang
user := &User{
    ModelDefault: toyorm.ModelDefault{ID: 1},
    Name: "bigpigeon",
    Age:  18,
    Sex:  "male",
}
_, err = toy.Model(&User{}).Debug().Save(&user)
// SELECT id,created_at FROM user WHERE id IN (?), args:[]interface {}{0x1}
// REPLACE INTO user(id,created_at,updated_at,name,age,sex) VALUES(?,?,?,?,?,?) , args:[]interface {}{0x1, time.Time{wall:0x0, ext:63651271470, loc:(*time.Location)(nil)}, time.Time{wall:0xbe8df4eb81a130c8, ext:232012178, loc:(*time.Location)(0x141af80)}, "bigpigeon", 18, "male"}

```

#### update

```golang
brick.Model(&User{}).Debug().Update(&User{
    Age: 4,
})
// UPDATE user SET updated_at=?,age=? WHERE deleted_at IS NULL, args:[]interface {}{time.Time{wall:0xbe8df4eb81b6c050, ext:233425327, loc:(*time.Location)(0x141af80)}, 4}
```


#### find

find one

```golang
var user User
_, err = toy.Model(&User{}).Debug().Find(&user}
// SELECT id,created_at,updated_at,deleted_at,name,age,sex FROM user WHERE deleted_at IS NULL LIMIT 1, args:[]interface {}(nil)
// print user format with json
/* {
  "ID": 1,
  "CreatedAt": "2018-01-11T12:47:01Z",
  "UpdatedAt": "2018-01-11T12:47:01Z",
  "DeletedAt": null,
  "Name": "bigpigeon",
  "Age": 4,
  "Sex": "male",
  "Detail": null,
  "Friends": null,
  "Blog": null
}*/
```

find multiple

```golang
var users []User
_, err = brick.Debug().Find(&users)
fmt.Printf("find users %s\n", JsonEncode(&users))

// SELECT id,created_at,updated_at,deleted_at,name,age,sex FROM user WHERE deleted_at IS NULL, args:[]interface {}(nil)
```

#### delete

delete with primary key

```golang
_, err = brick.Debug().Delete(&user)
// UPDATE user SET deleted_at=? WHERE id IN (?), args:[]interface {}{(*time.Time)(0xc4200f0520), 0x1}
```

delete with condition

```golang
_, err = brick.Debug().Where(toyorm.ExprEqual, Offsetof(User{}.Name), "bigpigeon").DeleteWithConditions()
// UPDATE user SET deleted_at=? WHERE name = ?, args:[]interface {}{(*time.Time)(0xc4200dbfa0), "bigpigeon"}
```

### ToyBrick

-----

use **toy.Model** will create a ToyBrick, you need use it to build grammar and operate the database

#### Where condition

affective update/find/delete operation

**usage**

where will clean old conditions and make new one

    brick.Where(<expr>, <Key>, [value])

conditions will copy conditions and clean old conditions

    brick.Conditions(<toyorm.Search>)

or & and condition will use or/and to link new condition when current condition is not nil

    brick.Or().Condition(<expr>, <Key>, [value])
    brick.And().Condition(<expr>, <Key>, [value])

or & and conditions will use or/and to link new conditions

    brick.Or().Conditions(<toyorm.Search>)
    brick.And().Conditions(<toyorm.Search>)

**SearchExpr**

SearchExpr        |  to sql      | example
------------------|--------------|:----------------
ExprAnd           | AND          | brick.Where(ExprAnd, Product{Name:"food one", Count: 4}) // WHERE name = "food one" AND Count = 4
ExprOr            | OR           | brick.Where(ExprOr, Product{Name:"food one", Count: 4}) // WHERE name = "food one" OR Count = "4"
ExprEqual         | =            | brick.Where(ExprEqual, OffsetOf(Product{}.Name), "food one") // WHERE name = "find one"
ExprNotEqual      | <>           | brick.Where(ExprNotEqual, OffsetOf(Product{}.Name), "food one") // WHERE name <> "find one"
ExprGreater       | >            | brick.Where(ExprGreater, OffsetOf(Product{}.Count), 3) // WHERE count > 3
ExprGreaterEqual  | >=           | brick.Where(ExprGreaterEqual, OffsetOf(Product{}.Count), 3) // WHERE count >= 3
ExprLess          | <            | brick.Where(ExprLess, OffsetOf(Product{}.Count), 3) // WHERE count < 3
ExprLessEqual     | <=           | brick.Where(ExprLessEqual, OffsetOf(Product{}.Count), 3) // WHERE count <= 3
ExprBetween       | Between      | brick.Where(ExprBetween, OffsetOf(Product{}.Count), [2]int{2,3}) // WHERE count BETWEEN 2 AND 3
ExprNotBetween    | NOT Between  | brick.Where(ExprNotBetween, OffsetOf(Product{}.Count), [2]int{2,3}) // WHERE count NOT BETWEEN 2 AND 3
ExprIn            | IN           | brick.Where(ExprIn, OffsetOf(Product{}.Count), []int{1, 2, 3}) // WHERE count IN (1,2,3)
ExprNotIn         | NOT IN       | brick.Where(ExprNotIn, OffsetOf(Product{}.Count), []int{1, 2, 3}) // WHERE count NOT IN (1,2,3)
ExprLike          | LIKE         | brick.Where(ExprLike, OffsetOf(Product{}.Name), "one") // WHERE name LIKE "one"
ExprNotLike       | NOT LIKE     | brick.Where(ExprNotLike, OffsetOf(Product{}.Name), "one") // WHERE name NOT LIKE "one"
ExprNull          | IS NULL      | brick.Where(ExprNull, OffsetOf(Product{}.DeletedAt)) // WHERE DeletedAt IS NULL
ExprNotNull       | IS NOT NULL  | brick.Where(ExprNotNull, OffsetOf(Product{}.DeletedAt)) // WHERE DeletedAt IS NOT NULL

**example**

single condition

```golang
brick = brick.Where(toyorm.ExprEqual, Offsetof(Product{}.Tag), "food")
// WHERE tag = "food"
```

combination condition

```golang
brick = brick.Where(toyorm.ExprEqual, Offsetof(Product{}.Count), 2).And().
    Condition(toyorm.ExprGreater, Offsetof(Product{}.Price), 3).Or().
    Condition(toyorm.ExprEqual, Offsetof(Product{}.Count), 4)
// WHERE count = 2 and price > 3 or count = 4
```

priority condition

```golang
brick.Where(toyorm.ExprGreater, Offsetof(Product{}.Price), 3).And().Conditions(
    brick.Where(toyorm.ExprEqual, Offsetof(Product{}.Count), 2).Or().
    Condition(toyorm.ExprEqual, Offsetof(Product{}.Count), 1).Search
)
// WHERE price > 3 and (count = 2 or count = 1)
```


```golang
brick.Conditions(
    brick.Where(toyorm.ExprEqual, Offsetof(Product{}.Count), 2).Or().
        Condition(toyorm.ExprEqual, Offsetof(Product{}.Count), 1).Search,
).And().Conditions(
    brick.Where(toyorm.ExprEqual, Offsetof(Product{}.Price), 3).Or().
        Condition(toyorm.ExprEqual, Offsetof(Product{}.Price), 4).Search,
)
// WHERE (count = ? OR count = ?) AND (price = ? OR price = ?)
```

limit & offset

```golang
brick := brick.Offset(2).Limit(2)
// LIMIT 2 OFFSET 2
```

order by

```golang
brick = brick.OrderBy(Offsetof(Product{}.Name))
// ORDER BY name
```


#### Transaction

start a transaction

```golang
brick = brick.Begin()
```

rollback all sql action

```golang
err = brick.Rollback()
```

commit all sql action

```golang
err = brick.Commit()
```

#### Debug

if Set debug all sql action will have log

```golang
brick = brick.Debug()
```


#### IgnoreMode

when I Update/Replace or Search with struct that have some zero value, did I update it ?


use IgnoreMode to differentiate what zero value should update

```golang
brick = brick.IgnoreMode(toyorm.IgnoreZero ^ toyorm.IgnoreZeroLen)
// ignore all zeor value but excloud zero len slice
// now value = []int(nil) will ignore
// but value = []int{} will update
```

---

mode              |  effective
------------------|---------------
IgnoreFalse       | ignore field type is bool and value is false
IgnoreZeroInt     | ignore field type is int/uint/uintptr(incloud their 16,32,64 bit type) and value is 0
IgnoreZeroFloat   | ignore field type is float32/float64 and value is 0.0
IgnoreZeroComplex | ignore field type is float32/float64 and value is 0.0
IgnoreNilString   | ignore field type is string and value is ""
IgnoreNilPoint    | ignore field type is point/map/slice and value is nil
IgnoreZeroLen     | ignore field type is map/array/slice and len = 0
IgnoreNullStruct  | ignore field type is struct and value is zero value struct e.g type A struct{A string,B int}, A{"", 0} will be ignore
IgnoreNil         | ignore with IgnoreNilPoint and IgnoreZeroLen
IgnoreZero        | ignore all of the above


#### Scope

// TODO

#### Thread safe

Thread safe if you comply with the following agreement

1. make sure **ToyBrick** object is read only, if you want to change it, create a new one

2. do not use **append** to change ToyBrick's slice data,use **make** and **copy** to clone new slice


### Full Feature Example

TODO