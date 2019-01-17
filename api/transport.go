package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vstarcevic/trader/trader"
)

// MakeHTTPHandler create handler for routes
func MakeHTTPHandler(svc Service) http.Handler {

	r := mux.NewRouter()

	r.HandleFunc("/contact", svc.contactCreateHandler).Methods("POST")
	r.HandleFunc("/checksubscription", svc.checkSubscriptionHandler).Methods("GET")
	r.HandleFunc("/findbyany/{search}", svc.findByAnyHandler).Methods("GET")
	r.HandleFunc("/contactsubscription", svc.contactSubscriptionHandler).Methods("POST")
	return r

}

func (s Service) contactCreateHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var contact trader.Contact
	err := decoder.Decode(&contact)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if contact.Clientid == 0 || contact.Broker == "" || contact.Country == "" || contact.Identifier == "" || contact.Language == "" {
		err := fmt.Errorf("Missing mandatory data")
		http.Error(w, err.Error(), 500)
		return
	}

	resp, err := s.ContactRepo.Create(&contact)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	response := ContactCreateResponse{Contactid: resp}

	out, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))

}

func (s Service) checkSubscriptionHandler(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	clientid := params.Get("clientid")
	subscriptionType := params.Get("subscriptiontype")
	result, err := s.SubscriptionRepo.CheckSubscription(clientid, subscriptionType)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	response := CheckSubscriptionResponse{SubscriptionExists: result}

	out, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))

}

func (s Service) findByAnyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	search := vars["search"]
	results, err := s.SubscriptionRepo.FindByAny(search)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var response []ContactSubscriptionResponse
	for _, result := range results {
		row := ContactSubscriptionResponse{
			ContactResponse: ContactResponse{
				Clientid:   strconv.Itoa(result.Contact.Clientid),
				Broker:     result.Contact.Broker,
				Country:    result.Contact.Country,
				Identifier: result.Contact.Identifier},
			SubscriptionResponse: SubscriptionResponse{
				SubscriptionType: result.SubscriptionType.Name,
				AccountCode:      result.Subscription.Account},
		}

		response = append(response, row)
	}

	out, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}

func (s Service) contactSubscriptionHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var subscrReq ContactSubscriptionRequest
	err := decoder.Decode(&subscrReq)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if subscrReq.Contactid == "" || subscrReq.SubscriptionType == "" {
		err := fmt.Errorf("Missing mandatory data")
		http.Error(w, err.Error(), 500)
		return
	}

	contact, err := s.ContactRepo.GetByID(subscrReq.Contactid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	subscrTypeid := s.SubscriptionRepo.GetSubscriptionTypeIDFromName(subscrReq.SubscriptionType)
	if subscrTypeid == 0 {
		err := fmt.Errorf("No subscription with that type")
		http.Error(w, err.Error(), 500)
		return
	}

	subscriptionid := s.SubscriptionRepo.GetFirstSubscriptionForType(strconv.Itoa(subscrTypeid))
	if subscriptionid == 0 {
		err := fmt.Errorf("No subscription")
		http.Error(w, err.Error(), 500)
		return
	}

	err = s.SubscriptionRepo.CreateContactSubscription(contact.Contactid, subscriptionid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
