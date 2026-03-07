use windows::Win32::{ Foundation::*, UI::WindowsAndMessaging::* };
use crate::{ CONFIG, config::MouseEvent, utils::press_keys };

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

		let mouse_event = match w_param.0 as u32 {
			WM_LBUTTONDOWN => Some(MouseEvent::LeftDown),
			WM_LBUTTONUP => Some(MouseEvent::LeftUp),
			WM_RBUTTONDOWN => Some(MouseEvent::RightDown),
			WM_RBUTTONUP => Some(MouseEvent::RightUp),
			WM_MBUTTONDOWN => Some(MouseEvent::MiddleDown),
			WM_MBUTTONUP => Some(MouseEvent::MiddleUp),
			WM_MOUSEWHEEL => Some(if info.mouseData as i32 >> 16 > 0
				{ MouseEvent::WheelUp } else { MouseEvent::WheelDown }),
			_ => None
		};

		if let Some(cfg) = CONFIG.get() {
			for rule in &cfg.rules {
				if rule.enabled
					&& let Some(event) = mouse_event
					&& rule.trigger.mouse.as_ref() == Some(&event)
					&& rule.trigger.region.contains(pt)
					&& let Some(keys) = &rule.action.send_keys {
					press_keys(keys);
					if rule.consume.unwrap_or(false) { return LRESULT(1) }
				}
			}
		}
    }
    unsafe { CallNextHookEx(None, n_code, w_param, l_param) }
}