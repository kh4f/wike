use windows::Win32::{
    Foundation::*,
    UI::{ WindowsAndMessaging::*, Input::KeyboardAndMouse::* },
};

fn create_input(v_key: VIRTUAL_KEY, key_up: bool) -> INPUT {
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

fn press_keys(keys: &[VIRTUAL_KEY]) {
    let mut inputs: Vec<INPUT> = Vec::new();
    inputs.extend(keys.iter().map(|&k| create_input(k, false)));
    inputs.extend(keys.iter().rev().map(|&k| create_input(k, true)));
    unsafe { SendInput(&inputs, std::mem::size_of::<INPUT>() as i32) };
}

unsafe extern "system" fn keyboard_proc(n_code: i32, w_param: WPARAM, l_param: LPARAM) -> LRESULT {
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

unsafe extern "system" fn mouse_proc(n_code: i32, w_param: WPARAM, l_param: LPARAM) -> LRESULT {
    if n_code >= 0 {
        let info = unsafe { &*(l_param.0 as *const MSLLHOOKSTRUCT) };
        let pt = info.pt;

        match w_param.0 as u32 {
            WM_MOUSEWHEEL => {
                let delta = (info.mouseData >> 16) as i16;
                press_keys(&[if delta > 0 { VK_VOLUME_UP } else { VK_VOLUME_DOWN }]);
				return LRESULT(1);
            }
            _ => (),
        }
    }
    unsafe { CallNextHookEx(None, n_code, w_param, l_param) }
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
	unsafe {
        let mouse_hook = SetWindowsHookExW(WH_MOUSE_LL, Some(mouse_proc), Some(HINSTANCE::default()), 0)?;
        let keyboard_hook = SetWindowsHookExW(WH_KEYBOARD_LL, Some(keyboard_proc), Some(HINSTANCE::default()), 0)?;

        let mut msg = MSG::default();
        while GetMessageW(&mut msg, None, 0, 0).into() {
            _ = TranslateMessage(&msg);
            DispatchMessageW(&msg);
        }

        UnhookWindowsHookEx(mouse_hook)?;
        UnhookWindowsHookEx(keyboard_hook)?;
        Ok(())
    }
}