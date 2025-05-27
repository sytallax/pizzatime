package dominos

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ApiEndpoint int

const (
	ApiEndpointStoreLocator ApiEndpoint = iota
	ApiEndpointStoreMenu
)

var urls = map[ApiEndpoint]string{
	ApiEndpointStoreLocator: "https://order.dominos.com/power/store-locator?s=%s&c=%s&type=delivery",
	ApiEndpointStoreMenu:    "https://order.dominos.com/power/store/%s/menu?lang=en&structured=true",
}

func (e ApiEndpoint) String() string {
	return urls[e]
}

type Customer struct {
	Name    string
	Address Address
}

type Order struct {
	Customer Customer
}

type Address struct {
	Street     string `json:"Street"`
	City       string `json:"City"`
	Region     string `json:"Region"`
	PostalCode int    `json:"PostalCode"`
}

func (a *Address) LineOne() string {
	return a.Street
}
func (a *Address) LineTwo() string {
	return strings.Join([]string{a.City, a.Region, strconv.Itoa(a.PostalCode)}, " ")
}

type Store struct {
	StoreID                           string `json:"StoreID"`
	IsOnlineCapable                   bool   `json:"IsOnlineCapable"`
	IsOnlineNow                       bool   `json:"IsOnlineNow"`
	AllowDeliveryOrders               bool   `json:"AllowDeliveryOrders"`
	ServiceMethodEstimatedWaitMinutes struct {
		Delivery struct {
			Min int `json:"Min"`
			Max int `json:"Max"`
		} `json:"Delivery"`
	} `json:"ServiceMethodEstimatedWaitMinutes"`
	ContactlessDelivery string `json:"ContactlessDelivery"`
	IsOpen              bool   `json:"IsOpen"`
	ServiceIsOpen       struct {
		Delivery bool `json:"Delivery"`
	} `json:"ServiceIsOpen"`
	Address Address `json:"Address"`
}

type storeLocatorResponse struct {
	Stores []Store `json:"Stores"`
}

func runStoreLocator(l1, l2 string) (storeLocatorResponse, error) {
	endpoint := fmt.Sprintf(ApiEndpointStoreLocator.String(), url.PathEscape(l1), url.PathEscape(l2))
	resp, err := http.Get(endpoint)
	if err != nil {
		return storeLocatorResponse{}, err
	}
	defer resp.Body.Close()

	var r storeLocatorResponse
	json.NewDecoder(resp.Body).Decode(&r)
	return r, nil
}

func (a *Address) GetNearestStore() (Store, error) {
	resp, err := runStoreLocator(a.LineOne(), a.LineTwo())
	if err != nil {
		return Store{}, err
	}

	for _, s := range resp.Stores {
		if s.IsOnlineCapable && s.IsOnlineNow && s.AllowDeliveryOrders && s.IsOpen && s.ServiceIsOpen.Delivery {
			return s, nil
		}
	}

	return Store{}, errors.New("Could not find open store")
}

type Menu struct {
	Categorization struct {
		Food struct {
			Categories []struct {
				Code        string   `json:"Code"`
				Description string   `json:"Description"`
				Products    []string `json:"Products"`
				Name        string   `json:"Name"`
			} `json:"Categories"`
			Code        string `json:"Code"`
			Description string `json:"Description"`
			Products    []any  `json:"Products"`
			Name        string `json:"Name"`
		} `json:"Food"`
	} `json:"Categorization"`
	Coupons map[string]struct {
		Code  string `json:"Code"`
		Name  string `json:"Name"`
		Price string `json:"Price"`
	} `json:"Coupons"`
	Flavors map[string]struct {
		BreadDipCombos map[string]struct {
			Code        string `json:"Code"`
			Description string `json:"Description"`
			Name        string `json:"Name"`
		} `json:"BreadDipCombos"`
		Pasta map[string]struct {
			Code        string `json:"Code"`
			Description string `json:"Description"`
			Name        string `json:"Name"`
		} `json:"Pasta"`
		Pizza map[string]struct {
			Code        string `json:"Code"`
			Description string `json:"Description"`
			Name        string `json:"Name"`
		} `json:"Pizza"`
		Wings map[string]struct {
			Code        string `json:"Code"`
			Description string `json:"Description"`
			Name        string `json:"Name"`
		} `json:"Wings"`
	} `json:"Flavors"`
	Products map[string]struct {
		Code        string   `json:"Code"`
		Name        string   `json:"Name"`
		ProductType string   `json:"ProductType"`
		Variants    []string `json:"Variants"`
	} `json:"Products"`
	Sides map[string]struct {
		Bread map[string]struct {
			Code        string `json:"Code"`
			Description string `json:"Description"`
			Name        string `json:"Name"`
		} `json:"Bread"`
		Dessert map[string]struct {
			Code        string `json:"Code"`
			Description string `json:"Description"`
			Name        string `json:"Name"`
		} `json:"Dessert"`
		GSalad map[string]struct {
			Code        string `json:"Code"`
			Description string `json:"Description"`
			Name        string `json:"Name"`
		} `json:"GSalad"`
		Tots map[string]struct {
			Code        string `json:"Code"`
			Description string `json:"Description"`
			Name        string `json:"Name"`
		} `json:"Tots"`
		Wings map[string]struct {
			Code        string `json:"Code"`
			Description string `json:"Description"`
			Name        string `json:"Name"`
		} `json:"Wings"`
	} `json:"Sides"`
	Sizes map[string]struct {
		Bread map[string]struct {
			Code string `json:"Code"`
			Name string `json:"Name"`
		} `json:"Bread"`
		Dessert map[string]struct {
			Code string `json:"Code"`
			Name string `json:"Name"`
		} `json:"Dessert"`
		Drinks map[string]struct {
			Code string `json:"Code"`
			Name string `json:"Name"`
		} `json:"Drinks"`
		Pizza map[string]struct {
			Code string `json:"Code"`
			Name string `json:"Name"`
		} `json:"Pizza"`
		Tots map[string]struct {
			Code string `json:"Code"`
			Name string `json:"Name"`
		} `json:"Tots"`
		Wings map[string]struct {
			Code string `json:"Code"`
			Name string `json:"Name"`
		} `json:"Wings"`
	} `json:"Sizes"`
	Toppings map[string]struct {
		Bread map[string]struct {
			Code string `json:"Code"`
			Name string `json:"Name"`
		} `json:"Bread"`
		Pasta map[string]struct {
			Code string `json:"Code"`
			Name string `json:"Name"`
		} `json:"Pasta"`
		Pizza map[string]struct {
			Code string `json:"Code"`
			Name string `json:"Name"`
		} `json:"Pizza"`
		Sandwich map[string]struct {
			Code string `json:"Code"`
			Name string `json:"Name"`
		} `json:"Sandwich"`
		Tots map[string]struct {
			Code string `json:"Code"`
			Name string `json:"Name"`
		} `json:"Tots"`
		Wings map[string]struct {
			Code string `json:"Code"`
			Name string `json:"Name"`
		} `json:"Wings"`
	} `json:"Toppings"`
	Variants map[string]struct {
		Code        string `json:"Code"`
		Name        string `json:"Name"`
		Price       string `json:"Price"`
		ProductCode string `json:"ProductCode"`
	} `json:"Variants"`
}

func (m *Menu) CouponCodes() []string {
	var keys []string
	for k := range m.Coupons {
		keys = append(keys, k)
	}
	return keys
}
func (m *Menu) ProductCodes() []string {
	var keys []string
	for k := range m.Products {
		keys = append(keys, k)
	}
	return keys
}
func (m *Menu) VariantCodes() []string {
	var keys []string
	for k := range m.Variants {
		keys = append(keys, k)
	}
	return keys
}

func (s *Store) GetMenu() (Menu, error) {
	sID, err := strconv.Atoi(s.StoreID)
	if err != nil {
		return Menu{}, err
	}

	endpoint := fmt.Sprintf(ApiEndpointStoreMenu.String(), strconv.Itoa(sID))
	resp, err := http.Get(endpoint)
	if err != nil {
		return Menu{}, err
	}
	defer resp.Body.Close()

	var m Menu
	json.NewDecoder(resp.Body).Decode(&m)
	return m, nil
}
