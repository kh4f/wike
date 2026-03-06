use windows::Win32::{
    Foundation::{HINSTANCE, LPARAM, LRESULT, WPARAM},
    UI::WindowsAndMessaging::*,
};

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

        if w_param.0 as u32 == WM_LBUTTONDOWN {
            println!("Left button down at ({}, {})", pt.x, pt.y);
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