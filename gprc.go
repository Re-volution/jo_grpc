package jo_grpc

import (
	"errors"
	"fmt"
	"reflect"
)

var negr = new(grpcStruct)

type grpcStruct struct {
	KIn   map[string][]int
	KOut  map[string][]int
	FName map[string]reflect.Value
}

func init() {
	negr.KIn = make(map[string][]int)
	negr.KOut = make(map[string][]int)
	negr.FName = make(map[string]reflect.Value)
}

// SetGPRC 线程不安全，不可以在运行过程中动态添加
func SetGPRC(strc interface{}) {
	var ty = reflect.TypeOf(strc)
	var val = reflect.ValueOf(strc)
	//fmt.Println(val.NumMethod(), ty.NumMethod())

	for n := 0; n < val.NumMethod(); n++ {
		var (
			method  = ty.Method(n)
			name    = method.Name
			methodT = method.Type
		)

		negr.FName[name] = val.Method(n)
		for i := 1; i < methodT.NumIn(); i++ {
			fmt.Println(methodT.In(i), methodT.In(i).Kind())
			negr.KIn[name] = append(negr.KIn[name], int(methodT.In(i).Kind()))
		}
		for i := 0; i < methodT.NumOut(); i++ {
			fmt.Println(methodT.Out(i), methodT.Out(i).Kind())
			negr.KOut[name] = append(negr.KOut[name], int(methodT.Out(i).Kind()))
		}
	}
}

func Call(name string, data ...interface{}) ([]interface{}, error) {
	if len(data) != len(negr.KIn[name]) {
		return nil, errors.New(fmt.Sprintf("param not enough,need:%d param length:%d", len(negr.KIn[name]), len(data)))
	}
	var temp []reflect.Value
	for _, v := range data {
		temp = append(temp, reflect.ValueOf(v))
	}
	f, ok := negr.FName[name]
	if !ok {
		return nil, errors.New("not found this func:" + name)
	}

	var resi []interface{}
	for _, v := range f.Call(temp) {
		resi = append(resi, v.Interface())
	}
	return resi, nil
}
