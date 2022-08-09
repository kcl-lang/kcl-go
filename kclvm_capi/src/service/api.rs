use protobuf::Message;

use crate::model::gpyrpc::*;
use crate::service::service::KclvmService;
use kclvm::utils::*;
use std::ffi::CString;
use std::os::raw::c_char;

#[allow(non_camel_case_types)]
type kclvm_service = KclvmService;

// Create an instance of KclvmService and return its pointer
#[no_mangle]
pub extern "C" fn kclvm_service_new(plugin_agent: u64) -> *mut kclvm_service {
    let mut serv = KclvmService::default();
    serv.plugin_agent = plugin_agent;
    Box::into_raw(Box::new(serv))
}

// Delete KclvmService
#[no_mangle]
pub extern "C" fn kclvm_service_delete(serv: *mut kclvm_service) {
    free_mut_ptr(serv);
}

// Free memory for string returned to the outside
#[no_mangle]
pub extern "C" fn kclvm_service_free_string(res: *mut c_char) {
    if !res.is_null() {
        unsafe { CString::from_raw(res) };
    }
}

// Provide the error information of the last call to the outside
#[no_mangle]
pub extern "C" fn kclvm_service_get_error_buffer(serv: *mut kclvm_service) -> *const c_char {
    let serv = mut_ptr_as_ref(serv);
    serv.kclvm_service_err_buffer.as_ptr() as *const i8
}

//Clear the error infomation buffer，so that external programs do not get too old error messages
#[no_mangle]
pub extern "C" fn kclvm_service_clear_error_buffer(serv: *mut kclvm_service) {
    let serv = mut_ptr_as_ref(serv);
    serv.kclvm_service_err_buffer = "\0".to_string();
}

/// Call kclvm service by C API
///
/// # Parameters
///
/// `serv`: [*mut kclvm_service]
///     The pointer of &\[[KclvmService]]
///
/// `call`: [*const c_char]
///     The C str of the name of the called service,
///     with the format "KclvmService.{MethodName}"
///
/// `args`: [*const c_char]
///     Arguments of the call serialized as protobuf byte sequence,
///     refer to kclvm/api/src/gpyrpc.proto for the specific definitions of arguments
///
/// # Returns
///
/// result: [*const c_char]
///     Result of the call serialized as protobuf byte sequence
#[no_mangle]
pub extern "C" fn kclvm_service_call(
    serv: *mut kclvm_service,
    call: *const c_char,
    args: *const c_char,
) -> *const c_char {
    let result = std::panic::catch_unwind(|| {
        let args = unsafe { std::ffi::CStr::from_ptr(args) }.to_bytes();
        let call = c2str(call);
        let call = _kclvm_get_service_fn_ptr_by_name(call);
        if call == 0 {
            panic!("null fn ptr");
        }
        let call = (&call as *const u64) as *const ()
            as *const fn(serv: *mut KclvmService, args: &[u8]) -> *const c_char;
        unsafe { (*call)(serv, args) }
    });

    match result {
        //todo uniform error handling
        Ok(result) => result,
        Err(panic_err) => {
            let mut err_message = if let Some(s) = panic_err.downcast_ref::<&str>() {
                s.to_string()
            } else if let Some(s) = panic_err.downcast_ref::<&String>() {
                (*s).clone()
            } else if let Some(s) = panic_err.downcast_ref::<String>() {
                (*s).clone()
            } else {
                "".to_string()
            };
            let serv_ref = mut_ptr_as_ref(serv);
            serv_ref.kclvm_service_err_buffer = err_message.clone();
            let c_string =
                std::ffi::CString::new(err_message.as_str()).expect("CString::new failed");
            let ptr = c_string.into_raw();
            ptr as *const i8
        }
    }
}

pub fn _kclvm_get_service_fn_ptr_by_name(name: &str) -> u64 {
    match name {
        "KclvmService.Ping" => ping as *const () as u64,
        "KclvmService.ExecProgram" => exec_program as *const () as u64,
        _ => panic!("unknown method name : {}", name),
    }
}

/// ping is used to test whether kclvm service is successfully imported
/// arguments and return results should be consistent
pub fn ping(serv: *mut KclvmService, args: &[u8]) -> *const c_char {
    let serv_ref = mut_ptr_as_ref(serv);
    let args = Ping_Args::parse_from_bytes(args).unwrap();
    let res = serv_ref.ping(&args);
    CString::new(res.write_to_bytes().unwrap())
        .unwrap()
        .into_raw()
}

/// exec_program provides users with the ability to execute KCL code
///
/// # Parameters
///
/// `serv`: [*mut kclvm_service]
///     The pointer of &\[[KclvmService]]
///
///
/// `args`: [&[u8]]
///     the items and compile parameters selected by the user in the KCLVM CLI
///     serialized as protobuf byte sequence
///
/// # Returns
///
/// result: [*const c_char]
///     Result of the call serialized as protobuf byte sequence
pub fn exec_program(serv: &mut KclvmService, args: &[u8]) -> *const c_char {
    let serv_ref = mut_ptr_as_ref(serv);
    let args = ExecProgram_Args::parse_from_bytes(args).unwrap();
    let res = serv_ref.exec_program(&args);
    let result_byte = match res {
        Ok(res) => match res.write_to_bytes() {
            Ok(bytes) => bytes,
            Err(err) => panic!("{}", err.to_string()),
        },
        Err(err) => panic!("{}", err.clone()),
    };
    CString::new(result_byte).unwrap().into_raw()
}
