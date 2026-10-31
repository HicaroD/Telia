package sdl

use "vendor::sdl::raw"
use "std::strings"

type ThreadEntry = raw::NativeThreadEntry

fn create_thread(entry ThreadEntry, name string, data rawptr) (*Thread, error) {
  converted, err := strings::to_c(name)
  if err != nil {
    return nil, err
  }
  defer strings::free_c(&converted)
  thread := raw::create_native_thread(entry, strings::c_view(&converted), data)
  if thread == nil {
    return nil, error("SDL create_thread failed")
  }
  return thread, nil
}

fn wait_thread(thread *Thread) i32 {
  return raw::wait_native_thread(thread)
}
