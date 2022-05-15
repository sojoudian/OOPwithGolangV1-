package main

import (
	"errors"
	"fmt"
	"log"
)

var (
	ErrorInSufficientBalance = errors.New("insufficient balance")
	ErrorOverdraftIncurred   = errors.New("overdraft incurred")
)

type Depositable interface {
	Deposit(float64)
}
type Withdrawable interface {
	Withdraw(float64) (float64, error)
}

// BankAccount emulating what a class would do
type BankAccount struct {
	Owner   string
	balance float64
}

func (ba BankAccount) Balance() float64 {
	return ba.balance
}

func (ba *BankAccount) Deposit(amount float64) {
	// ba.balance = ba.balance + amount
	ba.balance += amount
}

func (ba *BankAccount) Withdraw(amount float64) (float64, error) {
	if ba.balance < amount {
		return 0, ErrorInSufficientBalance
	}
	ba.balance -= amount
	return ba.balance, nil
}

// OverdraftableBankAccount +++++++++++++++++
type OverdraftableBankAccount struct {
	BankAccount
	Fee float64
}

func (oba *OverdraftableBankAccount) Withdraw(amount float64) (float64, error) {
	var err error
	if oba.balance < amount {
		oba.balance -= oba.Fee
		err = ErrorOverdraftIncurred
	}
	oba.balance -= amount
	return oba.balance, err
}

//func Transfer(debtor BankAccount, creditor BankAccount, amount float64) error {
//	return nil
//}
// after creatin the interfaces this function will bechange into:
func Transfer(debtor Withdrawable, creditor Depositable, amount float64) error {
	balance, err := debtor.Withdraw(amount)
	switch err {
	case ErrorInSufficientBalance:
		return err
	case ErrorOverdraftIncurred:
		log.Printf("debtor incurred overdraft, new balance is %.2F", balance)
	}
	creditor.Deposit(amount)
	return nil

}

func main() {
	accOne := &BankAccount{"Maz", 50}
	accOne.Deposit(20)
	fmt.Println("Balance:", accOne.Balance())

	a2 := &OverdraftableBankAccount{BankAccount{"Delphini", 100}, 20}
	a2.Deposit(30)
	fmt.Println("Balance for Delphini:", a2.balance)

	_, err := a2.Withdraw(150)
	if err != nil {
		log.Printf("Overdraft incurred, balance is now %.2f", a2.Balance())
	}

	_, err = accOne.Withdraw(150)
	if err != nil {
		log.Printf("Balance: %.2F", accOne.Balance())
	}

	//Also, we can have a2.BankAccount.Withdraw() which acts like accOne.BankAccount.Withdraw()
	// But now, we don't want that.

	// Put money back into a2 account
	a2.Deposit(100)
	//Now in most of the OOP languages we can do
	//err = Transfer(accOne, a2, 50)
	// it will return this error:
	//./main.go:82:25: cannot use a2 (variable of type OverdraftableBankAccount) as type BankAccount in argument to Transfer
	// That is called inheritance
	//You would be able to pass in a type hat is inherits
	//from another type
	// +++++++++++++++++++++++++++++++++ BUT YOU CAN DO THAT IN GO +++++++++++++++++++++++++++++++++
	//so what you have to do is instead of using inheritance you have to use composition AKA you have
	// to use interface, we will do this in line 14

	fmt.Println("Account 1 balance: ", accOne.Balance(), "Account 2 balance: ", a2.Balance())
	err = Transfer(accOne, a2, 100)
	if err != nil {
		log.Printf("Could not compelete transfer: %v", err)
	}

	err = Transfer(a2, accOne, 100)
	if err != nil {
		log.Printf("Could not compelete transfer: %v", err)
	}
	fmt.Println("Account 1 balance: ", accOne.Balance(), "Account 2 balance: ", a2.Balance())
}
