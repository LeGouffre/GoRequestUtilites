package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"unicode"
)

//CheckStruct define the params for checking data
type CheckStruct struct {
	Checkfunc func(string) (bool, string)
	Require   bool
}

//ParseBodyForm try to set the data request on the entity passing by reference
//return nil in success or an error
func ParseBodyForm(w http.ResponseWriter, r *http.Request, entity interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	err2 := json.Unmarshal(body, entity)
	fmt.Println(entity)
	if err2 != nil {
		return err2
	}
	return nil
}

// CheckForm Loop for check data with key => func check
func CheckForm(datas map[string]CheckStruct, body interface{}) (bool, map[string]string) {
	errorTable := make(map[string]string)
	r := reflect.ValueOf(body)
	for key, sFunc := range datas {
		tmp := UcFirst(key)
		f := r.FieldByName(tmp)
		if f.IsValid() == false && sFunc.Require == true {
			errorTable[key] = "Valeur non renseigner"
		} else if sFunc.Checkfunc != nil && f.IsValid() {
			tmp := strings.TrimSpace(f.String())
			res, err := sFunc.Checkfunc(tmp)
			if res == false {
				errorTable[key] = err
			}
		}
	}
	if len(errorTable) > 0 {
		return false, errorTable
	}
	return true, nil
}

//UcFirst func for to upper the first letter
func UcFirst(str1 string) string {
	for i, v := range str1 {
		return string(unicode.ToUpper(v)) + str1[i+1:]
	}
	return ("")
}

//SToS Struct To Struct
//dublicate informations of ent1 to ent2 will be a struct of data
// the duplication opeare with comparaison between the fiels of ent1 with ent2
func SToS(ent1 interface{}, ent2 interface{}) interface{} {
	vEnt1 := reflect.ValueOf(ent1)
	vEnt2 := reflect.ValueOf(ent2).Elem()
	var key string
	for c := 0; c < vEnt1.NumField(); c++ {
		key = vEnt1.Type().Field(c).Name
		field := vEnt2.FieldByName(key)
		if field.IsValid() && field.CanSet() {
			field.Set(vEnt1.Field(c))
		}
	}
	return ent2
}
