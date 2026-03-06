mod config;
mod utils;
mod hooks;

use std::sync::OnceLock;
use windows::Win32::{ Foundation::*, UI::WindowsAndMessaging::* };
use hooks::{ keyboard_proc, mouse_proc };
use config::{ Config, ScreenSize, ScreenRegion };

pub static SCREEN_SIZE: OnceLock<ScreenSize> = OnceLock::new();
pub static CONFIG: OnceLock<Config> = OnceLock::new();

fn main() -> Result<(), Box<dyn std::error::Error>> {
	unsafe {
		SCREEN_SIZE.set(ScreenSize {
			w: GetSystemMetrics(SM_CXSCREEN) as i16,
			h: GetSystemMetrics(SM_CYSCREEN) as i16
		}).ok();
		CONFIG.set(Config {
			volume_scroll_region: ScreenRegion::new(-1, -500, 50, 1000)
		}).ok();

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