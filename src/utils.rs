use windows::Win32::{ Foundation::*, UI::Input::KeyboardAndMouse::* };
use crate::config::ScreenRegion;

pub fn is_inside_region(pt: POINT, area: &ScreenRegion) -> bool {
	pt.x >= (area.x as i32) && pt.x <= (area.x + area.w) as i32 &&
	pt.y >= (area.y as i32) && pt.y <= (area.y + area.h) as i32
}

pub fn create_input(v_key: VIRTUAL_KEY, key_up: bool) -> INPUT {
    INPUT {
        r#type: INPUT_KEYBOARD,
        Anonymous: INPUT_0 {
            ki: KEYBDINPUT {
                wVk: v_key,
                wScan: 0,
                dwFlags: if key_up {
                    KEYEVENTF_KEYUP
                } else {
                    KEYBD_EVENT_FLAGS(0)
                },
                time: 0,
                dwExtraInfo: 0,
            },
        },
    }
}

pub fn press_keys(keys: &[VIRTUAL_KEY]) {
    let mut inputs: Vec<INPUT> = Vec::new();
    inputs.extend(keys.iter().map(|&k| create_input(k, false)));
    inputs.extend(keys.iter().rev().map(|&k| create_input(k, true)));
    unsafe { SendInput(&inputs, std::mem::size_of::<INPUT>() as i32) };
}