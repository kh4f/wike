use serde::Deserialize;
use windows::Win32::{ Foundation::POINT, UI::Input::KeyboardAndMouse::VIRTUAL_KEY };
use crate::SCREEN_SIZE;

pub struct Config {
	pub rules: Vec<Rule>,
}

pub struct Rule {
	pub name: Option<String>,
	pub enabled: bool,
    pub trigger: Trigger,
    pub action: Action,
    pub consume: Option<bool>,
}

pub struct Trigger {
	pub region: ScreenRegion,
    pub mouse: Option<MouseEvent>,
    pub key: Option<VIRTUAL_KEY>,
}

pub struct Action {
    pub keys: Option<Vec<VIRTUAL_KEY>>,
    pub cmd: Option<String>,
    pub open: Option<OpenAction>,
}

pub struct OpenAction {
    pub target: String,
    pub window_class: Option<String>,
}

#[derive(Deserialize, PartialEq, Copy, Clone)]
pub enum MouseEvent {
	LeftUp,
    LeftDown,
	RightUp,
    RightDown,
    MiddleUp,
    MiddleDown,
    WheelUp,
    WheelDown,
}

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