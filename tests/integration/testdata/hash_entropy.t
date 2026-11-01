package main
use "std::errors"
@[header="hash_entropy.h", default_cc="c"]
extern {
 fn hash_entropy_fails() i32
}
fn acquire() (u64[2], error) {
 seed := __key_seed() @try
 return seed, nil
}
fn main() i32 {
 seed, err := acquire()
 if hash_entropy_fails() != 0 {
  if err == nil or err.kind != errors::Error::RandomnessFailed {
   return 1
  }
  if seed[0] != 0 or seed[1] != 0 {
   return 2
  }
 } else {
  if err != nil {
   return 3
  }
  if seed[0] != 506097522914230528 or seed[1] != 1084818905618843912 {
   return 4
  }
 }
 return 0
}
