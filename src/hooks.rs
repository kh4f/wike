use std::process::Command;
use windows:: Win32::{ Foundation::*, UI:: WindowsAndMessaging::* };
use crate::{ CONFIG, config::{ Action, MouseEvent }, utils::{ open_or_focus_app, press_keys } };

pub unsafe extern "system" fn keyboard_proc(n_code: i32, w_param: WPARAM, l_param: LPARAM) -> LRESULT {
    if n_code >= 0 {
        let info = unsafe { &*(l_param.0 as *const KBDLLHOOKSTRUCT) };
        let mut pt = POINT::default();
        unsafe { GetCursorPos(&mut pt).ok(); };

        if handle_keyboard_event(info.vkCode, pt) { return LRESULT(1) }
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

		if let Some(event) = mouse_event && handle_mouse_event(event, pt) { return LRESULT(1) }
	}
	unsafe { CallNextHookEx(None, n_code, w_param, l_param) }
}

fn handle_mouse_event(mouse_event: MouseEvent, pt: POINT) -> bool {
    if let Some(cfg) = CONFIG.get() {
        for rule in &cfg.rules {
            if rule.enabled
                && let Some(event) = rule.trigger.mouse
                && event == mouse_event
                && rule.trigger.region.contains(pt)
            {
                execute_action(&rule.action);
                return rule.consume.unwrap_or(false)
            }
        }
    }
    false
}

fn handle_keyboard_event(vk_code: u32, pt: POINT) -> bool {
    if let Some(cfg) = CONFIG.get() {
        for rule in &cfg.rules {
            if rule.enabled
                && let Some(key) = rule.trigger.key
                && key.0 as u32 == vk_code
                && rule.trigger.region.contains(pt)
            {
                execute_action(&rule.action);
                return rule.consume.unwrap_or(false)
            }
        }
    }
    false
}

fn execute_action(action: &Action) {
    if let Some(keys) = &action.keys { press_keys(keys); }
    if let Some(cmd) = &action.cmd { Command::new(cmd).spawn().ok(); }
    if let Some(open) = &action.open { open_or_focus_app(open); }
}