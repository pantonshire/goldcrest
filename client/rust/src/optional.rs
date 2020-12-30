use crate::twitter1;

impl From<u64> for twitter1::OptFixed64 {
    fn from(x: u64) -> Self {
        twitter1::OptFixed64{val: x}
    }
}

impl From<twitter1::OptFixed64> for u64 {
    fn from(x: twitter1::OptFixed64) -> Self {
        x.val
    }
}

impl From<String> for twitter1::OptString {
    fn from(x: String) -> Self {
        twitter1::OptString{val: x}
    }
}

impl From<twitter1::OptString> for String {
    fn from(x: twitter1::OptString) -> Self {
        x.val
    }
}
