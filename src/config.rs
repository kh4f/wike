use windows::Win32::Foundation::POINT;
use crate::SCREEN_SIZE;

pub struct Config { pub volume_scroll_region: ScreenRegion }

pub struct ScreenSize { pub w: i16, pub h: i16 }

pub struct ScreenRegion { pub x: i16, pub y: i16, pub w: i16, pub h: i16 }

impl ScreenRegion {
	pub fn new(mut x: i16, mut y: i16, w: i16, h: i16) -> Self {
		if x < 0 { x += SCREEN_SIZE.get().unwrap().w; }
		if y < 0 { y += SCREEN_SIZE.get().unwrap().h; }
        Self { x, y, w, h }
    }

    pub fn contains(&self, pt: POINT) -> bool {
        pt.x >= self.x as i32 && pt.x <= (self.x + self.w) as i32 &&
        pt.y >= self.y as i32 && pt.y <= (self.y + self.h) as i32
    }
}