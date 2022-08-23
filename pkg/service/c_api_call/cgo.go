//go:build cgo
// +build cgo

package capicall

/*
#include <string.h>
#include <stdlib.h>
typedef struct kclvm_service kclvm_service;
kclvm_service * kclvm_service_new(void *f,long long plugin_agent){
	kclvm_service * (*new_service)(long long);
	new_service = (kclvm_service * (*)(long long))f;
	return new_service(plugin_agent);
}
void kclvm_service_delete(void *f,kclvm_service * c){
	void (*delete_service)(kclvm_service *);
	delete_service = (void (*)(kclvm_service *))f;
	return delete_service(c);
}
void kclvm_service_free_string(void *f,const char * res) {
	void (*free_string)(const char *);
	free_string = (void (*)(const char *))f;
	return free_string(res);
}
const char* kclvm_service_call(void *f,kclvm_service* c,const char * method,const char * args){
	const char* (*service_call)(kclvm_service*,const char *,const char *);
	service_call = (const char* (*)(kclvm_service*,const char *,const char *))f;
	return service_call(c,method,args);
}
const char* kclvm_service_get_error_buffer(void *f,kclvm_service* c){
	const char* (*get_error_buffer)(kclvm_service*);
	get_error_buffer = (const char* (*)(kclvm_service*))f;
	return get_error_buffer(c);
}
void kclvm_service_clear_error_buffer(void *f,kclvm_service* c){
	void (*clear_error_buffer)(kclvm_service*);
	clear_error_buffer = (void (*)(kclvm_service*))f;
	return clear_error_buffer(c);
}
*/
import "C"

import (
	"kusionstack.io/kclvm-go/pkg/3rdparty/dlopen"
)

var lib *dlopen.LibHandle = loadKclvmServiceCapiLib()

func NewKclvmService(pluginAgent C.longlong) *C.kclvm_service {
	f := "kclvm_service_new"

	newServ, _ := lib.GetSymbolPointer(f)

	return C.kclvm_service_new(newServ, pluginAgent)
}

func DeleteKclvmService(serv *C.kclvm_service) {
	f := "kclvm_service_delete"

	deleteServ, _ := lib.GetSymbolPointer(f)

	C.kclvm_service_delete(deleteServ, serv)
}

func KclvmServiceFreeString(str *C.char) {
	f := "kclvm_service_free_string"

	freeString, _ := lib.GetSymbolPointer(f)

	C.kclvm_service_free_string(freeString, str)
}

func KclvmServiceCall(serv *C.kclvm_service, method *C.char, args *C.char) *C.char {
	f := "kclvm_service_call"

	serviceCall, _ := lib.GetSymbolPointer(f)

	return C.kclvm_service_call(serviceCall, serv, method, args)
}

func GetKclvmServiceError(serv *C.kclvm_service) *C.char {
	f := "kclvm_service_get_error_buffer"

	getError, _ := lib.GetSymbolPointer(f)

	return C.kclvm_service_get_error_buffer(getError, serv)
}

func ClearKclvmServiceError(serv *C.kclvm_service) {
	f := "kclvm_service_clear_error_buffer"

	clearError, _ := lib.GetSymbolPointer(f)

	C.kclvm_service_get_error_buffer(clearError, serv)
}
