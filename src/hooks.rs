use windows::Win32::{
	Foundation::*,
	UI::{ Input::KeyboardAndMouse::*, WindowsAndMessaging::* }
};
use crate::{
	config::VOLUME_SCROLL_REGION,
	utils::{ is_inside_region, press_keys }
};

pub unsafe extern "system" fn keyboard_proc(n_code: i32, w_param: WPARAM, l_param: LPARAM) -> LRESULT {
    if n_code >= 0 {
        let info = unsafe { &*(l_param.0 as *const KBDLLHOOKSTRUCT) };
        let vk_code = info.vkCode;

        if w_param.0 as u32 == WM_KEYDOWN {
            println!("Key down: {}", vk_code);
        } else if w_param.0 as u32 == WM_KEYUP {
            println!("Key up: {}", vk_code);
        }
    }
    unsafe { CallNextHookEx(None, n_code, w_param, l_param) }
}

pub unsafe extern "system" fn mouse_proc(n_code: i32, w_param: WPARAM, l_param: LPARAM) -> LRESULT {
    if n_code >= 0 {
        let info = unsafe { &*(l_param.0 as *const MSLLHOOKSTRUCT) };
        let pt = info.pt;

        match w_param.0 as u32 {
            WM_MOUSEWHEEL => {
                let delta = (info.mouseData >> 16) as i16;
				if is_inside_region(pt, &VOLUME_SCROLL_REGION) {
					press_keys(&[if delta > 0 { VK_VOLUME_UP } else { VK_VOLUME_DOWN }]);
				}
				return LRESULT(1)
            }
            _ => (),
        }
    }
    unsafe { CallNextHookEx(None, n_code, w_param, l_param) }
}