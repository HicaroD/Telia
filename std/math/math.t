package math

#[default_cc="c"]
extern common {
  fn cos(x f64) f64
  fn sin(x f64) f64
  fn tan(x f64) f64
  fn acos(x f64) f64
  fn acosf(x f32) f32
  fn asin(x f64) f64
  fn atan(x f64) f64
  fn atan2(y f64, x f64) f64
  fn pow(y f64, x f64) f64
  fn powf(y f32, x f32) f32
}
