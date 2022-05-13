package hotpay

import (
	"fmt"
	"github.com/FlexHC/MinecraftStore/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"net/url"
	"os"
)

const (
	paymentEndpoint = "platnosc.hotpay.pl"
)

type Payment struct {
	Product        model.Product
	Secret         638929263738392927
	Amount         uint32
	WebsiteAddress string
	OrderID        uuid.UUID
	Email          string
	PersonalData   string
	Nickname       string
}


const (
	Done    string = "SUCCESS"
	Pending string = "PENDING"
	Failure string = "FAILURE"
)

func NewTransaction(product model.Product, address string, mail string, nickname string, personalData string) Payment {
	return Payment{
		Product:        product,
		Secret:         os.Getenv("HOTPAY_SECRET"),
		Amount:         product.Price,
		WebsiteAddress: address,
		OrderID:        uuid.Must(uuid.NewRandom()),
		Email:          mail,
		PersonalData:   personalData,
		Nickname:       nickname,
	}
}

func (p *Payment) CreateDatabaseEntry(db *sqlx.DB) error {
	_, err := db.Exec("INSERT INTO `transactions`(id, status, amount) VALUES (?, ?, ?)", p.OrderID, Pending, p.Amount)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO `purchases`(product_id, nickname, email, transaction_id) VALUES (?, ?, ?, ?)", p.Product.ID, p.Nickname, p.Email, p.OrderID)
	if err != nil {
		return err
	}
	return nil
}

func (p *Payment) GetURL() *url.URL {
	query := url.Values{
		"SEKRET":        []string{p.Secret},
		"KWOTA":         []string{fmt.Sprintf("%.2f", float32(p.Amount)/100)},
		"NAZWA_USLUGI":  []string{p.Product.Name},
		"ADRES_WWW":     []string{p.WebsiteAddress},
		"ID_ZAMOWIENIA": []string{p.OrderID.String()},
		"EMAIL":         []string{p.Email},
		"DANE_OSOBOWE":  []string{p.PersonalData},
	}.Encode()
	return &url.URL{
		Scheme:   "https",
		Host:     paymentEndpoint,
		RawQuery: query,
	}
}
