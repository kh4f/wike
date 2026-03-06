use windows::Win32::Foundation::POINT;

pub struct ScreenRegion { pub x: i16, pub y: i16, pub w: i16, pub h: i16 }

impl ScreenRegion {
    pub fn contains(&self, pt: POINT) -> bool {
        pt.x >= self.x as i32 && pt.x <= (self.x + self.w) as i32 &&
        pt.y >= self.y as i32 && pt.y <= (self.y + self.h) as i32
    }
}

pub const VOLUME_SCROLL_REGION: ScreenRegion = ScreenRegion { x: 1917, y: 600, w: 50, h: 1000 };