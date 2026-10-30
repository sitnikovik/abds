# Code Style Guide

> [!CAUTION]
> This is not a general Go style guide.
> This is the code style policy for this repository and its derived template-based projects

## How to read this document

**👀 Read this first**:

- Treat the top-level sections as project rules, not general Go advice.
- Read `❌ Bad` as an anti-pattern, `💀 Unacceptable` as a stronger anti-pattern, `✅ Good` as the default solution, and `⭐️ Better` as a refined version of the same idea.
- Some bad examples are **intentionally broken** to demonstrate failure modes. When that is the point, the section says so explicitly.
- Prefer the smallest example that proves the rule. Use larger multi-file examples only when package boundaries or interactions are part of the point.
- When adding new sections, follow the template and guidance in [docs/style/section-template.md](docs/style/section-template.md).

## Contents

- [How to read this document](#how-to-read-this-document)
- [Declarative style](#declarative-style)
  - [Immutability](#immutability)
  - [No Setters](#no-setters)
  - [No Getters](#no-getters)
  - [Functional options](#functional-options)
  - [One type — one constructor](#one-type--one-constructor)
  - [No logic in constructors](#no-logic-in-constructors)
  - [Custom types](#custom-types)
- [Vertical formatting](#vertical-formatting)
- [No newlines in functions](#no-newlines-in-functions)
- [No nil](#no-nil)
  - [Nullable](#nullable)
  - [Pointer structs](#pointer-structs)
  - [Errors can be nil](#errors-can-be-nil)
- [Document all](#document-all)
  - [What the function does](#what-the-function-does)
  - [What the field stores](#what-the-field-stores)
  - [Sentence structure](#sentence-structure)
  - [Multiple parameters / return values](#multiple-parameters--return-values)
  - [Single parameter / return values](#single-parameter--return-value)
  - [No arguments / return values](#no-arguments--return-values)
  - [Full example](#full-example)
  - [No comments](#no-comments)
- [Exporting](#exporting)
  - [Functions must be exported](#functions-must-be-exported)
  - [Any type must have a constructor](#any-type-must-have-a-constructor)
  - [Types must be exported](#types-must-be-exported)
  - [Fields must be unexported](#fields-must-be-unexported)
- [Test all](#test-all)
  - [Testify](#testify)
  - [Separate test packages](#separate-test-packages)
  - [Canonical test template](#canonical-test-template)
  - [Parallel tests](#parallel-tests)

## Declarative style

**📌 Rule**:

Favor a declarative programming style over imperative.
That means expressing the logic of computation without describing its control flow.

**🤔 Why**:

The document treats API design and data modeling as the primary tools for keeping code simple.
Declarative code tends to make invariants explicit and reduces the amount of stateful orchestration spread through the codebase.

### Immutability

> [!TIP]
> Immutability is the key to solve a lot of problems
that could be made in the projects!

**📌 Rule**:

The object cannot be changed after created!
If you really need to do it then create the new one
with the new value but specify it in a constructor.

**🤔 Why**:

Immutable objects are easier to reason about because their state does not change unexpectedly.
That reduces hidden side effects and makes the behavior of the code more predictable.

**🔗 Related rules**:

- [No setters](#no-setters)
- [Functional options](#functional-options)
- [No logic in constructors](#no-logic-in-constructors)

### No setters

> [!IMPORTANT]
> Setters are evil

**📌 Rule**:

Do not define setter methods that only assign values to struct fields.
The values of a struct should usually be fixed at creation time.
If you need a different value, create a new instance instead of mutating the existing one.

**🤔 Why**:

Setters make it possible to change the object from any place that has access to it.
That weakens encapsulation and makes the state of the value harder to predict.
The more mutable your structs are, the easier it is to introduce unexpected side effects.

**❌ Bad**:

The example below allows any caller to change the price after the order was created.
That means the original state of the object is no longer stable.

```go
package order

type Order struct {
    price int64
}

func NewOrder(price int64) *Order {
    return &Order{
        price: price,
    }
}

func (o *Order) GetPrice() int64 {
    return o.price
}

func (o *Order) SetPrice(price int64) {
    o.price = price
}
```

```go
package main

import (
    "fmt"
    "your_module/order"
)

func main() {
    ord := order.NewOrder(100)
    processOrder(ord)
    // Unexpectedly prints 200, but should be 100
    fmt.Println("Order price after processing:", ord.GetPrice())
}

func processOrder(o *order.Order) {
    o.SetPrice(200) // Oh, no! Some developer changed the price directly!
}
```

**✅ Good**:

The example below keeps the price immutable after construction.
If another price is needed, create another order value instead of mutating the existing one.

```go
package order

type Order struct {
    price int64
}

func NewOrder(price int64) *Order {
    return &Order{
        price: price,
    }
}

func (o *Order) Price() int64 {
    return o.price
}
```

```go
package main

import (
    "fmt"
    "your_module/order"
)

func main() {
    original := order.NewOrder(100)
    corrected := order.NewOrder(200)

    fmt.Println(original.Price())  // 100
    fmt.Println(corrected.Price()) // 200
}
```

**📝 Exceptions / Notes**:

This rule is a strong default, not a ban on every internal state change in every possible design.
But if you feel the need to add a setter, first assume that the struct design is wrong and revisit the API.

**🔗 Related rules**:

- [Immutability](#immutability)
- [No getters](#no-getters)
- [Functional options](#functional-options)
- [No logic in constructors](#no-logic-in-constructors)

### No getters

> [!IMPORTANT]
> Getters are evil

**📌 Rule**:

If you do not use setters - getters are useless.
If you really need to get some data from the struct - just omit `Get` prefix
in the method name. That is all!

**🤔 Why**:

If the method returns, for example, the `price` field - just name it `Price()`.
Why do you need `GetPrice()`?
Does the `Price() int64` not returns the price?
Why do you need `Get` prefix? Just avoid it.

**✅ Good**:

```go
package order

type Order struct {
    uuid string
    status string
    price int64
}

func NewOrder(
    uuid string,
    status string,
    price int64,
) *Order {
    return &Order{
        uuid: uuid,
        status: status,
        price: price,
    }
}

func (o *Order) UUID() string {
    return o.uuid
}

func (o *Order) Status() string {
    return o.status
}

func (o *Order) Price() int64 {
    return o.price
}
```

```go
package main

import (
    "fmt"
    "your_module/order"
)

func main() {
    ord1 := order.NewOrder("uuid1", "status1", 100)
    ord2 := order.NewOrder("uuid2", "status2", 200)
    fmt.Printf(
        "Order 1: %s, %s, %d\n",
        ord1.UUID(),
        ord1.Status(),
        ord1.Price(),
    ) // Prints uuid1, status1, 100
    fmt.Printf(
        "Order 2: %s, %s, %d\n",
        ord2.UUID(),
        ord2.Status(),
        ord2.Price(),
    ) // Prints uuid2, status2, 200
}

func processOrder(o *order.Order) {
    // There is no any method changes the input order!
}
```

Look! No `Get` prefix! Just `Price()` method that returns the price.
And we **do not have any setters**!
The `price` field is immutable after the struct created!
**We are sure that no one can change the price unexpectedly!**

**🔗 Related rules**:

- [No setters](#no-setters)
- [Immutability](#immutability)
- [Fields must be unexported](#fields-must-be-unexported)

### Functional options

> 💡 All the values could be set in a constructor

**📌 Rule**:

All! Remember that! Even there are a lot of fields in the struct,
they still could be specified in a constructor via **functional options**.

**🤔 Why**:

Functional options let the type keep one constructor while still supporting optional fields.
That keeps the API extensible without exploding the number of parameters or constructors.

**❌ Bad:**

File: *order.go*

```go
package order

type Order struct {
    uuid string
    status string
    userID string
    address string
    phone string
    price int64
    priority uint8
}

func NewOrder(
    uuid string,
    status string,
    userID string,
    address string,
    phone string,
    price int64,
    priority uint8,
) *Order {
    return &Order{
        uuid: uuid,
        status: status,
        userID: userID,
        address: address,
        phone: phone,
        price: price,
        priority: priority,
    }
}
```

**✅ Good:**

File: *order.go*

```go
package order

type Order struct {
    uuid string
    status string
    userID string
    address string
    phone string
    price int64
    priority uint8
}

func NewOrder(
    uuid string,
    status string,
    opts ...Option,
) *Order {
    o := &Order{
        uuid: uuid,
        status: status,
    }
    for _, opt := range opts {
        opt(o)
    }
    return o
}
```

File: *order.go*

```go
package order

type Option func(*Order)

func WithPrice(p int64) Option {
    return func(o *Order) {
        o.price = p
    }
}


func WithAddress(addr string) Option {
    return func(o *Order) {
        o.address = addr
    }
}

func WithPhone(phone string) Option {
    return func(o *Order) {
        o.phone = phone
    }
}

// etc...
```

File: *main.go*

```go
package main

import "github.com/you/repo/order"

func main() {
    fullOrder := order.NewOrder(
        "uuid1",
        "status1",
        order.WithPrice(100),
        order.WithAddress("Los Santos, Groove street"),
        order.WithPhone("88005553535"),
    )
    orderWithRequiredOnly := order.NewOrder(
        "uuid1",
        "status1",
    )
    orderWithPrice := order.NewOrder(
        "uuid1",
        "status1",
        order.WithPrice(100),
    )
    _, _, _ = fullOrder, orderWithRequiredOnly, orderWithPrice
}
```

Look! Here the type has many fields and only one constructor
that could be extended with functional options if we need
and allows us to create objects with different sets of values.

That way looks like ***overload*** in OOP languages like Java/C#/C++.

**🔗 Related rules**:

- [One type — one constructor](#one-type--one-constructor)
- [No logic in constructors](#no-logic-in-constructors)
- [No setters](#no-setters)

### One type — one constructor

**📌 Rule**:

According to functional options you may assume
that it may be optimized with making several constructors
for specific situations
to create objects with pre-defined set of arguments.

**🤔 Why**:

That is better than specifing a lof of arguments in the main constructor,
but one day you may have so many situations that you won't able to maintain.

**❌ Bad:**

```go
package order

type Order struct {
    uuid string
    status string
    userID string
    address string
    phone string
    price int64
    priority uint8
}

func NewOrder(
    uuid string,
    status string,
) *Order {
    return &Order{
        uuid: uuid,
        status: status,
    }
}

func NewOrderWithPrice(
    uuid string,
    status string,
    price int64,
) *Order {
    return &Order{
        uuid: uuid,
        status: status,
        price: price,
    }
}

func NewOrderWithAddress(
    uuid string,
    status string,
    addr string,
) *Order {
    return &Order{
        uuid: uuid,
        status: status,
        address: addr,
    }
}

func NewOrderWithPriceAndAddress(
    uuid string,
    status string,
    price int64,
    addr string,
) *Order {
    return &Order{
        uuid: uuid,
        status: status,
        price: price,
        address: addr,
    }
}
```

**💀 Unacceptable:**

File: *order.go*

```go
package order

type Order struct {
    uuid string
    status string
    userID string
    address string
    phone string
    price int64
    priority uint8
}

func NewOrder(
    uuid string,
    status string,
    userID string,
    address string,
    phone string,
    price int64,
    priority uint8,
) *Order {
    return &Order{
        uuid: uuid,
        status: status,
        userID: userID,
        address: address,
        phone: phone,
        price: price,
        priority: priority,
    }
}
```

**✅ Good:**

File: *order.go*

```go
package order

type Order struct {
    uuid string
    status string
    userID string
    address string
    phone string
    price int64
    priority uint8
}

func NewOrder(
    uuid string,
    status string,
    opts ...Option,
) *Order {
    o := &Order{
        uuid: uuid,
        status: status,
    }
    for _, opt := range opts {
        opt(o)
    }
    return o
}
```

File: *option.go*

```go
package order

type Option func(*Order)

func WithPrice(p int64) Option {
    return func(o *Order) {
        o.price = p
    }
}

func WithAddress(addr string) Option {
    return func(o *Order) {
        o.address = addr
    }
}

func WithPhone(phone string) Option {
    return func(o *Order) {
        o.phone = phone
    }
}

// etc...
```

**🔗 Related rules**:

- [Functional options](#functional-options)
- [No logic in constructors](#no-logic-in-constructors)
- [Any type must have a constructor](#any-type-must-have-a-constructor)

### No logic in constructors

**📌 Rule**:

Constructors are meant to initialize and return instances of types
and should not contain any logic or side effects.
They should simplify set up the struct with default values or provided parameters.

**🤔 Why**:

Especially do not return error! Error handling is also logic that is prohibited.

**❌ Bad:**

```go
package order

type Order struct {
    uuid string
    status string
    price int64
}

func NewOrder(
    uuid string,
    status string,
    price int64,
) *Order {
    o := &Order{
        uuid: uuid,
        status: status,
        price: price,
    }
    if price < 0 {
        o.price = 0
    } 
    if price == 0 && status == "" {
        o.status = "new"
    }
    return o
}
```

**💀 Unacceptable:**

```go
package order

import "errors"

type Order struct {
    uuid string
    status string
    price int64
}

func NewOrder(
    uuid string,
    status string,
    price int64,
) (*Order, error) {
    o := &Order{
        uuid: uuid,
        status: status,
        price: price,
    }
    if price < 0 {
        return nil, errors.New("empty price")
    } 
    if status == "" {
        return nil, errors.New("empty status")
    }
    return o
}
```

**✅ Good:**

```go
package order

type Order struct {
    uuid string
    status string
    price int64
}

func NewOrder(
    uuid string,
    status string,
    price int64,
) *Order {
    return &Order{
        uuid: uuid,
        status: status,
        price: price,
    }
}
```

Constructor can accept and move values from the arguments into the fields.
> 💡 This is what is precisely forms **encapsulation**

**🔗 Related rules**:

- [Immutability](#immutability)
- [Functional options](#functional-options)
- [One type — one constructor](#one-type--one-constructor)

### Custom types

**📌 Rule**:

Use type defintion to create custom types covers your business logic.
Especially when you denormalize your struct to optimize usage of your program.
In that case you may impair the readability of the project.

**🤔 Why**:

The main tip to catch that — you call fieilds with two words or more.

**❌ Bad:**

File: *order.go*

```go
package order

type Order struct {
    userID string
    userPhone string
    deliveryAddress string
    status string
    price uint64
}

func (o *Order) UserID() string {
    return o.userID
}

func (o *Order) UserPhone() string {
    return o.userPhone
}

func (o *Order) DeliveryAddress() string {
    return o.deliveryAddress
}

func (o *Order) Status() string {
    return o.status
}

func (o *Order) IsStatusNew() bool {
    return o.status == "new"
}

func (o *Order) IsStatusDelivered() bool {
    return o.status == "delivered"
}

func (o *Order) Price() uint64 {
    return o.price
}
```

**✅ Good:**

File: *order/order.go*

```go
package order

import (
    "github.com/you/repo/money/price"
    "github.com/you/repo/order/status"
    "github.com/you/repo/user"
)

type Order struct {
    delivery Delivery
    user user.Identity
    status status.Status
    price price.Price
}

func NewOrder(s status.Status, opts ...Option) *Order {
    o := &Order{
        status: s,
    }
    for _, opt := range opts {
        opt(o)
    }
    return o
}

func (o *Order) Delivery() Delivery {
    return o.delivery
}

func (o *Order) User() user.Identity {
    return o.user
}

func (o *Order) Status() status.Status {
    return o.status
}

func (o *Order) Price() price.Price {
    return o.price
}
```

File: *order/delivery.go*

```go
package order

import (
    "github.com/you/repo/delivery"
)

type Delivery struct {
    addr delivery.Address
}

func NewDelivery(addr delivery.Address) Delivery {
    return Delivery{
        addr: addr,
    }
}

func (d Delivery) Address() delivery.Address {
    return d.addr
}
```

File: *order/option.go*

```go
package order

import (
    "github.com/you/repo/money/price"
    "github.com/you/repo/user"
)
type Option func(*Order)

func WithDelivery(d Delivery) Option {
    return func(o *Order) {
        o.delivery = d
    }
}

func WithUser(usr user.Identity) Option {
    return func(o *Order) {
        o.user = usr
    }
}

func WithPrice(p price.Price) Option {
    return func(o *Order) {
        o.price = p
    }
}
```

File: *order/status/status.go*

> 💡 That way allows as to change initial type easy

```go
package status

type Status uint8

const (
    New = 0
    InProgress = 1
    Delivered = 2
)

func (s Status) IsNew() bool {
    return s == New
}

func (s Status) IsInProgress() bool {
    return s == InProgress
}

func (s Status) IsDelivered() bool {
    return s == Delivered
}
```

⚠️ To prevent creating unexpected `Status` we should omit constructor
and provide usage only with constants.

File: *user/identity.go*

```go
package user

import "github.com/you/repo/contact"

type Identity struct {
    phone contact.Phone
    id ID
}

func NewIdentity(id ID, opts ...IdentityOption) Identity {
    i := Identity{
        id: id,
    }
    for _, opt := range opts {
        opt(&i)
    }
    return i
}

func (i Identity) String() string {
    // That's can be useful in some logs or debugging.
    return i.id.String()
}

func (i Identity) Phone() contact.Phone {
    return i.phone
}

func (i Identity) ID() ID {
    return i.id
}
```

File: *user/identity_option.go*

```go
package user

import "github.com/you/repo/contact"

type IdentityOption func(i *Identity)

func WithPhone(p contact.Phone) IdentityOption {
    return func(i *Identity) {
        i.phone = p
    }
}
```

File: *user/id.go*

```go
package user

type ID string

func (i ID) String() string {
    return string(i)
}

// Here you can add some code covers logic of the identifiers
// for example checking if the id has new structure or
// if the id created after a year required.
```

File: *contact/phone.go*

```go
package contact

type Phone string

func (p Phone) String() string {
    return string(p)
}

// Here you can make logic for phones
// for example, extracting country code or get its provider name.
```

You may say that the "Good" code is too bigger
and complicated than "Bad". That's a mistake all developers
do once at least.
This way obligate you to think about the structure and maintainability
in the start of developing so you can cover with tests each point of
the main object or its secondary ones.
And think about that each type is not totally depends of others
so you can modify it quite easy.

**🔗 Related rules**:

- [No getters](#no-getters)
- [Types must be exported](#types-must-be-exported)
- [Fields must be unexported](#fields-must-be-unexported)

## Vertical formatting

**📌 Rule**:

Prefer vertical formatting for function calls, method calls, composite literals,
and other expressions that contain several parts.
If a call or expression contains arguments, options, or chained calls,
write each part on its own line.

**🤔 Why**:

Vertical formatting improves readability on small screens
and makes nested or extended expressions easier to scan.
It also makes diffs cleaner because adding, removing, or changing one argument
usually changes only one line.

**❌ Bad**:

The examples below pack too much information into one line.
That makes the code harder to scan and easier to miss during review.

```go
package main

import "your_module/order"

func main() {
    orderWithPrice := order.NewOrder("uuid1", "status1", order.WithPrice(100))
    _, _ = orderWithPrice, order.NewOrder("uuid2", "status2", order.WithPrice(200), order.WithAddress("Los Santos, Groov street"))
}
```

```go
package main

func main() {
    query := sql.Select("id", "email").From("users").Where("active = ?", true)
    _ = query
}
```

**✅ Good**:

The examples below give each argument or method call its own line.
That creates a stable visual structure and improves readability.

```go
package main

import "your_module/order"

func main() {
    orderWithPrice := order.NewOrder(
        "uuid1",
        "status1",
        order.WithPrice(100),
    )
    orderWithDetails := order.NewOrder(
        "uuid2",
        "status2",
        order.WithPrice(200),
        order.WithAddress("Los Santos, Groove street"),
    )
    _, _ = orderWithPrice, orderWithDetails
}
```

```go
package main

func main() {
    query := sql.
        Select("id", "email").
        From("users").
        Where("active = ?", true)
    _ = query
}
```

**📝 Exceptions / Notes**:

Calls without arguments may stay on one line.

If you are not sure which form to choose,
prefer the multiline version by default.

**🔗 Related rules**:

- [Functional options](#functional-options)
- [No newlines in functions](#no-newlines-in-functions)

## No newlines in functions

**📌 Rule**:

Avoid unnecessary newlines within function bodies to maintain readability and compactness.

**🤔 Why**:

This rule pushes you toward smaller functions with denser, easier-to-scan control flow.
If a function becomes hard to read without visual spacing, it is often a sign that the function should be smaller.

**❌ Bad:**

```go
package config

import "errors"

type Config struct {
    usr string
    pwd string
}

func (c Config) Empty() bool {
    return c.usr == "" && c.pwd == ""
}

func (c Config) Validate() error {
    if c.Empty() {
        return nil
    }

    if c.usr == "" {
        return errors.New("empty username")
    }

    if c.pwd == "" {
        return errors.New("empty password")
    }
    
    return nil
}
```

**✅ Good:**

```go
package config

import "errors"

type Config struct {
    usr string
    pwd string
}

func (c Config) Empty() bool {
    return c.usr == "" && c.pwd == ""
}

func (c Config) Validate() error {
    if c.Empty() {
        return nil
    }
    if c.usr == "" {
        return errors.New("empty username")
    }
    if c.pwd == "" {
        return errors.New("empty password")
    }
    return nil
}
```

First time you may decide that this rule is too harsh and excessive.
Please believe me! That way makes us to build very small functions and simplify the code.
Just try it and you will see the benefits soon!

## No nil

> [!TIP]
> Just do not use `nil` and your code will become much more reliable!

**📌 Rule**:

As you can know the `nil` is the same as `NULL` in Java, C++, or `None` in Python.
And we **totally prohibit** this usage in our projects.

**🤔 Why**:

`nil` pushes uncertainty into every consumer of a value.
This section explains how to model optionality and pointer-based APIs without normalizing `nil` checks across the whole codebase.

> [!CAUTION]
> Some `❌ Bad` examples in this section are intentionally broken and may not compile.
> This is deliberate: the point is to demonstrate how `nil`-driven APIs quickly become unsafe, awkward, or misleading.

### Nullable

**📌 Rule**:

Long time ago `NULL` was created to display that value was not set in RDBMS and soon it arrived to the code.
In that moment it seemed to be useful and friendly for developers
but now it looks like a **bad manner**.

**🤔 Why**:

Why you should use it?
Why you need to know that the values was set?
Can't you understand that the value was not set with empty value?
What if a user set precisely empty value?
You can say that you have to check that case.
OK, but think about maintanance in your code...

**❌ Bad:**

```go
package user

type User struct {
    name *string
}

func NewUser(name *string) *User {
    return &User{
        name: name,
    }
}

func (u *User) Name() *string {
    return u.name
}

// ...
```

```go
package main

func main() {
    // ...
    u := NewUser("John Doe")
    if name := u.Name(); name == "" { // will panic if the flag would be nil!
        // ...
    } else {
        // ...
    }
    // ...
}
```

Or how it will look like if we don't know whether the value is nil or not?

```go
package main

func main() {
    // ...
    var name *string
    if v != nil {
        name = &string(v)
    }
    u := NewUser(name)
    if nam := u.Name(); nam == "" { // will panic if the flag would be nil!
        // ...
    } else {
        // ...
    }
    // ...
}
```

You may say "just check for nil":

```go
package main

func main() {
    // ...
    var name *string
    if v != nil {
        name = &string(v)
    }
    if name := u.Name(); name != nil && name == "" {
        // ...
    } else if name == "" {
        // ...
    } else {
        // ...
    }
    // ...
}
```

It would be if you decice to add check to the predicate like this:

```go
package user

// ...
func (u *User) Name() string {
    if u.name == nil {
        return string
    }
    return u.name
}
// ...
```

But any solving will look like a crunch in your code
and at the same moment you lost the one "feature" of the `NULL`.
So why you still use it?

If your really need to check whether some value is really set by a user
you should do like this:

**✅ Good:**

File: *strings/string.go*.

```go
package strings

type String struct {
    v string
    ex bool
}

func NewString(v string, opts ...Option) String {
    s := String{
        v: v,
    }
    for _, opt := range opts {
        opt(&s)
    }
    return s
}

func (s String) Value() string {
    return s.v
}

func (s String) Empty() bool {
    return s.v != ""
}

func (s String) Exists() bool {
    return s.ex || s.v != ""
}
```

File: *strings/option.go*.

```go
package strings

type Option func(*String)

func WithExists(exists bool) Option {
    return func(s *String) {
        s.ex = exists
    }
}
```

```go
package main

func main() {
    // ...
    usr := NewUser(
        strings.NewString("John Doe"),
    )
    if nam := user.Name(); !nam.Exists() {
        
    } else if nam.Empty() {
        
    } else {

    }
    // ...    
}
// ...
```

What if we don't know whether the value is not empty? Then:

```go
package main

func main() {
    // ...
    usr := NewUser(
        strings.NewString(name),
        strings.WithExists(name != ""),
    )
    if nam := user.Name(); !nam.Exists() {
        // ...
    } else if nam.Empty() {
        // ...
    } else {
        // ...
    }
    // ...    
}
// ...
```

That way you delegate the logic to check for nil to `User`,
so you can never check for it and get rid of nil-pointer errors.
These errors look like childish mistakes that develovers just forgot to handle.

**🔗 Related rules**:

- [Pointer structs](#pointer-structs)
- [Errors can be nil](#errors-can-be-nil)
- [Custom types](#custom-types)

### Pointer structs

> [!TIP]
> Prevent and catch all nil-usage except of errors

**📌 Rule**:

If we store structs with a pointer and do not use declarative style,
we are guaranteed to meet panics with nil-pointer error.

Do not check objects and values for nil like that:

**🤔 Why**:

When a struct field may silently remain `nil`, every consumer has to defensively remember that fact.
That spreads uncertainty through the whole API and makes panics much easier to introduce.

**❌ Bad:**

File: **order/order.go**

```go
package order

type Order struct {
    User *User
}

func (o *Order) GetUserName() string {
    return o.User.Name
}
```

File: **order/user.go**

```go
package order

type User struct {
    Name string
}
```

File: **main.go**

```go
package main

import (
    "fmt"

    "github.com/you/repo/order"
)

func main() {
    // ...
    o := new(order.Order)
    if usr != nil {
        o.User = usr
    }
    // ...
    fmt.Println(o.User.Name) // will panic because o.User is nil
}
// ...
```

If we apply declarative style it will look like this:

**✅ Good:**

File: **order/order.go**

```go
package order

type Order struct {
    user *User
}

func NewOrder(usr *User) *Order {
    return &Order{
        user: usr,
    }
}

func (o *Order) User() *User {
    return o.user
}
```

File: **order/user.go**

```go
package order

type User struct {
    name string
}

func NewUser(name string) *User {
    return &User{
        name: name,
    }
}

func (u *User) Name() string {
    return u.name
}
```

File: **main.go**

```go
package main

import (
    "fmt"

    "github.com/you/repo/order"
)

func main() {
    // ...
    usr := order.NewUser("John Doe")
    o := order.NewOrder(usr)
    // ...
    fmt.Println(o.User().Name())
}
// ...
```

And you may say that if user is still nil it will panic.
Yes, it does. But this code is not able to fix all problems with nil.

> [!IMPORTANT]
> Pointer it is precisely a pointer and you should not store nil in the struct fields

Also consider usage without pointers.

**🔗 Related rules**:

- [Nullable](#nullable)
- [No setters](#no-setters)
- [Fields must be unexported](#fields-must-be-unexported)

### Errors can be nil

**📌 Rule**:

The one thing you can to check for `!= nil` is error.
But after you have read the upper paragraphs
you may think that we have to check errors something like that:

**🤔 Why**:

Errors are the idiomatic Go exception to the “no nil” rule.
Returning `err` explicitly keeps failure handling visible to the caller and aligned with Go tooling and conventions.

**❌ Bad**:

```go
package main

func main() {
    // ...
    value := cache.Get(ctx, "key")
    if value.Errored() {
        log.Fatalf("failed to get value from cache: %v", value.Err())
    }
    // handle the value...    
}
```

**✅ Good**:

And I was a fan of this style.
But one day I recognized that when get these values wrapps error inside itself — **it obligates you to know that this staff may be errored**.
And that is not good.

The best way you can do to check for errors is checking for nil.
Compare the code with this:

```go
package main

func main() {
    // ...
    value, err := cache.Get(ctx, "key")
    if err != nil {
        log.Fatalf("failed to get value from cache: %v", err)
    }
    // handle the value...    
}
```

This way obligates you to handle error
and many linters will alert you about that.

Maybe the first example is truly object oriented style
but it is not idiomatic in Go and the handling becomes not maintainable.

**🔗 Related rules**:

- [Nullable](#nullable)
- [No logic in constructors](#no-logic-in-constructors)

## Document all

> [!IMPORTANT]
> Even unexported! Even if the value is too obvious! Just document it!

**📌 Rule**:

Every package, function, method, field and type must have a clear
and concise comment explaining its purpose and usage.
Use [GoDoc conventions](https://go.dev/doc/comment) for documentation.

**🤔 Why**:

Documentation is part of the public contract of the code.
If a new contributor cannot understand what a symbol is for from its comment, the API is harder to use correctly.
What is obvious to one developer may be unclear to another one.
These doc comments also become a lightweight built-in documentation system when the project has no dedicated Markdown docs, wiki, or Confluence pages.

### What the function does

**📌 Rule**:

Describe what the function does and not how it does it.
The doc have to include a summary of the function's behavior,
its inputs, outputs, and any side effects.

### What the field stores

**📌 Rule**:

Describe what the field stores and how it is used.

### Sentence structure

**📌 Rule**:

The first sentence should be a short summary that starts with the name of the function.
**After a blank line**, provide more detailed information if necessary.
The comment should be written in full sentences and be clear
to someone unfamiliar with the code.

### Multiple parameters / return values

**📌 Rule**:

If the function has more than one parameter or returns values,
describe each of them clearly.
Specify the type and purpose of each parameter and return value.

**✅ Good**:

```go
package parser

// Tokenize splits the input string into a slice of substrings based on the specified delimiter.
//
// It returns a slice of substrings and an error if the operation fails.
//
// Parameters:
//   - input: The string to be tokenized.
//   - delimiter: The delimiter used to split the input string.
// 
// Example:
//
//  Tokenize("a,b,c", ",") // returns []string{"a", "b", "c"}, nil
func Tokenize(input string, delimiter string) ([]string, error) {
    return nil, nil
}
```

**How it looks in VSCode:**

![how it looks in VSCode](docs/style//img/godoc_02.png)

### Single parameter / return value

**📌 Rule**:

If the function has only one parameter and one return value,
a brief description is sufficient but if necessary, you can still elaborate.

**✅ Good**:

```go
package parser

type Type string

// ParseType parses the type from a commit subject string.
//
// The expected format has the type at the beginning of the subject,
// optionally followed by a scope in parentheses and a colon.
//
// For example:
//
//  "feat(ui): add new button" // type is "feat"
//  "fix: resolve issue" // type is "fix"
func ParseType(s string) Type {
    return ""
}
```

**How it looks in VSCode:**

![how it looks in VSCode](docs/style/img/godoc_01.png)

### No arguments / return values

**📌 Rule**:

If the function has no parameters or does not return any values,
a simple description of its purpose is sufficient.

**✅ Good**:

```go
package bootstrap

// Initialize sets up the necessary configurations for the application.
func Initialize() {}
```

### Full example

**📌 Rule**:

When a section needs a complete example, prefer one small exported type with documented unexported fields and documented public methods.

**✅ Good**:

```go
package user

// User describes an account in the customer area of an online store.
//
// User stores account data and exposes it through public methods.
// All fields are unexported to keep the state encapsulated.
type User struct {
    // id stores the stable internal identifier of the user.
    //
    // The value is used to reference the user in storage, logs,
    // and communication between application components.
    id string

    // email stores the email address of the user.
    //
    // The value is used for authentication, notifications,
    // and order-related communication.
    email string

    // firstName stores the personal name of the user.
    //
    // The value is used in the profile and other interface elements
    // that address the user directly.
    firstName string

    // lastName stores the family name of the user.
    //
    // The value is used together with firstName when the application
    // needs to show the full name of the user.
    lastName string
}

// NewUser creates a new User instance from the provided account data.
//
// NewUser only assigns the input values to the struct fields.
//
// Parameters:
//   - id: stable internal identifier of the user.
//   - email: email address used for authentication and notifications.
//   - firstName: personal name of the user.
//   - lastName: family name of the user.
//
// Returns:
//   - *User: initialized user instance.
func NewUser(
    id string,
    email string,
    firstName string,
    lastName string,
) *User {
    return &User{
        id: id,
        email: email,
        firstName: firstName,
        lastName: lastName,
    }
}

// ID returns the stable internal identifier of the user.
func (u *User) ID() string {
    return u.id
}

// Email returns the email address of the user.
func (u *User) Email() string {
    return u.email
}

// FullName returns the full display name of the user.
func (u *User) FullName() string {
    return u.firstName + " " + u.lastName
}
```

### No comments

**📌 Rule**:

Do not write comments inside function bodies.
If something matters enough to explain, put that explanation into the doc comment of the function, method, type, or field.

**🤔 Why**:

Inline comments inside the function body usually compensate for naming or structure that should already be clear from the API.
When the explanation sits in the doc comment instead, it helps both readers of the code and users of the symbol in editors and generated docs.

**❌ Bad**:

```go
package order

type Order struct {
    price int64
}

func (o *Order) ApplyDiscount(discount int64) int64 {
    // subtract the discount from the current price
    price := o.price - discount
    // return the discounted price to the caller
    return price
}
```

**✅ Good**:

```go
package order

type Order struct {
    price int64
}

// ApplyDiscount returns the discounted price of the order.
//
// ApplyDiscount calculates the result from the current order price
// and the provided discount value.
func (o *Order) ApplyDiscount(discount int64) int64 {
    price := o.price - discount
    return price
}
```

**🔗 Related rules**:

- [Document all](#document-all)
- [What the function does](#what-the-function-does)
- [Sentence structure](#sentence-structure)

## Exporting

**📌 Rule**:

Exporting rules define what should be visible outside the package and what must stay internal.

**🤔 Why**:

These rules keep package boundaries explicit and make APIs easier to reuse consistently.
Functions and types are exported so they can participate in signatures, composition, and reuse across packages.
Fields stay unexported so the package keeps control over invariants, naming, and future evolution.

### Functions must be exported

**📌 Rule**:

Every function or method is intended to be used outside its package and must be exported
(i.e., its name should start with an uppercase letter).
This makes the function accessible to other packages.

We do not allow any unexported functions in the codebase!

> ⚠️ If you dedicate a code in a function - it must be accessible from others outside the package!

**❌ Bad:**

```go
package file

import "os"

func read(path string) ([]byte, error) {
    return os.ReadFile(path)
}
```

**✅ Good:**

```go
package file

import "os"

// Read reads file by provided path and returns
// the slice of bytes read from the file
// or returns error if something went wrong.
func Read(path string) ([]byte, error) {
    return os.ReadFile(path)
}
```

**❌ Bad:**

```go
package main

import (
    "log"

    "your_module/file"
)

func main() {
    cfg, err := file.ReadConfig("config.json")
    if err != nil {
        log.Fatal(err)
    }
    // ...
}
```

```go
package file

import (
    "encoding/json"
    "os"
)

type Config struct{}

// ReadConfig reads and parses a config file.
func ReadConfig(path string) (*Config, error) {
    data, err := readFile(path)
    if err != nil {
        return nil, err
    }
    return parseConfig(data)
}

func readFile(path string) ([]byte, error) {
    return os.ReadFile(path)
}

func parseConfig(data []byte) (*Config, error) {
    var cfg Config
    err := json.Unmarshal(data, &cfg)
    return &cfg, err
}
```

**✅ Good:**

```go
package main

import (
    "log"
    "your_module/file"
)

func main() {
    cfg, err := file.ReadConfig("config.json")
    if err != nil {
        log.Fatal(err)
    }
    // ...
}
```

```go
package file

import (
    "encoding/json"
    "os"
)

type Config struct{}

// ReadConfig reads and parses a config file.
func ReadConfig(path string) (*Config, error) {
    data, err := Read(path)
    if err != nil {
        return nil, err
    }
    return ParseConfig(data)
}

// Read reads a file and returns its contents.
func Read(path string) ([]byte, error) {
    return os.ReadFile(path)
}

// ParseConfig parses config data from bytes.
func ParseConfig(data []byte) (*Config, error) {
    var cfg Config
    err := json.Unmarshal(data, &cfg)
    return &cfg, err
}
```

**⭐️ Better:**

```go
package main

import (
    "log"
    "your_module/file"
    "your_module/config"
)

func main() {
    data, err := file.Read("config.json")
    if err != nil {
        log.Fatal(err)
    }
    cfg, err := config.ParseConfig(data)
    if err != nil {
        log.Fatal(err)
    }
    // ...
}
```

```go
package file

import "os"

// Read reads a file and returns its contents.
func Read(path string) ([]byte, error) {
    return os.ReadFile(path)
}
```

```go
package config

import "encoding/json"

type Config struct{}

func ParseConfig(data []byte) (*Config, error) {
    var cfg Config
    err := json.Unmarshal(data, &cfg)
    return &cfg, err
}
```

**🔗 Related rules**:

- [Any type must have a constructor](#any-type-must-have-a-constructor)
- [Types must be exported](#types-must-be-exported)
- [Fields must be unexported](#fields-must-be-unexported)

### Any type must have a constructor

**📌 Rule**:

Each struct or type must have a constructor function that initializes and returns an instance of the struct.
The constructor should be named `New<TypeName>` and should set default values for the struct's fields if necessary.

### Types must be exported

**📌 Rule**:

All structs and types that are intended to be used outside their package must be exported
(i.e., their names should start with an uppercase letter).
This allows other packages to define them in outer function signatures.

**🤔 Why**:

But! **The fields of the structs must be unexported** to maintain encapsulation
and prevent direct access from outside the package.

### Fields must be unexported

**📌 Rule**:

**All struct fields must be unexported** to maintain encapsulation and prevent direct access from outside the package.

**🤔 Why**:

When defining structs, ensure that all fields are unexported (i.e., start with a lowercase letter).
This encapsulation helps maintain control over how data is accessed and modified.

## Test all

**📌 Rule**:

Every new feature or bugfix must include corresponding tests.
Use Go's testing package to write unit tests that cover various scenarios and edge cases.

**🤔 Why**:

The rest of this style guide pushes complexity into APIs, types, and contracts.
Tests are what make those rules safe to evolve over time.

### Testify

**📌 Rule**:

Use [testify package](https://github.com/stretchr/testify)
to make **asserts** in your tests.

**🤔 Why**:

> ❌ Never use `t.Error()` because it definitely damages the reading!

### No table-driven tests

**📌 Rule**:

Do not use table-driven tests.
Write every case as its own
`t.Run(name, func(t *testing.T) { ... })` block with inline setup, act, and assert.
Do not collect cases into a `tests := []struct{...}` slice
and do not define `args`/`want` helper structs.

**🤔 Why**:

Both VSCode and JetBrains Goland can generate unit-test templates in one click,
which boosts productivity, but both produce table-driven tests.
That is not the style we use in this project.
We write every `t.Run()` block by hand so that each case stays absolutely
independent from the others and can assert exactly what it needs, in its own way.

**❌ Bad:**

```go
package file

import "testing"

func TestName_String(t *testing.T) {
    tests := []struct {
        name string
        n    Name
        want string
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.n.String(); got != tt.want {
                t.Errorf("Name.String() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

**✅ Good:**

```go
package file_test

import (
    "testing"

    impl "github.com/you/repo/file"
    "github.com/stretchr/testify/assert"
)

func TestName_String(t *testing.T) {
    t.Parallel()

    t.Run("returns original file name", func(t *testing.T) {
        t.Parallel()
        n := impl.NewName("example.txt")
        assert.Equal(
            t,
            "example.txt",
            n.String(),
        )
    })
}
```

> 💡 Do not ignore `require` in testify package to ensure yourself
that some field is not nil pointer
or the implementation must have some functionality

**🔗 Related rules**:

- [Parallel tests](#parallel-tests)

### Canonical test template

**📌 Rule**:

Use one canonical test shape as the default starting point for ordinary unit tests.
Each case is its own `t.Run(name, func(t *testing.T) { ... })` block with inline setup, act, and assert.
Do not collect cases into a `tests := []struct{...}` slice and do not define `args`/`want` helper structs.

**🤔 Why**:

When every case is a self-contained `t.Run` block, the setup, the action, and the assertion sit next to each other.
That makes each case readable on its own, easier to review, copy, and extend without touching the others.

**✅ Good: only value returned**:

```go
package file_test

import (
    "testing"

    impl "github.com/you/repo/file"
    "github.com/stretchr/testify/assert"
)

func TestName_String(t *testing.T) {
    t.Parallel()

    t.Run("returns original file name", func(t *testing.T) {
        t.Parallel()
        n := impl.NewName("example.txt")
        assert.Equal(
            t,
            "example.txt",
            n.String(),
        )
    })
}
```

**✅ Good: value and error**:

```go
package file_test

import (
    "testing"

    impl "github.com/you/repo/file"
    "github.com/stretchr/testify/assert"
)

func TestParseName(t *testing.T) {
    t.Parallel()

    t.Run("parses valid file name", func(t *testing.T) {
        t.Parallel()
        got, gotErr := impl.ParseName("example.txt")
        assert.Equal(t, impl.NewName("example.txt"), got)
        assert.NoError(t, gotErr)
    })

    t.Run("returns error for invalid file name", func(t *testing.T) {
        t.Parallel()
        got, gotErr := impl.ParseName("")
        assert.Equal(t, impl.Name{}, got)
        assert.ErrorIs(t, gotErr, impl.ErrInvalidName)
    })
}
```

**✅ Good: bool flag only**:

```go
package file_test

import (
    "testing"

    impl "github.com/you/repo/file"
    "github.com/stretchr/testify/assert"
)

func TestName_HasExt(t *testing.T) {
    t.Parallel()

    t.Run("returns true when file has extension", func(t *testing.T) {
        t.Parallel()
        n := impl.NewName("example.txt")
        assert.Equal(t, true, n.HasExt())
    })

    t.Run("returns false when file has no extension", func(t *testing.T) {
        t.Parallel()
        n := impl.NewName("example")
        assert.Equal(t, false, n.HasExt())
    })
}
```

**🔗 Related rules**:

- [Testify](#testify)
- [Separate test packages](#separate-test-packages)
- [Parallel tests](#parallel-tests)

### Separate test packages

**📌 Rule**:

Add suffix `_test` to the name of the package you test.
It obligates you to import the implementation
and you have to declare a constructor for the type you test.

**🤔 Why**:

Testing through the public package boundary protects encapsulation and keeps you honest about the exported API.

**❌ Bad:**

```go
package file

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestName_String(t *testing.T) {
    t.Parallel()

    t.Run("some case", func(t *testing.T) {
        t.Parallel()
        n := Name{ // Allows to create with a literal
            value: "example.txt", // With unexported fields!
        }
        assert.Equal(
            t,
            "example.txt",
            n.String(),
        )
    })
}
```

**✅ Good:**

```go
package file_test

import (
    "testing"

    "github.com/you/repo/file"
    "github.com/stretchr/testify/assert"
)

func TestName_String(t *testing.T) {
    t.Parallel()

    t.Run("some case", func(t *testing.T) {
        t.Parallel()
        n := file.NewName("example.txt") // Obligates to use the constructor
        assert.Equal(
            t,
            "example.txt",
            n.String(),
        )
    })
}
```

Using unexported fields in tests is also prohibited.
But it still does not allow you to use this fields of the implementation that you test!

**🔗 Related rules**:

- [Fields must be unexported](#fields-must-be-unexported)
- [Any type must have a constructor](#any-type-must-have-a-constructor)

### Parallel tests

**📌 Rule**:

Use `t.Parallel()` in all test cases where it could be used.
But be careful to use it in the cases with asyncronious implementation
because it may bring your test flacky so you may spend a lot of time to find out
*yet another test you skip every second CI-job*

**📝 Exceptions / Notes**:

Parallelism is the default, but not when shared state, timing sensitivity, or asynchronous behavior makes the test unstable.
