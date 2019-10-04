package database

import (
	"payments/models"

	log "github.com/sirupsen/logrus"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// PaymentService holds the collection
type PaymentService struct {
	collection *mgo.Collection
}

// NewPaymentService creates a new payment service
func NewPaymentService(session *mgo.Session, db string, col string) *PaymentService {
	collection := session.DB(db).C(col)
	collection.EnsureIndex(paymentModelIndex())
	return &PaymentService{collection}
}

// CreatePayment creates a new payment
func (p *PaymentService) CreatePayment(u *models.Payment) error {
	payment, err := newPaymentModel(u)
	if err != nil {
		return err
	}
	return p.collection.Insert(&payment)
}

// GetPaymentsByIDs gets a list of payments
func (p *PaymentService) GetPaymentsByIDs(ids []string) ([]models.Payment, error) {
	output := []models.Payment{}

	for _, id := range ids {
		model := paymentModel{}
		err := p.collection.Find(bson.M{"some_other_id": id}).One(&model)
		if err != nil {
			log.Warnf("innexistend id requested %s", id)
		} else {
			payment := models.Payment{
				ID:             model.ID.Hex(),
				Type:           model.Type,
				Version:        model.Version,
				OrganisationID: model.OrganisationID,
				SomeOtherID:    model.SomeOtherID}

			output = append(output, payment)
		}
	}
	return output, nil
}

// DeletePayment deletes it
func (p *PaymentService) DeletePayment(id string) error {
	err := p.collection.Remove(bson.M{"some_other_id": id})
	if err != nil {
		return err
	}
	return nil
}

// UpdatePayment updates a payment by some id
func (p *PaymentService) UpdatePayment(u *models.Payment) error {

	filter := bson.D{{"some_other_id", u.SomeOtherID}}

	err := p.collection.Update(filter, u)
	if err != nil {
		return err
	}
	return nil
}
