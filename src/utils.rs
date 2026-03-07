use windows::{
	core::{HSTRING, w},
    Win32::UI::{
        Input::KeyboardAndMouse::{INPUT, INPUT_0, INPUT_KEYBOARD, KEYBD_EVENT_FLAGS, KEYBDINPUT, KEYEVENTF_KEYUP,SendInput, VIRTUAL_KEY},
        Shell::ShellExecuteW,
        WindowsAndMessaging::{FindWindowW, IsWindow, SW_RESTORE, SW_SHOW, SetForegroundWindow, ShowWindow},
    },
};
use crate::config::OpenAction;

pub fn open_or_focus_app(open_action: &OpenAction) {
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

pub fn press_keys(keys: &[VIRTUAL_KEY]) {
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
