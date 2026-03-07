mod config;
mod utils;
mod hooks;

use std::sync::OnceLock;
use windows::Win32::{ Foundation::*, UI::{ Input::KeyboardAndMouse::{ VK_VOLUME_DOWN, VK_VOLUME_UP }, WindowsAndMessaging::* } };
use hooks::{ keyboard_proc, mouse_proc };
use config::{ Config, ScreenSize, ScreenRegion, Rule, Trigger, Action, MouseEvent };

pub static SCREEN_SIZE: OnceLock<ScreenSize> = OnceLock::new();
pub static CONFIG: OnceLock<Config> = OnceLock::new();

fn main() -> Result<(), Box<dyn std::error::Error>> {
	unsafe {
		SCREEN_SIZE.set(ScreenSize {
			w: GetSystemMetrics(SM_CXSCREEN) as i16,
			h: GetSystemMetrics(SM_CYSCREEN) as i16
		}).ok();
		CONFIG.set(Config {
			rules: vec![
				Rule {
					name: Some("Volume Scroll (up)".to_string()),
					enabled: true,
					trigger: Trigger {
						region: ScreenRegion::new(-1, -500, 2, 500),
						mouse: Some(MouseEvent::WheelUp),
						kb: None,
					},
					action: Action {
						send_keys: Some(vec![VK_VOLUME_UP]),
						run_cmd: None,
					},
					consume: Some(true),
				},
				Rule {
					name: Some("Volume Scroll (down)".to_string()),
					enabled: true,
					trigger: Trigger {
						region: ScreenRegion::new(-1, -500, 2, 500),
						mouse: Some(MouseEvent::WheelDown),
						kb: None,
					},
					action: Action {
						send_keys: Some(vec![VK_VOLUME_DOWN]),
						run_cmd: None,
					},
					consume: Some(true),
				}
			]
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