use std::process::Command;
use windows::{
	core::{HSTRING, w},
	Win32::{
		Foundation::*,
		UI::{WindowsAndMessaging::*, Input::KeyboardAndMouse::*, Shell::ShellExecuteW},
	},
};
use crate::{ CONFIG, config::{ Action, MouseEvent, OpenAction } };

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

		if let Some(event) = get_mouse_event(w_param.0 as u32, info.mouseData)
			&& handle_mouse_event(event, pt)
		{ return LRESULT(1) }
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

fn get_mouse_event(msg: u32, mouse_data: u32) -> Option<MouseEvent> {
    match msg {
        WM_LBUTTONDOWN => Some(MouseEvent::LeftDown),
        WM_LBUTTONUP => Some(MouseEvent::LeftUp),
        WM_RBUTTONDOWN => Some(MouseEvent::RightDown),
        WM_RBUTTONUP => Some(MouseEvent::RightUp),
        WM_MBUTTONDOWN => Some(MouseEvent::MiddleDown),
        WM_MBUTTONUP => Some(MouseEvent::MiddleUp),
        WM_MOUSEWHEEL => Some(if mouse_data as i32 >> 16 > 0
            { MouseEvent::WheelUp } else { MouseEvent::WheelDown }),
        _ => None
    }
}

fn open_or_focus_app(open_action: &OpenAction) {
    unsafe {
        if let Some(w_class) = &open_action.window_class
            && let Ok(hwnd) = FindWindowW(&HSTRING::from(w_class), None)
            && IsWindow(Some(hwnd)).as_bool()
        {
            ShowWindow(hwnd, SW_RESTORE);
            SetForegroundWindow(hwnd);
            return;
        }
        ShellExecuteW(None, w!("open"), &HSTRING::from(&open_action.target), None, None, SW_SHOW);
    }
}

fn press_keys(keys: &[VIRTUAL_KEY]) {
    let mut inputs: Vec<INPUT> = Vec::new();
    inputs.extend(keys.iter().map(|&k| create_input(k, false)));
    inputs.extend(keys.iter().rev().map(|&k| create_input(k, true)));
    unsafe { SendInput(&inputs, std::mem::size_of::<INPUT>() as i32) };
}

fn create_input(v_key: VIRTUAL_KEY, key_up: bool) -> INPUT {
    INPUT {
        r#type: INPUT_KEYBOARD,
        Anonymous: INPUT_0 {
            ki: KEYBDINPUT {
                wVk: v_key,
                wScan: 0,
                dwFlags: if key_up { KEYEVENTF_KEYUP } else { KEYBD_EVENT_FLAGS(0) },
                time: 0,
                dwExtraInfo: 0,
            },
        },
    }
}